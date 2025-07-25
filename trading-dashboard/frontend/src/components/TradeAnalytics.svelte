<script>
	import { createEventDispatcher } from 'svelte';
	import { tradesStore } from '../stores/trades.js';

	const dispatch = createEventDispatcher();

	// Get trades from store
	$: trades = $tradesStore.trades || [];

	// Calculate analytics
	$: analytics = calculateAnalytics(trades);

	function calculateAnalytics(trades) {
		if (!trades || trades.length === 0) {
			return {
				totalTrades: 0,
				activeTrades: 0,
				closedTrades: 0,
				expiredTrades: 0,
				strategyCounts: {},
				sectorCounts: {},
				expiringThisWeek: 0,
				avgDaysToExpiry: 0,
				recentActivity: []
			};
		}

		const now = new Date();
		const oneWeekFromNow = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000);

		const analytics = {
			totalTrades: trades.length,
			activeTrades: trades.filter(t => t.status === 'active').length,
			closedTrades: trades.filter(t => t.status === 'closed').length,
			expiredTrades: trades.filter(t => t.status === 'expired').length,
			strategyCounts: {},
			sectorCounts: {},
			expiringThisWeek: 0,
			avgDaysToExpiry: 0,
			recentActivity: []
		};

		// Strategy and sector counts
		trades.forEach(trade => {
			analytics.strategyCounts[trade.strategy_type] = (analytics.strategyCounts[trade.strategy_type] || 0) + 1;
			analytics.sectorCounts[trade.sector] = (analytics.sectorCounts[trade.sector] || 0) + 1;

			// Check if expiring this week
			const expirationDate = new Date(trade.expiration_date);
			if (expirationDate <= oneWeekFromNow && expirationDate >= now && trade.status === 'active') {
				analytics.expiringThisWeek++;
			}
		});

		// Average days to expiry for active trades
		const activeTrades = trades.filter(t => t.status === 'active');
		if (activeTrades.length > 0) {
			const totalDays = activeTrades.reduce((sum, trade) => {
				const expirationDate = new Date(trade.expiration_date);
				const daysToExpiry = Math.ceil((expirationDate - now) / (1000 * 60 * 60 * 24));
				return sum + Math.max(0, daysToExpiry);
			}, 0);
			analytics.avgDaysToExpiry = Math.round(totalDays / activeTrades.length);
		}

		// Recent activity (last 7 days)
		const oneWeekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
		analytics.recentActivity = trades
			.filter(trade => new Date(trade.created_at) >= oneWeekAgo)
			.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
			.slice(0, 5);

		return analytics;
	}

	function getTopStrategies() {
		return Object.entries(analytics.strategyCounts)
			.sort(([,a], [,b]) => b - a)
			.slice(0, 5);
	}

	function getTopSectors() {
		return Object.entries(analytics.sectorCounts)
			.sort(([,a], [,b]) => b - a)
			.slice(0, 5);
	}

	function getStatusPercentage(status) {
		if (analytics.totalTrades === 0) return 0;
		const count = analytics[`${status}Trades`];
		return Math.round((count / analytics.totalTrades) * 100);
	}

	function formatDate(dateStr) {
		return new Date(dateStr).toLocaleDateString('en-US', { 
			month: 'short', 
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function exportAnalytics() {
		const data = {
			summary: analytics,
			topStrategies: getTopStrategies(),
			topSectors: getTopSectors(),
			exportDate: new Date().toISOString()
		};
		
		const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `trade-analytics-${new Date().toISOString().split('T')[0]}.json`;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
	}
</script>

<div class="trade-analytics">
	<div class="analytics-header">
		<h2>📊 Trade Analytics</h2>
		<button class="export-btn" on:click={exportAnalytics}>
			📤 Export
		</button>
	</div>

	<!-- Summary Cards -->
	<div class="summary-grid">
		<div class="summary-card total">
			<div class="card-icon">📈</div>
			<div class="card-content">
				<div class="card-value">{analytics.totalTrades}</div>
				<div class="card-label">Total Trades</div>
			</div>
		</div>

		<div class="summary-card active">
			<div class="card-icon">🟢</div>
			<div class="card-content">
				<div class="card-value">{analytics.activeTrades}</div>
				<div class="card-label">Active</div>
				<div class="card-percentage">{getStatusPercentage('active')}%</div>
			</div>
		</div>

		<div class="summary-card closed">
			<div class="card-icon">✅</div>
			<div class="card-content">
				<div class="card-value">{analytics.closedTrades}</div>
				<div class="card-label">Closed</div>
				<div class="card-percentage">{getStatusPercentage('closed')}%</div>
			</div>
		</div>

		<div class="summary-card expired">
			<div class="card-icon">⚫</div>
			<div class="card-content">
				<div class="card-value">{analytics.expiredTrades}</div>
				<div class="card-label">Expired</div>
				<div class="card-percentage">{getStatusPercentage('expired')}%</div>
			</div>
		</div>

		<div class="summary-card warning">
			<div class="card-icon">⚠️</div>
			<div class="card-content">
				<div class="card-value">{analytics.expiringThisWeek}</div>
				<div class="card-label">Expiring This Week</div>
			</div>
		</div>

		<div class="summary-card info">
			<div class="card-icon">📅</div>
			<div class="card-content">
				<div class="card-value">{analytics.avgDaysToExpiry}</div>
				<div class="card-label">Avg Days to Expiry</div>
			</div>
		</div>
	</div>

	<!-- Charts Section -->
	<div class="charts-section">
		<!-- Top Strategies -->
		<div class="chart-card">
			<h3>Top Strategies</h3>
			<div class="strategy-list">
				{#each getTopStrategies() as [strategy, count]}
					<div class="strategy-item">
						<div class="strategy-name">{strategy}</div>
						<div class="strategy-bar">
							<div 
								class="strategy-fill" 
								style="width: {(count / analytics.totalTrades) * 100}%"
							></div>
						</div>
						<div class="strategy-count">{count}</div>
					</div>
				{/each}
			</div>
		</div>

		<!-- Top Sectors -->
		<div class="chart-card">
			<h3>Top Sectors</h3>
			<div class="sector-list">
				{#each getTopSectors() as [sector, count]}
					<div class="sector-item">
						<div class="sector-name">{sector}</div>
						<div class="sector-bar">
							<div 
								class="sector-fill" 
								style="width: {(count / analytics.totalTrades) * 100}%"
							></div>
						</div>
						<div class="sector-count">{count}</div>
					</div>
				{/each}
			</div>
		</div>
	</div>

	<!-- Recent Activity -->
	<div class="activity-section">
		<h3>Recent Activity</h3>
		<div class="activity-list">
			{#each analytics.recentActivity as trade}
				<div class="activity-item">
					<div class="activity-ticker">{trade.ticker}</div>
					<div class="activity-strategy">{trade.strategy_type}</div>
					<div class="activity-sector">{trade.sector}</div>
					<div class="activity-date">{formatDate(trade.created_at)}</div>
				</div>
			{:else}
				<div class="empty-activity">
					<p>No recent activity in the last 7 days</p>
				</div>
			{/each}
		</div>
	</div>
</div>

<style>
	.trade-analytics {
		background: #1a1a1a;
		border-radius: 12px;
		padding: 24px;
		margin-bottom: 24px;
	}

	.analytics-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 24px;
	}

	.analytics-header h2 {
		margin: 0;
		color: #ffffff;
		font-size: 1.5rem;
		font-weight: 600;
	}

	.export-btn {
		background: linear-gradient(135deg, #4a90e2, #7b68ee);
		color: white;
		border: none;
		padding: 8px 16px;
		border-radius: 6px;
		cursor: pointer;
		font-size: 14px;
		font-weight: 500;
		transition: all 0.2s ease;
	}

	.export-btn:hover {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(74, 144, 226, 0.3);
	}

	.summary-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
		gap: 16px;
		margin-bottom: 32px;
	}

	.summary-card {
		background: #2a2a2a;
		border-radius: 8px;
		padding: 16px;
		display: flex;
		align-items: center;
		gap: 12px;
		border: 2px solid transparent;
		transition: all 0.2s ease;
	}

	.summary-card:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
	}

	.summary-card.total { border-color: #4a90e2; }
	.summary-card.active { border-color: #22c55e; }
	.summary-card.closed { border-color: #6b7280; }
	.summary-card.expired { border-color: #ef4444; }
	.summary-card.warning { border-color: #f59e0b; }
	.summary-card.info { border-color: #8b5cf6; }

	.card-icon {
		font-size: 24px;
	}

	.card-content {
		flex: 1;
	}

	.card-value {
		font-size: 24px;
		font-weight: 700;
		color: #ffffff;
		line-height: 1;
	}

	.card-label {
		font-size: 12px;
		color: #999;
		margin-top: 4px;
	}

	.card-percentage {
		font-size: 11px;
		color: #666;
		margin-top: 2px;
	}

	.charts-section {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 24px;
		margin-bottom: 32px;
	}

	.chart-card {
		background: #2a2a2a;
		border-radius: 8px;
		padding: 20px;
	}

	.chart-card h3 {
		margin: 0 0 16px 0;
		color: #ffffff;
		font-size: 1.1rem;
		font-weight: 600;
	}

	.strategy-list, .sector-list {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.strategy-item, .sector-item {
		display: grid;
		grid-template-columns: 1fr 2fr auto;
		align-items: center;
		gap: 12px;
	}

	.strategy-name, .sector-name {
		font-size: 14px;
		color: #cccccc;
		font-weight: 500;
	}

	.strategy-bar, .sector-bar {
		background: #444;
		border-radius: 4px;
		height: 8px;
		overflow: hidden;
	}

	.strategy-fill {
		background: linear-gradient(90deg, #4a90e2, #7b68ee);
		height: 100%;
		transition: width 0.3s ease;
	}

	.sector-fill {
		background: linear-gradient(90deg, #22c55e, #16a34a);
		height: 100%;
		transition: width 0.3s ease;
	}

	.strategy-count, .sector-count {
		font-size: 14px;
		font-weight: 600;
		color: #ffffff;
		min-width: 30px;
		text-align: right;
	}

	.activity-section {
		background: #2a2a2a;
		border-radius: 8px;
		padding: 20px;
	}

	.activity-section h3 {
		margin: 0 0 16px 0;
		color: #ffffff;
		font-size: 1.1rem;
		font-weight: 600;
	}

	.activity-list {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.activity-item {
		display: grid;
		grid-template-columns: 60px 1fr 1fr auto;
		align-items: center;
		gap: 12px;
		padding: 8px 0;
		border-bottom: 1px solid #333;
	}

	.activity-item:last-child {
		border-bottom: none;
	}

	.activity-ticker {
		font-weight: 700;
		color: #ffffff;
		font-size: 14px;
	}

	.activity-strategy {
		font-size: 13px;
		color: #cccccc;
	}

	.activity-sector {
		font-size: 13px;
		color: #999;
	}

	.activity-date {
		font-size: 12px;
		color: #666;
	}

	.empty-activity {
		text-align: center;
		padding: 20px;
		color: #666;
	}

	/* Responsive design */
	@media (max-width: 768px) {
		.summary-grid {
			grid-template-columns: repeat(2, 1fr);
		}

		.charts-section {
			grid-template-columns: 1fr;
		}

		.activity-item {
			grid-template-columns: 1fr;
			gap: 4px;
		}

		.activity-ticker, .activity-strategy {
			font-weight: 600;
		}
	}
</style> 