const API_BASE = "/v1";

export interface Persona {
  persona_id: string;
  user_role_display: string;
  llm_role_display: string;
}

export interface PromptTemplate {
  id: string;
  persona_id: string;
  version: number;
  template: string;
  variables: string[];
}

export interface Prompt {
  id: string;
  template_id: string;
  template_version: number;
  variable_values: Record<string, string>;
  content: string;
}

export const api = {
  // Personas
  async getPersonas(): Promise<Persona[]> {
    const res = await fetch(`${API_BASE}/personas`);
    return res.json();
  },

  async createPersona(persona: Omit<Persona, "persona_id">): Promise<Persona> {
    const res = await fetch(`${API_BASE}/personas`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(persona),
    });
    return res.json();
  },

  async updatePersona(id: string, persona: Omit<Persona, "persona_id">): Promise<Persona> {
    const res = await fetch(`${API_BASE}/personas/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(persona),
    });
    return res.json();
  },

  async deletePersona(id: string): Promise<void> {
    await fetch(`${API_BASE}/personas/${id}`, {
      method: "DELETE",
    });
  },

  // Templates
  async getTemplates(): Promise<PromptTemplate[]> {
    const res = await fetch(`${API_BASE}/templates`);
    return res.json();
  },

  async createTemplate(
    template: Omit<PromptTemplate, "id">
  ): Promise<PromptTemplate> {
    const res = await fetch(`${API_BASE}/templates`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(template),
    });
    return res.json();
  },

  async updateTemplate(id: string, template: Omit<PromptTemplate, "id">): Promise<PromptTemplate> {
    const res = await fetch(`${API_BASE}/templates/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(template),
    });
    return res.json();
  },

  async createTemplateVersion(id: string, template: Omit<PromptTemplate, "id" | "version">): Promise<PromptTemplate> {
    const res = await fetch(`${API_BASE}/templates/${id}/version`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(template),
    });
    return res.json();
  },

  async deleteTemplate(id: string, version: number): Promise<void> {
    await fetch(`${API_BASE}/templates/${id}?version=${version}`, {
      method: "DELETE",
    });
  },

  // Prompts
  async generatePrompt(
    templateId: string,
    values: Record<string, string>
  ): Promise<Prompt> {
    const res = await fetch(`${API_BASE}/generate-prompt`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        template_id: templateId,
        variable_values: values,
      }),
    });
    return res.json();
  },

  async getPrompts(): Promise<Prompt[]> {
    const res = await fetch(`${API_BASE}/prompts`);
    return res.json();
  },

  async updatePrompt(id: string, prompt: Omit<Prompt, "id">): Promise<Prompt> {
    const res = await fetch(`${API_BASE}/prompts/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(prompt),
    });
    return res.json();
  },

  async deletePrompt(id: string): Promise<void> {
    await fetch(`${API_BASE}/prompts/${id}`, {
      method: "DELETE",
    });
  },
};
