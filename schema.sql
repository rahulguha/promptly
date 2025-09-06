-- SQLite schema for Promptly application
-- Creates tables for profiles.

-- Profiles table - stores user-defined personas
CREATE TABLE IF NOT EXISTS profiles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    attributes TEXT, -- JSON blob for structured attributes
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);