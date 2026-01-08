# Chameleon Vitae - REST API Specification

> **Version:** 1.0.0  
> **Base URL:** `http://localhost:8080/v1`  
> **Authentication:** Bearer Token (Firebase JWT)

This document defines the complete REST API contract for Chameleon Vitae. All endpoints except `/health` require authentication via Firebase JWT tokens in the `Authorization` header.

---

## Table of Contents

1. [Authentication](#1-authentication)
2. [User Profile](#2-user-profile)
3. [Experiences](#3-experiences)
4. [Bullets](#4-bullets)
5. [Skills](#5-skills)
6. [Spoken Languages](#6-spoken-languages)
7. [Resume Engine](#7-resume-engine)
8. [Tools](#8-tools)
9. [Common Response Formats](#9-common-response-formats)

---

## Authentication Header

All authenticated endpoints require:

```http
Authorization: Bearer <firebase_jwt_token>
```

---

## 1. Authentication

### POST `/auth/sync`

Synchronize Firebase user with the local PostgreSQL database. Called after successful Firebase authentication on the frontend.

**Request Body:**

```json
{
  "firebase_uid": "string (from Firebase)",
  "email": "string (optional)",
  "name": "string (optional)",
  "picture": "string (optional, profile picture URL)"
}
```

**Response:** `201 Created`

```json
{
  "id": "uuid",
  "firebase_uid": "string",
  "email": "string",
  "name": "string",
  "created_at": "ISO8601 timestamp"
}
```

**Notes:**

- Creates user if not exists, updates if exists (upsert behavior).
- This is the first call after Firebase login.

---

## 2. User Profile

### GET `/me`

Get the authenticated user's complete profile.

**Response:** `200 OK`

```json
{
  "id": "uuid",
  "firebase_uid": "string",
  "picture_url": "string",
  "email": "string",
  "name": "string",
  "headline": "string",
  "summary": "string",
  "location": "string",
  "phone": "string",
  "website": "string",
  "linkedin_url": "string",
  "github_url": "string",
  "portfolio_url": "string",
  "preferred_language": "en | pt-br",
  "created_at": "ISO8601 timestamp",
  "updated_at": "ISO8601 timestamp"
}
```

### PATCH `/me`

Update user profile fields.

**Request Body:** (all fields optional)

```json
{
  "name": "string",
  "headline": "string",
  "summary": "string",
  "location": "string",
  "phone": "string",
  "website": "string",
  "linkedin_url": "string",
  "github_url": "string",
  "portfolio_url": "string",
  "preferred_language": "en | pt-br"
}
```

**Response:** `200 OK` (returns updated user object)

---

## 3. Experiences

### GET `/experiences`

List all experiences for the authenticated user.

**Query Parameters:**

| Parameter | Type   | Description                          |
| --------- | ------ | ------------------------------------ |
| `type`    | string | Filter by experience type (optional) |
| `limit`   | int    | Pagination limit (default: 50)       |
| `offset`  | int    | Pagination offset (default: 0)       |

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "uuid",
      "type": "work | education | certification | project | freelance | volunteer | open_source | hackathon | side_project | event_organization | publication | award",
      "title": "string",
      "organization": "string",
      "location": "string",
      "start_date": "YYYY-MM-DD",
      "end_date": "YYYY-MM-DD | null",
      "is_current": true,
      "description": "string",
      "url": "string",
      "metadata": {},
      "display_order": 0,
      "bullets": [
        {
          "id": "uuid",
          "content": "string",
          "impact_score": 75,
          "keywords": ["string"],
          "display_order": 0
        }
      ],
      "created_at": "ISO8601",
      "updated_at": "ISO8601"
    }
  ],
  "total": 15,
  "limit": 50,
  "offset": 0
}
```

### POST `/experiences`

Create a new experience.

**Request Body:**

```json
{
  "type": "work",
  "title": "Senior Software Engineer",
  "organization": "Tech Company Inc.",
  "location": "Remote",
  "start_date": "2022-01-15",
  "end_date": null,
  "is_current": true,
  "description": "Leading backend development...",
  "url": "https://techcompany.com",
  "metadata": {}
}
```

**Response:** `201 Created` (returns created experience)

### PUT `/experiences/{id}`

Update an existing experience.

**Path Parameters:**

- `id` (uuid): Experience ID

**Request Body:** Same as POST (all fields optional except `id`)

**Response:** `200 OK` (returns updated experience)

### DELETE `/experiences/{id}`

Delete an experience and all its bullets.

**Path Parameters:**

- `id` (uuid): Experience ID

**Response:** `204 No Content`

---

## 4. Bullets

Bullets are atomic units of experience that can be individually selected for resume tailoring.

### POST `/experiences/{experience_id}/bullets`

Add a new bullet to an experience.

**Path Parameters:**

- `experience_id` (uuid): Parent experience ID

**Request Body:**

```json
{
  "content": "Reduced API latency by 40% through database query optimization",
  "keywords": ["performance", "optimization", "database"],
  "display_order": 0
}
```

**Response:** `201 Created`

```json
{
  "id": "uuid",
  "experience_id": "uuid",
  "content": "string",
  "impact_score": 50,
  "keywords": ["string"],
  "metadata": {},
  "display_order": 0,
  "created_at": "ISO8601",
  "updated_at": "ISO8601"
}
```

**Note:** `impact_score` starts at 50 (neutral) and is recalculated by AI.

### PUT `/bullets/{id}`

Update an existing bullet.

**Path Parameters:**

- `id` (uuid): Bullet ID

**Request Body:**

```json
{
  "content": "Updated bullet content",
  "keywords": ["new", "keywords"],
  "display_order": 1
}
```

**Response:** `200 OK` (returns updated bullet)

### DELETE `/bullets/{id}`

Delete a bullet.

**Path Parameters:**

- `id` (uuid): Bullet ID

**Response:** `204 No Content`

### POST `/bullets/{id}/score`

Trigger AI recalculation of the bullet's impact score.

**Path Parameters:**

- `id` (uuid): Bullet ID

**Response:** `200 OK`

```json
{
  "id": "uuid",
  "content": "string",
  "impact_score": 85,
  "score_reasoning": "Contains quantifiable metrics (40%), demonstrates leadership..."
}
```

---

## 5. Skills

### GET `/skills`

List all skills for the authenticated user.

**Query Parameters:**

| Parameter  | Type   | Description                   |
| ---------- | ------ | ----------------------------- |
| `category` | string | Filter by category (optional) |

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Go",
      "category": "Programming Languages",
      "proficiency_level": 85,
      "years_of_experience": 3.5,
      "is_highlighted": true,
      "display_order": 0,
      "created_at": "ISO8601"
    }
  ],
  "total": 25
}
```

### POST `/skills/batch`

Batch create or update skills. Uses upsert logic based on `name`.

**Request Body:**

```json
{
  "skills": [
    {
      "name": "Go",
      "category": "Programming Languages",
      "proficiency_level": 85,
      "years_of_experience": 3.5,
      "is_highlighted": true
    },
    {
      "name": "PostgreSQL",
      "category": "Databases",
      "proficiency_level": 75,
      "years_of_experience": 5.0,
      "is_highlighted": false
    }
  ]
}
```

**Response:** `200 OK`

```json
{
  "created": 1,
  "updated": 1,
  "data": [...]
}
```

### DELETE `/skills/{id}`

Delete a skill.

**Path Parameters:**

- `id` (uuid): Skill ID

**Response:** `204 No Content`

---

## 6. Spoken Languages

### GET `/languages`

List all spoken languages for the authenticated user.

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "uuid",
      "language": "English",
      "proficiency": "native | fluent | advanced | intermediate | basic",
      "display_order": 0,
      "created_at": "ISO8601"
    }
  ]
}
```

### POST `/languages`

Add a spoken language.

**Request Body:**

```json
{
  "language": "Portuguese",
  "proficiency": "native",
  "display_order": 0
}
```

**Response:** `201 Created`

### DELETE `/languages/{id}`

Delete a spoken language.

**Response:** `204 No Content`

---

## 7. Resume Engine

### GET `/resumes`

List all generated resumes (history).

**Query Parameters:**

| Parameter | Type   | Description                    |
| --------- | ------ | ------------------------------ |
| `status`  | string | Filter by status (optional)    |
| `limit`   | int    | Pagination limit (default: 20) |
| `offset`  | int    | Pagination offset (default: 0) |

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "uuid",
      "job_title": "Senior Backend Engineer",
      "company_name": "Awesome Corp",
      "job_url": "https://linkedin.com/jobs/...",
      "target_language": "en",
      "score": 85,
      "status": "draft | generated | reviewed | submitted | interview | rejected | accepted",
      "created_at": "ISO8601",
      "updated_at": "ISO8601"
    }
  ],
  "total": 10
}
```

### POST `/resumes`

Create a new resume draft from a job description.

**Request Body:**

```json
{
  "job_description": "## Senior Backend Engineer\n\nWe are looking for...",
  "job_title": "Senior Backend Engineer",
  "company_name": "Awesome Corp",
  "job_url": "https://linkedin.com/jobs/12345",
  "target_language": "en"
}
```

**Response:** `201 Created`

```json
{
  "id": "uuid",
  "job_description": "string",
  "job_title": "string",
  "company_name": "string",
  "job_url": "string",
  "target_language": "en",
  "status": "draft",
  "created_at": "ISO8601"
}
```

### GET `/resumes/{id}`

Get a specific resume with all details.

**Response:** `200 OK`

```json
{
  "id": "uuid",
  "job_description": "string",
  "job_title": "string",
  "company_name": "string",
  "job_url": "string",
  "target_language": "en",
  "selected_bullets": ["uuid", "uuid"],
  "generated_content": {
    "summary": "Tailored professional summary...",
    "experiences": [
      {
        "experience_id": "uuid",
        "title": "string",
        "organization": "string",
        "bullets": [
          {
            "bullet_id": "uuid",
            "original_content": "string",
            "tailored_content": "string"
          }
        ]
      }
    ],
    "skills": ["Go", "PostgreSQL", "Docker"]
  },
  "pdf_url": "https://storage.../resume.pdf",
  "score": 85,
  "notes": "User notes",
  "status": "generated",
  "created_at": "ISO8601",
  "updated_at": "ISO8601"
}
```

### POST `/resumes/{id}/tailor`

Trigger AI to analyze the job description, select relevant bullets, and generate tailored content.

**Path Parameters:**

- `id` (uuid): Resume ID

**Request Body:** (optional overrides)

```json
{
  "max_bullets_per_experience": 5,
  "include_experience_types": ["work", "project", "open_source"],
  "highlight_skills": ["Go", "Kubernetes"]
}
```

**Response:** `200 OK`

```json
{
  "id": "uuid",
  "status": "generated",
  "score": 85,
  "selected_bullets": ["uuid", "uuid", "uuid"],
  "generated_content": { ... },
  "analysis": {
    "matched_keywords": ["golang", "microservices", "kubernetes"],
    "missing_keywords": ["terraform"],
    "recommendations": ["Consider adding cloud infrastructure experience"]
  }
}
```

### PATCH `/resumes/{id}/content`

Manually edit the generated content.

**Request Body:**

```json
{
  "generated_content": { ... },
  "notes": "Made adjustments to summary",
  "status": "reviewed"
}
```

**Response:** `200 OK`

### GET `/resumes/{id}/pdf`

Generate and download the PDF version of the resume.

**Query Parameters:**

| Parameter  | Type   | Description                                  |
| ---------- | ------ | -------------------------------------------- |
| `template` | string | Template name (default: "modern")            |
| `format`   | string | Paper format: "a4" or "letter" (default: a4) |

**Response:** `200 OK`

- Content-Type: `application/pdf`
- Content-Disposition: `attachment; filename="resume-{id}.pdf"`

---

## 8. Tools

### POST `/tools/parse-job`

Parse a job posting URL into LLM-friendly Markdown using Jina Reader.

**Request Body:**

```json
{
  "url": "https://linkedin.com/jobs/view/12345"
}
```

**Response:** `200 OK`

```json
{
  "url": "string",
  "title": "Senior Backend Engineer at Awesome Corp",
  "markdown": "## Senior Backend Engineer\n\n**Company:** Awesome Corp\n\n**Location:** Remote\n\n### Requirements\n\n- 5+ years of experience with Go...",
  "metadata": {
    "source": "linkedin.com",
    "fetched_at": "ISO8601"
  }
}
```

**Notes:**

- Uses Jina Reader API (`r.jina.ai`) under the hood.
- Supports LinkedIn, Gupy, Indeed, and other job boards.

---

## 9. Common Response Formats

### Success Response

```json
{
  "data": { ... },
  "message": "Operation successful"
}
```

### Error Response

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  }
}
```

### HTTP Status Codes

| Code | Description                              |
| ---- | ---------------------------------------- |
| 200  | OK - Request succeeded                   |
| 201  | Created - Resource created successfully  |
| 204  | No Content - Deletion successful         |
| 400  | Bad Request - Invalid input              |
| 401  | Unauthorized - Missing or invalid token  |
| 403  | Forbidden - Insufficient permissions     |
| 404  | Not Found - Resource doesn't exist       |
| 409  | Conflict - Resource already exists       |
| 422  | Unprocessable Entity - Validation failed |
| 429  | Too Many Requests - Rate limit exceeded  |
| 500  | Internal Server Error                    |

---

## Rate Limiting

| Endpoint Category | Rate Limit  |
| ----------------- | ----------- |
| Authentication    | 10 req/min  |
| Read operations   | 100 req/min |
| Write operations  | 30 req/min  |
| AI operations     | 10 req/min  |
| PDF generation    | 5 req/min   |

Rate limit headers are included in all responses:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1704672000
```

---

## Versioning

The API version is included in the URL path (`/v1/`). Breaking changes will result in a new version (`/v2/`). Non-breaking additions (new fields, new endpoints) will not increment the version.

---

Last Updated: 2026-01-08
