# Contributing to Chameleon Vitae

First off, thank you for considering contributing to Chameleon Vitae! 🦎

This document provides guidelines for contributing to the project. Following these guidelines helps communicate that you respect the time of the developers managing and developing this open source project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Workflow](#development-workflow)
- [Style Guides](#style-guides)
- [Architecture Guidelines](#architecture-guidelines)

## Code of Conduct

This project and everyone participating in it is governed by our commitment to creating a welcoming and inclusive environment. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.22+** - [Installation Guide](https://go.dev/doc/install)
- **Bun 1.0+** - [Installation Guide](https://bun.sh/)
- **Podman** - [Installation Guide](https://podman.io/docs/installation)
- **podman-compose** - [Installation Guide](https://github.com/containers/podman-compose#installation)

### Setting Up the Development Environment

1. **Fork the repository** on GitHub
2. **Clone your fork**:

   ```bash
   git clone https://github.com/YOUR_USERNAME/chameleon-vitae.git
   cd chameleon-vitae
   ```

3. **Add the upstream remote**:

   ```bash
   git remote add upstream https://github.com/SeltikHD/chameleon-vitae.git
   ```

4. **Copy the environment file**:

   ```bash
   cp .env.example .env
   ```

5. **Start the infrastructure**:

   ```bash
   podman-compose up -d
   ```

6. **Run the backend**:

   ```bash
   go run cmd/server/main.go
   ```

7. **Run the frontend** (in another terminal):

   ```bash
   cd frontend && bun install && bun dev
   ```

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Describe the behavior you observed and what you expected**
- **Include logs, screenshots, or code samples if applicable**
- **Specify your environment** (OS, Go version, Bun version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description of the proposed functionality**
- **Explain why this enhancement would be useful**
- **Consider the impact on the existing architecture**

### Pull Requests

1. **Create a feature branch** from `main`:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**, following our [Style Guides](#style-guides)
3. **Write or update tests** as needed
4. **Ensure all tests pass**:

   ```bash
   go test ./...
   cd frontend && bun test
   ```

5. **Commit your changes** using [Conventional Commits](#commit-messages)
6. **Push to your fork** and submit a pull request

## Development Workflow

### Branch Naming

- `feature/` - New features (e.g., `feature/linkedin-scraping`)
- `fix/` - Bug fixes (e.g., `fix/pdf-rendering-issue`)
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test additions or modifications
- `chore/` - Maintenance tasks

### Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```text
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Types:**

- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation only
- `style` - Formatting, missing semicolons, etc.
- `refactor` - Code change that neither fixes a bug nor adds a feature
- `perf` - Performance improvement
- `test` - Adding or correcting tests
- `chore` - Maintenance tasks

**Examples:**

```text
feat(core): add bullet selection algorithm
fix(pdf): resolve font rendering issue in Gotenberg
docs(readme): update installation instructions
refactor(adapters): simplify Groq API client
```

## Style Guides

### Go Code

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Write table-driven tests
- Document all exported functions and types

### TypeScript/Vue Code

- Follow the project's ESLint configuration
- Use TypeScript strict mode
- Follow Vue.js 3 Composition API patterns
- Use `<script setup>` syntax in Vue components

### Documentation

- Write in **English only**
- Use Markdown for all documentation
- Include code examples where appropriate
- Keep the README up to date

## Architecture Guidelines

Chameleon Vitae follows **Hexagonal Architecture** (Ports and Adapters). This is crucial to understand before contributing:

### Core Principles

1. **The Core Domain is Sacred**
   - The `internal/core/` directory contains pure business logic
   - **NEVER** import external adapters into the core
   - The core should only depend on its own interfaces (ports)

2. **Adapters are Interchangeable**
   - All external dependencies (database, AI, PDF engine) are adapters
   - Adapters implement the interfaces defined in the core
   - You should be able to swap adapters without touching core logic

3. **Dependency Direction**
   - Dependencies point **inward** toward the core
   - The core knows nothing about the outside world
   - Adapters depend on the core, not the other way around

### Directory Structure

```text
internal/
├── core/
│   ├── domain/        # Entities and value objects
│   ├── ports/         # Interfaces (input and output)
│   └── services/      # Use cases / application services
└── adapters/
    ├── primary/       # Input adapters (HTTP handlers)
    └── secondary/     # Output adapters (DB, AI, PDF)
```

### Before Submitting

Ask yourself:

- Does my code maintain the architectural boundaries?
- Have I written tests for my changes?
- Is my code documented?
- Does it follow the existing patterns in the codebase?

---

Thank you for contributing to Chameleon Vitae! 🎉
