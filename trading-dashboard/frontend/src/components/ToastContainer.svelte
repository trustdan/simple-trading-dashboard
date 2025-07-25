<script>
	import { toastStore } from '../stores/toast.js';
	import Toast from './Toast.svelte';
	
	$: toasts = $toastStore;
</script>

<div class="toast-container">
	{#each toasts as toast (toast.id)}
		<Toast
			message={toast.message}
			type={toast.type}
			duration={0}
			on:close={() => toastStore.remove(toast.id)}
		/>
	{/each}
</div>

<style>
	.toast-container {
		position: fixed;
		top: 20px;
		right: 20px;
		z-index: 10000;
		display: flex;
		flex-direction: column;
		gap: 8px;
		pointer-events: none;
	}
	
	.toast-container :global(.toast) {
		pointer-events: auto;
	}
	
	@media (max-width: 640px) {
		.toast-container {
			left: 20px;
			right: 20px;
		}
	}
</style> 