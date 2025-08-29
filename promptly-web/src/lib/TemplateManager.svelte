<script lang="ts">
	import { api, type PromptTemplate, type Persona } from './api.js';
	import { onMount } from 'svelte';
	import RichDropdown from './RichDropdown.svelte';

	let templates: PromptTemplate[] = [];
	let personas: Persona[] = [];
	let showForm = false;
	let editingTemplate: PromptTemplate | null = null;
	let creatingVersion = false;
	let newTemplate = {
		persona_id: '',
		template: '',
		variables: ['']
	};

	onMount(async () => {
		try {
			templates = await api.getTemplates();
			personas = await api.getPersonas();
		} catch (error) {
			console.error('Failed to load data:', error);
		}
	});

	function addVariable() {
		newTemplate.variables = [...newTemplate.variables, ''];
	}

	function removeVariable(index: number) {
		newTemplate.variables = newTemplate.variables.filter((_, i) => i !== index);
	}

	async function createTemplate() {
		try {
			const variables = newTemplate.variables.filter(v => v.trim());
			const created = await api.createTemplate({
				...newTemplate,
				variables
			});
			if (created) {
				templates = [...templates, created];
				resetForm();
			} else {
				console.error('Failed to create template: No data returned');
			}
		} catch (error) {
			console.error('Error creating template:', error);
		}
	}

	async function updateTemplate() {
		if (!editingTemplate) return;
		const variables = newTemplate.variables.filter(v => v.trim());
		
		let updated;
		if (creatingVersion) {
			// Create new version
			updated = await api.createTemplateVersion(editingTemplate.id, {
				persona_id: newTemplate.persona_id,
				template: newTemplate.template,
				variables
			});
			templates = [...templates, updated];
		} else {
			// Update existing template
			updated = await api.updateTemplate(editingTemplate.id, {
				persona_id: newTemplate.persona_id,
				version: editingTemplate.version,
				template: newTemplate.template,
				variables
			});
			templates = templates.map(t => t.id === updated.id && t.version === updated.version ? updated : t);
		}
		
		resetForm();
	}

	async function deleteTemplate(template: PromptTemplate) {
		const personaDisplay = getPersonaDisplay(template.persona_id);
		if (confirm(`Delete template v${template.version} for ${personaDisplay}?`)) {
			await api.deleteTemplate(template.id, template.version);
			templates = templates.filter(t => !(t.id === template.id && t.version === template.version));
		}
	}

	function editTemplate(template: PromptTemplate) {
		editingTemplate = template;
		creatingVersion = false;
		newTemplate = {
			persona_id: template.persona_id,
			template: template.template,
			variables: [...template.variables]
		};
		showForm = true;
	}

	function createNewVersion(template: PromptTemplate) {
		editingTemplate = template;
		creatingVersion = true;
		newTemplate = {
			persona_id: template.persona_id,
			template: template.template,
			variables: [...template.variables]
		};
		showForm = true;
	}

	function resetForm() {
		newTemplate = { persona_id: '', template: '', variables: [''] };
		editingTemplate = null;
		creatingVersion = false;
		showForm = false;
	}

	function getPersonaDisplay(personaId: string) {
		const persona = personas.find(p => p.persona_id === personaId);
		return persona ? `${persona.user_role_display} → ${persona.llm_role_display}` : 'Unknown';
	}
	
	$: personaOptions = personas ? personas.map(persona => ({
		value: persona.persona_id,
		display: `${persona.user_role_display} → ${persona.llm_role_display}`,
		meta: `${persona.user_role} → ${persona.llm_role}`,
		user_role: persona.user_role,
		llm_role: persona.llm_role
	})) : [];
</script>

<div class="template-manager">
	<h2>Templates</h2>
	
	<button onclick={() => showForm = !showForm}>
		{showForm ? 'Cancel' : 'Add Template'}
	</button>

	{#if showForm}
		<div class="form-header">
			<h3>
				{#if creatingVersion}
					Creating New Version of Template (v{editingTemplate?.version}) 
				{:else if editingTemplate}
					Editing Template (v{editingTemplate?.version})
				{:else}
					Create New Template
				{/if}
			</h3>
		</div>
		<form onsubmit={(e) => { e.preventDefault(); editingTemplate ? updateTemplate() : createTemplate(); }} class="template-form">
			<RichDropdown 
				items={personaOptions}
				bind:selectedValue={newTemplate.persona_id}
				placeholder="Select Persona"
			/>

			<textarea 
				bind:value={newTemplate.template} 
				placeholder="Template content (e.g., Review this {'{'}language{'}'} code: {'{'}code{'}'})"
				rows="4"
				required
			></textarea>

			<div class="variables-section">
				<h4>Variables</h4>
				{#each newTemplate.variables as variable, i}
					<div class="variable-input">
						<input bind:value={newTemplate.variables[i]} placeholder="Variable name" />
						<button type="button" onclick={() => removeVariable(i)}>Remove</button>
					</div>
				{/each}
				<button type="button" onclick={addVariable}>Add Variable</button>
			</div>

			<button type="submit">
				{#if creatingVersion}
					Create New Version
				{:else if editingTemplate}
					Update Template
				{:else}
					Create Template
				{/if}
			</button>
			{#if editingTemplate}
				<button type="button" onclick={resetForm}>Cancel Edit</button>
			{/if}
		</form>
	{/if}

	<div class="templates-list">
		{#each templates as template}
			<div class="template-item">
				<div class="template-actions">
					<button class="icon-btn edit-btn" onclick={() => editTemplate(template)} title="Edit">
						✏️
					</button>
					<button class="icon-btn version-btn" onclick={() => createNewVersion(template)} title="Create new Version">
						📝
					</button>
					<button class="icon-btn delete-btn" onclick={() => deleteTemplate(template)} title="Delete">
						🗑️
					</button>
				</div>
				<div class="template-display">
					<div class="template-version">v{template.version}</div>
					<!-- <div class="template-persona">{getPersonaDisplay(template.persona_id)}</div> -->
					<div class="template-preview">{template.template.slice(0, 80)}...</div>
				</div>
			</div>
		{/each}
	</div>
</div>

<style>
	.template-manager {
		margin: 20px;
	}

	.form-header {
		margin: 20px 0 10px 0;
	}

	.form-header h3 {
		color: #007cba;
		margin: 0;
		font-size: 18px;
	}
	
	.template-form {
		display: flex;
		flex-direction: column;
		gap: 15px;
		max-width: 600px;
		margin: 20px 0;
	}
	
	.template-form select,
	.template-form textarea {
		padding: 8px;
		border: 1px solid #ccc;
		border-radius: 4px;
	}
	
	.variables-section {
		border: 1px solid #eee;
		padding: 15px;
		border-radius: 4px;
	}
	
	.variable-input {
		display: flex;
		gap: 10px;
		margin-bottom: 10px;
	}
	
	.variable-input input {
		flex: 1;
		padding: 8px;
		border: 1px solid #ccc;
		border-radius: 4px;
	}
	
	.templates-list {
		margin-top: 20px;
	}
	
	.template-item {
		display: flex;
		align-items: center;
		gap: 15px;
		padding: 12px 15px;
		border-bottom: 1px solid #eee;
		background: white;
	}
	
	.template-item:hover {
		background: #f8f8f8;
	}
	
	.template-actions {
		display: flex;
		gap: 8px;
	}
	
	.template-display {
		flex: 1;
	}
	
	.template-version {
		font-size: 12px;
		font-weight: bold;
		color: #28a745;
		margin-bottom: 2px;
	}

	.template-persona {
		font-size: 14px;
		font-weight: bold;
		color: #007cba;
		margin-bottom: 4px;
	}
	
	.template-preview {
		font-size: 12px;
		color: #666;
		font-style: italic;
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

	.version-btn:hover {
		background: #d1ecf1;
		border-color: #bee5eb;
	}
	
	.delete-btn:hover {
		background: #f8d7da;
		border-color: #f5c6cb;
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
	
	button[type="button"] {
		background: #666;
	}
	
	button[type="button"]:hover {
		background: #444;
	}
</style>