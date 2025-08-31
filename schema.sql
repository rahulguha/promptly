-- SQLite schema for Promptly application
-- Creates tables for personas, prompt templates, and generated prompts

-- Personas table - stores user and LLM role definitions
CREATE TABLE IF NOT EXISTS personas (
    id TEXT PRIMARY KEY,
    user_role_display TEXT NOT NULL,
    llm_role_display TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Prompt templates table - stores reusable prompt templates with variables
CREATE TABLE IF NOT EXISTS prompt_templates (
    id TEXT NOT NULL,
    persona_id TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    variables TEXT NOT NULL, -- JSON array of variable names
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, version),

CREATE TABLE IF NOT EXISTS prompts (
    id TEXT PRIMARY KEY,
    template_id TEXT NOT NULL,
    template_version INTEGER NOT NULL,
