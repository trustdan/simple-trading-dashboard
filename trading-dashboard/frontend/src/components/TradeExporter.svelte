<script>
	import { createEventDispatcher } from 'svelte';
	import { tradesStore } from '../stores/trades.js';
	import { toastStore } from '../stores/toast.js';

	const dispatch = createEventDispatcher();

	// Get data from store
	$: trades = $tradesStore.trades || [];
	$: strategyTypes = $tradesStore.strategyTypes || [];

	let isExporting = false;
	let exportFormat = 'json'; // 'json', 'csv', 'excel'
	let exportRange = 'all'; // 'all', 'active', 'custom'
	let customStartDate = '';
	let customEndDate = '';

	// Export formats
	const formats = [
		{ value: 'json', label: 'JSON', description: 'Machine-readable format' },
		{ value: 'csv', label: 'CSV', description: 'Spreadsheet compatible' },
		{ value: 'excel', label: 'Excel', description: 'Microsoft Excel format' }
	];

	// Export ranges
	const ranges = [
		{ value: 'all', label: 'All Trades', description: 'Export all trades' },
		{ value: 'active', label: 'Active Only', description: 'Only active trades' },
		{ value: 'custom', label: 'Date Range', description: 'Custom date range' }
	];

	function getFilteredTrades() {
		let filteredTrades = [...trades];

		switch (exportRange) {
			case 'active':
				filteredTrades = filteredTrades.filter(trade => trade.status === 'active');
				break;
			case 'custom':
				if (customStartDate && customEndDate) {
					const startDate = new Date(customStartDate);
					const endDate = new Date(customEndDate);
					filteredTrades = filteredTrades.filter(trade => {
						const tradeDate = new Date(trade.entry_date);
						return tradeDate >= startDate && tradeDate <= endDate;
					});
				}
				break;
		}

		return filteredTrades;
	}

	function generateAnalytics(trades) {
		const analytics = {
			summary: {
				totalTrades: trades.length,
				activeTrades: trades.filter(t => t.status === 'active').length,
				closedTrades: trades.filter(t => t.status === 'closed').length,
				expiredTrades: trades.filter(t => t.status === 'expired').length
			},
			strategies: {},
			sectors: {},
			monthlyBreakdown: {},
			exportMetadata: {
				exportDate: new Date().toISOString(),
				exportRange,
				totalRecords: trades.length
			}
		};

		// Strategy breakdown
		trades.forEach(trade => {
			analytics.strategies[trade.strategy_type] = (analytics.strategies[trade.strategy_type] || 0) + 1;
			analytics.sectors[trade.sector] = (analytics.sectors[trade.sector] || 0) + 1;

			// Monthly breakdown
			const month = new Date(trade.entry_date).toISOString().slice(0, 7); // YYYY-MM
			analytics.monthlyBreakdown[month] = (analytics.monthlyBreakdown[month] || 0) + 1;
		});

		return analytics;
	}

	function exportAsJSON() {
		const filteredTrades = getFilteredTrades();
		const analytics = generateAnalytics(filteredTrades);
		
		const exportData = {
			trades: filteredTrades,
			analytics,
			strategyTypes,
			exportInfo: {
				format: 'json',
				timestamp: new Date().toISOString(),
				version: '1.0'
			}
		};

		const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' });
		downloadFile(blob, `trades-export-${getDateString()}.json`);
	}

	function exportAsCSV() {
		const filteredTrades = getFilteredTrades();
		
		// CSV headers
		const headers = [
			'ID', 'Ticker', 'Sector', 'Strategy Type', 'Entry Date', 'Expiration Date',
			'Target Price', 'Stop Loss', 'Status', 'Notes', 'Created At', 'Updated At'
		];

		// CSV rows
		const rows = filteredTrades.map(trade => [
			trade.id,
			trade.ticker,
			trade.sector,
			trade.strategy_type,
			formatDateForCSV(trade.entry_date),
			formatDateForCSV(trade.expiration_date),
			trade.target_price || '',
			trade.stop_loss || '',
			trade.status,
			`"${(trade.notes || '').replace(/"/g, '""')}"`, // Escape quotes
			formatDateForCSV(trade.created_at),
			formatDateForCSV(trade.updated_at)
		]);

		// Combine headers and rows
		const csvContent = [headers, ...rows]
			.map(row => row.join(','))
			.join('\n');

		const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
		downloadFile(blob, `trades-export-${getDateString()}.csv`);
	}

	function exportAsExcel() {
		// For Excel export, we'll create a more structured CSV that Excel can handle better
		const filteredTrades = getFilteredTrades();
		const analytics = generateAnalytics(filteredTrades);

		// Create multiple sheets as CSV sections
		let excelContent = '';

		// Trades sheet
		excelContent += 'TRADES\n';
		excelContent += 'ID,Ticker,Sector,Strategy Type,Entry Date,Expiration Date,Target Price,Stop Loss,Status,Notes,Created At,Updated At\n';
		
		filteredTrades.forEach(trade => {
			excelContent += [
				trade.id,
				trade.ticker,
				trade.sector,
				trade.strategy_type,
				formatDateForCSV(trade.entry_date),
				formatDateForCSV(trade.expiration_date),
				trade.target_price || '',
				trade.stop_loss || '',
				trade.status,
				`"${(trade.notes || '').replace(/"/g, '""')}"`,
				formatDateForCSV(trade.created_at),
				formatDateForCSV(trade.updated_at)
			].join(',') + '\n';
		});

		// Analytics sheet
		excelContent += '\n\nANALYTICS SUMMARY\n';
		excelContent += 'Metric,Value\n';
		excelContent += `Total Trades,${analytics.summary.totalTrades}\n`;
		excelContent += `Active Trades,${analytics.summary.activeTrades}\n`;
		excelContent += `Closed Trades,${analytics.summary.closedTrades}\n`;
		excelContent += `Expired Trades,${analytics.summary.expiredTrades}\n`;

		// Strategy breakdown
		excelContent += '\n\nSTRATEGY BREAKDOWN\n';
		excelContent += 'Strategy,Count\n';
		Object.entries(analytics.strategies).forEach(([strategy, count]) => {
			excelContent += `${strategy},${count}\n`;
		});

		// Sector breakdown
		excelContent += '\n\nSECTOR BREAKDOWN\n';
		excelContent += 'Sector,Count\n';
		Object.entries(analytics.sectors).forEach(([sector, count]) => {
			excelContent += `${sector},${count}\n`;
		});

		const blob = new Blob([excelContent], { type: 'text/csv;charset=utf-8;' });
		downloadFile(blob, `trades-export-${getDateString()}.csv`);
	}

	function formatDateForCSV(dateStr) {
		if (!dateStr) return '';
		return new Date(dateStr).toISOString().split('T')[0];
	}

	function getDateString() {
		return new Date().toISOString().split('T')[0];
	}

	function downloadFile(blob, filename) {
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = filename;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
	}

	async function handleExport() {
		if (isExporting) return;

		const filteredTrades = getFilteredTrades();
		if (filteredTrades.length === 0) {
			toastStore.error('No trades to export with current filters');
			return;
		}

		isExporting = true;

		try {
			switch (exportFormat) {
				case 'json':
					exportAsJSON();
					break;
				case 'csv':
					exportAsCSV();
					break;
				case 'excel':
					exportAsExcel();
					break;
			}

			toastStore.success(`Successfully exported ${filteredTrades.length} trades as ${exportFormat.toUpperCase()}`);
		} catch (error) {
			console.error('Export error:', error);
			toastStore.error(`Failed to export: ${error.message}`);
		} finally {
			isExporting = false;
		}
	}

	// Validate custom date range
	$: isCustomRangeValid = exportRange !== 'custom' || (customStartDate && customEndDate && customStartDate <= customEndDate);
	$: filteredCount = getFilteredTrades().length;
</script>

<div class="trade-exporter">
	<div class="exporter-header">
		<h3>📤 Export Trades</h3>
		<div class="trade-count">
			{filteredCount} trade{filteredCount !== 1 ? 's' : ''} selected
		</div>
	</div>

	<div class="export-options">
		<!-- Format Selection -->
		<div class="option-group">
			<label class="group-label">Export Format</label>
			<div class="format-grid">
				{#each formats as format}
					<label class="format-option">
						<input 
							type="radio" 
							bind:group={exportFormat} 
							value={format.value}
							disabled={isExporting}
						/>
						<div class="format-card" class:selected={exportFormat === format.value}>
							<div class="format-name">{format.label}</div>
							<div class="format-desc">{format.description}</div>
						</div>
					</label>
				{/each}
			</div>
		</div>

		<!-- Range Selection -->
		<div class="option-group">
			<label class="group-label">Export Range</label>
			<div class="range-options">
				{#each ranges as range}
					<label class="range-option">
						<input 
							type="radio" 
							bind:group={exportRange} 
							value={range.value}
							disabled={isExporting}
						/>
						<div class="range-content">
							<span class="range-name">{range.label}</span>
							<span class="range-desc">{range.description}</span>
						</div>
					</label>
				{/each}
			</div>
		</div>

		<!-- Custom Date Range -->
		{#if exportRange === 'custom'}
			<div class="option-group">
				<label class="group-label">Date Range</label>
				<div class="date-range">
					<div class="date-input">
						<label for="start-date">Start Date</label>
						<input 
							id="start-date"
							type="date" 
							bind:value={customStartDate}
							disabled={isExporting}
						/>
					</div>
					<div class="date-separator">to</div>
					<div class="date-input">
						<label for="end-date">End Date</label>
						<input 
							id="end-date"
							type="date" 
							bind:value={customEndDate}
							disabled={isExporting}
						/>
					</div>
				</div>
			</div>
		{/if}
	</div>

	<!-- Export Button -->
	<div class="export-actions">
		<button 
			class="export-btn"
			class:loading={isExporting}
			disabled={isExporting || !isCustomRangeValid || filteredCount === 0}
			on:click={handleExport}
		>
			{#if isExporting}
				<div class="loading-spinner"></div>
				Exporting...
			{:else}
				📤 Export {exportFormat.toUpperCase()}
			{/if}
		</button>
	</div>
</div>

<style>
	.trade-exporter {
		background: #1a1a1a;
		border-radius: 12px;
		padding: 24px;
		margin-bottom: 24px;
	}

	.exporter-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 24px;
	}

	.exporter-header h3 {
		margin: 0;
		color: #ffffff;
		font-size: 1.25rem;
		font-weight: 600;
	}

	.trade-count {
		background: #2a2a2a;
		padding: 6px 12px;
		border-radius: 6px;
		font-size: 14px;
		color: #cccccc;
		font-weight: 500;
	}

	.export-options {
		display: flex;
		flex-direction: column;
		gap: 24px;
		margin-bottom: 24px;
	}

	.option-group {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.group-label {
		font-size: 16px;
		font-weight: 600;
		color: #ffffff;
	}

	.format-grid {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 12px;
	}

	.format-option {
		cursor: pointer;
	}

	.format-option input {
		display: none;
	}

	.format-card {
		background: #2a2a2a;
		border: 2px solid #444;
		border-radius: 8px;
		padding: 16px;
		text-align: center;
		transition: all 0.2s ease;
		cursor: pointer;
	}

	.format-card:hover {
		border-color: #666;
		transform: translateY(-1px);
	}

	.format-card.selected {
		border-color: #4a90e2;
		background: rgba(74, 144, 226, 0.1);
	}

	.format-name {
		font-size: 16px;
		font-weight: 600;
		color: #ffffff;
		margin-bottom: 4px;
	}

	.format-desc {
		font-size: 12px;
		color: #999;
	}

	.range-options {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.range-option {
		display: flex;
		align-items: center;
		gap: 12px;
		cursor: pointer;
		padding: 12px;
		background: #2a2a2a;
		border-radius: 6px;
		transition: background-color 0.2s ease;
	}

	.range-option:hover {
		background: #333;
	}

	.range-option input {
		accent-color: #4a90e2;
	}

	.range-content {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.range-name {
		font-size: 14px;
		font-weight: 600;
		color: #ffffff;
	}

	.range-desc {
		font-size: 12px;
		color: #999;
	}

	.date-range {
		display: flex;
		align-items: end;
		gap: 16px;
	}

	.date-input {
		display: flex;
		flex-direction: column;
		gap: 4px;
		flex: 1;
	}

	.date-input label {
		font-size: 12px;
		color: #cccccc;
		font-weight: 500;
	}

	.date-input input {
		padding: 8px 12px;
		background: #2a2a2a;
		border: 2px solid #444;
		border-radius: 6px;
		color: #ffffff;
		font-size: 14px;
	}

	.date-input input:focus {
		outline: none;
		border-color: #4a90e2;
	}

	.date-separator {
		padding: 8px 0;
		color: #666;
		font-weight: 500;
	}

	.export-actions {
		display: flex;
		justify-content: center;
	}

	.export-btn {
		background: linear-gradient(135deg, #4a90e2, #7b68ee);
		color: white;
		border: none;
		padding: 12px 32px;
		border-radius: 8px;
		cursor: pointer;
		font-size: 16px;
		font-weight: 600;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		gap: 8px;
		min-width: 160px;
		justify-content: center;
	}

	.export-btn:hover:not(:disabled) {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(74, 144, 226, 0.3);
	}

	.export-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		transform: none;
		box-shadow: none;
	}

	.export-btn.loading {
		background: #666;
	}

	.loading-spinner {
		width: 16px;
		height: 16px;
		border: 2px solid transparent;
		border-top: 2px solid currentColor;
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	/* Responsive design */
	@media (max-width: 768px) {
		.format-grid {
			grid-template-columns: 1fr;
		}

		.date-range {
			flex-direction: column;
			align-items: stretch;
		}

		.date-separator {
			text-align: center;
		}
	}
</style> 