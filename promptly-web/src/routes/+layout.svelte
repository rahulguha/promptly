<script>
	import '../app.css';
</script>

<div class="app">
	<header>
		<div class="header-content">
			<img src="/promptly-logo.png" alt="Promptly" class="logo" />
			<div class="header-text">
				<h1>Promptly</h1> 
				<p>opensourced Prompt Management System</p>
			</div>
		</div>
	</header>
	
	<main>
		<slot />
	</main>
</div>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
		background: #fafafa;
	}
	
	.app {
		min-height: 100vh;
	}
	
	header {
		background: #ebdf8f;
		color: rgb(9, 9, 9);
		padding: 2px 2px;
	}
	
	.header-content {
		display: flex;
		align-items: center;
		justify-content: flex-start;
		gap: 12px;
		max-width: 1200px;
		margin: 0 auto;
	}
	
	.logo {
		width: 100px;
		height: 100px;
		object-fit: contain;
	}
	
	.header-text {
		text-align: center;
		flex: 1;
	}
	
	header h1 {
		margin: 0;
		font-size: 2em;
	}
	
	header p {
		margin: 5px 0 0 0;
		opacity: 0.9;
	}
	
	main {
		max-width: 1200px;
		margin: 0 auto;
		padding: 5px;
	}
</style>