<script>
	import { createEventDispatcher, onMount } from 'svelte';
	import { tradesStore } from '../stores/trades.js';
	import { toastStore } from '../stores/toast.js';

	const dispatch = createEventDispatcher();

	// Props
	export let isOpen = false;
	export let trade = null; // If editing existing trade
	export let selectedDate = null; // If creating new trade for specific date
	export let selectedSector = null; // If creating new trade for specific sector

	// Form data
	let formData = {
		ticker: '',
		sector: selectedSector || 'Technology',
		strategy_type: 'Long Call',
		entry_date: selectedDate || new Date().toISOString().split('T')[0],
		expiration_date: '',
		target_price: '',
		stop_loss: '',
		notes: ''
	};

	// Form state
	let isLoading = false;
	let errors = {};
	let strategyTypes = [];

	// Sector options
	const sectors = [
		'Technology', 'Healthcare', 'Financial', 'Consumer Discretionary',
		'Communication Services', 'Industrials', 'Consumer Staples',
		'Energy', 'Utilities', 'Real Estate', 'Materials'
	];

	// Load strategy types and populate form
	onMount(async () => {
		try {
			strategyTypes = await tradesStore.loadStrategyTypes();
			populateFormFromTrade();
		} catch (error) {
			toastStore.error('Failed to load strategy types');
			console.error('Error loading strategy types:', error);
		}
	});
	
	// Watch for changes to trade prop and update form
	$: if (trade) {
		populateFormFromTrade();
	}
	
	// Reset form when modal opens for new trade
	$: if (isOpen && !trade) {
		resetForm();
	}
	
	function populateFormFromTrade() {
		if (trade && typeof trade === 'object') {
			formData = {
				ticker: trade.ticker || '',
				sector: trade.sector || 'Technology',
				strategy_type: trade.strategy_type || 'Long Call',
				entry_date: trade.entry_date ? trade.entry_date.split('T')[0] : '',
				expiration_date: trade.expiration_date ? trade.expiration_date.split('T')[0] : '',
				target_price: trade.target_price || '',
				stop_loss: trade.stop_loss || '',
				notes: trade.notes || ''
			};
		}
	}
	
	function resetForm() {
		formData = {
			ticker: '',
			sector: selectedSector || 'Technology',
			strategy_type: 'Long Call',
			entry_date: selectedDate || new Date().toISOString().split('T')[0],
			expiration_date: '',
			target_price: '',
			stop_loss: '',
			notes: ''
		};
		errors = {};
	}

	// Form validation
	function validateForm() {
		errors = {};

		if (!formData.ticker.trim()) {
			errors.ticker = 'Ticker is required';
		} else if (!/^[A-Z]{1,5}$/.test(formData.ticker.trim().toUpperCase())) {
			errors.ticker = 'Ticker must be 1-5 uppercase letters';
		}

		if (!formData.entry_date) {
			errors.entry_date = 'Entry date is required';
		}

		if (!formData.expiration_date) {
			errors.expiration_date = 'Expiration date is required';
		} else if (new Date(formData.expiration_date) <= new Date(formData.entry_date)) {
			errors.expiration_date = 'Expiration must be after entry date';
		}

		if (formData.target_price && isNaN(parseFloat(formData.target_price))) {
			errors.target_price = 'Target price must be a valid number';
		}

		if (formData.stop_loss && isNaN(parseFloat(formData.stop_loss))) {
			errors.stop_loss = 'Stop loss must be a valid number';
		}

		return Object.keys(errors).length === 0;
	}

	// Handle form submission
	async function handleSubmit() {
		if (!validateForm()) {
			toastStore.error('Please fix form errors before submitting');
			return;
		}

		isLoading = true;

		try {
			// Prepare request data
			const requestData = {
				ticker: formData.ticker.trim().toUpperCase(),
				sector: formData.sector,
				strategy_type: formData.strategy_type,
				entry_date: new Date(formData.entry_date + 'T00:00:00Z'),
				expiration_date: new Date(formData.expiration_date + 'T00:00:00Z'),
				target_price: formData.target_price ? parseFloat(formData.target_price) : null,
				stop_loss: formData.stop_loss ? parseFloat(formData.stop_loss) : null,
				notes: formData.notes.trim()
			};

			let result;
			if (trade && typeof trade === 'object' && trade.id) {
				// Update existing trade
				result = await window.go.main.App.UpdateTrade(trade.id, requestData);
				toastStore.success('Trade updated successfully!');
			} else {
				// Create new trade
				result = await window.go.main.App.CreateTrade(requestData);
				toastStore.success('Trade created successfully!');
			}

			// Close modal and notify parent first
			close();
			dispatch('trade-saved', result);

		} catch (error) {
			console.error('Error saving trade:', error);
			toastStore.error(`Failed to ${trade ? 'update' : 'create'} trade: ${error.message || error}`);
		} finally {
			isLoading = false;
		}
	}

	// Close modal
	function close() {
		isOpen = false;
		dispatch('close');
	}

	// Handle escape key
	function handleKeydown(event) {
		if (event.key === 'Escape') {
			close();
		}
	}

	// Auto-format ticker input
	function handleTickerInput(event) {
		formData.ticker = event.target.value.toUpperCase();
	}

	// Get strategy color
	function getStrategyColor(strategyName) {
		const strategy = strategyTypes.find(s => s.name === strategyName);
		return strategy?.color_hex || '#4a90e2';
	}
</script>

<!-- Modal backdrop -->
{#if isOpen}
	<div class="modal-backdrop" on:click={close} on:keydown={handleKeydown}>
		<div class="modal-container" on:click|stopPropagation>
			<div class="modal-header">
				<h2>{trade ? 'Edit Trade' : 'New Trade'}</h2>
				<button class="close-btn" on:click={close}>
					<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<line x1="18" y1="6" x2="6" y2="18"></line>
						<line x1="6" y1="6" x2="18" y2="18"></line>
					</svg>
				</button>
			</div>

			<form class="modal-body" on:submit|preventDefault={handleSubmit}>
				<div class="form-row">
					<div class="form-group">
						<label for="ticker">Ticker Symbol *</label>
						<input
							id="ticker"
							type="text"
							bind:value={formData.ticker}
							on:input={handleTickerInput}
							placeholder="AAPL"
							maxlength="5"
							class:error={errors.ticker}
							disabled={isLoading}
						/>
						{#if errors.ticker}
							<span class="error-message">{errors.ticker}</span>
						{/if}
					</div>

					<div class="form-group">
						<label for="sector">Sector</label>
						<select
							id="sector"
							bind:value={formData.sector}
							disabled={isLoading}
						>
							{#each sectors as sector}
								<option value={sector}>{sector}</option>
							{/each}
						</select>
					</div>
				</div>

				<div class="form-group">
					<label for="strategy">Strategy Type</label>
					<select
						id="strategy"
						bind:value={formData.strategy_type}
						disabled={isLoading}
						style="border-left: 4px solid {getStrategyColor(formData.strategy_type)}"
					>
						{#each strategyTypes as strategy}
							<option value={strategy.name}>{strategy.name} - {strategy.description}</option>
						{/each}
					</select>
				</div>

				<div class="form-row">
					<div class="form-group">
						<label for="entry_date">Entry Date *</label>
						<input
							id="entry_date"
							type="date"
							bind:value={formData.entry_date}
							class:error={errors.entry_date}
							disabled={isLoading}
						/>
						{#if errors.entry_date}
							<span class="error-message">{errors.entry_date}</span>
						{/if}
					</div>

					<div class="form-group">
						<label for="expiration_date">Expiration Date *</label>
						<input
							id="expiration_date"
							type="date"
							bind:value={formData.expiration_date}
							class:error={errors.expiration_date}
							disabled={isLoading}
						/>
						{#if errors.expiration_date}
							<span class="error-message">{errors.expiration_date}</span>
						{/if}
					</div>
				</div>

				<div class="form-row">
					<div class="form-group">
						<label for="target_price">Target Price</label>
						<input
							id="target_price"
							type="number"
							step="0.01"
							min="0"
							bind:value={formData.target_price}
							placeholder="0.00"
							class:error={errors.target_price}
							disabled={isLoading}
						/>
						{#if errors.target_price}
							<span class="error-message">{errors.target_price}</span>
						{/if}
					</div>

					<div class="form-group">
						<label for="stop_loss">Stop Loss</label>
						<input
							id="stop_loss"
							type="number"
							step="0.01"
							min="0"
							bind:value={formData.stop_loss}
							placeholder="0.00"
							class:error={errors.stop_loss}
							disabled={isLoading}
						/>
						{#if errors.stop_loss}
							<span class="error-message">{errors.stop_loss}</span>
						{/if}
					</div>
				</div>

				<div class="form-group">
					<label for="notes">Notes</label>
					<textarea
						id="notes"
						bind:value={formData.notes}
						placeholder="Optional notes about this trade..."
						rows="3"
						disabled={isLoading}
					></textarea>
				</div>

				<div class="modal-footer">
					<button type="button" class="btn-secondary" on:click={close} disabled={isLoading}>
						Cancel
					</button>
					<button type="submit" class="btn-primary" disabled={isLoading}>
						{#if isLoading}
							<span class="loading-spinner"></span>
							{trade ? 'Updating...' : 'Creating...'}
						{:else}
							{trade ? 'Update Trade' : 'Create Trade'}
						{/if}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<style>
	.modal-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: rgba(0, 0, 0, 0.75);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		backdrop-filter: blur(4px);
		animation: fadeIn 0.2s ease-out;
	}

	.modal-container {
		background: #1a1a1a;
		border-radius: 12px;
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
		width: 90%;
		max-width: 600px;
		max-height: 90vh;
		overflow: hidden;
		animation: slideIn 0.3s ease-out;
		border: 1px solid #333;
	}

	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 24px;
		border-bottom: 1px solid #333;
		background: linear-gradient(135deg, #2a2a2a, #1a1a1a);
	}

	.modal-header h2 {
		margin: 0;
		color: #ffffff;
		font-size: 24px;
		font-weight: 600;
	}

	.close-btn {
		background: none;
		border: none;
		color: #999;
		cursor: pointer;
		padding: 8px;
		border-radius: 6px;
		transition: all 0.2s ease;
	}

	.close-btn:hover {
		color: #ffffff;
		background: #333;
	}

	.modal-body {
		padding: 24px;
		max-height: calc(90vh - 140px);
		overflow-y: auto;
	}

	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 16px;
	}

	.form-group {
		margin-bottom: 20px;
	}

	.form-group label {
		display: block;
		margin-bottom: 6px;
		color: #cccccc;
		font-weight: 500;
		font-size: 14px;
	}

	.form-group input,
	.form-group select,
	.form-group textarea {
		width: 100%;
		padding: 12px;
		background: #2a2a2a;
		border: 2px solid #444;
		border-radius: 6px;
		color: #ffffff;
		font-size: 14px;
		transition: all 0.2s ease;
		box-sizing: border-box;
	}

	.form-group input:focus,
	.form-group select:focus,
	.form-group textarea:focus {
		outline: none;
		border-color: #4a90e2;
		box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.1);
	}

	.form-group input.error,
	.form-group select.error,
	.form-group textarea.error {
		border-color: #ef4444;
		box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
	}

	.form-group input:disabled,
	.form-group select:disabled,
	.form-group textarea:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.error-message {
		display: block;
		color: #ef4444;
		font-size: 12px;
		margin-top: 4px;
	}

	.modal-footer {
		display: flex;
		gap: 12px;
		justify-content: flex-end;
		padding: 24px;
		border-top: 1px solid #333;
		background: #1a1a1a;
	}

	.btn-secondary,
	.btn-primary {
		padding: 12px 24px;
		border-radius: 6px;
		font-weight: 600;
		font-size: 14px;
		cursor: pointer;
		transition: all 0.2s ease;
		border: none;
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.btn-secondary {
		background: #333;
		color: #cccccc;
	}

	.btn-secondary:hover:not(:disabled) {
		background: #444;
		color: #ffffff;
	}

	.btn-primary {
		background: linear-gradient(135deg, #4a90e2, #357abd);
		color: white;
	}

	.btn-primary:hover:not(:disabled) {
		background: linear-gradient(135deg, #357abd, #2968a3);
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(74, 144, 226, 0.3);
	}

	.btn-secondary:disabled,
	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		transform: none;
		box-shadow: none;
	}

	.loading-spinner {
		width: 16px;
		height: 16px;
		border: 2px solid transparent;
		border-top: 2px solid currentColor;
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	@keyframes fadeIn {
		from { opacity: 0; }
		to { opacity: 1; }
	}

	@keyframes slideIn {
		from {
			opacity: 0;
			transform: translateY(-20px) scale(0.95);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	/* Responsive design */
	@media (max-width: 640px) {
		.modal-container {
			width: 95%;
			margin: 20px;
		}

		.form-row {
			grid-template-columns: 1fr;
		}

		.modal-header,
		.modal-body,
		.modal-footer {
			padding: 16px;
		}
	}
</style> 