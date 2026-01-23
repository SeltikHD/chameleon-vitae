// Package groq provides an AI adapter using the Groq API with Llama models.
package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
	"github.com/SeltikHD/chameleon-vitae/internal/core/ports"
)

const (
	baseURL          = "https://api.groq.com/openai/v1"
	defaultMaxTokens = 4096
)

// Config holds Groq API configuration.
type Config struct {
	// APIKey is the Groq API key.
	APIKey string

	// ModelGeneration is the model used for content generation (summary, tailoring).
	ModelGeneration string

	// ModelAnalysis is the model used for analysis tasks (job analysis, scoring).
	ModelAnalysis string

	// MaxRetries is the maximum number of retries on rate limit errors.
	MaxRetries int

	// Timeout is the HTTP request timeout.
	Timeout time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		ModelGeneration: "llama-3.3-70b-versatile",
		ModelAnalysis:   "meta-llama/llama-4-scout-17b-16e-instruct",
		MaxRetries:      3,
		Timeout:         60 * time.Second,
	}
}

// Client implements ports.AIProvider using the Groq API.
type Client struct {
	config     Config
	httpClient *http.Client
}

// New creates a new Groq API client.
func New(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("groq: API key is required")
	}

	if cfg.ModelGeneration == "" {
		cfg.ModelGeneration = DefaultConfig().ModelGeneration
	}
	if cfg.ModelAnalysis == "" {
		cfg.ModelAnalysis = DefaultConfig().ModelAnalysis
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = DefaultConfig().MaxRetries
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultConfig().Timeout
	}

	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}, nil
}

// AnalyzeJob analyzes a job description and extracts key requirements.
func (c *Client) AnalyzeJob(ctx context.Context, req ports.AnalyzeJobRequest) (*ports.JobAnalysis, error) {
	prompt := fmt.Sprintf(`Analyze the following job description and extract key information.

Job Description:
%s

Provide a JSON response with the following structure:
{
  "title": "extracted job title",
  "company": "company name if found",
  "required_skills": ["list", "of", "required", "skills"],
  "preferred_skills": ["list", "of", "nice-to-have", "skills"],
  "keywords": ["important", "keywords", "from", "description"],
  "seniority_level": "junior/mid/senior/lead/executive",
  "years_experience": null or number,
  "summary": "brief 2-3 sentence summary of the role"
}

Respond ONLY with valid JSON, no additional text.`, req.JobDescription)

	response, err := c.chatCompletion(ctx, c.config.ModelAnalysis, prompt, 0.3)
	if err != nil {
		return nil, fmt.Errorf("groq: analyze job failed: %w", err)
	}

	var result struct {
		Title           string   `json:"title"`
		Company         string   `json:"company"`
		RequiredSkills  []string `json:"required_skills"`
		PreferredSkills []string `json:"preferred_skills"`
		Keywords        []string `json:"keywords"`
		SeniorityLevel  string   `json:"seniority_level"`
		YearsExperience *int     `json:"years_experience"`
		Summary         string   `json:"summary"`
	}

	if err := json.Unmarshal([]byte(cleanJSON(response)), &result); err != nil {
		log.Printf("json to parse: %s", cleanJSON(response))
		return nil, fmt.Errorf("groq: failed to parse job analysis: %w", err)
	}

	return &ports.JobAnalysis{
		Title:           result.Title,
		Company:         result.Company,
		RequiredSkills:  result.RequiredSkills,
		PreferredSkills: result.PreferredSkills,
		Keywords:        result.Keywords,
		SeniorityLevel:  result.SeniorityLevel,
		YearsExperience: result.YearsExperience,
		Summary:         result.Summary,
	}, nil
}

// SelectBullets selects the most relevant bullets for a job description.
func (c *Client) SelectBullets(ctx context.Context, req ports.SelectBulletsRequest) (*ports.BulletSelection, error) {
	// Build bullets list for prompt.
	var bulletsText strings.Builder
	for i, bullet := range req.AvailableBullets {
		fmt.Fprintf(&bulletsText, "%d. [ID: %s] %s\n", i+1, bullet.ID, bullet.Content)
	}

	prompt := fmt.Sprintf(`You are an expert resume consultant. Select the most relevant experience bullets for this job.

JOB REQUIREMENTS:
- Title: %s
- Company: %s
- Required Skills: %s
- Preferred Skills: %s
- Keywords: %s
- Summary: %s

AVAILABLE BULLETS:
%s

Select up to %d bullets that best match this job. Prioritize:
1. Direct skill matches
2. Quantifiable achievements
3. Relevant industry experience
4. Leadership/impact indicators

Respond with JSON:
{
  "selected_bullet_ids": ["id1", "id2", ...],
  "reasoning": "Brief explanation of selection strategy"
}

Respond ONLY with valid JSON.`,
		req.JobAnalysis.Title,
		req.JobAnalysis.Company,
		strings.Join(req.JobAnalysis.RequiredSkills, ", "),
		strings.Join(req.JobAnalysis.PreferredSkills, ", "),
		strings.Join(req.JobAnalysis.Keywords, ", "),
		req.JobAnalysis.Summary,
		bulletsText.String(),
		req.MaxBullets,
	)

	response, err := c.chatCompletion(ctx, c.config.ModelAnalysis, prompt, 0.3)
	if err != nil {
		return nil, fmt.Errorf("groq: select bullets failed: %w", err)
	}

	var result struct {
		SelectedBulletIDs []string `json:"selected_bullet_ids"`
		Reasoning         string   `json:"reasoning"`
	}

	if err := json.Unmarshal([]byte(cleanJSON(response)), &result); err != nil {
		log.Printf("json to parse: %s", cleanJSON(response))
		return nil, fmt.Errorf("groq: failed to parse bullet selection: %w", err)
	}

	return &ports.BulletSelection{
		SelectedBulletIDs: result.SelectedBulletIDs,
		Reasoning:         result.Reasoning,
	}, nil
}

// TailorBullet rewrites a bullet to better match job requirements.
func (c *Client) TailorBullet(ctx context.Context, req ports.TailorBulletRequest) (*ports.TailoredBulletResult, error) {
	prompt := fmt.Sprintf(`You are an expert Resume Writer and STAR Method Specialist. Your task is to optimize a specific experience bullet point.

ORIGINAL BULLET:
%s

TARGET CONTEXT:
- Job Title: %s
- Required Skills: %s
- Keywords: %s

TASK INSTRUCTIONS:
1. **Analyze & Polish:** First, check the original bullet for grammar and clarity. Fix any errors.
2. **STAR Method Check:** Does the bullet follow the STAR method (Situation, Task, **Action**, **Result**)?
   - *If YES (it has a clear action and quantifiable result):* Keep the structure close to the original. Do not rewrite unnecessary parts.
   - *If NO (it is vague, e.g., "Worked on API"):* Rewrite it to include a specific **Action** and a measurable **Result** (e.g., "Architected a REST API handling **10k requests/sec**").
3. **Keyword Integration:** Naturally weave in the provided keywords if they fit the context.
4. **Style:** Write strictly in %s.

SMART BOLDING (CRITICAL):
Apply **bold** markdown syntax to specific high-value terms. Use bolding for:
- **Hard Skills/Tech Stack:** (e.g., **Go**, **PostgreSQL**, **Docker**)
- **Quantifiable Metrics:** (e.g., **30%% reduction**, **500ms**, **$1M revenue**)
- **Strong Action Verbs:** (e.g., **Orchestrated**, **Deployed**, **Optimized**)
*Constraint:* Limit to 3-5 bolded terms per bullet to ensure readability.

Response format (JSON ONLY):
{
  "tailored_content": "The optimized bullet string with **markdown** formatting",
  "keywords": ["list", "of", "keywords", "used"]
}`,
		req.Bullet.Content,
		req.JobAnalysis.Title,
		strings.Join(req.JobAnalysis.RequiredSkills, ", "),
		strings.Join(req.JobAnalysis.Keywords, ", "),
		req.Style,
	)

	response, err := c.chatCompletion(ctx, c.config.ModelGeneration, prompt, 0.7)
	if err != nil {
		return nil, fmt.Errorf("groq: tailor bullet failed: %w", err)
	}

	var result struct {
		TailoredContent string   `json:"tailored_content"`
		Keywords        []string `json:"keywords"`
	}

	if err := json.Unmarshal([]byte(cleanJSON(response)), &result); err != nil {
		log.Printf("json to parse: %s", cleanJSON(response))
		return nil, fmt.Errorf("groq: failed to parse tailored bullet: %w", err)
	}

	return &ports.TailoredBulletResult{
		OriginalID:      req.Bullet.ID,
		TailoredContent: result.TailoredContent,
		Keywords:        result.Keywords,
	}, nil
}

// GenerateSummary generates a professional summary tailored to the job.
func (c *Client) GenerateSummary(ctx context.Context, req ports.GenerateSummaryRequest) (*ports.SummaryResult, error) {
	// Build context from user and bullets.
	userName := "Professional"
	if req.User.Name != nil {
		userName = *req.User.Name
	}

	var bulletsContext strings.Builder
	for _, bullet := range req.SelectedBullets {
		fmt.Fprintf(&bulletsContext, "- %s\n", bullet.Content)
	}

	prompt := fmt.Sprintf(`Generate a professional summary for a resume application.

CANDIDATE INFO:
- Name: %s
- Headline: %s
- Current Summary: %s

KEY ACHIEVEMENTS (selected for this job):
%s

TARGET JOB:
- Title: %s
- Company: %s
- Required Skills: %s
- Summary: %s

Write a compelling 3-4 sentence professional summary that:
1. Highlights relevant experience and skills
2. Incorporates key achievements
3. Aligns with the target job requirements
4. Uses confident, professional language
5. Is written in %s

SMART BOLDING (REQUIRED):
Apply **bold** markdown syntax to highlight:
- Years of experience (e.g., **7+ years**)
- Key technical domains (e.g., **distributed systems**, **machine learning**)
- Core competencies (e.g., **architecting**, **scaling**, **leading teams**)
- Notable achievements or metrics (e.g., **Fortune 500**, **$10M revenue**)
Use sparingly - maximum 4-6 bold terms in the summary to maintain readability.

Respond with JSON:
{
  "summary": "the generated professional summary with **bold** highlights"
}

Respond ONLY with valid JSON.`,
		userName,
		stringPtr(req.User.Headline),
		stringPtr(req.User.Summary),
		bulletsContext.String(),
		req.JobAnalysis.Title,
		req.JobAnalysis.Company,
		strings.Join(req.JobAnalysis.RequiredSkills, ", "),
		req.JobAnalysis.Summary,
		req.TargetLanguage,
	)

	response, err := c.chatCompletion(ctx, c.config.ModelGeneration, prompt, 0.8)
	if err != nil {
		return nil, fmt.Errorf("groq: generate summary failed: %w", err)
	}

	var result struct {
		Summary string `json:"summary"`
	}

	if err := json.Unmarshal([]byte(cleanJSON(response)), &result); err != nil {
		log.Printf("json to parse: %s", cleanJSON(response))
		return nil, fmt.Errorf("groq: failed to parse summary: %w", err)
	}

	return &ports.SummaryResult{
		Summary: result.Summary,
	}, nil
}

// ScoreMatch calculates a match score between resume and job.
func (c *Client) ScoreMatch(ctx context.Context, req ports.ScoreMatchRequest) (*domain.MatchScore, error) {
	// Build skills list.
	var skillsList strings.Builder
	for _, skill := range req.UserSkills {
		fmt.Fprintf(&skillsList, "- %s (proficiency: %d%%)\n", skill.Name, skill.ProficiencyLevel.Int())
	}

	// Build experiences from resume content.
	var experiencesText strings.Builder
	if req.Resume != nil {
		fmt.Fprintf(&experiencesText, "Summary: %s\n\n", req.Resume.Summary)
		for _, exp := range req.Resume.Experiences {
			fmt.Fprintf(&experiencesText, "%s at %s:\n", exp.Title, exp.Organization)
			for _, bullet := range exp.Bullets {
				fmt.Fprintf(&experiencesText, "  - %s\n", bullet.TailoredContent)
			}
		}
	}

	prompt := fmt.Sprintf(`Score how well this resume matches the job requirements.

JOB REQUIREMENTS:
- Title: %s
- Required Skills: %s
- Preferred Skills: %s
- Years Experience: %v
- Summary: %s

CANDIDATE SKILLS:
%s

RESUME CONTENT:
%s

Analyze the match and provide a score from 0-100 based on:
1. Skill alignment (40%%)
2. Experience relevance (30%%)
3. Seniority fit (15%%)
4. Keyword coverage (15%%)

Respond with JSON:
{
  "score": 85,
  "breakdown": {
    "skills": 90,
    "experience": 80,
    "seniority": 85,
    "keywords": 75
  },
  "explanation": "Brief explanation of the score"
}

Respond ONLY with valid JSON.`,
		req.JobAnalysis.Title,
		strings.Join(req.JobAnalysis.RequiredSkills, ", "),
		strings.Join(req.JobAnalysis.PreferredSkills, ", "),
		req.JobAnalysis.YearsExperience,
		req.JobAnalysis.Summary,
		skillsList.String(),
		experiencesText.String(),
	)

	response, err := c.chatCompletion(ctx, c.config.ModelAnalysis, prompt, 0.2)
	if err != nil {
		return nil, fmt.Errorf("groq: score match failed: %w", err)
	}

	var result struct {
		Score int `json:"score"`
	}

	if err := json.Unmarshal([]byte(cleanJSON(response)), &result); err != nil {
		log.Printf("json to parse: %s", cleanJSON(response))
		return nil, fmt.Errorf("groq: failed to parse match score: %w", err)
	}

	score, err := domain.NewMatchScore(result.Score)
	if err != nil {
		return nil, fmt.Errorf("groq: invalid score value: %w", err)
	}

	return &score, nil
}

// Close releases any resources held by the AI provider.
func (c *Client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// chatCompletion sends a chat completion request to Groq API.
func (c *Client) chatCompletion(ctx context.Context, model, prompt string, temperature float64) (string, error) {
	reqBody := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens":  defaultMaxTokens,
		"temperature": temperature,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff.
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(backoff):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/chat/completions", bytes.NewReader(body))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %w", err)
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = fmt.Errorf("rate limited (attempt %d/%d)", attempt+1, c.config.MaxRetries+1)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
		}

		var response struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(cleanJSON(string(respBody))), &response); err != nil {
			log.Printf("json to parse: %s", cleanJSON(string(respBody)))
			return "", fmt.Errorf("failed to parse response: %w", err)
		}

		if len(response.Choices) == 0 {
			return "", fmt.Errorf("no choices in response")
		}

		content := response.Choices[0].Message.Content

		// Clean JSON from markdown code blocks if present.
		content = strings.TrimSpace(content)
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)

		return content, nil
	}

	return "", fmt.Errorf("max retries exceeded: %w", lastErr)
}

// stringPtr safely dereferences a string pointer.
func stringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// cleanJSON extracts JSON content from a string that may contain markdown formatting.
func cleanJSON(input string) string {
	// Starts the JSON object
	start := strings.Index(input, "{")
	// Finds the end of the JSON object
	end := strings.LastIndex(input, "}")

	// If both start and end are found, extract the JSON substring
	if start != -1 && end != -1 && start < end {
		return input[start : end+1]
	}

	// If not found, return the original input
	return input
}
