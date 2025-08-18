<script lang="ts">
	import PersonaManager from '$lib/PersonaManager.svelte';
	import TemplateManager from '$lib/TemplateManager.svelte';
	import PromptGenerator from '$lib/PromptGenerator.svelte';

	let activeTab = 'personas';

	function setActiveTab(tab: string) {
		activeTab = tab;
	}
</script>

<div class="prompt-management">
	<div class="tab-navigation">
		<button 
			class="tab-button {activeTab === 'personas' ? 'active' : ''}"
			onclick={() => setActiveTab('personas')}
		>
			Personas
		</button>
		<button 
			class="tab-button {activeTab === 'templates' ? 'active' : ''}"
			onclick={() => setActiveTab('templates')}
		>
			Templates
		</button>
		<button 
			class="tab-button {activeTab === 'generate' ? 'active' : ''}"
			onclick={() => setActiveTab('generate')}
		>
			Generate Prompts
		</button>
	</div>

	<div class="tab-content">
		{#if activeTab === 'personas'}
			<PersonaManager />
		{:else if activeTab === 'templates'}
			<TemplateManager />
		{:else if activeTab === 'generate'}
			<PromptGenerator />
		{/if}
	</div>
</div>

<style>
	.prompt-management {
		background: white;
		border-radius: 8px;
		padding: 20px;
		box-shadow: 0 2px 4px rgba(0,0,0,0.1);
	}
	
	.tab-content {
		margin-top: 20px;
	}
</style>