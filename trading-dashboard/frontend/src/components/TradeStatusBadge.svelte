<script>
	import { createEventDispatcher } from 'svelte';

	export let trade;
	export let showDropdown = false;
	export let size = 'small'; // 'small', 'medium', 'large'

	const dispatch = createEventDispatcher();

	let isDropdownOpen = false;

	const statusConfig = {
		active: {
			label: 'Active',
			color: '#22c55e',
			bgColor: 'rgba(34, 197, 94, 0.1)',
			icon: '●'
		},
		closed: {
			label: 'Closed',
			color: '#6b7280',
			bgColor: 'rgba(107, 114, 128, 0.1)',
			icon: '✓'
		},
		expired: {
			label: 'Expired',
			color: '#ef4444',
			bgColor: 'rgba(239, 68, 68, 0.1)',
			icon: '⚫'
		}
	};

	function handleStatusChange(newStatus) {
		if (newStatus !== trade.status) {
			dispatch('status-change', {
				trade,
				newStatus,
				oldStatus: trade.status
			});
		}
		isDropdownOpen = false;
	}

	function toggleDropdown(event) {
		if (showDropdown) {
			event.stopPropagation();
			isDropdownOpen = !isDropdownOpen;
		}
	}

	function closeDropdown() {
		isDropdownOpen = false;
	}

	// Close dropdown when clicking outside
	function handleClickOutside(event) {
		if (isDropdownOpen && !event.target.closest('.status-badge-container')) {
			closeDropdown();
		}
	}

	$: currentStatus = statusConfig[trade.status] || statusConfig.active;
	$: sizeClass = `badge-${size}`;
</script>

<svelte:window on:click={handleClickOutside} />

<div class="status-badge-container" class:interactive={showDropdown}>
	<button
		class="status-badge {sizeClass}"
		class:clickable={showDropdown}
		style="color: {currentStatus.color}; background-color: {currentStatus.bgColor}; border-color: {currentStatus.color};"
		on:click={toggleDropdown}
		title="{showDropdown ? 'Click to change status' : `Status: ${currentStatus.label}`}"
	>
		<span class="status-icon">{currentStatus.icon}</span>
		<span class="status-label">{currentStatus.label}</span>
		{#if showDropdown}
			<span class="dropdown-arrow">▼</span>
		{/if}
	</button>

	{#if showDropdown && isDropdownOpen}
		<div class="status-dropdown">
			{#each Object.entries(statusConfig) as [statusKey, config]}
				<button
					class="status-option"
					class:current={statusKey === trade.status}
					style="color: {config.color};"
					on:click={() => handleStatusChange(statusKey)}
				>
					<span class="status-icon">{config.icon}</span>
					<span class="status-label">{config.label}</span>
					{#if statusKey === trade.status}
						<span class="checkmark">✓</span>
					{/if}
				</button>
			{/each}
		</div>
	{/if}
</div>

<style>
	.status-badge-container {
		position: relative;
		display: inline-block;
	}

	.status-badge {
		display: inline-flex;
		align-items: center;
		gap: 4px;
		border: 1px solid;
		border-radius: 12px;
		font-weight: 500;
		font-family: system-ui, -apple-system, sans-serif;
		transition: all 0.2s ease;
		background: none;
		cursor: default;
		user-select: none;
	}

	.status-badge.clickable {
		cursor: pointer;
	}

	.status-badge.clickable:hover {
		transform: translateY(-1px);
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.badge-small {
		padding: 2px 6px;
		font-size: 0.75rem;
	}

	.badge-medium {
		padding: 4px 8px;
		font-size: 0.875rem;
	}

	.badge-large {
		padding: 6px 12px;
		font-size: 1rem;
	}

	.status-icon {
		font-size: 0.8em;
	}

	.dropdown-arrow {
		font-size: 0.6em;
		margin-left: 2px;
		transition: transform 0.2s ease;
	}

	.status-badge.clickable:hover .dropdown-arrow {
		transform: rotate(180deg);
	}

	.status-dropdown {
		position: absolute;
		top: 100%;
		left: 0;
		z-index: 1000;
		background: #1a1a1a;
		border: 1px solid #333;
		border-radius: 8px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
		min-width: 120px;
		overflow: hidden;
		animation: slideDown 0.2s ease-out;
	}

	.status-option {
		display: flex;
		align-items: center;
		gap: 8px;
		width: 100%;
		padding: 8px 12px;
		background: none;
		border: none;
		text-align: left;
		cursor: pointer;
		transition: background-color 0.2s ease;
		font-size: 0.875rem;
	}

	.status-option:hover {
		background-color: rgba(255, 255, 255, 0.05);
	}

	.status-option.current {
		background-color: rgba(74, 144, 226, 0.1);
	}

	.checkmark {
		margin-left: auto;
		color: #4a90e2;
		font-size: 0.8rem;
	}

	@keyframes slideDown {
		from {
			opacity: 0;
			transform: translateY(-4px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	/* Responsive adjustments */
	@media (max-width: 640px) {
		.badge-small {
			padding: 1px 4px;
			font-size: 0.7rem;
		}

		.status-dropdown {
			min-width: 100px;
		}

		.status-option {
			padding: 6px 8px;
			font-size: 0.8rem;
		}
	}
</style> 