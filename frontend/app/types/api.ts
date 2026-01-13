/**
 * API Type Definitions for Chameleon Vitae.
 *
 * These interfaces strictly match the Golang struct/JSON responses
 * from the backend API. Any changes to backend DTOs MUST be reflected here.
 *
 * @see internal/adapters/primary/http/dto.go
 */

// ============================================================================
// Common Response Types
// ============================================================================

/** Standard error detail for field validation. */
export interface ErrorDetail {
  field: string
  message: string
}

/** Error body structure. */
export interface ErrorBody {
  code: string
  message: string
  details?: ErrorDetail[]
}

/** Standard error response format. */
export interface ErrorResponse {
  error: ErrorBody
}

/** Generic success response. */
export interface SuccessResponse<T = unknown> {
  data?: T
  message?: string
}

/** Paginated list response. */
export interface PaginatedResponse<T> {
  data: T[]
  total: number
  limit: number
  offset: number
}

// ============================================================================
// Auth DTOs
// ============================================================================

/** Request body for user synchronization after Firebase login. */
export interface SyncUserRequest {
  firebase_uid: string
  email?: string | null
  name?: string | null
  picture?: string | null
}

/** Response after user synchronization. */
export interface SyncUserResponse {
  id: string
  firebase_uid: string
  email?: string | null
  name?: string | null
  created_at: string
}

// ============================================================================
// User DTOs
// ============================================================================

/** User profile response from GET /me. */
export interface UserResponse {
  id: string
  firebase_uid: string
  picture_url?: string | null
  email?: string | null
  name?: string | null
  headline?: string | null
  summary?: string | null
  location?: string | null
  phone?: string | null
  website?: string | null
  linkedin_url?: string | null
  github_url?: string | null
  portfolio_url?: string | null
  preferred_language: string
  created_at: string
  updated_at: string
}

/** Request body for updating user profile (PATCH /me). */
export interface UpdateUserRequest {
  name?: string | null
  headline?: string | null
  summary?: string | null
  location?: string | null
  phone?: string | null
  website?: string | null
  linkedin_url?: string | null
  github_url?: string | null
  portfolio_url?: string | null
  preferred_language?: string | null
}

// ============================================================================
// Experience DTOs
// ============================================================================

/** Experience types matching backend enum. */
export type ExperienceType =
  | 'work'
  | 'education'
  | 'project'
  | 'volunteer'
  | 'certification'
  | 'award'

/** Experience response from API. */
export interface ExperienceResponse {
  id: string
  type: ExperienceType
  title: string
  organization: string
  location?: string | null
  start_date: string
  end_date?: string | null
  is_current: boolean
  description?: string | null
  url?: string | null
  metadata?: Record<string, unknown> | null
  display_order: number
  bullets?: BulletResponse[]
  created_at: string
  updated_at: string
}

/** Request body for creating an experience. */
export interface CreateExperienceRequest {
  type: ExperienceType
  title: string
  organization: string
  location?: string | null
  start_date: string
  end_date?: string | null
  is_current?: boolean
  description?: string | null
  url?: string | null
  metadata?: Record<string, unknown> | null
}

/** Request body for updating an experience. */
export interface UpdateExperienceRequest {
  type?: ExperienceType | null
  title?: string | null
  organization?: string | null
  location?: string | null
  start_date?: string | null
  end_date?: string | null
  is_current?: boolean | null
  description?: string | null
  url?: string | null
  metadata?: Record<string, unknown> | null
  display_order?: number | null
}

/** Paginated list of experiences. */
export type ListExperiencesResponse = PaginatedResponse<ExperienceResponse>

// ============================================================================
// Bullet DTOs
// ============================================================================

/** Bullet (atomic experience unit) response from API. */
export interface BulletResponse {
  id: string
  experience_id: string
  content: string
  impact_score: number
  keywords: string[]
  metadata?: Record<string, unknown> | null
  display_order: number
  created_at: string
  updated_at: string
}

/** Request body for creating a bullet. */
export interface CreateBulletRequest {
  content: string
  keywords?: string[]
  display_order?: number
}

/** Request body for updating a bullet. */
export interface UpdateBulletRequest {
  content?: string | null
  keywords?: string[]
  display_order?: number | null
}

/** Response after recalculating bullet score. */
export interface ScoreBulletResponse {
  id: string
  content: string
  impact_score: number
  score_reasoning?: string
}

// ============================================================================
// Education DTOs
// ============================================================================

/** Education entry response from API. */
export interface EducationResponse {
  id: string
  institution: string
  degree: string
  field_of_study: string
  location?: string | null
  start_date?: string | null
  end_date?: string | null
  gpa?: string | null
  honors?: string[]
  display_order: number
  created_at: string
  updated_at: string
}

/** Request body for creating an education entry. */
export interface CreateEducationRequest {
  institution: string
  degree: string
  field_of_study: string
  location?: string | null
  start_date?: string | null
  end_date?: string | null
  gpa?: string | null
  honors?: string[]
  display_order?: number
}

/** Request body for updating an education entry. */
export interface UpdateEducationRequest {
  institution?: string | null
  degree?: string | null
  field_of_study?: string | null
  location?: string | null
  start_date?: string | null
  end_date?: string | null
  gpa?: string | null
  honors?: string[]
  display_order?: number | null
}

/** List of education entries response. */
export interface ListEducationResponse {
  data: EducationResponse[]
  total: number
}

// ============================================================================
// Project DTOs
// ============================================================================

/** Project bullet response from API. */
export interface ProjectBulletResponse {
  id: string
  project_id: string
  content: string
  display_order: number
  created_at: string
}

/** Project response from API. */
export interface ProjectResponse {
  id: string
  name: string
  description?: string
  tech_stack: string[]
  url?: string
  repository_url?: string
  start_date?: string | null
  end_date?: string | null
  display_order: number
  bullets?: ProjectBulletResponse[]
  created_at: string
  updated_at: string
}

/** Request body for creating a project. */
export interface CreateProjectRequest {
  name: string
  description?: string
  tech_stack?: string[]
  url?: string
  repository_url?: string
  start_date?: string | null
  end_date?: string | null
  display_order?: number
  bullets?: string[]
}

/** Request body for updating a project. */
export interface UpdateProjectRequest {
  name?: string | null
  description?: string | null
  tech_stack?: string[]
  url?: string | null
  repository_url?: string | null
  start_date?: string | null
  end_date?: string | null
  display_order?: number | null
}

/** Request body for creating a project bullet. */
export interface CreateProjectBulletRequest {
  content: string
  display_order?: number
}

/** List of projects response. */
export interface ListProjectsResponse {
  data: ProjectResponse[]
  total: number
}

// ============================================================================
// Skill DTOs
// ============================================================================

/** Skill response from API. */
export interface SkillResponse {
  id: string
  name: string
  category?: string | null
  proficiency_level: number
  years_of_experience?: number | null
  is_highlighted: boolean
  display_order: number
  created_at: string
}

/** Skill input for batch operations. */
export interface SkillInput {
  name: string
  category?: string | null
  proficiency_level?: number | null
  years_of_experience?: number | null
  is_highlighted?: boolean | null
}

/** Request body for batch skill upsert. */
export interface BatchUpsertSkillsRequest {
  skills: SkillInput[]
}

/** Response for batch skill upsert. */
export interface BatchUpsertSkillsResponse {
  created: number
  updated: number
  data: SkillResponse[]
}

/** List of skills response. */
export interface ListSkillsResponse {
  data: SkillResponse[]
  total: number
}

// ============================================================================
// Spoken Language DTOs
// ============================================================================

/** Spoken language proficiency levels. */
export type LanguageProficiency = 'native' | 'fluent' | 'professional' | 'intermediate' | 'beginner'

/** Spoken language response from API. */
export interface SpokenLanguageResponse {
  id: string
  language: string
  proficiency: LanguageProficiency
  display_order: number
  created_at: string
}

/** Request body for creating a spoken language. */
export interface CreateSpokenLanguageRequest {
  language: string
  proficiency: LanguageProficiency
  display_order?: number
}

/** List of spoken languages response. */
export interface ListSpokenLanguagesResponse {
  data: SpokenLanguageResponse[]
}

// ============================================================================
// Resume DTOs
// ============================================================================

/** Resume status values matching backend domain. */
export type ResumeStatus =
  | 'draft'
  | 'generated'
  | 'reviewed'
  | 'submitted'
  | 'interview'
  | 'rejected'
  | 'accepted'

/** Tailored bullet in generated content. */
export interface TailoredBulletDTO {
  bullet_id: string
  original_content: string
  tailored_content: string
}

/** Tailored experience in generated content. */
export interface TailoredExperienceDTO {
  experience_id: string
  title: string
  organization: string
  start_date: string
  end_date?: string | null
  is_current: boolean
  bullets: TailoredBulletDTO[]
}

/** AI analysis of resume match. */
export interface ResumeAnalysisDTO {
  matched_keywords: string[]
  missing_keywords: string[]
  recommendations: string[]
}

/** AI-generated resume content. */
export interface ResumeContentDTO {
  summary: string
  experiences: TailoredExperienceDTO[]
  skills: string[]
  analysis?: ResumeAnalysisDTO | null
}

/** Full resume response from API. */
export interface ResumeResponse {
  id: string
  job_title?: string
  company_name?: string
  job_url?: string
  job_description?: string
  target_language: string
  selected_bullets?: string[]
  generated_content?: ResumeContentDTO | null
  pdf_url?: string
  score?: number
  notes?: string
  status: ResumeStatus
  created_at: string
  updated_at: string
}

/** Resume list item (without full content). */
export interface ResumeListItem {
  id: string
  job_title: string
  company_name: string
  job_url?: string | null
  target_language: string
  score?: number | null
  status: ResumeStatus
  created_at: string
  updated_at: string
}

/** Request body for creating a resume. */
export interface CreateResumeRequest {
  job_description: string
  job_title?: string
  company_name?: string
  job_url?: string
  target_language?: string
}

/** Request body for tailoring a resume. */
export interface TailorResumeRequest {
  max_bullets_per_job?: number
}

/** Response after tailoring a resume. */
export interface TailorResumeResponse {
  id: string
  status: ResumeStatus
  score: number
  selected_bullets: string[]
  generated_content: Record<string, unknown>
  analysis?: ResumeAnalysisDTO | null
}

/** Request body for updating resume status/content. */
export interface UpdateResumeContentRequest {
  status: ResumeStatus
  notes?: string | null
}

/** Paginated list of resumes. */
export type ListResumesResponse = PaginatedResponse<ResumeResponse>

// ============================================================================
// Tools DTOs
// ============================================================================

/** Request body for parsing a job URL. */
export interface ParseJobURLRequest {
  url: string
}

/** Metadata about parsed job. */
export interface ParseJobMetadata {
  source: string
  fetched_at: string
}

/** Response from job URL parsing. */
export interface ParseJobURLResponse {
  url: string
  title: string
  markdown: string
  metadata?: ParseJobMetadata | null
}

// ============================================================================
// AI Analysis DTOs
// ============================================================================

/** Request body for analyzing a job description. */
export interface AnalyzeJobRequest {
  job_description: string
  job_url?: string
  target_language?: string
}

/** Response from job analysis with AI. */
export interface AnalyzeJobResponse {
  title?: string
  company?: string
  position?: string
  required_skills: string[]
  nice_to_have: string[]
  keywords: string[]
  seniority_level?: string
  experience_level?: string
  years_experience?: number
  summary?: string
}
