# 🦎 Chameleon Vitae

> **"The CV that adapts."**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Vue Version](https://img.shields.io/badge/Vue.js-3.0-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![Nuxt](https://img.shields.io/badge/Nuxt-3-00DC82?style=flat&logo=nuxt.js)](https://nuxt.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Architecture](https://img.shields.io/badge/Architecture-Hexagonal-orange)](https://alistair.cockburn.us/hexagonal-architecture/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Chameleon Vitae** is an open-source resume engineering system designed to beat Applicant Tracking Systems (ATS). Instead of maintaining a single static resume, the system maintains a database of atomic "experience bullets" and uses Artificial Intelligence (LLM) to assemble, in real-time, the perfect version of your profile for each specific job posting.

## ✨ Key Features

- **Atomic Modularity:** Your experiences are not fixed text blocks. They are split into independent bullets (topics) in the database.
- **AI Tailoring (via Groq):** Analyzes the job description (Markdown) and selects/rewrites the most relevant bullets from your history.
- **ATS-Friendly PDF Generation:** High-fidelity rendering using Headless Chrome (via Gotenberg), ensuring robots can read the file correctly.
- **Hexagonal Architecture:** Domain core isolated from external dependencies, facilitating testing and swapping AI providers.

## 🛠️ Tech Stack

| Layer          | Technology        | Rationale                                                   |
| :------------- | :---------------- | :---------------------------------------------------------- |
| **Backend**    | **Golang + Chi**  | Performance, strong typing, and idiomatic simplicity.       |
| **Frontend**   | **Vue.js + Nuxt** | Static Site Generation (SSG) and reactive components.       |
| **Database**   | **PostgreSQL 17** | Relational robustness and JSONB support for metadata.       |
| **Infra**      | **Podman**        | Secure (rootless) containers and lightweight orchestration. |
| **AI/LLM**     | **Groq API**      | Llama 3/4 inference at extreme speed (LPU).                 |
| **PDF Engine** | **Gotenberg 8**   | Dockerized API for reliable HTML -> PDF conversion.         |

## 🏗️ Project Structure (Hexagonal Architecture)

```text
chameleon-vitae/
├── cmd/
│   └── server/              # Application entrypoint
│       └── main.go
├── internal/
│   ├── core/                # 🔒 PURE DOMAIN — NO EXTERNAL DEPENDENCIES
│   │   ├── domain/          # Entities, Value Objects, Domain Errors
│   │   ├── ports/           # Interfaces (Input & Output Ports)
│   │   └── services/        # Application Services / Use Cases
│   └── adapters/
│       ├── primary/         # Input Adapters (HTTP handlers, CLI)
│       │   └── http/        # Chi router handlers
│       └── secondary/       # Output Adapters (implementations)
│           ├── postgres/    # Database adapter
│           ├── groq/        # AI provider adapter
│           └── gotenberg/   # PDF engine adapter
├── pkg/                     # Shared utilities (can be imported by adapters)
├── frontend/                # Nuxt.js application
└── deploy/                  # Infrastructure (Compose, Dockerfiles)
    └── postgres/
        └── init/            # Database initialization scripts
```

## ⚡ Getting Started

### Prerequisites

- **Go 1.22+** — [Installation Guide](https://go.dev/doc/install)
- **Bun** - [Installation Guide](https://bun.sh/)
- **Podman** and **podman-compose** — [Installation Guide](https://podman.io/docs/installation)

### Quick Start

**1. Clone the repository:**

```bash
git clone https://github.com/SeltikHD/chameleon-vitae.git
cd chameleon-vitae
```

**2. Set up environment variables:**

```bash
cp .env.example .env
# Edit .env with your Groq API key
```

**3. Start the infrastructure (Database + PDF Engine):**

```bash
podman-compose up -d
```

**4. Run the backend:**

```bash
go run cmd/server/main.go
```

**5. Run the frontend** (in another terminal):

```bash
cd frontend
bun install
bun dev
```

### Environment Variables

| Variable                | Description                          | Default                     |
| ----------------------- | ------------------------------------ | --------------------------- |
| `GROQ_API_KEY`          | Your Groq API key                    | (required)                  |
| `POSTGRES_HOST`         | PostgreSQL host                      | `localhost`                 |
| `POSTGRES_PORT`         | PostgreSQL port                      | `5432`                      |
| `POSTGRES_USER`         | PostgreSQL user                      | `chameleon`                 |
| `POSTGRES_PASSWORD`     | PostgreSQL password                  | `changeme_in_production`    |
| `POSTGRES_DB`           | PostgreSQL database                  | `chameleon_vitae`           |
| `GOTENBERG_URL`         | Gotenberg service URL                | `http://localhost:3000`     |

## 📖 Documentation

- [Architecture Decisions (ADRs)](DECISIONS.md)
- [Contributing Guide](CONTRIBUTING.md)
- [GitHub Copilot Instructions](.github/copilot-instructions.md)

## 🗺️ Roadmap

- [ ] **MVP:** Resume generation based on Job Description (Markdown)
- [ ] LinkedIn profile scraping integration
- [ ] "Auto-Apply" module (automatic submission)
- [ ] Local LLM support via Ollama
- [ ] Resume template marketplace
- [ ] Analytics dashboard for job applications

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) before submitting a Pull Request.

## 📄 License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

---

Built with ☕ and a healthy disdain for HR forms.
