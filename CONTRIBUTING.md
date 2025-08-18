# Contributing to Promptly

Thank you for your interest in contributing to Promptly! This document provides guidelines for contributing to the project.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/yourusername/promptly.git
   cd promptly
   ```

3. Set up the development environment:
   ```bash
   # Backend
   go mod download
   go build -o promptly cmd/promptly/main.go

   # Frontend
   cd promptly-web
   npm install
   ```

## Development Workflow

### Backend (Go)
- Follow standard Go conventions and formatting (`go fmt`)
- Write tests for new functionality
- Run tests: `go test ./...`
- Build: `go build -o promptly cmd/promptly/main.go`

### Frontend (SvelteKit)
- Use TypeScript for new code
- Follow existing component patterns
- Test in browser with: `npm run dev`
- Build: `npm run build`

## Making Changes

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the coding standards
3. Test your changes thoroughly
4. Commit with clear, descriptive messages
5. Push to your fork and create a Pull Request

## Pull Request Guidelines

- Provide a clear description of the changes
- Reference any related issues
- Ensure all tests pass
- Keep changes focused and atomic
- Update documentation if needed

## Code Style

### Go
- Use `go fmt` for formatting
- Follow effective Go guidelines
- Use meaningful variable and function names
- Add comments for exported functions

### TypeScript/Svelte
- Use TypeScript for type safety
- Follow existing component structure
- Use descriptive component and variable names
- Keep components focused and reusable

## Project Structure

- `cmd/promptly/` - Main CLI application
- `internal/` - Internal Go packages (models, routes, storage)
- `promptly-web/` - SvelteKit frontend application
- `data/` - JSON data files for development

## Reporting Issues

When reporting issues, please include:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Environment details (Go version, Node version, OS)

## Questions?

Feel free to open an issue for discussion before starting work on major features.

## License

By contributing to Promptly, you agree that your contributions will be licensed under the Apache License 2.0.