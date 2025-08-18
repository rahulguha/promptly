# API Documentation

Promptly provides a RESTful API for managing personas, prompt templates, and generated prompts.

## Base URL

```
http://localhost:8080/v1
```

## Authentication

No authentication required for local development.

## Endpoints

### Health Check

```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "service": "promptly"
}
```

### Personas

Personas define user and LLM roles for different contexts.

#### Get All Personas
```http
GET /v1/personas
```

**Response:**
```json
[
  {
    "persona_id": "323a1004-1526-4ee9-b9bc-3ba5cfbdc9b8",
    "user_role": "high_schooler",
    "user_role_display": "High School Student",
    "llm_role": "teacher",
    "llm_role_display": "Patient High School Teacher"
  }
]
```

#### Get Persona by ID
```http
GET /v1/personas/{id}
```

#### Create Persona
```http
POST /v1/personas
Content-Type: application/json

{
  "user_role": "developer",
  "user_role_display": "Software Developer", 
  "llm_role": "code_reviewer",
  "llm_role_display": "Senior Code Reviewer"
}
```

#### Update Persona
```http
PUT /v1/personas/{id}
Content-Type: application/json

{
  "user_role": "developer",
  "user_role_display": "Software Developer",
  "llm_role": "code_reviewer", 
  "llm_role_display": "Senior Code Reviewer"
}
```

#### Delete Persona
```http
DELETE /v1/personas/{id}
```

### Templates

Prompt templates with variable placeholders linked to personas.

#### Get All Templates
```http
GET /v1/templates
```

**Response:**
```json
[
  {
    "id": "7407b2d4-1448-40cb-a628-dc5775aa3268",
    "persona_id": "323a1004-1526-4ee9-b9bc-3ba5cfbdc9b8",
    "template": "Please review this {{code_type}} code for {{focus_area}}.",
    "variables": ["code_type", "focus_area"]
  }
]
```

#### Get Template by ID
```http
GET /v1/templates/{id}
```

#### Create Template
```http
POST /v1/templates
Content-Type: application/json

{
  "persona_id": "323a1004-1526-4ee9-b9bc-3ba5cfbdc9b8",
  "template": "Please review this {{code_type}} code for {{focus_area}}.",
  "variables": ["code_type", "focus_area"]
}
```

#### Update Template
```http
PUT /v1/templates/{id}
Content-Type: application/json

{
  "persona_id": "323a1004-1526-4ee9-b9bc-3ba5cfbdc9b8",
  "template": "Please review this {{code_type}} code for {{focus_area}}.",
  "variables": ["code_type", "focus_area"]
}
```

#### Delete Template
```http
DELETE /v1/templates/{id}
```

### Prompts

Generated prompts from templates with variable substitution.

#### Get All Prompts
```http
GET /v1/prompts
```

**Response:**
```json
[
  {
    "id": "9a0a0de1-af62-44af-a62c-ab0be14780ca",
    "template_id": "7407b2d4-1448-40cb-a628-dc5775aa3268",
    "values": {
      "code_type": "JavaScript",
      "focus_area": "performance"
    },
    "content": "User is a Software Developer and wants LLM to play the role of Senior Code Reviewer. Please review this JavaScript code for performance."
  }
]
```

#### Get Prompt by ID
```http
GET /v1/prompts/{id}
```

#### Create Prompt
```http
POST /v1/prompts
Content-Type: application/json

{
  "template_id": "7407b2d4-1448-40cb-a628-dc5775aa3268",
  "values": {
    "code_type": "JavaScript",
    "focus_area": "performance"
  },
  "content": "Generated prompt content..."
}
```

#### Update Prompt
```http
PUT /v1/prompts/{id}
Content-Type: application/json

{
  "template_id": "7407b2d4-1448-40cb-a628-dc5775aa3268",
  "values": {
    "code_type": "JavaScript", 
    "focus_area": "performance"
  },
  "content": "Updated prompt content..."
}
```

#### Delete Prompt
```http
DELETE /v1/prompts/{id}
```

### Generate Prompt from Template

Generate a final prompt by substituting template variables.

```http
POST /v1/generate-prompt
Content-Type: application/json

{
  "template_id": "7407b2d4-1448-40cb-a628-dc5775aa3268",
  "values": {
    "code_type": "JavaScript",
    "focus_area": "performance"
  }
}
```

**Response:**
```json
{
  "id": "9a0a0de1-af62-44af-a62c-ab0be14780ca",
  "template_id": "7407b2d4-1448-40cb-a628-dc5775aa3268", 
  "values": {
    "code_type": "JavaScript",
    "focus_area": "performance"
  },
  "content": "User is a Software Developer and wants LLM to play the role of Senior Code Reviewer. Please review this JavaScript code for performance."
}
```

## Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message describing what went wrong"
}
```

**Status Codes:**
- `200` - Success
- `201` - Created
- `400` - Bad Request (invalid input)
- `404` - Not Found
- `500` - Internal Server Error

## Variable Substitution

Templates use `{{variable_name}}` syntax. When generating prompts:
1. Persona context is automatically prepended
2. Variables are replaced with provided values
3. Missing variables remain as placeholders