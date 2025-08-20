# Promptly

![Build Status](https://github.com/rahulguha/promptly/actions/workflows/ci.yml/badge.svg)

# Promptly

A prompt management application with persona-based templating for AI interactions. Create reusable prompt templates with variable substitution and manage different personas for various use cases. User can play a role and may want LLM to play another role. For example, an user may want to play role of a **_software developer_** and may want LLM to play role of a **_code reviewer_**.

## Features

![alt text](image.png)

- **Persona Management**: Define user and LLM roles for different contexts
- **Template System**: Create reusable prompt templates with `{{variable}}` placeholders
- **Template Version Management**: Maintain versions of templates and edit them separately. Track generated prompts
- **Prompt Generation**: Generate final prompts by substituting template variables
- **Web Interface**: Modern SvelteKit frontend for easy management
- **REST API**: Full CRUD operations for all entities
- **Storage**:
  - **JsonStore**: Simple file-based persistence
  - **SQLite store**

## Use Case

- **Generic RAG Application**: User adds their artifacts and later wants to query that corpus. However they want to use sophesticated prompts to customize response. Example: A High School student may want to generate Mock Tests for his preparation.
- **Interviewee**: User adds their artifacts and later wants to query that corpus. They want to simulate interview process using LLMs ability to play a role of an interviewer.
- **Pair Programming**: User may want to use LLM for pair programming by asking LLM to play that role.

## Quick Start

### Prerequisites

- Go 1.22+
- Node.js 18+

### Installation

1. Clone the repository:

```bash
git clone https://github.com/rahulguha/promptly.git
cd promptly
```

2. Build the Go backend:

```bash
go build -o promptly cmd/promptly/main.go
```

3. Start the API server:

```bash
./promptly serve
```

4. Set up the web interface (in a new terminal):

```bash
cd promptly-web
npm install
npm run dev
```

The API will be available at `http://localhost:8080` and the web interface at `http://localhost:5173`.

## Usage

### CLI Options

```bash
# Start server on default port 8080
./promptly serve

# Custom port and data path
./promptly serve --port 3000 --data ./custom-data/prompts.json

# Help
./promptly --help
```

### Workflow

1. **Create Personas**: Define user roles (e.g., "Developer") and corresponding LLM roles (e.g., "Code Reviewer")
2. **Build Templates**: Create prompt templates linked to personas with variable placeholders like `{{project_name}}`
3. **Generate Prompts**: Substitute variables with actual values to create final prompts

### API Endpoints

- `GET/POST/PUT/DELETE /v1/personas` - Manage user/LLM role definitions
- `GET/POST/PUT/DELETE /v1/templates` - Manage prompt templates
- `GET/POST/PUT/DELETE /v1/prompts` - Manage generated prompts
- `POST /v1/generate-prompt` - Generate prompts from templates
- `GET /health` - Health check

## Development

### Backend

```bash
# Run tests
go test ./...

# Build
go build -o promptly cmd/promptly/main.go
```

### Frontend

```bash
cd promptly-web

# Development server
npm run dev

# Build for production
npm run build

# Preview build
npm run preview
```

## Configuration

The application looks for configuration in:

- Command line flags (`--port`, `--data`)
- Config file (`config.yaml` or `.promptly`)
- Environment variables

Default data directory: `data/` (creates `prompts.json`, `prompt_template.json`, `persona.json`)

## Documentation

- [API Documentation](API.md) - Complete REST API reference
- [Examples](examples/) - Usage examples and demo data

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

Licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.
