<script>
	export let sector = '';
	export let date = new Date();
	export let trades = [];

	import { createEventDispatcher } from 'svelte';
	import TradeStatusBadge from './TradeStatusBadge.svelte';
	import TradeContextMenu from './TradeContextMenu.svelte';
	
	const dispatch = createEventDispatcher();

	// Context menu state
	let contextMenuVisible = false;
	let contextMenuX = 0;
	let contextMenuY = 0;
	let contextMenuTrade = null;

	// Group trades by strategy type for visual stacking
	$: groupedTrades = groupTradesByStrategy(trades);

	function groupTradesByStrategy(trades) {
		const groups = {};
		trades.forEach(trade => {
			if (!groups[trade.strategy_type]) {
				groups[trade.strategy_type] = [];
			}
			groups[trade.strategy_type].push(trade);
		});
		return groups;
	}

	function getStrategyColor(strategyType) {
		// Default colors based on strategy categories
		const colorMap = {
			// Bullish strategies (green tones)
			'Long Call': '#22c55e',
			'Bull Call Spread': '#16a34a', 
			'Cash-Secured Put': '#15803d',
			
			// Bearish strategies (red tones)
			'Long Put': '#ef4444',
			'Bear Put Spread': '#dc2626',
			'Covered Call': '#b91c1c',
			
			// Neutral strategies (purple tones)
			'Iron Condor': '#8b5cf6',
			'Butterfly Spread': '#7c3aed',
			'Straddle': '#6d28d9'
		};
		
		return colorMap[strategyType] || '#4a90e2';
	}

	function handleCellClick() {
		dispatch('cellClick', {
			sector,
			date,
			trades
		});
	}

	function handleTradeClick(trade, event) {
		event.stopPropagation();
		dispatch('tradeClick', {
			trade,
			sector,
			date
		});
	}

	function formatTicker(ticker) {
		return ticker.toUpperCase();
	}

	function isTradeExpiring(trade) {
		const expirationDate = new Date(trade.expiration_date);
		const cellDate = new Date(date);
		const diffDays = Math.ceil((expirationDate - cellDate) / (1000 * 60 * 60 * 24));
		return diffDays <= 7 && diffDays >= 0; // Expiring within a week
	}

	function isTradeStarting(trade) {
		const entryDate = new Date(trade.entry_date);
		const cellDate = new Date(date);
		return entryDate.toDateString() === cellDate.toDateString();
	}

	function getDaysToExpiration(trade) {
		const expirationDate = new Date(trade.expiration_date);
		const cellDate = new Date(date);
		return Math.ceil((expirationDate - cellDate) / (1000 * 60 * 60 * 24));
	}

	function handleStatusChange(event) {
		const { trade, newStatus } = event.detail;
		dispatch('status-change', { trade, newStatus });
	}

	function handleRightClick(trade, event) {
		event.preventDefault();
		event.stopPropagation();
		
		contextMenuTrade = trade;
		contextMenuX = event.clientX;
		contextMenuY = event.clientY;
		contextMenuVisible = true;
	}

	function closeContextMenu() {
		contextMenuVisible = false;
		contextMenuTrade = null;
	}

	function handleContextMenuEdit(event) {
		dispatch('tradeClick', { trade: event.detail });
	}

	function handleContextMenuStatusChange(event) {
		dispatch('status-change', event.detail);
	}

	function handleContextMenuDelete(event) {
		dispatch('delete-trade', event.detail);
	}

	function handleContextMenuDuplicate(event) {
		dispatch('duplicate-trade', event.detail);
	}
</script>

<div 
	class="trade-cell" 
	class:has-trades={trades.length > 0}
	class:clickable={trades.length > 0}
	on:click={handleCellClick}
	role="button"
	tabindex="0"
	on:keydown={(e) => e.key === 'Enter' && handleCellClick()}
>
	{#if trades.length > 0}
		<div class="trades-container">
			{#each Object.entries(groupedTrades) as [strategyType, strategyTrades]}
				<div class="strategy-group">
					{#each strategyTrades as trade}
						{@const daysToExpiry = getDaysToExpiration(trade)}
						<div 
							class="trade-item"
							class:expiring={isTradeExpiring(trade)}
							class:starting={isTradeStarting(trade)}
							class:active={trade.status === 'active'}
							class:closed={trade.status === 'closed'}
							class:expired={trade.status === 'expired'}
							style="background-color: {getStrategyColor(trade.strategy_type)}; border-color: {getStrategyColor(trade.strategy_type)}; opacity: {trade.status === 'active' ? 1 : 0.6};"
							on:click={(e) => handleTradeClick(trade, e)}
							on:contextmenu={(e) => handleRightClick(trade, e)}
							title="{trade.ticker} - {trade.strategy_type} ({trade.status}) - Right-click for options"
						>
							<div class="trade-header">
								<div class="trade-ticker">
									{formatTicker(trade.ticker)}
								</div>
								<TradeStatusBadge 
									{trade} 
									showDropdown={true} 
									size="small"
									on:status-change={handleStatusChange}
								/>
							</div>
							
							<div class="trade-details">
								<div class="strategy-name">{trade.strategy_type}</div>
								{#if daysToExpiry >= 0}
									<div class="expiry-info" class:urgent={daysToExpiry <= 3}>
										{daysToExpiry}d
									</div>
								{:else}
									<div class="expiry-info expired">EXPIRED</div>
								{/if}
							</div>

							{#if isTradeStarting(trade)}
								<div class="start-indicator">🚀</div>
							{/if}
							{#if isTradeExpiring(trade)}
								<div class="expiry-indicator">⚠️</div>
							{/if}
						</div>
					{/each}
				</div>
			{/each}
			
			{#if trades.length > 3}
				<div class="trade-overflow">
					+{trades.length - 3}
				</div>
			{/if}
		</div>
	{:else}
		<div class="empty-cell">
			<!-- Empty cell for potential trade placement -->
		</div>
	{/if}
</div>

<style>
	.trade-cell {
		min-height: 60px;
		padding: 4px;
		border-right: 1px solid rgba(255, 255, 255, 0.05);
		position: relative;
		transition: all 0.2s ease;
		cursor: default;
		background: transparent;
	}

	.trade-cell.clickable {
		cursor: pointer;
	}

	.trade-cell.clickable:hover {
		background: rgba(255, 255, 255, 0.05);
	}

	.trade-cell:focus {
		outline: 2px solid #4a90e2;
		outline-offset: -2px;
	}

	.trades-container {
		display: flex;
		flex-direction: column;
		gap: 2px;
		height: 100%;
	}

	.strategy-group {
		display: flex;
		flex-direction: column;
		gap: 1px;
	}

	.trade-item {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 2px 4px;
		border-radius: 3px;
		border: 1px solid;
		font-size: 0.7rem;
		font-weight: 600;
		color: white;
		min-height: 16px;
		position: relative;
		transition: all 0.2s ease;
		cursor: pointer;
		overflow: hidden;
	}

	.trade-item:hover {
		transform: scale(1.05);
		z-index: 10;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
	}

	.trade-item.expiring {
		animation: blinkExpiring 2s infinite;
		border-width: 2px;
	}

	.trade-item.starting {
		box-shadow: 0 0 8px rgba(255, 255, 255, 0.5);
		border-width: 2px;
	}

	@keyframes blinkExpiring {
		0%, 50% { opacity: 1; }
		51%, 100% { opacity: 0.6; }
	}

	.trade-ticker {
		font-weight: 700;
		letter-spacing: 0.5px;
		flex: 1;
		text-align: center;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.expiry-indicator {
		position: absolute;
		top: -2px;
		right: -2px;
		width: 8px;
		height: 8px;
		background: #ff4444;
		border-radius: 50%;
		color: white;
		font-size: 0.6rem;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: bold;
		border: 1px solid white;
	}

	.start-indicator {
		position: absolute;
		top: -2px;
		left: -2px;
		width: 8px;
		height: 8px;
		background: #22c55e;
		border-radius: 50%;
		color: white;
		font-size: 0.8rem;
		display: flex;
		align-items: center;
		justify-content: center;
		border: 1px solid white;
	}

	.trade-overflow {
		background: rgba(255, 255, 255, 0.1);
		color: #ccc;
		padding: 2px 4px;
		border-radius: 3px;
		font-size: 0.6rem;
		text-align: center;
		border: 1px dashed rgba(255, 255, 255, 0.3);
		margin-top: 2px;
	}

	.empty-cell {
		height: 100%;
		min-height: 52px;
		background: transparent;
		border-radius: 3px;
		transition: background-color 0.2s ease;
	}

	.trade-cell:hover .empty-cell {
		background: rgba(74, 144, 226, 0.1);
		border: 1px dashed rgba(74, 144, 226, 0.3);
	}

	/* Responsive adjustments */
	@media (max-width: 1200px) {
		.trade-item {
			font-size: 0.6rem;
			padding: 1px 2px;
		}
		
		.trade-ticker {
			font-size: 0.6rem;
		}
	}
</style> 