<script>
	import { createEventDispatcher, onMount } from 'svelte';
	
	export let message = '';
	export let type = 'success'; // 'success', 'error', 'warning', 'info'
	export let duration = 3000; // Auto-hide duration in ms
	export let showIcon = true;
	
	const dispatch = createEventDispatcher();
	let visible = false;
	let timeoutId;
	
	// Show toast on mount
	onMount(() => {
		visible = true;
		
		if (duration > 0) {
			timeoutId = setTimeout(() => {
				hide();
			}, duration);
		}
	});
	
	function hide() {
		visible = false;
		setTimeout(() => {
			dispatch('close');
		}, 300); // Wait for animation to complete
	}
	
	function handleClick() {
		if (timeoutId) {
			clearTimeout(timeoutId);
		}
		hide();
	}
	
	// Icon based on type
	$: icon = {
		success: '✓',
		error: '✕',
		warning: '⚠',
		info: 'ℹ'
	}[type];
</script>

<div 
	class="toast toast-{type}" 
	class:visible 
	on:click={handleClick}
	role="alert"
	aria-live="polite"
>
	{#if showIcon}
		<div class="toast-icon">
			{icon}
		</div>
	{/if}
	
	<div class="toast-message">
		{message}
	</div>
	
	<button class="toast-close" on:click={handleClick} aria-label="Close notification">
		×
	</button>
</div>

<style>
	.toast {
		position: fixed;
		top: 20px;
		right: 20px;
		z-index: 10000;
		
		display: flex;
		align-items: center;
		gap: 12px;
		
		min-width: 300px;
		max-width: 500px;
		padding: 16px 20px;
		
		background: rgba(0, 0, 0, 0.9);
		border: 1px solid;
		border-radius: 8px;
		
		color: white;
		font-size: 14px;
		font-weight: 500;
		
		backdrop-filter: blur(10px);
		cursor: pointer;
		
		transform: translateX(100%);
		opacity: 0;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	}
	
	.toast.visible {
		transform: translateX(0);
		opacity: 1;
	}
	
	.toast-success {
		background: rgba(34, 197, 94, 0.15);
		border-color: rgba(34, 197, 94, 0.5);
		color: #22c55e;
	}
	
	.toast-error {
		background: rgba(239, 68, 68, 0.15);
		border-color: rgba(239, 68, 68, 0.5);
		color: #ef4444;
	}
	
	.toast-warning {
		background: rgba(245, 158, 11, 0.15);
		border-color: rgba(245, 158, 11, 0.5);
		color: #f59e0b;
	}
	
	.toast-info {
		background: rgba(59, 130, 246, 0.15);
		border-color: rgba(59, 130, 246, 0.5);
		color: #3b82f6;
	}
	
	.toast-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 24px;
		height: 24px;
		font-size: 16px;
		font-weight: bold;
		flex-shrink: 0;
	}
	
	.toast-message {
		flex: 1;
		line-height: 1.4;
	}
	
	.toast-close {
		background: none;
		border: none;
		color: inherit;
		font-size: 20px;
		cursor: pointer;
		padding: 0;
		width: 24px;
		height: 24px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 4px;
		transition: background-color 0.2s;
		flex-shrink: 0;
	}
	
	.toast-close:hover {
		background: rgba(255, 255, 255, 0.1);
	}
	
	.toast:hover {
		transform: translateX(-4px);
	}
	
	@media (max-width: 640px) {
		.toast {
			left: 20px;
			right: 20px;
			min-width: auto;
			transform: translateY(-100%);
		}
		
		.toast.visible {
			transform: translateY(0);
		}
		
		.toast:hover {
			transform: translateY(-2px);
		}
	}
</style> 