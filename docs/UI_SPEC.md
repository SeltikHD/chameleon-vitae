# UI Specification — Chameleon Vitae

> **Frontend Sitemap and Component Structure**

This document outlines all pages, their routes, and the high-level content structure for the Chameleon Vitae frontend application.

---

## Table of Contents

1. [Public Routes](#1-public-routes)
2. [Dashboard Routes (Protected)](#2-dashboard-routes-protected)
3. [Component Hierarchy](#3-component-hierarchy)
4. [Page Specifications](#4-page-specifications)

---

## 1. Public Routes

Routes accessible without authentication.

| Route | Page | Description |
|-------|------|-------------|
| `/` | Landing | Marketing page with hero, features, CTA |
| `/login` | Login | Authentication form (Email, Google, GitHub) |

---

## 2. Dashboard Routes (Protected)

Routes requiring Firebase authentication. Wrapped in `dashboard` layout with sidebar navigation.

| Route | Page | Description |
|-------|------|-------------|
| `/dashboard` | Overview | Stats cards, recent activity, quick actions |
| `/dashboard/resumes` | Resume List | Table of generated resumes with filters |
| `/dashboard/resumes/new` | New Resume | Job description input, URL parsing, resume generation |
| `/dashboard/resumes/[id]` | Resume Workbench | Split-view editor with live PDF preview |
| `/dashboard/experiences` | Experience Timeline | Vertical timeline of work/education entries |
| `/dashboard/skills` | Skill Matrix | Tag-based skill management with categories |
| `/dashboard/profile` | Profile Settings | Static user info, preferences, account |

---

## 3. Component Hierarchy

```text
app/
├── layouts/
│   ├── default.vue          # Public pages (minimal header)
│   └── dashboard.vue         # Protected pages (sidebar + header)
│
├── components/
│   ├── landing/
│   │   ├── HeroSection.vue
│   │   ├── FeatureGrid.vue
│   │   └── CallToAction.vue
│   │
│   ├── dashboard/
│   │   ├── Sidebar.vue
│   │   ├── Header.vue
│   │   ├── StatsCard.vue
│   │   └── ActivityFeed.vue
│   │
│   ├── resume/
│   │   ├── ResumeCard.vue
│   │   ├── ResumeEditor.vue
│   │   ├── ResumePreview.vue
│   │   └── JobDescriptionInput.vue
│   │
│   ├── experience/
│   │   ├── ExperienceCard.vue
│   │   ├── ExperienceTimeline.vue
│   │   └── BulletItem.vue
│   │
│   └── shared/
│       ├── Logo.vue
│       ├── ThemeToggle.vue
│       └── EmptyState.vue
│
└── pages/
    ├── index.vue             # Landing
    ├── login.vue             # Login
    └── dashboard/
        ├── index.vue         # Overview
        ├── resumes/
        │   ├── index.vue     # List
        │   ├── new.vue       # Create
        │   └── [id].vue      # Workbench
        ├── experiences.vue   # Timeline
        ├── skills.vue        # Matrix
        └── profile.vue       # Settings
```

---

## 4. Page Specifications

### 4.1 Landing Page (`/`)

**Layout:** `default`

**Sections:**

1. **Hero Section**
   - Headline: "The CV That Adapts"
   - Subheadline: Brief value proposition
   - CTA Button: "Get Started Free" → `/login`
   - Background: Gradient with subtle chameleon motif

2. **Feature Grid** (3 columns)
   - Feature 1: "AI-Powered Tailoring" (Violet accent)
   - Feature 2: "ATS-Optimized PDFs" (Emerald accent)
   - Feature 3: "Atomic Experience Bullets" (Emerald accent)

3. **How It Works** (3 steps)
   - Step 1: Import your experiences
   - Step 2: Paste a job description
   - Step 3: Download your tailored resume

4. **Call to Action**
   - Final CTA with social proof
   - "Join X professionals using Chameleon Vitae"

**Components:** `HeroSection`, `FeatureGrid`, `CallToAction`

---

### 4.2 Login Page (`/login`)

**Layout:** `default`

**Content:**

1. **Auth Card** (centered)
   - Logo + App name
   - Email/Password form fields
   - "Sign In" button (primary)
   - Divider: "or continue with"
   - Social buttons: Google, GitHub
   - Link: "Don't have an account? Sign up"

**Components:** Uses Nuxt UI `UCard`, `UInput`, `UButton`, `UDivider`

---

### 4.3 Dashboard Overview (`/dashboard`)

**Layout:** `dashboard`

**Content:**

1. **Header**
   - Page title: "Dashboard"
   - Quick action button: "New Resume" (primary)

2. **Stats Grid** (4 cards)
   - Total Resumes (count)
   - Active Applications (count)
   - Experience Bullets (count)
   - Match Score Average (percentage)

3. **Recent Activity**
   - List of last 5 generated resumes
   - Each item: Title, company, date, match score badge

4. **Quick Actions**
   - "Create Resume from URL"
   - "Import from LinkedIn" (disabled, coming soon)

**Components:** `StatsCard`, `ActivityFeed`

---

### 4.4 Resume List (`/dashboard/resumes`)

**Layout:** `dashboard`

**Content:**

1. **Header**
   - Title: "My Resumes"
   - Filters: Status (All, Draft, Generated), Date range
   - Search input
   - "New Resume" button

2. **Resume Table** (UTable)
   - Columns: Title, Company, Score, Status, Created, Actions
   - Row actions: View, Duplicate, Delete

3. **Empty State**
   - Illustration
   - "No resumes yet. Create your first tailored CV!"
   - CTA: "Create Resume"

**Components:** `ResumeCard`, `EmptyState`

---

### 4.5 New Resume (`/dashboard/resumes/new`)

**Layout:** `dashboard`

**Content:**

1. **Header**
   - Title: "Create New Resume"
   - Breadcrumb: Dashboard > Resumes > New

2. **Job Input Card**
   - Tab 1: "Paste Job Description" (textarea)
   - Tab 2: "Import from URL" (input + fetch button)
   - "Analyze Job" button (secondary/violet)

3. **Analysis Preview** (shown after analysis)
   - Extracted: Job Title, Company, Required Skills (badges)
   - Matched bullets preview
   - "Generate Resume" button (primary)

**Components:** `JobDescriptionInput`

---

### 4.6 Resume Workbench (`/dashboard/resumes/[id]`)

**Layout:** `dashboard` (full-width variant)

**Content:**

1. **Split View**
   - Left Panel (50%): Resume Editor
     - Sections: Contact, Summary, Experience, Skills, Education
     - Drag-and-drop bullet reordering
     - Inline editing of text
   - Right Panel (50%): PDF Preview
     - Live-updating iframe/embed
     - Template selector dropdown
     - Download button

2. **Toolbar**
   - Template selector
   - "Regenerate with AI" button
   - "Download PDF" button
   - Match score indicator

**Components:** `ResumeEditor`, `ResumePreview`

---

### 4.7 Experience Timeline (`/dashboard/experiences`)

**Layout:** `dashboard`

**Content:**

1. **Header**
   - Title: "My Experiences"
   - Filter by type: All, Work, Education, Projects
   - "Add Experience" button

2. **Timeline View**
   - Vertical timeline with date markers
   - Each entry: Type icon, Title, Organization, Duration
   - Expand to see bullets
   - Edit/Delete actions

3. **Experience Detail Modal**
   - Edit form for experience
   - Bullet list with add/edit/delete
   - Impact score per bullet

**Components:** `ExperienceTimeline`, `ExperienceCard`, `BulletItem`

---

### 4.8 Skill Matrix (`/dashboard/skills`)

**Layout:** `dashboard`

**Content:**

1. **Header**
   - Title: "Skills & Competencies"
   - "Add Skill" button

2. **Category Sections**
   - Programming Languages
   - Frameworks & Libraries
   - Tools & Platforms
   - Soft Skills

3. **Skill Tags**
   - Each skill as a badge with proficiency indicator
   - Click to edit, hover to see usage count
   - Drag to reorder or change category

**Components:** Uses Nuxt UI `UBadge`, `UCard`

---

### 4.9 Profile Settings (`/dashboard/profile`)

**Layout:** `dashboard`

**Content:**

1. **Header**
   - Title: "Profile Settings"

2. **Sections**
   - **Personal Information**
     - Name, Headline, Location, Photo
   - **Contact Details**
     - Email, Phone, Website, LinkedIn, GitHub
   - **Preferences**
     - Default language, Preferred template
   - **Account**
     - Change password, Delete account

3. **Save Actions**
   - "Save Changes" button
   - "Cancel" button

**Components:** Uses Nuxt UI form components

---

## 5. Navigation Structure

### Sidebar Menu (Dashboard Layout)

```text
┌─────────────────────────────┐
│  🦎 Chameleon Vitae         │
├─────────────────────────────┤
│  📊 Dashboard               │
│  📄 Resumes                 │
│  💼 Experiences             │
│  🏷️ Skills                  │
│  👤 Profile                 │
├─────────────────────────────┤
│  ⚙️ Settings                │
│  🚪 Sign Out                │
└─────────────────────────────┘
```

---

## 6. Responsive Behavior

| Breakpoint | Sidebar | Content |
|------------|---------|---------|
| Desktop (lg+) | Fixed, expanded (240px) | Full width minus sidebar |
| Tablet (md) | Collapsed (64px icons only) | Expanded |
| Mobile (< md) | Hidden (hamburger toggle) | Full width |

---

## 7. Mock Data Requirements

For the "Dumb UI" phase, the following mock data structures are needed:

```typescript
// User
const mockUser = {
  id: 'user-1',
  name: 'John Developer',
  email: 'john@example.com',
  avatar: '/avatars/default.png',
  headline: 'Senior Software Engineer'
}

// Resumes
const mockResumes = [
  { id: '1', title: 'Frontend Developer', company: 'TechCorp', score: 92, status: 'generated', createdAt: '2026-01-08' },
  { id: '2', title: 'Full Stack Engineer', company: 'StartupXYZ', score: 87, status: 'draft', createdAt: '2026-01-07' }
]

// Experiences
const mockExperiences = [
  { id: '1', type: 'work', title: 'Senior Engineer', org: 'BigTech Inc', start: '2022-01', end: null, bullets: [...] },
  { id: '2', type: 'education', title: 'BSc Computer Science', org: 'State University', start: '2016-09', end: '2020-05' }
]

// Skills
const mockSkills = [
  { id: '1', name: 'TypeScript', category: 'Languages', proficiency: 'expert' },
  { id: '2', name: 'Vue.js', category: 'Frameworks', proficiency: 'advanced' }
]
```

---

*Document Version: 1.0*  
*Last Updated: 2026-01-09*
