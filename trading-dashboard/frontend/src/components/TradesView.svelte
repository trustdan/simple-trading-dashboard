<script>
	import TradeCell from './TradeCell.svelte';
	import TradeModal from './TradeModal.svelte';
	import TradeFilters from './TradeFilters.svelte';
	import TradeAnalytics from './TradeAnalytics.svelte';
	import TradeHeatMap from './TradeHeatMap.svelte';
	import TradeExporter from './TradeExporter.svelte';
	import { onMount } from 'svelte';
	import { tradesStore } from '../stores/trades.js';
	import { toastStore } from '../stores/toast.js';

	let loading = false;
	let currentCenterDate = new Date(); // Changed from currentWeekStart to currentCenterDate
	
	// Modal state
	let isModalOpen = false;
	let selectedTradeForEdit = null;
	let selectedDateForNew = null;
	let selectedSectorForNew = null;

	// Filter state
	let filters = {
		status: 'all',
		strategy: 'all',
		sector: 'all',
		search: ''
	};
	let filteredTrades = []; // Initialize with empty array
	
	// Table state
	let allTrades = []; // All trades regardless of date filter

	// View state
	let currentView = 'grid'; // 'grid', 'analytics', 'heatmap'
	
	// Memoization for expensive operations
	let tradesCache = new Map();
	let lastFilters = null;
	let lastTradesLength = 0;
	
	// Get trades and date columns from store with safe defaults
	$: trades = $tradesStore?.trades || [];
	$: dateColumns = $tradesStore?.dateColumns || [];
	$: sectors = $tradesStore?.sectors || [];

	// Apply filters to trades with memoization (only if store is ready)
	$: if ($tradesStore && trades) {
		// Only recompute filtered trades if trades or filters actually changed
		const filtersKey = JSON.stringify(filters);
		const tradesChanged = trades.length !== lastTradesLength || lastFilters !== filtersKey;
		
		if (tradesChanged) {
			filteredTrades = applyFilters(trades, filters);
			lastFilters = filtersKey;
			lastTradesLength = trades.length;
			
			// Clear cache when trades change
			tradesCache.clear();
		}
	}

	// Set the current center date to today
	onMount(() => {
		const today = new Date();
		currentCenterDate = new Date(today); // Start centered on today
		
		generateDateColumns();
		loadTrades();
	});

	function generateDateColumns() {
		const columns = [];
		const centerDate = new Date(currentCenterDate);
		
		// Generate 42 days total: 14 days back + today + 27 days forward (6 weeks total)
		const startDate = new Date(centerDate);
		startDate.setDate(centerDate.getDate() - 14); // 2 weeks back
		
		for (let i = 0; i < 42; i++) {
			const date = new Date(startDate);
			date.setDate(startDate.getDate() + i);
			columns.push(date);
		}
		
		tradesStore.setDateColumns(columns);
	}

	async function loadTrades() {
		loading = true;
		try {
			const centerDate = new Date(currentCenterDate);
			const startDate = new Date(centerDate);
			startDate.setDate(centerDate.getDate() - 14); // 2 weeks back
			const endDate = new Date(centerDate);
			endDate.setDate(centerDate.getDate() + 27); // ~4 weeks forward
			
			await tradesStore.loadTradesByDateRange(startDate, endDate);
			
			// Also load all trades for the table (wider date range)
			await loadAllTrades();
		} catch (error) {
			console.error('Failed to load trades:', error);
			toastStore.error('Failed to load trades');
		} finally {
			loading = false;
		}
	}
	
	async function loadAllTrades() {
		try {
			const startDate = new Date();
			startDate.setMonth(startDate.getMonth() - 6); // 6 months ago
			const endDate = new Date();
			endDate.setMonth(endDate.getMonth() + 6); // 6 months ahead
			
			const allTradesData = await window['go']['main']['App']['GetTrades'](startDate, endDate);
			allTrades = allTradesData || [];
		} catch (error) {
			console.error('Failed to load all trades:', error);
			allTrades = [];
		}
	}

	function navigateWeek(direction) {
		const newCenter = new Date(currentCenterDate);
		newCenter.setDate(currentCenterDate.getDate() + (direction * 7));
		currentCenterDate = newCenter;
		
		generateDateColumns();
		loadTrades();
	}

	function navigateDay(direction) {
		const newCenter = new Date(currentCenterDate);
		newCenter.setDate(currentCenterDate.getDate() + direction);
		currentCenterDate = newCenter;
		
		generateDateColumns();
		loadTrades();
	}

	function goToToday() {
		const today = new Date();
		currentCenterDate = new Date(today);
		generateDateColumns();
		loadTrades();
		
		// Scroll the current day into view
		setTimeout(() => {
			const todayElements = document.querySelectorAll('.date-header.today');
			if (todayElements.length > 0) {
				todayElements[0].scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'center' });
			}
		}, 100);
	}

	// Modal handlers
	function openNewTradeModal(date = null, sector = null) {
		selectedTradeForEdit = null;
		selectedDateForNew = date ? date.toISOString().split('T')[0] : null;
		selectedSectorForNew = sector;
		isModalOpen = true;
	}

	function openEditTradeModal(trade) {
		selectedTradeForEdit = trade;
		selectedDateForNew = null;
		selectedSectorForNew = null;
		isModalOpen = true;
	}

	function closeModal() {
		isModalOpen = false;
		selectedTradeForEdit = null;
		selectedDateForNew = null;
		selectedSectorForNew = null;
	}

	async function handleTradeSaved() {
		closeModal();
		// Force refresh trades immediately after saving
		await loadTrades();
	}
	
	// Delete trade functionality
	async function deleteTrade(trade) {
		if (!confirm(`Are you sure you want to delete the ${trade.strategy_type} trade for ${trade.ticker}?`)) {
			return;
		}
		
		try {
			await window['go']['main']['App']['DeleteTrade'](trade.id);
			toastStore.success('Trade deleted successfully!');
			await loadTrades();
		} catch (error) {
			console.error('Failed to delete trade:', error);
			toastStore.error('Failed to delete trade');
		}
	}
	
	// Edit trade functionality 
	function editTrade(trade) {
		openEditTradeModal(trade);
	}

	function formatDateColumn(date) {
		const today = new Date();
		const isToday = date.toDateString() === today.toDateString();
		
		return {
			dayName: date.toLocaleDateString('en-US', { weekday: 'short' }),
			dayNumber: date.getDate(),
			month: date.toLocaleDateString('en-US', { month: 'short' }),
			isToday,
			isWeekend: date.getDay() === 0 || date.getDay() === 6
		};
	}

	function applyFilters(trades, filters) {
		return trades.filter(trade => {
			// Status filter
			if (filters.status !== 'all' && trade.status !== filters.status) {
				return false;
			}

			// Strategy filter
			if (filters.strategy !== 'all' && trade.strategy_type !== filters.strategy) {
				return false;
			}

			// Sector filter
			if (filters.sector !== 'all' && trade.sector !== filters.sector) {
				return false;
			}

			// Search filter
			if (filters.search.trim() !== '') {
				const searchTerm = filters.search.toLowerCase();
				const ticker = trade.ticker.toLowerCase();
				if (!ticker.includes(searchTerm)) {
					return false;
				}
			}

			return true;
		});
	}

	function getTradesForCell(sector, date) {
		// Create a cache key for memoization
		const cacheKey = `${sector}-${date.getTime()}`;
		
		// Return cached result if available
		if (tradesCache.has(cacheKey)) {
			return tradesCache.get(cacheKey);
		}
		
		// Convert date once for comparison
		const cellTime = date.getTime();
		
		const result = filteredTrades.filter(trade => {
			// Quick sector check first (most trades will fail this)
			if (trade.sector !== sector) {
				return false;
			}
			
			// Parse dates only once per trade (cache in trade object if needed)
			if (!trade._entryTime) {
				trade._entryTime = new Date(trade.entry_date).getTime();
			}
			if (!trade._expirationTime) {
				trade._expirationTime = new Date(trade.expiration_date).getTime();
			}
			
			return cellTime >= trade._entryTime && cellTime <= trade._expirationTime;
		});
		
		// Cache the result
		tradesCache.set(cacheKey, result);
		return result;
	}

	function handleFiltersChange(event) {
		filters = event.detail;
	}

	async function handleTradeStatusChange(event) {
		const { trade, newStatus } = event.detail;
		
		try {
			await window['go']['main']['App']['UpdateTradeStatus'](trade.id, newStatus);
			toastStore.success(`Trade status updated to ${newStatus}`);
			// Reload trades to get updated data
			loadTrades();
		} catch (error) {
			console.error('Error updating trade status:', error);
			toastStore.error(`Failed to update trade status: ${error.message || error}`);
		}
	}
</script>

<div class="trades-view">
	<header class="trades-header">
		<div class="header-content">
			<h1>Options Trading Calendar</h1>
			<p>Track your options positions across time and sectors</p>
		</div>
		
		<div class="header-controls">
			<div class="navigation-group">
				<button class="nav-button week-nav" on:click={() => navigateWeek(-1)}>
					← Week
				</button>
				<button class="nav-button day-nav" on:click={() => navigateDay(-1)}>
					← Day
				</button>
				
				<button class="today-button" on:click={goToToday}>
					Today
				</button>
				
				<button class="nav-button day-nav" on:click={() => navigateDay(1)}>
					Day →
				</button>
				<button class="nav-button week-nav" on:click={() => navigateWeek(1)}>
					Week →
				</button>
			</div>
			
			<div class="view-switcher">
				<button 
					class="view-btn" 
					class:active={currentView === 'grid'}
					on:click={() => currentView = 'grid'}
				>
					📅 Grid
				</button>
				<button 
					class="view-btn" 
					class:active={currentView === 'analytics'}
					on:click={() => currentView = 'analytics'}
				>
					📊 Analytics
				</button>
				<button 
					class="view-btn" 
					class:active={currentView === 'heatmap'}
					on:click={() => currentView = 'heatmap'}
				>
					🔥 Heat Map
				</button>
			</div>

			<button class="new-trade-button" on:click={() => openNewTradeModal()}>
				+ New Trade
			</button>
		</div>
	</header>

	{#if loading}
		<div class="loading-overlay">
			<div class="loading-spinner"></div>
			<p>Loading trades...</p>
		</div>
	{/if}

	<main class="trades-content">
		<!-- Filters (only show for grid view) -->
		{#if currentView === 'grid'}
			<TradeFilters bind:filters on:filters-change={handleFiltersChange} />
		{/if}

		<!-- Analytics View -->
		{#if currentView === 'analytics'}
			<TradeAnalytics />
			<TradeExporter />
		{:else if currentView === 'heatmap'}
			<TradeHeatMap 
				trades={filteredTrades} 
				{dateColumns} 
				{sectors}
				on:cell-click={(event) => {
					const { trades } = event.detail;
					if (trades.length > 0) {
						openEditTradeModal(trades[0]);
					}
				}}
			/>
		{:else}
			<!-- Grid View -->
			<div class="trade-grid">
			<!-- Header row with dates -->
			<div class="grid-header">
				<div class="sector-header">Sector</div>
				{#each dateColumns || [] as date}
					{@const formatted = formatDateColumn(date)}
					<div class="date-header" class:today={formatted.isToday} class:weekend={formatted.isWeekend}>
						<div class="day-name">{formatted.dayName}</div>
						<div class="day-number">{formatted.dayNumber}</div>
						<div class="month">{formatted.month}</div>
					</div>
				{/each}
			</div>

			<!-- Grid body with sector rows -->
			<div class="grid-body">
				{#each sectors || [] as sector}
					<div class="sector-row">
						<div class="sector-label">
							{sector}
						</div>
						
						{#each dateColumns || [] as date}
							{@const cellTrades = getTradesForCell(sector, date)}
							<TradeCell 
								{sector}
								{date}
								trades={cellTrades}
								on:cellClick={(event) => {
									const detail = event.detail;
									if (detail.trades.length > 0) {
										openEditTradeModal(detail.trades[0]);
									} else {
										openNewTradeModal(detail.date, detail.sector);
									}
								}}
								on:tradeClick={(event) => openEditTradeModal(event.detail.trade)}
								on:status-change={handleTradeStatusChange}
							/>
						{/each}
					</div>
				{/each}
			</div>
			</div>

			{#if filteredTrades.length === 0 && trades.length > 0 && !loading}
				<div class="empty-state">
					<div class="empty-icon">🔍</div>
					<h3>No trades match your filters</h3>
					<p>Try adjusting your filter criteria to see more trades</p>
				</div>
			{:else if trades.length === 0 && !loading}
				<div class="empty-state">
					<div class="empty-icon">📊</div>
					<h3>No trades found</h3>
					<p>Create your first options trade to see it on the calendar</p>
					<button class="create-trade-button" on:click={() => openNewTradeModal()}>
						+ Add Trade
					</button>
				</div>
			{/if}
		{/if}
	</main>

	<!-- Trades Table Section -->
	{#if currentView === 'grid' && allTrades.length > 0}
		<section class="trades-table-section">
			<div class="table-header">
				<h2>All Trades</h2>
				<p>Manage your existing trades</p>
			</div>
			
			<div class="trades-table-container">
				<table class="trades-table">
					<thead>
						<tr>
							<th>Ticker</th>
							<th>Strategy</th>
							<th>Sector</th>
							<th>Entry Date</th>
							<th>Expiration</th>
							<th>Target Price</th>
							<th>Stop Loss</th>
							<th>Status</th>
							<th>Actions</th>
						</tr>
					</thead>
					<tbody>
						{#each allTrades as trade}
							<tr>
								<td class="ticker-cell">{trade.ticker}</td>
								<td class="strategy-cell">
									<span class="strategy-badge">{trade.strategy_type}</span>
								</td>
								<td>{trade.sector}</td>
								<td>{new Date(trade.entry_date).toLocaleDateString()}</td>
								<td>{new Date(trade.expiration_date).toLocaleDateString()}</td>
								<td>{trade.target_price ? `$${trade.target_price.toFixed(2)}` : '-'}</td>
								<td>{trade.stop_loss ? `$${trade.stop_loss.toFixed(2)}` : '-'}</td>
								<td>
									<span class="status-badge status-{trade.status}">{trade.status}</span>
								</td>
								<td class="actions-cell">
									<button 
										class="action-btn edit-btn" 
										on:click={() => editTrade(trade)}
										title="Edit trade"
									>
										✏️
									</button>
									<button 
										class="action-btn delete-btn" 
										on:click={() => deleteTrade(trade)}
										title="Delete trade"
									>
										❌
									</button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</section>
	{/if}
</div>

<!-- Trade Modal -->
<TradeModal 
	bind:isOpen={isModalOpen}
	trade={selectedTradeForEdit}
	selectedDate={selectedDateForNew}
	selectedSector={selectedSectorForNew}
	on:close={closeModal}
	on:trade-saved={handleTradeSaved}
/>

<style>
	.trades-view {
		min-height: 100vh;
		background: linear-gradient(135deg, #1a1a1a, #2d2d2d);
		color: #ffffff;
	}

	.trades-header {
		padding: 24px;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
		display: flex;
		justify-content: space-between;
		align-items: center;
		flex-wrap: wrap;
		gap: 20px;
	}

	.header-content h1 {
		font-size: 2.5rem;
		font-weight: 700;
		margin-bottom: 8px;
		background: linear-gradient(135deg, #4a90e2, #7b68ee);
		background-clip: text;
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	.header-content p {
		font-size: 1.1rem;
		color: #aaa;
		margin: 0;
	}

	.header-controls {
		display: flex;
		gap: 12px;
		align-items: center;
	}

	.nav-button {
		background: rgba(255, 255, 255, 0.1);
		border: 1px solid rgba(255, 255, 255, 0.2);
		color: white;
		padding: 12px 20px;
		border-radius: 8px;
		cursor: pointer;
		transition: all 0.3s ease;
		font-size: 0.9rem;
	}

	.nav-button:hover {
		background: rgba(255, 255, 255, 0.15);
		border-color: rgba(74, 144, 226, 0.5);
	}

	.navigation-group {
		display: flex;
		gap: 8px;
		align-items: center;
	}

	.day-nav {
		font-size: 0.85rem;
		padding: 10px 16px;
		background: rgba(255, 255, 255, 0.08);
	}

	.week-nav {
		font-size: 0.9rem;
		padding: 12px 20px;
	}

	.today-button {
		background: linear-gradient(135deg, #4a90e2, #7b68ee);
		border: none;
		color: white;
		padding: 12px 24px;
		border-radius: 8px;
		cursor: pointer;
		transition: all 0.3s ease;
		font-weight: 600;
	}

	.today-button:hover {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(74, 144, 226, 0.3);
	}

	.new-trade-button {
		background: linear-gradient(135deg, #4a90e2, #7b68ee);
		border: none;
		color: white;
		padding: 12px 24px;
		border-radius: 8px;
		cursor: pointer;
		transition: all 0.3s ease;
		font-weight: 600;
		font-size: 0.9rem;
	}

	.new-trade-button:hover {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(74, 144, 226, 0.3);
	}

	.view-switcher {
		display: flex;
		background: #2a2a2a;
		border-radius: 8px;
		padding: 4px;
		gap: 2px;
	}

	.view-btn {
		background: none;
		border: none;
		color: #cccccc;
		padding: 8px 16px;
		border-radius: 6px;
		cursor: pointer;
		font-size: 14px;
		font-weight: 500;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		gap: 6px;
	}

	.view-btn:hover {
		color: #ffffff;
		background: rgba(255, 255, 255, 0.05);
	}

	.view-btn.active {
		background: linear-gradient(135deg, #4a90e2, #7b68ee);
		color: white;
		transform: translateY(-1px);
		box-shadow: 0 2px 8px rgba(74, 144, 226, 0.3);
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

	.trades-content {
		padding: 24px;
		max-width: 100%;
		overflow-x: auto;
	}

	.trade-grid {
		min-width: 2760px; /* 200px + (42 * 80px) = 3560px, but let's be conservative */
		background: rgba(255, 255, 255, 0.02);
		border-radius: 12px;
		overflow: hidden;
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.grid-header {
		display: grid;
		grid-template-columns: 200px repeat(42, 80px);
		background: rgba(255, 255, 255, 0.05);
		border-bottom: 2px solid rgba(255, 255, 255, 0.1);
	}

	.sector-header {
		padding: 16px;
		font-weight: 600;
		font-size: 1.1rem;
		display: flex;
		align-items: center;
		background: rgba(74, 144, 226, 0.1);
		border-right: 1px solid rgba(255, 255, 255, 0.1);
	}

	.date-header {
		padding: 12px 8px;
		text-align: center;
		border-right: 1px solid rgba(255, 255, 255, 0.05);
		transition: background-color 0.2s ease;
	}

	.date-header.today {
		background: rgba(74, 144, 226, 0.2);
		color: #4a90e2;
		font-weight: 600;
	}

	.date-header.weekend {
		background: rgba(255, 255, 255, 0.02);
		color: #888;
	}

	.day-name {
		font-size: 0.8rem;
		text-transform: uppercase;
		letter-spacing: 1px;
		margin-bottom: 4px;
	}

	.day-number {
		font-size: 1.2rem;
		font-weight: 600;
		margin-bottom: 2px;
	}

	.month {
		font-size: 0.7rem;
		opacity: 0.7;
		text-transform: uppercase;
	}

	.grid-body {
		display: grid;
		grid-template-rows: repeat(auto, 1fr);
	}

	.sector-row {
		display: grid;
		grid-template-columns: 200px repeat(42, 80px);
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
	}

	.sector-row:last-child {
		border-bottom: none;
	}

	.sector-label {
		padding: 20px 16px;
		font-weight: 500;
		font-size: 0.9rem;
		display: flex;
		align-items: center;
		background: rgba(255, 255, 255, 0.02);
		border-right: 1px solid rgba(255, 255, 255, 0.1);
		color: #ccc;
	}

	.empty-state {
		text-align: center;
		padding: 80px 24px;
		color: #888;
	}

	.empty-icon {
		font-size: 4rem;
		margin-bottom: 24px;
		opacity: 0.5;
	}

	.empty-state h3 {
		font-size: 1.5rem;
		margin-bottom: 12px;
		color: #ccc;
	}

	.empty-state p {
		font-size: 1.1rem;
		margin-bottom: 32px;
		max-width: 400px;
		margin-left: auto;
		margin-right: auto;
	}

	.create-trade-button {
		background: linear-gradient(135deg, #4a90e2, #7b68ee);
		color: white;
		border: none;
		padding: 16px 32px;
		font-size: 1.1rem;
		font-weight: 600;
		border-radius: 8px;
		cursor: pointer;
		transition: all 0.3s ease;
	}

	.create-trade-button:hover {
		transform: translateY(-2px);
		box-shadow: 0 8px 16px rgba(74, 144, 226, 0.3);
	}

	@media (max-width: 1200px) {
		.trades-header {
			flex-direction: column;
			align-items: flex-start;
		}

		.header-controls {
			width: 100%;
			justify-content: center;
		}

		.trade-grid {
			min-width: 2400px;
		}
	}

	/* Trades Table Styles */
	.trades-table-section {
		padding: 24px;
		border-top: 1px solid rgba(255, 255, 255, 0.1);
		background: rgba(0, 0, 0, 0.2);
	}
	
	.table-header {
		margin-bottom: 20px;
	}
	
	.table-header h2 {
		font-size: 1.8rem;
		font-weight: 600;
		margin-bottom: 8px;
		color: #ffffff;
	}
	
	.table-header p {
		color: #aaa;
		margin: 0;
	}
	
	.trades-table-container {
		background: rgba(255, 255, 255, 0.05);
		border-radius: 12px;
		overflow: hidden;
		border: 1px solid rgba(255, 255, 255, 0.1);
	}
	
	.trades-table {
		width: 100%;
		border-collapse: collapse;
	}
	
	.trades-table th {
		background: rgba(74, 144, 226, 0.2);
		color: #ffffff;
		padding: 16px 12px;
		text-align: left;
		font-weight: 600;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}
	
	.trades-table td {
		padding: 12px;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
		color: #cccccc;
	}
	
	.trades-table tbody tr:hover {
		background: rgba(255, 255, 255, 0.03);
	}
	
	.ticker-cell {
		font-weight: 600;
		font-family: monospace;
		color: #4a90e2;
	}
	
	.strategy-badge {
		background: linear-gradient(135deg, #4a90e2, #7b68ee);
		color: white;
		padding: 4px 8px;
		border-radius: 12px;
		font-size: 0.8rem;
		font-weight: 500;
	}
	
	.status-badge {
		padding: 4px 8px;
		border-radius: 12px;
		font-size: 0.8rem;
		font-weight: 500;
		text-transform: uppercase;
	}
	
	.status-active {
		background: rgba(34, 197, 94, 0.2);
		color: #22c55e;
	}
	
	.status-closed {
		background: rgba(156, 163, 175, 0.2);
		color: #9ca3af;
	}
	
	.status-expired {
		background: rgba(239, 68, 68, 0.2);
		color: #ef4444;
	}
	
	.actions-cell {
		text-align: center;
	}
	
	.action-btn {
		background: none;
		border: none;
		padding: 8px;
		margin: 0 2px;
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.2s ease;
		font-size: 1rem;
	}
	
	.action-btn:hover {
		background: rgba(255, 255, 255, 0.1);
		transform: scale(1.1);
	}
	
	.edit-btn:hover {
		background: rgba(74, 144, 226, 0.2);
	}
	
	.delete-btn:hover {
		background: rgba(239, 68, 68, 0.2);
	}

	@media (max-width: 768px) {
		.trades-view {
			padding: 12px;
		}

		.header-content h1 {
			font-size: 2rem;
		}

		.nav-button, .today-button {
			padding: 10px 16px;
			font-size: 0.8rem;
		}
		
		.trades-table-section {
			padding: 16px;
		}
		
		.trades-table-container {
			overflow-x: auto;
		}
		
		.trades-table {
			min-width: 800px;
		}
	}
</style> 