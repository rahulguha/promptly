# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Promptly is a sophisticated prompt management application with both a Go backend and SvelteKit frontend. It provides REST APIs for managing personas, prompt templates, and generated prompts with a web interface for interactive use. The system supports user authentication, multi-tenant data isolation, and LLM integration for prompt evaluation.

## Architecture

### Backend (Go)
- **Main entry**: `cmd/promptly/main.go` - CLI with Cobra, serves HTTP API with Gin
- **Models**: `internal/models/` - Core data structures (Persona, PromptTemplate, Prompt, Profile, Intent)
- **Storage**: Multiple storage backends with factory pattern:
  - `internal/storage/jsonstore/` - JSON file-based storage
  - `internal/storage/sqlite/` - SQLite database storage (default)
  - `internal/storage/factory.go` - Storage abstraction layer
- **API**: `internal/routes/` - RESTful handlers for all entities with session-based authentication
- **Authentication**: `internal/api/` - OAuth-based authentication with Google
- **Tracking**: `internal/tracking/` - User activity tracking with AWS DynamoDB
- **LLM Integration**: `internal/llm/` - Prompt evaluation with OpenAI and Claude APIs
- **Configuration**: `internal/config/` - Viper-based configuration management
- **Data files**: `data/` directory contains SQLite databases and JSON files for persistence

### Frontend (SvelteKit)
- **Location**: `promptly-web/` directory
- **Framework**: SvelteKit 2.0 with TypeScript and Vite
- **API client**: `src/lib/api.ts` - TypeScript interfaces and HTTP client functions
- **Components**: Svelte components for managing personas, templates, and prompt generation

### Multi-Tenant Architecture
- User-specific SQLite databases created per authenticated user
- Session-based authentication with Google OAuth
- DB middleware creates user-specific database connections
- Profile-based data isolation within user databases

## Development Commands

### Go Backend
```bash
# Build the CLI
go build -o promptly cmd/promptly/main.go

# Run the server (default: SQLite storage, port 8080)
./promptly serve

# Run with custom configuration
./promptly serve --port 3000 --storage sqlite --db ./custom/path.db

# Run with OpenAI LLM provider
./promptly serve --llm-provider openai

# Run with Claude LLM provider
./promptly serve --llm-provider claude

# Run tests
go test ./...
```

### Command Line Options
- `--port` or `-p`: Port to run the server on (default: 8080)
- `--storage` or `-s`: Storage type, either `json` or `sqlite` (default: `sqlite`)
- `--data` or `-d`: Path to data file for json storage (default: `data/prompts.json`)
- `--db` or `-b`: Path to database file for sqlite storage (default: `data/promptly.db`)
- `--llm-provider`: LLM provider to use (`openai` or `claude`)

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

## API Structure

The backend exposes a versioned REST API at `/v1/` with these endpoints:

### Core Resources
- `/v1/personas` - CRUD operations for user/LLM role definitions
- `/v1/templates` - CRUD operations for prompt templates with variables
  - `POST /:id/version` - Create new template versions
- `/v1/prompts` - CRUD operations for generated prompts
  - `POST /evaluate` - Evaluate prompts using LLM providers
- `/v1/profiles` - CRUD operations for user profiles
- `/v1/generate-prompt` - Generate final prompts from templates + values

### Additional Endpoints
- `/v1/intents` - Get user intents for analytics
- `/v1/track/users` - Track user data for analytics
- `/v1/track/activity` - Track user activity for analytics
- `/v1/api/auth/*` - OAuth authentication routes (login, callback, me, logout)
- `/health` - Health check endpoint

## Key Patterns

### Data Flow
1. **Personas** define user and LLM roles (e.g., "Developer" → "Code Reviewer")
2. **Templates** contain parameterized prompts linked to personas with `{{variable}}` placeholders
3. **Prompts** are generated from templates by substituting variables with actual values
4. **LLM Evaluation** provides automated grading and improvement suggestions for prompts
5. Persona context is automatically prepended to templates during prompt generation

### Storage Pattern
- Multi-backend storage support (SQLite default, JSON fallback)
- User-specific database isolation for multi-tenancy
- UUID-based entity identification throughout the system
- Auto-generation of UUIDs for new entities
- Concurrent read/write protection via database transactions

### Authentication & Sessions
- Google OAuth integration for user authentication
- Session-based state management with Gin sessions
- CORS configured for frontend development (localhost:5175)
- User-specific database connections via middleware

### LLM Integration
- Pluggable evaluator interface supporting multiple providers
- OpenAI GPT-3.5-turbo integration for prompt evaluation
- Claude API integration for prompt evaluation
- Structured JSON responses for evaluation results

### Error Handling
- Consistent HTTP status codes (400 for bad requests, 404 for not found, 500 for server errors)
- UUID validation on all ID parameters
- JSON binding validation for request bodies
- Proper error wrapping with context using pkg/errors

## Environment Configuration

### Required Environment Variables
- `OPENAI_API_KEY` - Required for OpenAI LLM evaluation
- `ANTHROPIC_API_KEY` - Required for Claude LLM evaluation
- `GOOGLE_CLIENT_ID` - Required for OAuth authentication
- `GOOGLE_CLIENT_SECRET` - Required for OAuth authentication
- `SESSION_SECRET` - Required for session encryption

### DynamoDB Configuration (Optional)
- `DYNAMODB_REGION` - AWS region for DynamoDB
- `DYNAMODB_TABLE_NAME` - DynamoDB table for user tracking
- `DYNAMODB_ACTIVITY_TABLE_NAME` - DynamoDB table for activity tracking

## Testing

- Uses `testify` package for Go testing
- Test files located alongside source files (e.g., `jsonstore_test.go`)
- Run all tests: `go test ./...`
- Frontend testing framework to be determined

## Dependencies

### Go Backend
- **Web Framework**: Gin with CORS and session middleware
- **CLI**: Cobra with Viper configuration
- **Database**: SQLite with modernc.org/sqlite driver
- **Authentication**: OAuth2 with Google provider
- **LLM APIs**: OpenAI and Anthropic SDKs
- **AWS**: AWS SDK for DynamoDB tracking
- **Testing**: Testify for unit tests

### Frontend
- **Framework**: SvelteKit 2.0 with Svelte 5.0
- **Build Tool**: Vite 6.x
- **Language**: TypeScript 5.x
- **Adapter**: SvelteJS auto adapter

## important-instruction-reminders
Do what has been asked; nothing more, nothing less.
NEVER create files unless they're absolutely necessary for achieving your goal.
ALWAYS prefer editing an existing file to creating a new one.
NEVER proactively create documentation files (*.md) or README files. Only create documentation files if explicitly requested by the User.