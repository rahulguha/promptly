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
    id TEXT PRIMARY KEY,
    persona_id TEXT NOT NULL,
    template TEXT NOT NULL,
    variables TEXT NOT NULL, -- JSON array of variable names
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (persona_id) REFERENCES personas(id) ON DELETE CASCADE
);

-- Prompts table - stores generated prompts from templates
CREATE TABLE IF NOT EXISTS prompts (
    id TEXT PRIMARY KEY,
    template_id TEXT NOT NULL,
    variable_values TEXT NOT NULL, -- JSON object with variable values
    content TEXT NOT NULL, -- Final generated prompt content
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (template_id) REFERENCES prompt_templates(id) ON DELETE CASCADE
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_personas_user_role ON personas(user_role_display);
CREATE INDEX IF NOT EXISTS idx_personas_llm_role ON personas(llm_role_display);
CREATE INDEX IF NOT EXISTS idx_templates_persona ON prompt_templates(persona_id);
CREATE INDEX IF NOT EXISTS idx_prompts_template ON prompts(template_id);

-- Sample data migration from existing JSON files
INSERT OR IGNORE INTO personas (id, user_role_display, llm_role_display) VALUES 
('323a1004-1526-4ee9-b9bc-3ba5cfbdc9b8', 'High School Student at 9th grade', 'Patient High School Teacher'),
('e44f89bf-dbd9-4609-b5d3-56ec4fe67a0c', 'High School Student', 'Mock Tester who generates Questions'),
('f27738a8-c2d9-4342-92ff-4514933107dc', 'High School Student', 'Hard High School Tester'),
('235658ed-f107-4e2a-836b-88ac3d7f75c2', 'College Student', 'College Professor'),
('2d989125-7e4b-4f09-b833-836cca653103', 'College Student', 'College Buddy');

INSERT OR IGNORE INTO prompt_templates (id, persona_id, template, variables) VALUES
('7407b2d4-1448-40cb-a628-dc5775aa3268', '323a1004-1526-4ee9-b9bc-3ba5cfbdc9b8', 
 'User is a High School Student and wants LLM to play the role of Patient High School Teacher. Please teach for {{grade}} grade.', 
 '["grade"]'),
('a56dc324-21d6-413a-a02d-2f72844fe169', 'e44f89bf-dbd9-4609-b5d3-56ec4fe67a0c', 
 'User is a High School Student and wants LLM to play the role of Mock Tester who generates Questions. Generate {{numQ}} questions with answers. Make them multiple choice. Around 25% should be easy, 40% medium hard and remaining hard questions. The student is in {{grade}} grade.', 
 '["numQ", "grade"]'),
('0112d01b-f7c9-425d-877f-ee9a724378e7', '235658ed-f107-4e2a-836b-88ac3d7f75c2', 
 'User is a College Student and wants LLM to play the role of College Professor. Please teach for {{level}} level. Also discuss details and be concise.', 
 '["level"]');