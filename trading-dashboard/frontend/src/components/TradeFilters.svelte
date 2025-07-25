<script>
	import { createEventDispatcher } from 'svelte';
	import { tradesStore } from '../stores/trades.js';

	const dispatch = createEventDispatcher();

	// Filter state
	export let filters = {
		status: 'all',
		strategy: 'all',
		sector: 'all',
		search: ''
	};

	let strategyTypes = [];
	let isExpanded = false;

	// Get strategy types from store
	$: strategyTypes = $tradesStore.strategyTypes || [];

	const sectors = [
		'Technology', 'Healthcare', 'Financial', 'Consumer Discretionary',
		'Communication Services', 'Industrials', 'Consumer Staples',
		'Energy', 'Utilities', 'Real Estate', 'Materials'
	];

	const statusOptions = [
		{ value: 'all', label: 'All Status' },
		{ value: 'active', label: 'Active' },
		{ value: 'closed', label: 'Closed' },
		{ value: 'expired', label: 'Expired' }
	];

	// Emit filter changes
	function updateFilters() {
		dispatch('filters-change', filters);
	}

	function clearFilters() {
		filters = {
			status: 'all',
			strategy: 'all',
			sector: 'all',
			search: ''
		};
		updateFilters();
	}

	function toggleExpanded() {
		isExpanded = !isExpanded;
	}

	// Reactive filter updates
	$: if (filters) {
		updateFilters();
	}

	// Count active filters
	$: activeFilterCount = Object.entries(filters).filter(([key, value]) => {
		if (key === 'search') return value.trim() !== '';
		return value !== 'all';
	}).length;
</script>

<div class="trade-filters">
	<div class="filters-header">
		<button class="toggle-filters" on:click={toggleExpanded}>
			<span class="filter-icon">🔍</span>
			<span class="filter-text">Filters</span>
			{#if activeFilterCount > 0}
				<span class="filter-count">{activeFilterCount}</span>
			{/if}
			<span class="expand-arrow" class:expanded={isExpanded}>▼</span>
		</button>

		{#if activeFilterCount > 0}
			<button class="clear-filters" on:click={clearFilters}>
				Clear All
			</button>
		{/if}
	</div>

	{#if isExpanded}
		<div class="filters-content">
			<div class="filter-row">
				<!-- Search -->
				<div class="filter-group">
					<label for="search-filter">Search</label>
					<div class="search-input-container">
						<input
							id="search-filter"
							type="text"
							placeholder="Search by ticker..."
							bind:value={filters.search}
							class="search-input"
						/>
						<span class="search-icon">🔍</span>
					</div>
				</div>

				<!-- Status Filter -->
				<div class="filter-group">
					<label for="status-filter">Status</label>
					<select id="status-filter" bind:value={filters.status} class="filter-select">
						{#each statusOptions as option}
							<option value={option.value}>{option.label}</option>
						{/each}
					</select>
				</div>
			</div>

			<div class="filter-row">
				<!-- Strategy Filter -->
				<div class="filter-group">
					<label for="strategy-filter">Strategy</label>
					<select id="strategy-filter" bind:value={filters.strategy} class="filter-select">
						<option value="all">All Strategies</option>
						{#each strategyTypes as strategy}
							<option value={strategy.name}>{strategy.name}</option>
						{/each}
					</select>
				</div>

				<!-- Sector Filter -->
				<div class="filter-group">
					<label for="sector-filter">Sector</label>
					<select id="sector-filter" bind:value={filters.sector} class="filter-select">
						<option value="all">All Sectors</option>
						{#each sectors as sector}
							<option value={sector}>{sector}</option>
						{/each}
					</select>
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	.trade-filters {
		background: #1a1a1a;
		border: 1px solid #333;
		border-radius: 8px;
		margin-bottom: 16px;
	}

	.filters-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 12px 16px;
	}

	.toggle-filters {
		display: flex;
		align-items: center;
		gap: 8px;
		background: none;
		border: none;
		color: #cccccc;
		cursor: pointer;
		font-size: 14px;
		transition: color 0.2s ease;
	}

	.toggle-filters:hover {
		color: #ffffff;
	}

	.filter-icon {
		font-size: 16px;
	}

	.filter-count {
		background: #4a90e2;
		color: white;
		border-radius: 10px;
		padding: 2px 6px;
		font-size: 12px;
		font-weight: 600;
		min-width: 18px;
		text-align: center;
	}

	.expand-arrow {
		font-size: 12px;
		transition: transform 0.2s ease;
	}

	.expand-arrow.expanded {
		transform: rotate(180deg);
	}

	.clear-filters {
		background: #ef4444;
		color: white;
		border: none;
		padding: 6px 12px;
		border-radius: 4px;
		cursor: pointer;
		font-size: 12px;
		transition: background-color 0.2s ease;
	}

	.clear-filters:hover {
		background: #dc2626;
	}

	.filters-content {
		padding: 0 16px 16px;
		border-top: 1px solid #333;
		animation: slideDown 0.2s ease-out;
	}

	.filter-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 16px;
		margin-top: 16px;
	}

	.filter-group {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.filter-group label {
		color: #cccccc;
		font-size: 12px;
		font-weight: 500;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.search-input-container {
		position: relative;
	}

	.search-input,
	.filter-select {
		width: 100%;
		padding: 8px 12px;
		background: #2a2a2a;
		border: 1px solid #444;
		border-radius: 4px;
		color: #ffffff;
		font-size: 14px;
		transition: border-color 0.2s ease;
	}

	.search-input {
		padding-right: 32px;
	}

	.search-input:focus,
	.filter-select:focus {
		outline: none;
		border-color: #4a90e2;
	}

	.search-icon {
		position: absolute;
		right: 8px;
		top: 50%;
		transform: translateY(-50%);
		color: #666;
		font-size: 14px;
		pointer-events: none;
	}

	.filter-select option {
		background: #2a2a2a;
		color: #ffffff;
	}

	@keyframes slideDown {
		from {
			opacity: 0;
			transform: translateY(-8px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	/* Responsive design */
	@media (max-width: 640px) {
		.filter-row {
			grid-template-columns: 1fr;
			gap: 12px;
		}

		.filters-header {
			padding: 10px 12px;
		}

		.filters-content {
			padding: 0 12px 12px;
		}
	}
</style> 