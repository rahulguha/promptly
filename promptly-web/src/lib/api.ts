// In development, Vite's proxy will forward requests from /v1 to the API server.
// In production, we assume the API is being served on the same host.
const API_BASE = "/v1";

export interface Persona {
  persona_id: string;
  user_role_display: string;
  llm_role_display: string;
}

export interface PromptTemplate {
  id: string;
  name: string;
  persona_id: string;
  version: number;
  meta_role: string;
  task: string;
  answer_guideline: string;
  template: string;
  variables: string[];
}

export interface Prompt {
  id: string;
  name: string;
  template_id: string;
  template_version: number;
  variable_values: Record<string, string>;
  content: string;
}

// Define a type for our API request options that allows a structured body.
type ApiRequestOptions = Omit<RequestInit, "body"> & {
  body?: unknown;
};

async function apiRequest<T>(
  endpoint: string,
  options: ApiRequestOptions = {}
): Promise<T> {
  const headers = new Headers(options.headers);
  headers.set("Accept", "application/json");

  const fetchOptions: RequestInit = {
    method: options.method,
    headers: options.headers,
    mode: options.mode,
    credentials: options.credentials,
    cache: options.cache,
    redirect: options.redirect,
    referrerPolicy: options.referrerPolicy,
    signal: options.signal,
  };

  if (options.body) {
    fetchOptions.body = JSON.stringify(options.body);
    headers.set("Content-Type", "application/json");
  }

  fetchOptions.headers = headers;

  const res = await fetch(`${API_BASE}${endpoint}`, fetchOptions);

  if (!res.ok) {
    let message = `HTTP error! status: ${res.status}`;
    try {
      const body = await res.json();
      if (body.error) {
        message = body.error;
      }
    } catch (e) {
      // Not a JSON error response, or no body
    }
    throw new Error(message);
  }

  // For DELETE requests, we don't expect a JSON body in the success response.
  if (options.method?.toUpperCase() === "DELETE" || res.status === 204) {
    return undefined as T;
  }

  return res.json();
}

export const api = {
  // Personas
  getPersonas(): Promise<Persona[]> {
    return apiRequest("/personas");
  },

  createPersona(persona: Omit<Persona, "persona_id">): Promise<Persona> {
    return apiRequest("/personas", { method: "POST", body: persona });
  },

  updatePersona(
    id: string,
    persona: Omit<Persona, "persona_id">
  ): Promise<Persona> {
    return apiRequest(`/personas/${id}`, { method: "PUT", body: persona });
  },

  deletePersona(id: string): Promise<void> {
    return apiRequest(`/personas/${id}`, { method: "DELETE" });
  },

  // Templates
  getTemplates(): Promise<PromptTemplate[]> {
    return apiRequest("/templates");
  },

  createTemplate(
    template: Omit<PromptTemplate, "id">
  ): Promise<PromptTemplate> {
    return apiRequest("/templates", { method: "POST", body: template });
  },

  updateTemplate(
    id: string,
    template: Omit<PromptTemplate, "id">
  ): Promise<PromptTemplate> {
    return apiRequest(`/templates/${id}`, { method: "PUT", body: template });
  },

  createTemplateVersion(
    id: string,
    template: Omit<PromptTemplate, "id" | "version">
  ): Promise<PromptTemplate> {
    return apiRequest(`/templates/${id}/version`, {
      method: "POST",
      body: template,
    });
  },

  deleteTemplate(id: string, version: number): Promise<void> {
    return apiRequest(`/templates/${id}?version=${version}`, {
      method: "DELETE",
    });
  },

  // Prompts
  generatePrompt(
    templateId: string,
    name: string,
    values: Record<string, string>
  ): Promise<Prompt> {
    return apiRequest("/generate-prompt", {
      method: "POST",
      body: {
        template_id: templateId,
        name: name,
        variable_values: values,
      },
    });
  },

  getPrompts(): Promise<Prompt[]> {
    return apiRequest("/prompts");
  },

  updatePrompt(id: string, prompt: Omit<Prompt, "id">): Promise<Prompt> {
    return apiRequest(`/prompts/${id}`, { method: "PUT", body: prompt });
  },

  deletePrompt(id: string): Promise<void> {
    return apiRequest(`/prompts/${id}`, { method: "DELETE" });
  },
};