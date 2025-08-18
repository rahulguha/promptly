<script lang="ts">
	import { api, type Persona } from './api.js';
	import { onMount } from 'svelte';
	import RichDropdown from './RichDropdown.svelte';

	let personas: Persona[] = [];
	let showForm = false;
	let editingPersona: Persona | null = null;
	let newPersona = {
		user_role: '',
		user_role_display: '',
		llm_role: '',
		llm_role_display: ''
	};
	
	let showUserRoleFields = false;
	let showLLMRoleFields = false;
	
	$: userRoleOptions = uniqueUserRoles.map(role => {
		const persona = personas.find(p => p.user_role === role);
		return {
			value: role,
			display: persona?.user_role_display || role,
			meta: `Role: ${role}`,
			user_role: role,
			llm_role: 'user'
		};
	});
	
	$: llmRoleOptions = uniqueLLMRoles.map(role => {
		const persona = personas.find(p => p.llm_role === role);
		return {
			value: role,
			display: persona?.llm_role_display || role,
			meta: `Role: ${role}`,
			user_role: 'assistant', 
			llm_role: role
		};
	});
	
	$: uniqueUserRoles = [...new Set(personas.map(p => p.user_role))];
	$: uniqueLLMRoles = [...new Set(personas.map(p => p.llm_role))];

	onMount(async () => {
		personas = await api.getPersonas();
	});

	async function createPersona() {
		const created = await api.createPersona(newPersona);
		personas = [...personas, created];
		resetForm();
	}

	async function updatePersona() {
		if (!editingPersona) return;
		const updated = await api.updatePersona(editingPersona.persona_id, newPersona);
		personas = personas.map(p => p.persona_id === updated.persona_id ? updated : p);
		resetForm();
	}

	async function deletePersona(persona: Persona) {
		if (confirm(`Delete persona: ${persona.user_role_display} → ${persona.llm_role_display}?`)) {
			await api.deletePersona(persona.persona_id);
			personas = personas.filter(p => p.persona_id !== persona.persona_id);
		}
	}

	function editPersona(persona: Persona) {
		editingPersona = persona;
		newPersona = {
			user_role: persona.user_role,
			user_role_display: persona.user_role_display,
			llm_role: persona.llm_role,
			llm_role_display: persona.llm_role_display
		};
		showForm = true;
	}

	function resetForm() {
		newPersona = { user_role: '', user_role_display: '', llm_role: '', llm_role_display: '' };
		editingPersona = null;
		showForm = false;
		showUserRoleFields = false;
		showLLMRoleFields = false;
	}
	
	function handleUserRoleChange(value: string) {
		if (value === 'ADD_NEW') {
			showUserRoleFields = true;
			newPersona.user_role = '';
			newPersona.user_role_display = '';
		} else {
			showUserRoleFields = false;
			newPersona.user_role = value;
			// Auto-fill display name from existing persona
			const existing = personas.find(p => p.user_role === value);
			if (existing) {
				newPersona.user_role_display = existing.user_role_display;
			}
		}
	}
	
	function handleLLMRoleChange(value: string) {
		if (value === 'ADD_NEW') {
			showLLMRoleFields = true;
			newPersona.llm_role = '';
			newPersona.llm_role_display = '';
		} else {
			showLLMRoleFields = false;
			newPersona.llm_role = value;
			// Auto-fill display name from existing persona
			const existing = personas.find(p => p.llm_role === value);
			if (existing) {
				newPersona.llm_role_display = existing.llm_role_display;
			}
		}
	}
</script>

<div class="persona-manager">
	<h2>Personas</h2>
	
	<button onclick={() => showForm = !showForm}>
		{showForm ? 'Cancel' : 'Add Persona'}
	</button>

	{#if showForm}
		<form onsubmit={(e) => { e.preventDefault(); editingPersona ? updatePersona() : createPersona(); }} class="persona-form">
			<div class="role-section">
				<label>User Role:</label>
				<RichDropdown 
					items={userRoleOptions}
					bind:selectedValue={newPersona.user_role}
					placeholder="Select User Role"
					allowAddNew={true}
					on:select={(e) => handleUserRoleChange(e.detail.value)}
					on:addNew={() => handleUserRoleChange('ADD_NEW')}
				/>
				
				{#if showUserRoleFields}
					<input bind:value={newPersona.user_role} placeholder="User Role (e.g., developer)" required />
					<input bind:value={newPersona.user_role_display} placeholder="User Role Display (e.g., Software Developer)" required />
				{:else if newPersona.user_role}
					<input bind:value={newPersona.user_role_display} placeholder="User Role Display" required />
				{/if}
			</div>

			<div class="role-section">
				<label>LLM Role:</label>
				<RichDropdown 
					items={llmRoleOptions}
					bind:selectedValue={newPersona.llm_role}
					placeholder="Select LLM Role"
					allowAddNew={true}
					on:select={(e) => handleLLMRoleChange(e.detail.value)}
					on:addNew={() => handleLLMRoleChange('ADD_NEW')}
				/>
				
				{#if showLLMRoleFields}
					<input bind:value={newPersona.llm_role} placeholder="LLM Role (e.g., code_reviewer)" required />
					<input bind:value={newPersona.llm_role_display} placeholder="LLM Role Display (e.g., Senior Code Reviewer)" required />
				{:else if newPersona.llm_role}
					<input bind:value={newPersona.llm_role_display} placeholder="LLM Role Display" required />
				{/if}
			</div>

			<button type="submit">{editingPersona ? 'Update' : 'Create'} Persona</button>
			{#if editingPersona}
				<button type="button" onclick={resetForm}>Cancel Edit</button>
			{/if}
		</form>
	{/if}

	<div class="personas-list">
		{#each personas as persona}
			<div class="persona-item">
				<div class="persona-actions">
					<button class="icon-btn edit-btn" onclick={() => editPersona(persona)} title="Edit">
						✏️
					</button>
					<button class="icon-btn delete-btn" onclick={() => deletePersona(persona)} title="Delete">
						🗑️
					</button>
				</div>
				<div class="persona-display">
					{persona.user_role_display} → {persona.llm_role_display}
				</div>
			</div>
		{/each}
	</div>
</div>

<style>
	.persona-manager {
		margin: 20px;
	}
	
	.persona-form {
		display: flex;
		flex-direction: column;
		gap: 10px;
		max-width: 400px;
		margin: 20px 0;
	}
	
	.persona-form input,
	.persona-form select {
		padding: 8px;
		border: 1px solid #ccc;
		border-radius: 4px;
	}
	
	.role-section {
		border: 1px solid #eee;
		padding: 15px;
		border-radius: 4px;
		background: #fafafa;
	}
	
	.role-section label {
		display: block;
		font-weight: bold;
		margin-bottom: 8px;
	}
	
	.role-section select {
		width: 100%;
		margin-bottom: 10px;
	}
	
	.role-section input {
		width: 100%;
		margin-bottom: 8px;
	}
	
	.personas-list {
		margin-top: 20px;
	}
	
	.persona-item {
		display: flex;
		align-items: center;
		gap: 15px;
		padding: 8px 15px;
		border-bottom: 1px solid #eee;
		background: white;
	}
	
	.persona-item:hover {
		background: #f8f8f8;
	}
	
	.persona-actions {
		display: flex;
		gap: 8px;
	}
	
	.persona-display {
		font-size: 14px;
		color: #333;
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