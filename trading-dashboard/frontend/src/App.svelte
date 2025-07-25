<script>
	import MarketView from './components/MarketView.svelte';
	import TradesView from './components/TradesView.svelte';
	import Navigation from './components/Navigation.svelte';
	import './app.css'

	let currentView = 'market'; // 'market' or 'trades'

	function handleViewChange(event) {
		currentView = event.detail;
	}
</script>

<div class="app">
	<Navigation {currentView} on:viewChange={handleViewChange} />
	
	<main class="app-main">
		{#if currentView === 'market'}
			<MarketView />
		{:else if currentView === 'trades'}
			<TradesView />
		{/if}
	</main>
</div>

<style>
	.app {
		min-height: 100vh;
		background: linear-gradient(135deg, #1a1a1a, #2d2d2d);
	}

	.app-main {
		min-height: calc(100vh - 80px);
	}
</style>
