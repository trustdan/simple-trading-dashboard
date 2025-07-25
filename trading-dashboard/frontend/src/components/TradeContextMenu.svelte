<script>
	import { createEventDispatcher } from 'svelte';

	export let trade;
	export let x = 0;
	export let y = 0;
	export let visible = false;

	const dispatch = createEventDispatcher();

	let menuElement;

	// Close menu when clicking outside
	function handleClickOutside(event) {
		if (visible && menuElement && !menuElement.contains(event.target)) {
			close();
		}
	}

	function close() {
		visible = false;
		dispatch('close');
	}

	function handleEdit() {
		dispatch('edit', trade);
		close();
	}

	function handleStatusChange(newStatus) {
		dispatch('status-change', { trade, newStatus });
		close();
	}

	function handleDelete() {
		if (confirm(`Are you sure you want to delete the ${trade.ticker} trade?`)) {
			dispatch('delete', trade);
		}
		close();
	}

	function handleDuplicate() {
		dispatch('duplicate', trade);
		close();
	}

	// Position menu to stay within viewport
	function getMenuPosition() {
		if (!menuElement) return { x, y };
		
		const rect = menuElement.getBoundingClientRect();
		const viewportWidth = window.innerWidth;
		const viewportHeight = window.innerHeight;
		
		let adjustedX = x;
		let adjustedY = y;
		
		// Adjust horizontal position if menu would overflow
		if (x + rect.width > viewportWidth) {
			adjustedX = viewportWidth - rect.width - 10;
		}
		
		// Adjust vertical position if menu would overflow
		if (y + rect.height > viewportHeight) {
			adjustedY = viewportHeight - rect.height - 10;
		}
		
		return { x: adjustedX, y: adjustedY };
	}

	$: if (visible && menuElement) {
		const position = getMenuPosition();
		menuElement.style.left = `${position.x}px`;
		menuElement.style.top = `${position.y}px`;
	}
</script>

<svelte:window on:click={handleClickOutside} />

{#if visible}
	<div 
		bind:this={menuElement}
		class="context-menu"
		style="left: {x}px; top: {y}px;"
		on:click|stopPropagation
	>
		<div class="menu-header">
			<span class="trade-ticker">{trade.ticker}</span>
			<span class="trade-strategy">{trade.strategy_type}</span>
		</div>
		
		<div class="menu-divider"></div>
		
		<button class="menu-item" on:click={handleEdit}>
			<span class="menu-icon">✏️</span>
			Edit Trade
		</button>
		
		<button class="menu-item" on:click={handleDuplicate}>
			<span class="menu-icon">📋</span>
			Duplicate Trade
		</button>
		
		<div class="menu-divider"></div>
		
		<div class="menu-section">
			<span class="menu-section-title">Change Status</span>
			<button 
				class="menu-item status-item" 
				class:current={trade.status === 'active'}
				on:click={() => handleStatusChange('active')}
			>
				<span class="status-indicator active">●</span>
				Active
			</button>
			<button 
				class="menu-item status-item" 
				class:current={trade.status === 'closed'}
				on:click={() => handleStatusChange('closed')}
			>
				<span class="status-indicator closed">✓</span>
				Closed
			</button>
			<button 
				class="menu-item status-item" 
				class:current={trade.status === 'expired'}
				on:click={() => handleStatusChange('expired')}
			>
				<span class="status-indicator expired">⚫</span>
				Expired
			</button>
		</div>
		
		<div class="menu-divider"></div>
		
		<button class="menu-item danger" on:click={handleDelete}>
			<span class="menu-icon">🗑️</span>
			Delete Trade
		</button>
	</div>
{/if}

<style>
	.context-menu {
		position: fixed;
		z-index: 10000;
		background: #1a1a1a;
		border: 1px solid #333;
		border-radius: 8px;
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
		min-width: 200px;
		padding: 4px 0;
		font-size: 14px;
		animation: contextMenuSlide 0.15s ease-out;
	}

	.menu-header {
		padding: 8px 12px;
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.trade-ticker {
		font-weight: 700;
		color: #ffffff;
		font-size: 16px;
	}

	.trade-strategy {
		font-size: 12px;
		color: #999;
	}

	.menu-divider {
		height: 1px;
		background: #333;
		margin: 4px 0;
	}

	.menu-section {
		padding: 4px 0;
	}

	.menu-section-title {
		display: block;
		padding: 4px 12px;
		font-size: 11px;
		font-weight: 600;
		color: #666;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.menu-item {
		display: flex;
		align-items: center;
		gap: 8px;
		width: 100%;
		padding: 8px 12px;
		background: none;
		border: none;
		color: #cccccc;
		text-align: left;
		cursor: pointer;
		transition: background-color 0.15s ease;
		font-size: 14px;
	}

	.menu-item:hover {
		background: rgba(255, 255, 255, 0.05);
		color: #ffffff;
	}

	.menu-item.current {
		background: rgba(74, 144, 226, 0.1);
		color: #4a90e2;
	}

	.menu-item.danger {
		color: #ef4444;
	}

	.menu-item.danger:hover {
		background: rgba(239, 68, 68, 0.1);
		color: #ff6b6b;
	}

	.status-item {
		padding-left: 20px;
		font-size: 13px;
	}

	.menu-icon {
		font-size: 16px;
		width: 20px;
		text-align: center;
	}

	.status-indicator {
		font-size: 12px;
		width: 16px;
		text-align: center;
	}

	.status-indicator.active {
		color: #22c55e;
	}

	.status-indicator.closed {
		color: #6b7280;
	}

	.status-indicator.expired {
		color: #ef4444;
	}

	@keyframes contextMenuSlide {
		from {
			opacity: 0;
			transform: translateY(-4px) scale(0.95);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}
</style> 