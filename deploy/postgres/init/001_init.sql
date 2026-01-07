-- ============================================================================
-- Chameleon Vitae - Database Initialization Script
-- ============================================================================
-- This script runs automatically when the PostgreSQL container starts
-- with an empty data volume.
-- ============================================================================

-- Enable UUID extension for generating UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable pg_trgm for fuzzy text search (useful for matching skills)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- ============================================================================
-- Core Tables
-- ============================================================================

-- Users table (Firebase authentication)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    firebase_uid VARCHAR(128) UNIQUE NOT NULL,
    email VARCHAR(255),
    name VARCHAR(255),
    headline VARCHAR(255),
    summary TEXT,
    location VARCHAR(255),
    phone VARCHAR(50),
    website VARCHAR(512),
    linkedin_url VARCHAR(512),
    github_url VARCHAR(512),
    portfolio_url VARCHAR(512),
    preferred_language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Experiences (work, education, volunteer, freelance, etc.)
CREATE TABLE IF NOT EXISTS experiences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN (
        'work',
        'education',
        'certification',
        'project',
        'freelance',
        'volunteer',
        'open_source',
        'hackathon',
        'side_project',
        'event_organization',
        'publication',
        'award'
    )),
    title VARCHAR(255) NOT NULL,
    organization VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    start_date DATE NOT NULL,
    end_date DATE,
    is_current BOOLEAN DEFAULT FALSE,
    description TEXT,
    url VARCHAR(512),
    metadata JSONB DEFAULT '{}',
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Experience bullets (atomic units of experience for AI selection)
-- Each bullet represents a single achievement/responsibility that can be
-- independently selected and rewritten for tailored resumes.
CREATE TABLE IF NOT EXISTS bullets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    experience_id UUID NOT NULL REFERENCES experiences(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    -- impact_score: AI-calculated metric (0-100) indicating the strength/impact
    -- of this bullet point. Used for prioritization during resume generation.
    -- Higher scores = more impressive achievements (quantifiable results, leadership, etc.)
    impact_score INTEGER DEFAULT 50 CHECK (impact_score >= 0 AND impact_score <= 100),
    keywords TEXT[] DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Skills with categories and proficiency
CREATE TABLE IF NOT EXISTS skills (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    proficiency_level INTEGER DEFAULT 50 CHECK (proficiency_level >= 0 AND proficiency_level <= 100),
    years_of_experience NUMERIC(4, 1),
    is_highlighted BOOLEAN DEFAULT FALSE,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, name)
);

-- Languages spoken by the user (distinct from programming languages)
CREATE TABLE IF NOT EXISTS spoken_languages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language VARCHAR(100) NOT NULL,
    proficiency VARCHAR(50) NOT NULL CHECK (proficiency IN (
        'native',
        'fluent',
        'advanced',
        'intermediate',
        'basic'
    )),
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, language)
);

-- Generated resumes (tailored to specific job applications)
CREATE TABLE IF NOT EXISTS resumes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    job_description TEXT NOT NULL,
    job_title VARCHAR(255),
    company_name VARCHAR(255),
    job_url VARCHAR(512),
    target_language VARCHAR(10) DEFAULT 'en',
    selected_bullets UUID[] DEFAULT '{}',
    generated_content JSONB,
    pdf_url VARCHAR(512),
    score INTEGER DEFAULT 0 CHECK (score >= 0 AND score <= 100),
    notes TEXT,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN (
        'draft',
        'generated',
        'reviewed',
        'submitted',
        'interview',
        'rejected',
        'accepted'
    )),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================================
-- Indexes
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_users_firebase_uid ON users(firebase_uid);
CREATE INDEX IF NOT EXISTS idx_experiences_user_id ON experiences(user_id);
CREATE INDEX IF NOT EXISTS idx_experiences_type ON experiences(type);
CREATE INDEX IF NOT EXISTS idx_experiences_user_type ON experiences(user_id, type);
CREATE INDEX IF NOT EXISTS idx_bullets_experience_id ON bullets(experience_id);
CREATE INDEX IF NOT EXISTS idx_bullets_keywords ON bullets USING GIN(keywords);
CREATE INDEX IF NOT EXISTS idx_bullets_impact_score ON bullets(impact_score DESC);
CREATE INDEX IF NOT EXISTS idx_skills_user_id ON skills(user_id);
CREATE INDEX IF NOT EXISTS idx_skills_name_trgm ON skills USING GIN(name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_skills_category ON skills(category);
CREATE INDEX IF NOT EXISTS idx_spoken_languages_user_id ON spoken_languages(user_id);
CREATE INDEX IF NOT EXISTS idx_resumes_user_id ON resumes(user_id);
CREATE INDEX IF NOT EXISTS idx_resumes_status ON resumes(status);

-- ============================================================================
-- Triggers for updated_at
-- ============================================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_experiences_updated_at
    BEFORE UPDATE ON experiences
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bullets_updated_at
    BEFORE UPDATE ON bullets
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_resumes_updated_at
    BEFORE UPDATE ON resumes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- Comments
-- ============================================================================

COMMENT ON TABLE users IS 'User accounts linked to Firebase authentication';
COMMENT ON COLUMN users.firebase_uid IS 'Unique identifier from Firebase Authentication';
COMMENT ON COLUMN users.headline IS 'Professional headline (e.g., "Senior Software Engineer")';
COMMENT ON COLUMN users.summary IS 'Professional summary for the resume header';
COMMENT ON COLUMN users.preferred_language IS 'ISO 639-1 language code for resume generation';

COMMENT ON TABLE experiences IS 'Work experiences, projects, education, certifications, and other resume entries';
COMMENT ON COLUMN experiences.type IS 'Type of experience: work, education, certification, project, freelance, volunteer, open_source, hackathon, side_project, event_organization, publication, award';
COMMENT ON COLUMN experiences.display_order IS 'Order in which experiences appear within their type';

COMMENT ON TABLE bullets IS 'Atomic experience bullets for AI-powered resume tailoring';
COMMENT ON COLUMN bullets.impact_score IS 'AI-calculated impact score (0-100) for prioritization. Higher = more impressive.';
COMMENT ON COLUMN bullets.keywords IS 'Keywords extracted from the bullet for job matching';

COMMENT ON TABLE skills IS 'User skills with proficiency levels';
COMMENT ON COLUMN skills.category IS 'Skill category (e.g., Programming Languages, Frameworks, Tools)';
COMMENT ON COLUMN skills.is_highlighted IS 'Whether to feature this skill prominently';

COMMENT ON TABLE spoken_languages IS 'Natural languages spoken by the user';

COMMENT ON TABLE resumes IS 'Generated resumes tailored to specific job applications';
COMMENT ON COLUMN resumes.score IS 'AI-calculated match score with job description (0-100)';
COMMENT ON COLUMN resumes.target_language IS 'Language for resume generation (e.g., en, pt-br)';
