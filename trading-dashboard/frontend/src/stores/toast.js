import { writable } from 'svelte/store';

// Toast store for managing notifications
function createToastStore() {
	const { subscribe, update } = writable([]);
	
	return {
		subscribe,
		
		// Add a new toast
		add: (message, type = 'success', duration = 3000) => {
			const id = Date.now() + Math.random();
			const toast = {
				id,
				message,
				type,
				duration
			};
			
			update(toasts => [...toasts, toast]);
			
			// Auto-remove toast after duration
			if (duration > 0) {
				setTimeout(() => {
					update(toasts => toasts.filter(t => t.id !== id));
				}, duration);
			}
			
			return id;
		},
		
		// Remove a specific toast
		remove: (id) => {
			update(toasts => toasts.filter(t => t.id !== id));
		},
		
		// Clear all toasts
		clear: () => {
			update(() => []);
		},
		
		// Convenience methods
		success: (message, duration) => {
			return createToastStore().add(message, 'success', duration);
		},
		
		error: (message, duration) => {
			return createToastStore().add(message, 'error', duration);
		},
		
		warning: (message, duration) => {
			return createToastStore().add(message, 'warning', duration);
		},
		
		info: (message, duration) => {
			return createToastStore().add(message, 'info', duration);
		}
	};
}

export const toastStore = createToastStore();

// Global toast functions for easy access
export const toast = {
	success: (message, duration = 3000) => toastStore.add(message, 'success', duration),
	error: (message, duration = 4000) => toastStore.add(message, 'error', duration),
	warning: (message, duration = 3500) => toastStore.add(message, 'warning', duration),
	info: (message, duration = 3000) => toastStore.add(message, 'info', duration)
}; 