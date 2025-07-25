<script>
	import { createEventDispatcher } from 'svelte';
	
	export let trades = [];
	export let dateColumns = [];
	export let sectors = [];

	const dispatch = createEventDispatcher();

	// Heat map modes
	let heatMapMode = 'density'; // 'density', 'expiration', 'activity'

	// Calculate heat map data
	$: heatMapData = calculateHeatMapData(trades, dateColumns, sectors, heatMapMode);

	function calculateHeatMapData(trades, dateColumns, sectors, mode) {
		if (!trades || !dateColumns || !sectors) return {};

		const data = {};
		
		// Initialize data structure
		sectors.forEach(sector => {
			data[sector] = {};
			dateColumns.forEach(date => {
				const dateKey = date.toISOString().split('T')[0];
				data[sector][dateKey] = {
					count: 0,
					intensity: 0,
					trades: []
				};
			});
		});

		// Populate data based on mode
		trades.forEach(trade => {
			const entryDate = new Date(trade.entry_date);
			const expirationDate = new Date(trade.expiration_date);

			dateColumns.forEach(date => {
				const dateKey = date.toISOString().split('T')[0];
				const cellDate = new Date(date);

				// Check if trade is active on this date
				if (cellDate >= entryDate && cellDate <= expirationDate && trade.sector) {
					const cell = data[trade.sector]?.[dateKey];
					if (cell) {
						cell.count++;
						cell.trades.push(trade);

						// Calculate intensity based on mode
						switch (mode) {
							case 'density':
								cell.intensity = cell.count;
								break;
							case 'expiration':
								const daysToExpiry = Math.ceil((expirationDate - cellDate) / (1000 * 60 * 60 * 24));
								if (daysToExpiry <= 7 && daysToExpiry >= 0) {
									cell.intensity += 10; // High intensity for expiring trades
								} else {
									cell.intensity += 1;
								}
								break;
							case 'activity':
								const createdDate = new Date(trade.created_at);
								const daysSinceCreated = Math.ceil((cellDate - createdDate) / (1000 * 60 * 60 * 24));
								if (daysSinceCreated <= 3) {
									cell.intensity += 5; // High intensity for new trades
								} else {
									cell.intensity += 1;
								}
								break;
						}
					}
				}
			});
		});

		return data;
	}

	function getMaxIntensity() {
		let max = 0;
		Object.values(heatMapData).forEach(sectorData => {
			Object.values(sectorData).forEach(cellData => {
				max = Math.max(max, cellData.intensity);
			});
		});
		return max || 1;
	}

	function getCellIntensity(sector, date) {
		const dateKey = date.toISOString().split('T')[0];
		const cellData = heatMapData[sector]?.[dateKey];
		if (!cellData || cellData.intensity === 0) return 0;
		
		return cellData.intensity / getMaxIntensity();
	}

	function getCellColor(sector, date) {
		const intensity = getCellIntensity(sector, date);
		if (intensity === 0) return 'rgba(255, 255, 255, 0.02)';

		switch (heatMapMode) {
			case 'density':
				return `rgba(74, 144, 226, ${0.2 + intensity * 0.8})`; // Blue gradient
			case 'expiration':
				return `rgba(245, 158, 11, ${0.2 + intensity * 0.8})`; // Orange gradient
			case 'activity':
				return `rgba(34, 197, 94, ${0.2 + intensity * 0.8})`; // Green gradient
			default:
				return `rgba(74, 144, 226, ${0.2 + intensity * 0.8})`;
		}
	}

	function getCellTooltip(sector, date) {
		const dateKey = date.toISOString().split('T')[0];
		const cellData = heatMapData[sector]?.[dateKey];
		if (!cellData || cellData.count === 0) return `${sector} - ${dateKey}: No trades`;

		const modeText = {
			density: 'trades',
			expiration: 'expiration pressure',
			activity: 'activity level'
		}[heatMapMode];

		return `${sector} - ${dateKey}\n${cellData.count} trades\n${modeText}: ${cellData.intensity}`;
	}

	function handleCellClick(sector, date) {
		const dateKey = date.toISOString().split('T')[0];
		const cellData = heatMapData[sector]?.[dateKey];
		if (cellData && cellData.trades.length > 0) {
			dispatch('cell-click', {
				sector,
				date,
				trades: cellData.trades
			});
		}
	}

	function formatDate(date) {
		return {
			dayName: date.toLocaleDateString('en-US', { weekday: 'short' }),
			dayNumber: date.getDate(),
			month: date.toLocaleDateString('en-US', { month: 'short' }),
			isToday: date.toDateString() === new Date().toDateString(),
			isWeekend: date.getDay() === 0 || date.getDay() === 6
		};
	}
</script>

<div class="trade-heatmap">
	<div class="heatmap-header">
		<h3>🔥 Trade Heat Map</h3>
		<div class="mode-selector">
			<button 
				class="mode-btn" 
				class:active={heatMapMode === 'density'}
				on:click={() => heatMapMode = 'density'}
			>
				Density
			</button>
			<button 
				class="mode-btn" 
				class:active={heatMapMode === 'expiration'}
				on:click={() => heatMapMode = 'expiration'}
			>
				Expiration
			</button>
			<button 
				class="mode-btn" 
				class:active={heatMapMode === 'activity'}
				on:click={() => heatMapMode = 'activity'}
			>
				Activity
			</button>
		</div>
	</div>

	<div class="heatmap-legend">
		<span class="legend-label">
			{#if heatMapMode === 'density'}
				Trade Count
			{:else if heatMapMode === 'expiration'}
				Expiration Pressure
			{:else}
				Recent Activity
			{/if}
		</span>
		<div class="legend-gradient" style="background: linear-gradient(90deg, 
			{heatMapMode === 'density' ? 'rgba(74, 144, 226, 0.2), rgba(74, 144, 226, 1)' : 
			 heatMapMode === 'expiration' ? 'rgba(245, 158, 11, 0.2), rgba(245, 158, 11, 1)' : 
			 'rgba(34, 197, 94, 0.2), rgba(34, 197, 94, 1)'})"></div>
		<span class="legend-label">High</span>
	</div>

	<div class="heatmap-grid">
		<!-- Header row with dates -->
		<div class="grid-header">
			<div class="sector-header">Sector</div>
			{#each dateColumns as date}
				{@const formatted = formatDate(date)}
				<div class="date-header" class:today={formatted.isToday} class:weekend={formatted.isWeekend}>
					<div class="day-name">{formatted.dayName}</div>
					<div class="day-number">{formatted.dayNumber}</div>
					<div class="month">{formatted.month}</div>
				</div>
			{/each}
		</div>

		<!-- Grid body with sector rows -->
		<div class="grid-body">
			{#each sectors as sector}
				<div class="sector-row">
					<div class="sector-label">
						{sector}
					</div>
					
					{#each dateColumns as date}
						<div 
							class="heat-cell"
							style="background-color: {getCellColor(sector, date)}"
							title={getCellTooltip(sector, date)}
							on:click={() => handleCellClick(sector, date)}
						>
							{#if getCellIntensity(sector, date) > 0}
								<div class="cell-indicator">
									{heatMapData[sector]?.[date.toISOString().split('T')[0]]?.count || ''}
								</div>
							{/if}
						</div>
					{/each}
				</div>
			{/each}
		</div>
	</div>
</div>

<style>
	.trade-heatmap {
		background: #1a1a1a;
		border-radius: 12px;
		padding: 20px;
		margin-bottom: 24px;
	}

	.heatmap-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 16px;
	}

	.heatmap-header h3 {
		margin: 0;
		color: #ffffff;
		font-size: 1.25rem;
		font-weight: 600;
	}

	.mode-selector {
		display: flex;
		background: #2a2a2a;
		border-radius: 6px;
		padding: 2px;
	}

	.mode-btn {
		background: none;
		border: none;
		color: #cccccc;
		padding: 6px 12px;
		border-radius: 4px;
		cursor: pointer;
		font-size: 12px;
		font-weight: 500;
		transition: all 0.2s ease;
	}

	.mode-btn:hover {
		color: #ffffff;
		background: rgba(255, 255, 255, 0.05);
	}

	.mode-btn.active {
		background: #4a90e2;
		color: white;
	}

	.heatmap-legend {
		display: flex;
		align-items: center;
		gap: 8px;
		margin-bottom: 16px;
		font-size: 12px;
		color: #999;
	}

	.legend-gradient {
		width: 100px;
		height: 12px;
		border-radius: 6px;
	}

	.heatmap-grid {
		border: 1px solid #333;
		border-radius: 8px;
		overflow: hidden;
	}

	.grid-header {
		display: grid;
		grid-template-columns: 140px repeat(21, 1fr);
		background: #2a2a2a;
		border-bottom: 1px solid #333;
	}

	.sector-header {
		padding: 12px 8px;
		font-weight: 600;
		color: #ffffff;
		font-size: 14px;
		border-right: 1px solid #333;
		background: #333;
	}

	.date-header {
		padding: 8px 4px;
		text-align: center;
		font-size: 11px;
		color: #cccccc;
		border-right: 1px solid rgba(255, 255, 255, 0.1);
		min-width: 0;
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
		font-weight: 500;
	}

	.day-number {
		font-size: 13px;
		font-weight: 600;
		margin: 2px 0;
	}

	.month {
		font-size: 10px;
		opacity: 0.7;
	}

	.grid-body {
		background: #1a1a1a;
	}

	.sector-row {
		display: grid;
		grid-template-columns: 140px repeat(21, 1fr);
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
	}

	.sector-row:last-child {
		border-bottom: none;
	}

	.sector-label {
		padding: 12px 8px;
		font-weight: 500;
		color: #cccccc;
		font-size: 13px;
		border-right: 1px solid #333;
		background: #2a2a2a;
		display: flex;
		align-items: center;
	}

	.heat-cell {
		min-height: 40px;
		border-right: 1px solid rgba(255, 255, 255, 0.05);
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		justify-content: center;
		position: relative;
	}

	.heat-cell:hover {
		transform: scale(1.05);
		z-index: 10;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-radius: 4px;
	}

	.cell-indicator {
		font-size: 11px;
		font-weight: 600;
		color: white;
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
	}

	/* Responsive design */
	@media (max-width: 1200px) {
		.grid-header,
		.sector-row {
			grid-template-columns: 120px repeat(21, 1fr);
		}

		.date-header,
		.heat-cell {
			min-width: 30px;
		}
	}

	@media (max-width: 768px) {
		.heatmap-header {
			flex-direction: column;
			gap: 12px;
			align-items: stretch;
		}

		.mode-selector {
			justify-content: center;
		}

		.grid-header,
		.sector-row {
			grid-template-columns: 100px repeat(21, 1fr);
		}

		.sector-label {
			font-size: 11px;
			padding: 8px 4px;
		}

		.date-header {
			padding: 6px 2px;
			font-size: 10px;
		}

		.heat-cell {
			min-height: 30px;
		}
	}
</style> 