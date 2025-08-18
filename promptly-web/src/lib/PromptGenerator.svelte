<script lang="ts">
	import { api, type PromptTemplate, type Prompt, type Persona } from './api.js';
	import { onMount } from 'svelte';
	import RichDropdown from './RichDropdown.svelte';

	let templates: PromptTemplate[] = [];
	let personas: Persona[] = [];
	let generatedPrompts: Prompt[] = [];
	let selectedTemplateId = '';
	let variables: Record<string, string> = {};
	let selectedTemplate: PromptTemplate | null = null;

	onMount(async () => {
		try {
			const [templatesData, personasData, promptsData] = await Promise.all([
				api.getTemplates(),
				api.getPersonas(),
				api.getPrompts()
			]);
			templates = templatesData;
			personas = personasData;
			generatedPrompts = promptsData;
			console.log('Loaded personas:', personas.length);
		} catch (error) {
			console.error('Failed to load data:', error);
		}
	});

	$: {
		if (selectedTemplateId) {
			selectedTemplate = templates.find(t => t.id === selectedTemplateId) || null;
			if (selectedTemplate) {
				// Initialize variables object
				variables = {};
				selectedTemplate.variables.forEach(v => {
					if (v !== 'user_role_display' && v !== 'llm_role_display') {
						variables[v] = '';
					}
				});
			}
		}
	}

	async function generatePrompt(e: Event) {
		e.preventDefault();
		if (!selectedTemplateId) return;
		
		const prompt = await api.generatePrompt(selectedTemplateId, variables);
		generatedPrompts = [prompt, ...generatedPrompts];
		variables = {};
		selectedTemplateId = '';
	}

	function getPersonaDisplay(personaId: string) {
		if (!personas || personas.length === 0) {
			return 'Loading...';
		}
		
		const persona = personas.find(p => p.persona_id === personaId);
		if (!persona) {
			console.warn('Persona not found for ID:', personaId, 'Available IDs:', personas.map(p => p.persona_id));
		}
		return persona ? `${persona.user_role_display} → ${persona.llm_role_display}` : 'Unknown';
	}

	function getTemplateDisplay(templateId: string) {
		const template = templates.find(t => t.id === templateId);
		return template ? getPersonaDisplay(template.persona_id) : 'Unknown';
	}
	
	$: templateOptions = templates.map(template => {
		const persona = personas.find(p => p.persona_id === template.persona_id);
		return {
			value: template.id,
			display: persona ? `${persona.user_role_display} → ${persona.llm_role_display}` : 'Unknown Persona',
			meta: template.template.slice(0, 60) + '...',
			user_role: persona?.user_role || 'unknown',
			llm_role: persona?.llm_role || 'unknown'
		};
	});
</script>

<div class="prompt-generator">
	<h2>Generate Prompts</h2>
	
	<div class="generator-form">
		<RichDropdown 
			items={templateOptions}
			bind:selectedValue={selectedTemplateId}
			placeholder="Select Template"
		/>

		{#if selectedTemplate}
			<div class="template-preview">
				<h4>Template Preview</h4>
				<pre>{selectedTemplate.template}</pre>
			</div>

			<div class="variables-form">
				<h4>Variables</h4>
				{#each selectedTemplate.variables as variable}
					{#if variable !== 'user_role_display' && variable !== 'llm_role_display'}
						<div class="variable-input">
							<label>{variable}:</label>
							<input bind:value={variables[variable]} placeholder={`Enter ${variable}`} />
						</div>
					{/if}
				{/each}
			</div>

			<button onclick={(e) => generatePrompt(e)}>Generate Prompt</button>
		{/if}
	</div>

	<div class="generated-prompts">
		<h3>Generated Prompts</h3>
		{#each generatedPrompts as prompt}
			<div class="prompt-card">
				<div class="prompt-meta">
					<small><strong>Template:</strong> {getTemplateDisplay(prompt.template_id)}</small>
					<small><strong>ID:</strong> {prompt.id}</small>
				</div>
				<div class="prompt-content">
					<strong>Final Content:</strong>
					<pre>{prompt.content}</pre>
				</div>
				{#if Object.keys(prompt.values).length > 0}
					<div class="prompt-values">
						<strong>Values Used:</strong>
						{#each Object.entries(prompt.values) as [key, value]}
							<span class="value-tag">{key}: {value}</span>
						{/each}
					</div>
				{/if}
			</div>
		{/each}
	</div>
</div>

<style>
	.prompt-generator {
		margin: 20px;
	}
	
	.generator-form {
		max-width: 600px;
		margin-bottom: 30px;
	}
	
	.generator-form :global(.rich-dropdown) {
		margin-bottom: 15px;
	}
	
	.template-preview {
		background: #f5f5f5;
		padding: 15px;
		border-radius: 4px;
		margin: 15px 0;
	}
	
	.template-preview pre {
		margin: 0;
		font-size: 12px;
		white-space: pre-wrap;
	}
	
	.variables-form {
		border: 1px solid #eee;
		padding: 15px;
		border-radius: 4px;
		margin: 15px 0;
	}
	
	.variable-input {
		display: flex;
		align-items: center;
		gap: 10px;
		margin-bottom: 10px;
	}
	
	.variable-input label {
		min-width: 100px;
		font-weight: bold;
	}
	
	.variable-input input {
		flex: 1;
		padding: 8px;
		border: 1px solid #ccc;
		border-radius: 4px;
	}
	
	.generated-prompts {
		margin-top: 30px;
	}
	
	.prompt-card {
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 20px;
		margin-bottom: 15px;
		background: white;
	}
	
	.prompt-meta {
		display: flex;
		justify-content: space-between;
		margin-bottom: 15px;
		color: #666;
	}
	
	.prompt-content pre {
		background: #f8f8f8;
		padding: 15px;
		border-radius: 4px;
		border-left: 4px solid #007cba;
		white-space: pre-wrap;
		font-size: 14px;
	}
	
	.prompt-values {
		margin-top: 15px;
	}
	
	.value-tag {
		display: inline-block;
		background: #e1f5fe;
		padding: 4px 8px;
		border-radius: 4px;
		margin: 2px;
		font-size: 12px;
	}
	
	button {
		padding: 10px 15px;
		background: #007cba;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
	}
	
	button:hover {
		background: #005a87;
	}
</style>