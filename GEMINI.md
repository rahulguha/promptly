# Project Overview

This is a Go project with a SvelteKit frontend. The project, named "Promptly", is a prompt management application with profile and persona-based templating for AI interactions. It allows users to create reusable personas and prompt templates with variable substitution for various use cases.

The backend is a Go application using the Gin web framework. It uses Cobra for command-line argument parsing and Viper for configuration management. It supports two types of storage: a simple JSON file-based storage and a SQLite database. The backend also integrates with AWS DynamoDB for user tracking.

The frontend is a SvelteKit application, which communicates with the Go backend via a REST API.

# Building and Running

## Backend

To build the backend, run the following command:

```bash
go build -o promptly cmd/promptly/main.go
```

To run the backend, use the `serve` command:

```bash
./promptly serve
```

The server can be configured with the following flags:

*   `--port` or `-p`: Port to run the server on (default: 8080)
*   `--storage` or `-s`: Storage type, either `json` or `sqlite` (default: `sqlite`)
*   `--data` or `-d`: Path to data file for json storage (default: `data/prompts.json`)
*   `--db` or `-b`: Path to database file for sqlite storage (default: `data/promptly.db`)

## Frontend

To run the frontend in development mode, first navigate to the `promptly-web` directory, then run the following commands:

```bash
npm install
npm run dev
```

# Development Conventions

## Backend

The backend code is organized into several packages:

*   `cmd/promptly`: The main application entry point.
*   `internal/api`: Handles authentication and user tracking.
*   `internal/config`: Manages application configuration.
*   `internal/models`: Defines the data models.
*   `internal/routes`: Defines the API routes and handlers.
*   `internal/storage`: Implements the storage layer, with sub-packages for `jsonstore` and `sqlite`.
*   `internal/tracking`: Handles user activity tracking with DynamoDB.

The backend uses the `testify` package for testing. Tests can be run with the `go test ./...` command.

## Frontend

The frontend code is located in the `promptly-web` directory. It is a SvelteKit application. The main source code is in the `src` directory.

## API

The API is versioned under `/v1`. The main routes are:

*   `/v1/personas`: Manage personas
*   `/v1/templates`: Manage prompt templates
*   `/v1/prompts`: Manage prompts
*   `/v1/profiles`: Manage user profiles
*   `/v1/generate-prompt`: Generate a prompt from a template
*   `/v1/track/users`: Track user data
*   `/v1/track/activity`: Track user activity
*   `/v1/api/auth`: Authentication routes

A health check endpoint is available at `/health`.
