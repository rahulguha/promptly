# Examples

This directory contains example data and usage scenarios for Promptly.

## Demo Data

The `data/` directory already contains sample personas and templates for educational use cases:

### Personas
- **High School Student → Patient Teacher**: For learning assistance
- **High School Student → Mock Tester**: For practice questions  
- **College Student → Professor**: For advanced learning

### Templates
- Teaching templates with grade-level variables
- Question generation with difficulty and quantity controls
- Detailed explanations for college-level topics

## API Usage Examples

### 1. Create a Code Review Persona

```bash
curl -X POST http://localhost:8080/v1/personas \
  -H "Content-Type: application/json" \
  -d '{
    "user_role": "developer",
    "user_role_display": "Software Developer",
    "llm_role": "code_reviewer", 
    "llm_role_display": "Senior Code Reviewer"
  }'
```

### 2. Create a Code Review Template

```bash
curl -X POST http://localhost:8080/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "persona_id": "YOUR_PERSONA_ID_HERE",
    "template": "Please review this {{language}} code for {{focus}}. Pay attention to {{concerns}}.",
    "variables": ["language", "focus", "concerns"]
  }'
```

### 3. Generate a Prompt

```bash
curl -X POST http://localhost:8080/v1/generate-prompt \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": "YOUR_TEMPLATE_ID_HERE",
    "values": {
      "language": "JavaScript",
      "focus": "performance optimization",
      "concerns": "memory leaks and async patterns"
    }
  }'
```

## Getting Started

1. Start the server: `./promptly serve`
2. Open the web interface: http://localhost:5173
3. Try creating personas, templates, and generating prompts
4. Use the existing educational examples as reference

## Common Use Cases

- **Code Review**: Developer seeking feedback on code quality
- **Learning**: Student getting explanations at appropriate level
- **Writing**: Author getting editing suggestions for different audiences
- **Interview Prep**: Candidate practicing with role-specific questions