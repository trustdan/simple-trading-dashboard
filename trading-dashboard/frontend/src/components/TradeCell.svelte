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
	
	// Tooltip state
	let tooltipVisible = false;
	let tooltipX = 0;
	let tooltipY = 0;
	let tooltipTrade = null;

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
	
	// Tooltip functions using immediate DOM manipulation
	function showTooltip(trade, event) {
		// Clean up any existing tooltip first
		const existingTooltip = document.querySelector('.trade-tooltip-overlay');
		if (existingTooltip) {
			existingTooltip.remove();
		}
		
		// Calculate position with screen boundaries in mind
		const tooltipWidth = 280;
		const tooltipHeight = 250;
		const padding = 15;
		
		let x = Number(event.clientX) + padding;
		let y = Number(event.clientY) - padding;
		
		// Check right edge
		if (x + tooltipWidth > window.innerWidth) {
			x = Number(event.clientX) - tooltipWidth - padding;
		}
		
		// Check bottom edge  
		if (y + tooltipHeight > window.innerHeight) {
			y = Number(event.clientY) - tooltipHeight + padding;
		}
		
		// Check top edge
		if (y < 0) {
			y = padding;
		}
		
		// Check left edge
		if (x < 0) {
			x = padding;
		}
		
		// Create tooltip immediately
		const tooltip = document.createElement('div');
		tooltip.className = 'trade-tooltip-overlay';
		tooltip.style.cssText = `
			position: fixed !important;
			left: ${x}px !important;
			top: ${y}px !important;
			z-index: 2147483647 !important;
			background: rgba(15, 15, 15, 0.98) !important;
			border: 2px solid #4a90e2 !important;
			border-radius: 12px !important;
			padding: 16px !important;
			color: #ffffff !important;
			font-size: 0.9rem !important;
			font-family: system-ui, -apple-system, sans-serif !important;
			box-shadow: 0 12px 48px rgba(0, 0, 0, 0.9), 0 0 0 1px rgba(255, 255, 255, 0.1) !important;
			max-width: 280px !important;
			min-width: 250px !important;
			pointer-events: none !important;
			backdrop-filter: blur(20px) !important;
			transform: translateZ(0) !important;
			will-change: transform !important;
		`;
		
		tooltip.innerHTML = `
			<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; padding-bottom: 8px; border-bottom: 1px solid rgba(255, 255, 255, 0.2);">
				<span style="font-weight: 700; font-size: 1.2rem; color: #4a90e2; font-family: monospace;">${trade.ticker}</span>
				<span style="padding: 3px 8px; border-radius: 12px; font-size: 0.7rem; font-weight: 600; text-transform: uppercase; ${getStatusStyle(trade.status)}">${trade.status}</span>
			</div>
			<div style="font-weight: 600; color: #7b68ee; margin-bottom: 8px; font-size: 1rem;">${trade.strategy_type}</div>
			<div style="color: #cccccc; font-size: 0.85rem; margin-bottom: 10px;">${trade.sector}</div>
			<div style="margin-bottom: 10px;">
				<div style="display: flex; justify-content: space-between; margin-bottom: 3px;">
					<span style="color: #aaaaaa; font-size: 0.8rem;">Entry:</span> 
					<span style="color: #ffffff;">${formatDate(trade.entry_date)}</span>
				</div>
				<div style="display: flex; justify-content: space-between; margin-bottom: 3px;">
					<span style="color: #aaaaaa; font-size: 0.8rem;">Expiry:</span> 
					<span style="color: #ffffff;">${formatDate(trade.expiration_date)}</span>
				</div>
			</div>
			${trade.target_price || trade.stop_loss ? `
				<div style="margin-bottom: 10px;">
					${trade.target_price ? `
						<div style="display: flex; justify-content: space-between; margin-bottom: 3px;">
							<span style="color: #aaaaaa; font-size: 0.8rem;">Target:</span> 
							<span style="color: #22c55e; font-weight: 600;">${formatCurrency(trade.target_price)}</span>
						</div>
					` : ''}
					${trade.stop_loss ? `
						<div style="display: flex; justify-content: space-between; margin-bottom: 3px;">
							<span style="color: #aaaaaa; font-size: 0.8rem;">Stop:</span> 
							<span style="color: #ef4444; font-weight: 600;">${formatCurrency(trade.stop_loss)}</span>
						</div>
					` : ''}
				</div>
			` : ''}
			${trade.notes ? `<div style="background: rgba(255, 255, 255, 0.08); padding: 8px; border-radius: 6px; font-size: 0.8rem; color: #cccccc; margin-bottom: 10px; font-style: italic; border-left: 3px solid #4a90e2;">${trade.notes}</div>` : ''}
			<div style="text-align: center; color: #888888; font-size: 0.75rem; border-top: 1px solid rgba(255, 255, 255, 0.1); padding-top: 8px; font-style: italic;">💡 Click to edit</div>
		`;
		
		document.body.appendChild(tooltip);
		
		// Store reference for cleanup
		tooltipTrade = trade;
		tooltipVisible = true;
	}
	
	function hideTooltip() {
		tooltipVisible = false;
		tooltipTrade = null;
		
		// Clean up tooltip
		const existingTooltip = document.querySelector('.trade-tooltip-overlay');
		if (existingTooltip) {
			existingTooltip.remove();
		}
	}
	
	function getStatusStyle(status) {
		switch(status) {
			case 'active':
				return 'background: rgba(34, 197, 94, 0.2); color: #22c55e; border: 1px solid #22c55e;';
			case 'closed':
				return 'background: rgba(156, 163, 175, 0.2); color: #9ca3af; border: 1px solid #9ca3af;';
			case 'expired':
				return 'background: rgba(239, 68, 68, 0.2); color: #ef4444; border: 1px solid #ef4444;';
			default:
				return 'background: rgba(156, 163, 175, 0.2); color: #9ca3af; border: 1px solid #9ca3af;';
		}
	}
	
	function formatCurrency(value) {
		return value ? `$${value.toFixed(2)}` : '-';
	}
	
	function formatDate(dateString) {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric'
		});
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
							on:mouseenter={(e) => showTooltip(trade, e)}
							on:mouseleave={hideTooltip}
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



<!-- Context Menu -->
<TradeContextMenu 
	bind:visible={contextMenuVisible}
	trade={contextMenuTrade}
	x={contextMenuX}
	y={contextMenuY}
	on:close={closeContextMenu}
	on:edit={handleContextMenuEdit}
	on:status-change={handleContextMenuStatusChange}
	on:delete={handleContextMenuDelete}
	on:duplicate={handleContextMenuDuplicate}
/> 