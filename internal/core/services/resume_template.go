// Package services contains the application services (use cases).
package services

import (
	"fmt"
	"html"
	"slices"
	"strings"

	"github.com/SeltikHD/chameleon-vitae/internal/core/domain"
)

// ResumeTemplateData contains all data needed to render a resume.
type ResumeTemplateData struct {
	User        *domain.User
	Resume      *domain.Resume
	Education   []domain.Education
	Projects    []domain.Project
	Languages   []domain.SpokenLanguage
	Skills      []domain.Skill
	FontSize    int    // Base font size in pt (11, 10, or 9)
	ShowSummary bool   // Whether to show the professional summary
	Locale      Locale // Locale for internationalization (defaults to en-US)
}

// JakeResumeTemplate implements the Jake's Resume format.
// This is the gold standard for developer resumes:
// - Single page, dense, ATS-friendly
// - Sections: Header → Education → Experience → Projects → Technical Skills
// - Clean typography with clear visual hierarchy
type JakeResumeTemplate struct{}

// NewJakeResumeTemplate creates a new Jake's Resume template.
func NewJakeResumeTemplate() *JakeResumeTemplate {
	return &JakeResumeTemplate{}
}

// Render generates the HTML for the resume.
func (t *JakeResumeTemplate) Render(data ResumeTemplateData) string {
	if data.FontSize == 0 {
		data.FontSize = 11
	}

	// Initialize i18n with the specified locale (defaults to en-US).
	i18n := NewI18n(data.Locale)

	var sb strings.Builder

	// Document head
	sb.WriteString(t.renderHead(data))

	// Body content
	sb.WriteString(`<body>`)
	sb.WriteString(`<div class="resume-container">`)

	// Header section
	sb.WriteString(t.renderHeader(data.User))

	// Professional Summary section (optional - after header, before education)
	if data.ShowSummary {
		summary := ""
		if data.Resume.GeneratedContent != nil && data.Resume.GeneratedContent.Summary != "" {
			summary = data.Resume.GeneratedContent.Summary
		} else if data.User != nil && data.User.Summary != nil && *data.User.Summary != "" {
			summary = *data.User.Summary
		}
		if summary != "" {
			sb.WriteString(t.renderSummary(summary, i18n))
		}
	}

	// Education section (always first in Jake's Resume)
	if len(data.Education) > 0 {
		sb.WriteString(t.renderEducation(data.Education, i18n))
	}

	// Technical Skills section
	if data.Resume.GeneratedContent != nil && len(data.Resume.GeneratedContent.Skills) > 0 {
		sb.WriteString(t.renderSkills(data.Resume.GeneratedContent.Skills, data.Skills, i18n))
	}

	// Experience section
	if data.Resume.GeneratedContent != nil && len(data.Resume.GeneratedContent.Experiences) > 0 {
		sb.WriteString(t.renderExperience(data.Resume.GeneratedContent.Experiences, i18n))
	}

	// Projects section (buffer section - can be dropped for one-page fit)
	if len(data.Projects) > 0 {
		sb.WriteString(t.renderProjects(data.Projects, i18n))
	}

	// Languages section (if any)
	if len(data.Languages) > 0 {
		sb.WriteString(t.renderLanguages(data.Languages, i18n))
	}

	sb.WriteString(`</div>`)
	sb.WriteString(`</body></html>`)

	return sb.String()
}

// renderHead generates the HTML head with Jake's Resume CSS.
func (t *JakeResumeTemplate) renderHead(data ResumeTemplateData) string {
	userName := "Resume"
	if data.User != nil {
		userName = data.User.GetDisplayName()
	}

	lang := "en"
	if data.Resume != nil {
		lang = data.Resume.TargetLanguage
	}

	baseFontSize := data.FontSize
	if baseFontSize == 0 {
		baseFontSize = 11
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="%s">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Resume - %s</title>
    <style>
        /* Jake's Resume CSS - ATS-friendly, single-page optimized */

        /* Reset and base */
        *, *::before, *::after {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
			font-family: Arial, sans-serif;
			line-height: 1.6;
            --font-family: 'Times New Roman', Times, serif;
            font-size: %dpt;
            --line-height: 1.15;
            color: #000;
            background: #fff;
        }

        .resume-container {
            max-width: 8.5in;
            margin: 0 auto;
            padding: 0.3in 0.4in;
        }

        /* Header */
        .resume-header {
            text-align: center;
            margin-bottom: 8pt;
            border-bottom: 1pt solid #000;
            padding-bottom: 4pt;
        }

        .resume-name {
            font-size: 18pt;
            font-weight: bold;
            letter-spacing: 0.5pt;
            text-transform: uppercase;
            margin-bottom: 4pt;
        }

		.resume-contact {
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
			width: 100%%;
			font-size: 9pt;
			color: #333;
		}

        .resume-contact a {
            color: #000;
            text-decoration: none;
        }

        .resume-contact a:hover {
            text-decoration: underline;
        }

        .contact-separator {
            margin: 0 6pt;
        }

        /* Section styling */
        .resume-section {
            margin-bottom: 8pt;
        }

        .section-title {
            font-size: 11pt;
            font-weight: bold;
            text-transform: uppercase;
            letter-spacing: 1pt;
            border-bottom: 1pt solid #000;
            padding-bottom: 2pt;
            margin-bottom: 4pt;
        }

        /* Professional Summary */
        .summary-section {
            margin-bottom: 10pt;
        }

        .summary-text {
            margin: 0;
            text-align: justify;
            line-height: 1.3;
        }

        /* Entry (Education, Experience, Project) */
        .resume-entry {
            margin-bottom: 6pt;
        }

        .entry-header {
            display: flex;
            justify-content: space-between;
            align-items: baseline;
        }

        .entry-title {
            font-weight: bold;
        }

        .entry-location {
            font-style: italic;
            font-size: 10pt;
        }

        .entry-subheader {
            display: flex;
            justify-content: space-between;
            align-items: baseline;
            font-style: italic;
        }

        .entry-subtitle {
            font-style: italic;
        }

        .entry-date {
            font-size: 10pt;
        }

        /* Bullets */
        .entry-bullets {
            list-style-type: disc;
            margin-left: 18pt;
            margin-top: 2pt;
        }

        .entry-bullets li {
            margin-bottom: 1pt;
            text-align: justify;
        }

        /* Projects specific */
        .project-header {
            display: flex;
            align-items: baseline;
            gap: 8pt;
        }

        .project-name {
            font-weight: bold;
        }

        .project-tech {
            font-style: italic;
            font-size: 10pt;
        }

        .project-link {
            font-size: 9pt;
            color: #0066cc;
            text-decoration: none;
            margin-left: 4pt;
        }

        .project-link:hover {
            text-decoration: underline;
        }

        .project-links {
            font-size: 9pt;
        }

        .project-links a {
            color: #000;
            text-decoration: none;
        }

        /* Skills section */
        .skills-list {
            margin: 0;
            padding: 0;
            list-style: none;
        }

        .skills-row {
            margin-bottom: 2pt;
        }

        .skill-category {
            font-weight: bold;
        }

        .skill-items {
            font-weight: normal;
        }

        /* Languages section */
        .languages-list {
            display: flex;
            flex-wrap: wrap;
            gap: 12pt;
        }

        .language-item {
            font-size: 10pt;
        }

        .language-name {
            font-weight: bold;
        }

        .language-level {
            font-style: italic;
        }

        /* Honors/GPA inline */
        .education-honors {
            font-style: italic;
            font-size: 10pt;
        }

        /* Print optimization */
        @media print {
            body {
                -webkit-print-color-adjust: exact;
                print-color-adjust: exact;
            }

            .resume-container {
                padding: 0;
            }

            @page {
                size: letter;
                margin: 0.3in 0.4in;
            }
        }
    </style>
</head>
`, lang, html.EscapeString(userName), baseFontSize)
}

// renderHeader generates the header section with name and contact info.
func (t *JakeResumeTemplate) renderHeader(user *domain.User) string {
	if user == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<header class="resume-header">`)
	fmt.Fprintf(&sb, `<h1 class="resume-name">%s</h1>`, html.EscapeString(user.GetDisplayName()))

	// Build contact line
	var contacts []string

	if user.Phone != nil && *user.Phone != "" {
		contacts = append(contacts, html.EscapeString(*user.Phone))
	}

	if user.Email != nil && *user.Email != "" {
		contacts = append(contacts, fmt.Sprintf(`<a href="mailto:%s">%s</a>`,
			html.EscapeString(*user.Email),
			html.EscapeString(*user.Email)))
	}

	if user.LinkedInURL != nil && *user.LinkedInURL != "" {
		// Extract username from LinkedIn URL if possible
		linkedIn := extractURLDisplay(*user.LinkedInURL, "linkedin.com/in/")
		contacts = append(contacts, fmt.Sprintf(`<a href="%s">%s</a>`,
			html.EscapeString(*user.LinkedInURL),
			html.EscapeString(linkedIn)))
	}

	if user.GitHubURL != nil && *user.GitHubURL != "" {
		// Extract username from GitHub URL if possible
		github := extractURLDisplay(*user.GitHubURL, "github.com/")
		contacts = append(contacts, fmt.Sprintf(`<a href="%s">%s</a>`,
			html.EscapeString(*user.GitHubURL),
			html.EscapeString(github)))
	}

	if user.PortfolioURL != nil && *user.PortfolioURL != "" {
		contacts = append(contacts, fmt.Sprintf(`<a href="%s">%s</a>`,
			html.EscapeString(*user.PortfolioURL),
			html.EscapeString(extractDomain(*user.PortfolioURL))))
	}

	if len(contacts) > 0 {
		sb.WriteString(`<p class="resume-contact">`)
		sb.WriteString(strings.Join(contacts, `<span class="contact-separator">|</span>`))
		sb.WriteString(`</p>`)
	}

	sb.WriteString(`</header>`)
	return sb.String()
}

// renderSummary generates the professional summary section.
func (t *JakeResumeTemplate) renderSummary(summary string, i18n *I18n) string {
	if summary == "" {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<section class="resume-section summary-section">`)
	sb.WriteString(fmt.Sprintf(`<h2 class="section-title">%s</h2>`, html.EscapeString(i18n.T(KeyProfessionalSummary))))
	sb.WriteString(`<p class="summary-text">`)
	sb.WriteString(renderMarkdownBold(summary))
	sb.WriteString(`</p>`)
	sb.WriteString(`</section>`)
	return sb.String()
}

// renderEducation generates the education section.
func (t *JakeResumeTemplate) renderEducation(education []domain.Education, i18n *I18n) string {
	if len(education) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<section class="resume-section">`)
	sb.WriteString(fmt.Sprintf(`<h2 class="section-title">%s</h2>`, html.EscapeString(i18n.T(KeyEducation))))

	for _, edu := range education {
		sb.WriteString(`<div class="resume-entry">`)

		// First line: Institution | Location
		sb.WriteString(`<div class="entry-header">`)
		fmt.Fprintf(&sb, `<span class="entry-title">%s</span>`, html.EscapeString(edu.Institution))
		if edu.Location != nil && *edu.Location != "" {
			fmt.Fprintf(&sb, `<span class="entry-location">%s</span>`, html.EscapeString(*edu.Location))
		}
		sb.WriteString(`</div>`)

		// Second line: Degree, Field of Study | Dates
		sb.WriteString(`<div class="entry-subheader">`)
		degree := edu.Degree
		if edu.FieldOfStudy != nil && *edu.FieldOfStudy != "" {
			degree += " in " + *edu.FieldOfStudy
		}
		fmt.Fprintf(&sb, `<span class="entry-subtitle">%s</span>`, html.EscapeString(degree))
		fmt.Fprintf(&sb, `<span class="entry-date">%s</span>`, formatEducationDateRangeLocalized(edu.StartDate, edu.EndDate, i18n))
		sb.WriteString(`</div>`)

		// Honors/GPA if present
		var extras []string
		if edu.GPA != nil && *edu.GPA != "" {
			extras = append(extras, i18n.T(KeyGPA)+": "+*edu.GPA)
		}
		if len(edu.Honors) > 0 {
			extras = append(extras, strings.Join(edu.Honors, ", "))
		}
		if len(extras) > 0 {
			fmt.Fprintf(&sb, `<div class="education-honors">%s</div>`, html.EscapeString(strings.Join(extras, " | ")))
		}

		sb.WriteString(`</div>`)
	}

	sb.WriteString(`</section>`)
	return sb.String()
}

// renderExperience generates the experience section.
func (t *JakeResumeTemplate) renderExperience(experiences []domain.TailoredExperience, i18n *I18n) string {
	if len(experiences) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<section class="resume-section">`)
	sb.WriteString(fmt.Sprintf(`<h2 class="section-title">%s</h2>`, html.EscapeString(i18n.T(KeyExperience))))

	for _, exp := range experiences {
		sb.WriteString(`<div class="resume-entry">`)

		// First line: Title | Dates
		sb.WriteString(`<div class="entry-header">`)
		fmt.Fprintf(&sb, `<span class="entry-title">%s</span>`, html.EscapeString(exp.Title))
		dateStr := formatExperienceDateRangeLocalized(exp.StartDate, exp.EndDate, exp.IsCurrent, i18n)
		fmt.Fprintf(&sb, `<span class="entry-date">%s</span>`, html.EscapeString(dateStr))
		sb.WriteString(`</div>`)

		// Second line: Organization
		sb.WriteString(`<div class="entry-subheader">`)
		fmt.Fprintf(&sb, `<span class="entry-subtitle">%s</span>`, html.EscapeString(exp.Organization))
		sb.WriteString(`</div>`)

		// Bullets
		if len(exp.Bullets) > 0 {
			sb.WriteString(`<ul class="entry-bullets">`)
			for _, bullet := range exp.Bullets {
				content := bullet.TailoredContent
				if content == "" {
					content = bullet.OriginalContent
				}
				fmt.Fprintf(&sb, `<li>%s</li>`, renderMarkdownBold(content))
			}
			sb.WriteString(`</ul>`)
		}

		sb.WriteString(`</div>`)
	}

	sb.WriteString(`</section>`)
	return sb.String()
}

// renderProjects generates the projects section.
func (t *JakeResumeTemplate) renderProjects(projects []domain.Project, i18n *I18n) string {
	if len(projects) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<section class="resume-section">`)
	sb.WriteString(fmt.Sprintf(`<h2 class="section-title">%s</h2>`, html.EscapeString(i18n.T(KeyProjects))))

	for _, proj := range projects {
		sb.WriteString(`<div class="resume-entry">`)

		// Project header: Name | Tech Stack | Links | Date
		sb.WriteString(`<div class="entry-header">`)
		sb.WriteString(`<div class="project-header">`)
		fmt.Fprintf(&sb, `<span class="project-name">%s</span>`, html.EscapeString(proj.Name))
		if len(proj.TechStack) > 0 {
			fmt.Fprintf(&sb, `<span class="project-tech">| %s</span>`,
				html.EscapeString(strings.Join(proj.TechStack, ", ")))
		}
		// Discrete project links
		if proj.RepositoryURL != nil && *proj.RepositoryURL != "" {
			fmt.Fprintf(&sb, `<a href="%s" class="project-link">[Source]</a>`,
				html.EscapeString(*proj.RepositoryURL))
		}
		if proj.URL != nil && *proj.URL != "" {
			fmt.Fprintf(&sb, `<a href="%s" class="project-link">[Demo]</a>`,
				html.EscapeString(*proj.URL))
		}
		sb.WriteString(`</div>`)
		dateStr := formatProjectDateRangeLocalized(proj.StartDate, proj.EndDate, i18n)
		if dateStr != "" {
			fmt.Fprintf(&sb, `<span class="entry-date">%s</span>`, html.EscapeString(dateStr))
		}
		sb.WriteString(`</div>`)

		// Bullets
		if len(proj.Bullets) > 0 {
			sb.WriteString(`<ul class="entry-bullets">`)
			for _, bullet := range proj.Bullets {
				fmt.Fprintf(&sb, `<li>%s</li>`, renderMarkdownBold(bullet.Content))
			}
			sb.WriteString(`</ul>`)
		}

		sb.WriteString(`</div>`)
	}

	sb.WriteString(`</section>`)
	return sb.String()
}

// renderSkills generates the technical skills section in key-value format.
func (t *JakeResumeTemplate) renderSkills(selectedSkills []string, userSkills []domain.Skill, i18n *I18n) string {
	if len(selectedSkills) == 0 {
		return ""
	}

	// Group skills by category
	categorySkills := make(map[string][]string)
	skillCategories := make(map[string]string) // skill name -> category

	// Build skill lookup from user skills
	for _, skill := range userSkills {
		category := "Other"
		if skill.Category != nil && *skill.Category != "" {
			category = *skill.Category
		}
		skillCategories[strings.ToLower(skill.Name)] = category
	}

	// Group selected skills by category
	for _, skillName := range selectedSkills {
		category := skillCategories[strings.ToLower(skillName)]
		if category == "" {
			category = "Other"
		}
		categorySkills[category] = append(categorySkills[category], skillName)
	}

	// Define category order
	categoryOrder := []string{"Languages", "Frameworks", "Tools", "Databases", "Cloud", "Other"}

	var sb strings.Builder
	sb.WriteString(`<section class="resume-section">`)
	sb.WriteString(fmt.Sprintf(`<h2 class="section-title">%s</h2>`, html.EscapeString(i18n.T(KeyTechnicalSkills))))
	sb.WriteString(`<ul class="skills-list">`)

	for _, category := range categoryOrder {
		skills, exists := categorySkills[category]
		if !exists || len(skills) == 0 {
			continue
		}
		sb.WriteString(`<li class="skills-row">`)
		fmt.Fprintf(&sb, `<span class="skill-category">%s:</span> `, html.EscapeString(category))
		fmt.Fprintf(&sb, `<span class="skill-items">%s</span>`, html.EscapeString(strings.Join(skills, ", ")))
		sb.WriteString(`</li>`)
	}

	// Handle any remaining categories not in the predefined order
	for category, skills := range categorySkills {
		found := slices.Contains(categoryOrder, category)
		if !found && len(skills) > 0 {
			sb.WriteString(`<li class="skills-row">`)
			fmt.Fprintf(&sb, `<span class="skill-category">%s:</span> `, html.EscapeString(category))
			fmt.Fprintf(&sb, `<span class="skill-items">%s</span>`, html.EscapeString(strings.Join(skills, ", ")))
			sb.WriteString(`</li>`)
		}
	}

	sb.WriteString(`</ul>`)
	sb.WriteString(`</section>`)
	return sb.String()
}

// renderLanguages generates the spoken languages section.
func (t *JakeResumeTemplate) renderLanguages(languages []domain.SpokenLanguage, i18n *I18n) string {
	if len(languages) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<section class="resume-section">`)
	sb.WriteString(fmt.Sprintf(`<h2 class="section-title">%s</h2>`, html.EscapeString(i18n.T(KeyLanguages))))
	sb.WriteString(`<div class="languages-list">`)

	for _, lang := range languages {
		sb.WriteString(`<span class="language-item">`)
		fmt.Fprintf(&sb, `<span class="language-name">%s</span>`, html.EscapeString(lang.Language))
		fmt.Fprintf(&sb, ` (<span class="language-level">%s</span>)`, html.EscapeString(i18n.FormatProficiencyLevel(string(lang.Proficiency))))
		sb.WriteString(`</span>`)
	}

	sb.WriteString(`</div>`)
	sb.WriteString(`</section>`)
	return sb.String()
}

// Helper functions

// renderMarkdownBold converts markdown **bold** syntax to HTML <strong> tags.
// It first escapes HTML in the input, then converts **text** to <strong>text</strong>.
func renderMarkdownBold(text string) string {
	// First escape HTML to prevent XSS
	escaped := html.EscapeString(text)

	// Replace **text** with <strong>text</strong>
	// Use a simple state machine approach
	var result strings.Builder
	i := 0
	for i < len(escaped) {
		// Check for ** at current position
		if i+1 < len(escaped) && escaped[i] == '*' && escaped[i+1] == '*' {
			// Find closing **
			closeIdx := strings.Index(escaped[i+2:], "**")
			if closeIdx != -1 {
				// Found matching **
				boldContent := escaped[i+2 : i+2+closeIdx]
				result.WriteString("<strong>")
				result.WriteString(boldContent)
				result.WriteString("</strong>")
				i = i + 2 + closeIdx + 2 // Skip past closing **
				continue
			}
		}
		result.WriteByte(escaped[i])
		i++
	}
	return result.String()
}

func extractURLDisplay(url, prefix string) string {
	// Try to extract meaningful part from URL
	_, after, ok := strings.Cut(url, prefix)
	if ok {
		remaining := after
		// Remove trailing slashes
		remaining = strings.TrimSuffix(remaining, "/")
		return prefix + remaining
	}
	return url
}

func extractDomain(url string) string {
	// Remove protocol
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	// Remove path
	if idx := strings.Index(url, "/"); idx != -1 {
		url = url[:idx]
	}
	return url
}

func formatEducationDateRangeLocalized(startDate, endDate *domain.Date, i18n *I18n) string {
	if startDate == nil && endDate == nil {
		return ""
	}

	format := func(d *domain.Date) string {
		if d == nil || d.IsZero() {
			return ""
		}
		return i18n.FormatDate(d.Time)
	}

	start := format(startDate)
	end := format(endDate)

	if end == "" {
		end = i18n.T(KeyPresent)
	}

	if start == "" {
		return end
	}

	return start + " – " + end
}

func formatExperienceDateRangeLocalized(startDate string, endDate *string, isCurrent bool, i18n *I18n) string {
	if startDate == "" {
		return ""
	}

	start := i18n.FormatDateString(startDate)
	end := i18n.T(KeyPresent)
	if !isCurrent && endDate != nil && *endDate != "" {
		end = i18n.FormatDateString(*endDate)
	}

	return start + " – " + end
}

func formatProjectDateRangeLocalized(startDate, endDate *domain.Date, i18n *I18n) string {
	if startDate == nil && endDate == nil {
		return ""
	}

	format := func(d *domain.Date) string {
		if d == nil || d.IsZero() {
			return ""
		}
		return i18n.FormatDate(d.Time)
	}

	start := format(startDate)
	end := format(endDate)

	if start == "" && end == "" {
		return ""
	}

	if end == "" {
		return start
	}

	if start == "" {
		return end
	}

	return start + " – " + end
}
