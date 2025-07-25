<script>
	import { createEventDispatcher } from 'svelte';
	
	export let currentView = 'market'; // 'market' or 'trades'
	
	const dispatch = createEventDispatcher();
	
	function switchView(view) {
		if (view !== currentView) {
			currentView = view;
			dispatch('viewChange', view);
		}
	}
	
	// Keyboard navigation
	function handleKeyDown(event, view) {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			switchView(view);
		}
	}
</script>

<nav class="main-nav">
	<div class="nav-container">
		<div class="nav-brand">
			<span class="brand-icon">📈</span>
			<span class="brand-text">Trading Dashboard</span>
		</div>
		
		<div class="nav-buttons">
			<button 
				class="nav-button"
				class:active={currentView === 'market'}
				on:click={() => switchView('market')}
				on:keydown={(e) => handleKeyDown(e, 'market')}
				aria-pressed={currentView === 'market'}
			>
				<span class="button-icon">🎯</span>
				<span class="button-text">Market Sentiment</span>
			</button>
			
			<button 
				class="nav-button"
				class:active={currentView === 'trades'}
				on:click={() => switchView('trades')}
				on:keydown={(e) => handleKeyDown(e, 'trades')}
				aria-pressed={currentView === 'trades'}
			>
				<span class="button-icon">📊</span>
				<span class="button-text">Trade Management</span>
			</button>
		</div>
		
		<div class="nav-info">
			<div class="status-indicator" class:online={true}>
				<span class="status-dot"></span>
				<span class="status-text">Live</span>
			</div>
		</div>
	</div>
</nav>

<style>
	.main-nav {
		background: rgba(0, 0, 0, 0.8);
		backdrop-filter: blur(10px);
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
		position: sticky;
		top: 0;
		z-index: 100;
	}

	.nav-container {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 12px 24px;
		max-width: 1400px;
		margin: 0 auto;
	}

	.nav-brand {
		display: flex;
		align-items: center;
		gap: 12px;
		font-weight: 700;
		font-size: 1.2rem;
		color: #ffffff;
	}

	.brand-icon {
		font-size: 1.5rem;
	}

	.brand-text {
		background: linear-gradient(135deg, #4a90e2, #7b68ee);
		background-clip: text;
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	.nav-buttons {
		display: flex;
		gap: 8px;
		background: rgba(255, 255, 255, 0.05);
		padding: 4px;
		border-radius: 12px;
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.nav-button {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 12px 20px;
		border: none;
		background: transparent;
		color: #ccc;
		border-radius: 8px;
		cursor: pointer;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		font-weight: 500;
		font-size: 0.9rem;
		position: relative;
		overflow: hidden;
	}

	.nav-button::before {
		content: '';
		position: absolute;
		top: 0;
		left: -100%;
		width: 100%;
		height: 100%;
		background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
		transition: left 0.5s;
	}

	.nav-button:hover::before {
		left: 100%;
	}

	.nav-button:hover {
		color: #ffffff;
		background: rgba(255, 255, 255, 0.1);
		transform: translateY(-1px);
	}

	.nav-button.active {
		background: linear-gradient(135deg, #4a90e2, #7b68ee);
		color: #ffffff;
		font-weight: 600;
		box-shadow: 0 4px 12px rgba(74, 144, 226, 0.3);
	}

	.nav-button.active:hover {
		transform: translateY(-1px);
		box-shadow: 0 6px 16px rgba(74, 144, 226, 0.4);
	}

	.button-icon {
		font-size: 1.1rem;
	}

	.button-text {
		white-space: nowrap;
	}

	.nav-info {
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.status-indicator {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 6px 12px;
		border-radius: 20px;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		font-size: 0.8rem;
		color: #ccc;
	}

	.status-indicator.online {
		color: #22c55e;
		border-color: rgba(34, 197, 94, 0.3);
		background: rgba(34, 197, 94, 0.1);
	}

	.status-dot {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: currentColor;
		animation: pulse 2s infinite;
	}

	@keyframes pulse {
		0%, 100% { opacity: 1; }
		50% { opacity: 0.5; }
	}

	.status-text {
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 1px;
	}

	/* Responsive design */
	@media (max-width: 768px) {
		.nav-container {
			padding: 8px 16px;
			flex-wrap: wrap;
			gap: 12px;
		}

		.nav-brand {
			font-size: 1rem;
		}

		.brand-text {
			display: none;
		}

		.nav-buttons {
			flex: 1;
			justify-content: center;
		}

		.nav-button {
			padding: 10px 16px;
			font-size: 0.8rem;
		}

		.button-text {
			display: none;
		}

		.button-icon {
			font-size: 1.2rem;
		}

		.nav-info {
			gap: 8px;
		}

		.status-text {
			display: none;
		}
	}

	@media (max-width: 480px) {
		.nav-container {
			flex-direction: column;
			gap: 8px;
		}

		.nav-buttons {
			width: 100%;
		}
	}
</style> 