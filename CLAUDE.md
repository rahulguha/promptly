# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Promptly is a prompt management application with both a Go backend and SvelteKit frontend. It provides REST APIs for managing personas, prompt templates, and generated prompts with a web interface for interactive use.

## Architecture

### Backend (Go)
- **Main entry**: `cmd/promptly/main.go` - CLI with Cobra, serves HTTP API with Gin
- **Models**: `internal/models/prompt.go` - Core data structures (Persona, PromptTemplate, Prompt)
- **Storage**: `internal/storage/jsonstore/` - JSON file-based storage with concurrent access safety
- **API**: `internal/routes/` - RESTful handlers for all entities
- **Data files**: `data/` directory contains JSON files for persistence

### Frontend (SvelteKit)
- **Location**: `promptly-web/` directory
- **Framework**: SvelteKit with TypeScript and Vite
- **API client**: `src/lib/api.ts` - TypeScript interfaces and HTTP client functions
- **Components**: Svelte components for managing personas, templates, and prompt generation

### API Structure
The backend exposes a versioned REST API at `/v1/` with these endpoints:
- `/v1/personas` - CRUD operations for user/LLM role definitions
- `/v1/templates` - CRUD operations for prompt templates with variables
- `/v1/prompts` - CRUD operations for generated prompts
- `/v1/generate-prompt` - Generate final prompts from templates + values

## Development Commands

### Go Backend
```bash
# Build the CLI
go build -o promptly cmd/promptly/main.go

# Run the server (default port 8080)
./promptly serve

# Run with custom port and data path
./promptly serve --port 3000 --data ./custom-data/prompts.json

# Run tests
go test ./...
```

### Frontend
```bash
cd promptly-web

# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Key Patterns

### Data Flow
1. **Personas** define user and LLM roles (e.g., "Developer" → "Code Reviewer")
2. **Templates** contain parameterized prompts linked to personas with `{{variable}}` placeholders
3. **Prompts** are generated from templates by substituting variables with actual values
4. Persona context is automatically prepended to templates during prompt generation

### Storage Pattern
- JSON file storage with concurrent read/write protection via mutex
- Separate files: `prompts.json`, `prompt_template.json`, `persona.json`
- UUID-based entity identification throughout the system
- Auto-generation of UUIDs for new entities

### Error Handling
- Consistent HTTP status codes (400 for bad requests, 404 for not found, 500 for server errors)
- UUID validation on all ID parameters
- JSON binding validation for request bodies