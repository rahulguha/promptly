<script lang="ts">
	import { api, type PromptTemplate, type Prompt, type Persona } from './api.js';
	import { onMount } from 'svelte';
	import RichDropdown from './RichDropdown.svelte';

	let templates: PromptTemplate[] = [];
	let personas: Persona[] = [];
	let generatedPrompts: Prompt[] = [];
	let selectedTemplateId = '';
	let selectedTemplateVersion = 1;
	let variables: Record<string, string> = {};
	let selectedTemplate: PromptTemplate | null = null;
	let editingPrompt: Prompt | null = null;
	let showEditForm = false;
	let filterTemplateId = '';

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
			console.log('Loaded personas:', personas?.length || 0);
		} catch (error) {
			console.error('Failed to load data:', error);
		}
	});

	$: {
		if (selectedTemplateId) {
			// Parse the compound value (id:version)
			const [templateId, versionStr] = selectedTemplateId.split(':');
			const version = parseInt(versionStr);
			selectedTemplateVersion = version;
			
			selectedTemplate = templates.find(t => t.id === templateId && t.version === version) || null;
			if (selectedTemplate) {
				// Initialize variables object
				variables = {};
				selectedTemplate.variables.forEach(v => {
					if (v !== 'user_role_display' && v !== 'llm_role_display') {
						variables[v] = '';
					}
				});
			}
		} else {
			selectedTemplate = null;
			selectedTemplateVersion = 1;
			variables = {};
		}
	}

	async function generatePrompt(e: Event) {
		e.preventDefault();
		if (!selectedTemplateId) return;
		
		try {
			if (editingPrompt) {
				await updatePrompt();
			} else {
				// Parse template ID from the compound value
				const [templateId] = selectedTemplateId.split(':');
				const prompt = await api.generatePrompt(templateId, variables);
				if (prompt) {
					generatedPrompts = [prompt, ...generatedPrompts];
					variables = {};
					selectedTemplateId = '';
				} else {
					console.error('Failed to generate prompt: No data returned');
				}
			}
		} catch (error) {
			console.error('Error generating prompt:', error);
		}
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

	function getTemplateDisplay(templateId: string, templateVersion?: number) {
		const template = templates.find(t => t.id === templateId && (templateVersion ? t.version === templateVersion : true));
		if (!template) return 'Unknown';
		const personaDisplay = getPersonaDisplay(template.persona_id);
		return `${personaDisplay} (v${template.version})`;
	}

	function showTemplateDetails(templateId: string, templateVersion: number) {
		const template = templates.find(t => t.id === templateId && t.version === templateVersion);
		if (!template) {
			alert('Template not found');
			return;
		}
		
		const persona = personas.find(p => p.persona_id === template.persona_id);
		const personaDisplay = persona ? `${persona.user_role_display} → ${persona.llm_role_display}` : 'Unknown Persona';
		
		const details = `Template Details:
		
Persona: ${personaDisplay}
Version: ${template.version}
Variables: ${template.variables.join(', ')}

Template Content:
${template.template}`;
		
		alert(details);
	}
	
	$: templateOptions = templates ? templates.map(template => {
		const persona = personas ? personas.find(p => p.persona_id === template.persona_id) : null;
		return {
			value: `${template.id}:${template.version}`,
			display: persona ? `${persona.user_role_display} → ${persona.llm_role_display} (v${template.version})` : 'Unknown Persona',
			meta: template.template.slice(0, 60) + '...',
			user_role: persona?.user_role || 'unknown',
			llm_role: persona?.llm_role || 'unknown'
		};
	}) : [];

	$: filterOptions = [
		{ value: '', display: 'All Templates', meta: 'Show all prompts', user_role: 'all', llm_role: 'all' },
		...templateOptions
	];

	$: filteredPrompts = filterTemplateId 
		? generatedPrompts.filter(p => {
			// Handle both compound format (id:version) and plain id for backward compatibility
			if (filterTemplateId.includes(':')) {
				const [templateId, versionStr] = filterTemplateId.split(':');
				const version = parseInt(versionStr);
				return p.template_id === templateId && p.template_version === version;
			} else {
				return p.template_id === filterTemplateId;
			}
		})
		: generatedPrompts;

	async function deletePrompt(prompt: Prompt) {
		if (confirm('Delete this prompt?')) {
			await api.deletePrompt(prompt.id);
			generatedPrompts = generatedPrompts.filter(p => p.id !== prompt.id);
		}
	}

	function editPrompt(prompt: Prompt) {
		editingPrompt = prompt;
		selectedTemplateId = `${prompt.template_id}:${prompt.template_version}`;
		variables = { ...prompt.variable_values || {} };
		showEditForm = true;
	}

	async function updatePrompt() {
		if (!editingPrompt) return;
		
		const template = templates.find(t => t.id === editingPrompt.template_id && t.version === editingPrompt.template_version);
		if (!template) return;

		// Generate new content with updated variables
		let content = template.template;
		for (const [variable, value] of Object.entries(variables)) {
			const placeholder = "{{" + variable + "}}";
			content = content.replace(new RegExp(placeholder.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'g'), value);
		}

		const updatedPrompt = await api.updatePrompt(editingPrompt.id, {
			template_id: editingPrompt.template_id,
			variable_values: variables,
			content: content
		});

		generatedPrompts = generatedPrompts.map(p => p.id === updatedPrompt.id ? updatedPrompt : p);
		resetEditForm();
	}

	function resetEditForm() {
		editingPrompt = null;
		showEditForm = false;
		variables = {};
		selectedTemplateId = '';
	}
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

			<button onclick={(e) => generatePrompt(e)}>
				{editingPrompt ? 'Update Prompt' : 'Generate Prompt'}
			</button>
			{#if editingPrompt}
				<button type="button" onclick={resetEditForm}>Cancel Edit</button>
			{/if}
		{/if}
	</div>

	<div class="generated-prompts">
		<div class="prompts-header">
			<h3>Generated Prompts</h3>
			<div class="filter-section">
				<RichDropdown 
					items={filterOptions}
					bind:selectedValue={filterTemplateId}
					placeholder="Filter by template"
				/>
			</div>
		</div>
		{#each filteredPrompts as prompt}
			<div class="prompt-card">
				<div class="prompt-header">
					<div class="prompt-meta">
						<small><strong>Template:</strong> 
							<button 
								class="template-link" 
								onclick={() => showTemplateDetails(prompt.template_id, prompt.template_version)}
								title="Click to view template details"
							>
								{getTemplateDisplay(prompt.template_id, prompt.template_version)}
							</button>
						</small>
						<small><strong>ID:</strong> {prompt.id}</small>
					</div>
					<div class="prompt-actions">
						<button class="icon-btn edit-btn" onclick={() => editPrompt(prompt)} title="Edit">
							✏️
						</button>
						<button class="icon-btn delete-btn" onclick={() => deletePrompt(prompt)} title="Delete">
							🗑️
						</button>
					</div>
				</div>
				<div class="prompt-content">
					<strong>Final Content:</strong>
					<pre>{prompt.content}</pre>
				</div>
				{#if prompt.variable_values && Object.keys(prompt.variable_values).length > 0}
					<div class="prompt-values">
						<strong>Values Used:</strong>
						{#each Object.entries(prompt.variable_values) as [key, value]}
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
	
	.prompts-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 20px;
	}
	
	.filter-section {
		min-width: 250px;
	}
	
	.prompt-card {
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 20px;
		margin-bottom: 15px;
		background: white;
	}
	
	.prompt-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 15px;
	}
	
	.prompt-meta {
		display: flex;
		flex-direction: column;
		gap: 4px;
		color: #666;
	}
	
	.prompt-actions {
		display: flex;
		gap: 8px;
	}
	
	.icon-btn {
		padding: 6px;
		font-size: 16px;
		border-radius: 4px;
		background: #f8f9fa;
		border: 1px solid #dee2e6;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
	}
	
	.icon-btn:hover {
		background: #e9ecef;
	}
	
	.edit-btn:hover {
		background: #d4edda;
		border-color: #c3e6cb;
	}
	
	.delete-btn:hover {
		background: #f8d7da;
		border-color: #f5c6cb;
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

	.template-link {
		background: none;
		border: none;
		color: #007cba;
		text-decoration: underline;
		cursor: pointer;
		padding: 0;
		font-size: inherit;
		font-family: inherit;
	}

	.template-link:hover {
		background: none;
		color: #005a87;
		text-decoration: none;
	}
</style>