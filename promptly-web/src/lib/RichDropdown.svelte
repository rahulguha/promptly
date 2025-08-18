<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	
	export let items: any[] = [];
	export let selectedValue: string = '';
	export let placeholder: string = 'Select...';
	export let displayKey: string = 'display';
	export let valueKey: string = 'value';
	export let iconKey: string = 'icon';
	export let metaKey: string = 'meta';
	export let allowAddNew: boolean = false;
	
	let isOpen = false;
	let dropdownElement: HTMLDivElement;
	
	const dispatch = createEventDispatcher();
	
	$: selectedItem = items.find(item => item[valueKey] === selectedValue);
	
	function selectItem(item: any) {
		selectedValue = item[valueKey];
		isOpen = false;
		dispatch('select', item);
	}
	
	function addNew() {
		isOpen = false;
		dispatch('addNew');
	}
	
	function toggleDropdown() {
		isOpen = !isOpen;
	}
	
	function handleClickOutside(event: MouseEvent) {
		if (dropdownElement && !dropdownElement.contains(event.target as Node)) {
			isOpen = false;
		}
	}
	
	function getRoleIcon(userRole: string, llmRole: string) {
		const roleIcons: Record<string, string> = {
			'developer': '💻',
			'writer': '✍️',
			'analyst': '📊',
			'student': '🎓',
			'teacher': '👨‍🏫',
			'code_reviewer': '🔍',
			'assistant': '🤖',
			'editor': '✏️',
			'tutor': '📚',
			'mentor': '🧭'
		};
		return `${roleIcons[userRole] || '👤'} → ${roleIcons[llmRole] || '🤖'}`;
	}
</script>

<svelte:window on:click={handleClickOutside} />

<div class="rich-dropdown" bind:this={dropdownElement}>
	<button class="dropdown-trigger" onclick={toggleDropdown} type="button">
		<div class="selected-content">
			{#if selectedItem}
				<span class="icon">{selectedItem[iconKey] || getRoleIcon(selectedItem.user_role, selectedItem.llm_role)}</span>
				<div class="text">
					<div class="main">{selectedItem[displayKey]}</div>
					{#if selectedItem[metaKey]}
						<div class="meta">{selectedItem[metaKey]}</div>
					{/if}
				</div>
			{:else}
				<span class="placeholder">{placeholder}</span>
			{/if}
		</div>
		<span class="arrow" class:open={isOpen}>▼</span>
	</button>
	
	{#if isOpen}
		<div class="dropdown-menu">
			{#each items as item}
				<button 
					class="dropdown-item" 
					onclick={() => selectItem(item)}
					type="button"
				>
					<span class="icon">{item[iconKey] || getRoleIcon(item.user_role, item.llm_role)}</span>
					<div class="text">
						<div class="main">{item[displayKey]}</div>
						{#if item[metaKey]}
							<div class="meta">{item[metaKey]}</div>
						{/if}
					</div>
				</button>
			{/each}
			
			{#if allowAddNew}
				<div class="separator"></div>
				<button class="dropdown-item add-new" onclick={addNew} type="button">
					<span class="icon">➕</span>
					<div class="text">
						<div class="main">Add New</div>
					</div>
				</button>
			{/if}
		</div>
	{/if}
</div>

<style>
	.rich-dropdown {
		position: relative;
		width: 100%;
	}
	
	.dropdown-trigger {
		width: 100%;
		padding: 12px;
		border: 1px solid #ccc;
		border-radius: 6px;
		background: white;
		display: flex;
		align-items: center;
		justify-content: space-between;
		cursor: pointer;
		text-align: left;
	}
	
	.dropdown-trigger:hover {
		border-color: #007cba;
	}
	
	.selected-content {
		display: flex;
		align-items: center;
		gap: 10px;
		flex: 1;
	}
	
	.icon {
		font-size: 18px;
		min-width: 24px;
	}
	
	.text {
		flex: 1;
	}
	
	.main {
		font-weight: 500;
		color: #333;
	}
	
	.meta {
		font-size: 12px;
		color: #666;
		margin-top: 2px;
	}
	
	.placeholder {
		color: #999;
	}
	
	.arrow {
		transition: transform 0.2s;
		color: #666;
	}
	
	.arrow.open {
		transform: rotate(180deg);
	}
	
	.dropdown-menu {
		position: absolute;
		top: 100%;
		left: 0;
		right: 0;
		background: white;
		border: 1px solid #ccc;
		border-radius: 6px;
		box-shadow: 0 4px 12px rgba(0,0,0,0.15);
		z-index: 1000;
		max-height: 300px;
		overflow-y: auto;
	}
	
	.dropdown-item {
		width: 100%;
		padding: 12px;
		border: none;
		background: none;
		display: flex;
		align-items: center;
		gap: 10px;
		cursor: pointer;
		text-align: left;
	}
	
	.dropdown-item:hover {
		background: #f8f9fa;
	}
	
	.add-new {
		border-top: 1px solid #eee;
		color: #007cba;
		font-weight: 500;
	}
	
	.add-new:hover {
		background: #e3f2fd;
	}
	
	.separator {
		height: 1px;
		background: #eee;
		margin: 5px 0;
	}
</style>