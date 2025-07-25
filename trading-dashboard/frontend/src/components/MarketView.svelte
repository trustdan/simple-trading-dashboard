<script>
	import Dial from './Dial.svelte';
	import ToastContainer from './ToastContainer.svelte';
	import { onMount } from 'svelte';
	import { marketStore } from '../stores/market.js';
	import { toast } from '../stores/toast.js';

	let loading = false;
	let saving = false;
	let lastSaved = null;
	let autoSaveTimeout;

	// Get sectors from the store
	$: sectors = $marketStore.sectors;
	$: overallRating = $marketStore.overallRating;
	$: sectorRatings = $marketStore.sectorRatings;

	onMount(async () => {
		await loadLatestRating();
		
		// Start auto-save watcher
		const unsubscribe = marketStore.subscribe(state => {
			if (autoSaveTimeout) clearTimeout(autoSaveTimeout);
			autoSaveTimeout = setTimeout(() => {
				if (!saving && !loading) {
					autoSave();
				}
			}, 2000); // Auto-save after 2 seconds of inactivity
		});
		
		return unsubscribe;
	});

	async function loadLatestRating() {
		loading = true;
		try {
			await marketStore.loadLatestRating();
			toast.info('Market sentiment loaded');
		} catch (error) {
			console.error('Failed to load latest rating:', error);
			toast.error('Failed to load market data');
		} finally {
			loading = false;
		}
	}

	async function saveRatings(showToast = true) {
		saving = true;
		try {
			const result = await marketStore.saveRating();
			lastSaved = new Date();
			if (showToast) {
				toast.success('Market sentiment saved successfully!');
			}
			return result;
		} catch (error) {
			console.error('Failed to save ratings:', error);
			if (showToast) {
				toast.error('Failed to save market sentiment');
			}
			throw error;
		} finally {
			saving = false;
		}
	}

	async function autoSave() {
		try {
			await saveRatings(false); // Don't show toast for auto-save
		} catch (error) {
			// Silently fail auto-save, user can manually save
		}
	}

	function handleOverallChange(value) {
		marketStore.setOverallRating(value);
	}

	function handleSectorChange(sector, value) {
		marketStore.setSectorRating(sector, value);
	}

	function formatLastSaved(date) {
		if (!date) return '';
		const now = new Date();
		const diff = now - date;
		
		if (diff < 60000) return 'Just now';
		if (diff < 3600000) return `${Math.floor(diff / 60000)} minutes ago`;
		if (diff < 86400000) return `${Math.floor(diff / 3600000)} hours ago`;
		return date.toLocaleDateString();
	}
</script>

<div class="market-view">
	<header class="market-header">
		<h1>Market Sentiment Dashboard</h1>
		<p>Set your market sentiment ratings from -3 (very bearish) to +3 (very bullish)</p>
	</header>

	{#if loading}
		<div class="loading-overlay">
			<div class="loading-spinner"></div>
			<p>Loading...</p>
		</div>
	{/if}

	<main class="market-content">
		<!-- Overall Market Rating -->
		<section class="overall-section">
			<h2>Overall Market Sentiment</h2>
			<div class="overall-dial">
				<Dial
					value={overallRating}
					size={300}
					label="Overall Market"
					onValueChange={handleOverallChange}
				/>
			</div>
		</section>

		<!-- Sector Ratings -->
		<section class="sectors-section">
			<h2>Sector Sentiment</h2>
			<div class="sector-grid">
				{#each sectors as sector}
					<div class="sector-item">
						<Dial
							value={sectorRatings[sector] || 0}
							size={200}
							label={sector}
							onValueChange={(value) => handleSectorChange(sector, value)}
						/>
					</div>
				{/each}
			</div>
		</section>

		<!-- Save Button -->
		<section class="actions-section">
			<button 
				class="save-button" 
				class:saving 
				class:pulse={saving}
				on:click={() => saveRatings(true)}
				disabled={loading || saving}
			>
				<span class="save-button-text">
					{saving ? 'Saving...' : 'Save Ratings'}
				</span>
				<span class="save-button-icon" class:saving>
					💾
				</span>
			</button>
			
			{#if lastSaved}
				<div class="last-saved">
					Last saved: {formatLastSaved(lastSaved)}
				</div>
			{/if}
		</section>
	</main>
	
	<!-- Toast notifications -->
	<ToastContainer />
</div>

<style>
	.market-view {
		min-height: 100vh;
		background: linear-gradient(135deg, #1a1a1a, #2d2d2d);
		color: #ffffff;
		padding: 24px;
	}

	.market-header {
		text-align: center;
		margin-bottom: 32px;
	}

	.market-header h1 {
		font-size: 2.5rem;
		font-weight: 700;
		margin-bottom: 12px;
		background: linear-gradient(135deg, #4a90e2, #7b68ee);
		background-clip: text;
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	.market-header p {
		font-size: 1.1rem;
		color: #aaa;
		max-width: 600px;
		margin: 0 auto;
	}

	.loading-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.7);
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.loading-spinner {
		width: 50px;
		height: 50px;
		border: 4px solid #333;
		border-top: 4px solid #4a90e2;
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin-bottom: 16px;
	}

	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}

	.market-content {
		max-width: 1400px;
		margin: 0 auto;
	}

	.overall-section {
		text-align: center;
		margin-bottom: 48px;
		padding: 32px;
		background: rgba(255, 255, 255, 0.05);
		border-radius: 16px;
		backdrop-filter: blur(10px);
	}

	.overall-section h2 {
		font-size: 1.8rem;
		margin-bottom: 24px;
		color: #4a90e2;
	}

	.overall-dial {
		display: flex;
		justify-content: center;
	}

	.sectors-section {
		margin-bottom: 48px;
	}

	.sectors-section h2 {
		font-size: 1.8rem;
		text-align: center;
		margin-bottom: 32px;
		color: #4a90e2;
	}

	.sector-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 32px;
		padding: 0 16px;
	}

	.sector-item {
		display: flex;
		justify-content: center;
		padding: 24px;
		background: rgba(255, 255, 255, 0.03);
		border-radius: 12px;
		border: 1px solid rgba(255, 255, 255, 0.1);
		transition: all 0.3s ease;
	}

	.sector-item:hover {
		background: rgba(255, 255, 255, 0.05);
		border-color: rgba(74, 144, 226, 0.3);
		transform: translateY(-2px);
	}

	.actions-section {
		text-align: center;
		margin-top: 32px;
	}

	.save-button {
		background: linear-gradient(135deg, #4a90e2, #7b68ee);
		color: white;
		border: none;
		padding: 16px 32px;
		font-size: 1.1rem;
		font-weight: 600;
		border-radius: 12px;
		cursor: pointer;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		min-width: 180px;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 12px;
		position: relative;
		overflow: hidden;
	}

	.save-button::before {
		content: '';
		position: absolute;
		top: 0;
		left: -100%;
		width: 100%;
		height: 100%;
		background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
		transition: left 0.5s;
	}

	.save-button:hover:not(:disabled)::before {
		left: 100%;
	}

	.save-button:hover:not(:disabled) {
		transform: translateY(-2px);
		box-shadow: 0 12px 24px rgba(74, 144, 226, 0.4);
	}

	.save-button:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.save-button.saving {
		animation: saveGlow 2s infinite;
	}

	.save-button.pulse {
		animation: saveGlow 2s infinite, buttonPulse 0.5s ease-in-out;
	}

	.save-button-icon {
		transition: transform 0.3s ease;
	}

	.save-button-icon.saving {
		animation: iconSpin 1s linear infinite;
	}

	@keyframes saveGlow {
		0%, 100% {
			box-shadow: 0 0 20px rgba(74, 144, 226, 0.3);
		}
		50% {
			box-shadow: 0 0 30px rgba(74, 144, 226, 0.6);
		}
	}

	@keyframes buttonPulse {
		0% { transform: scale(1); }
		50% { transform: scale(1.05); }
		100% { transform: scale(1); }
	}

	@keyframes iconSpin {
		from { transform: rotate(0deg); }
		to { transform: rotate(360deg); }
	}

	.last-saved {
		margin-top: 12px;
		font-size: 0.9rem;
		color: #888;
		text-align: center;
	}

	@media (max-width: 768px) {
		.market-view {
			padding: 16px;
		}

		.market-header h1 {
			font-size: 2rem;
		}

		.overall-section {
			padding: 24px 16px;
		}

		.sector-grid {
			grid-template-columns: 1fr;
			gap: 24px;
		}

		.overall-dial :global(.dial-container) {
			transform: scale(0.8);
		}
	}
</style> 