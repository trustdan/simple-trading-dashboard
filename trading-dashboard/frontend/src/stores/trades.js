import { writable } from 'svelte/store';

// Market sectors list (same as market store)
const SECTORS = [
	'Basic Materials',
	'Communication Services',
	'Consumer Cyclical',
	'Consumer Defensive',
	'Energy',
	'Financial Services',
	'Healthcare',
	'Industrials',
	'Real Estate',
	'Technology',
	'Utilities'
];

// Create the trades store
function createTradesStore() {
	const { subscribe, set, update } = writable({
		trades: [],
		strategyTypes: [],
		sectors: SECTORS,
		dateColumns: [],
		loading: false,
		selectedTrade: null
	});

	return {
		subscribe,
		
		// Set date columns for calendar view
		setDateColumns: (columns) => update(state => ({
			...state,
			dateColumns: columns
		})),
		
		// Load trades by date range
		loadTradesByDateRange: async (startDate, endDate) => {
			update(state => ({ ...state, loading: true }));
			
			try {
				// Add a timeout to prevent infinite loading
				const timeoutPromise = new Promise((_, reject) => 
					setTimeout(() => reject(new Error('Request timeout')), 10000)
				);
				
				// Call the Wails backend API with timeout
				const trades = await Promise.race([
					window['go']['main']['App']['GetActiveTradesByDateRange'](startDate, endDate),
					timeoutPromise
				]);
				
				update(state => ({
					...state,
					trades: trades || [],
					loading: false
				}));
				
				console.log(`Loaded ${trades?.length || 0} trades for date range`);
				return trades;
			} catch (error) {
				console.error('Failed to load trades:', error);
				update(state => ({ 
					...state, 
					trades: [],
					loading: false 
				}));
				throw error;
			}
		},
		
		// Load all strategy types
		loadStrategyTypes: async () => {
			try {
				// Add a timeout to prevent infinite loading
				const timeoutPromise = new Promise((_, reject) => 
					setTimeout(() => reject(new Error('Request timeout')), 10000)
				);
				
				// Call the Wails backend API with timeout
				const strategies = await Promise.race([
					window['go']['main']['App']['GetStrategyTypes'](),
					timeoutPromise
				]);
				
				update(state => ({
					...state,
					strategyTypes: strategies || []
				}));
				
				console.log(`Loaded ${strategies?.length || 0} strategy types`);
				return strategies;
			} catch (error) {
				console.error('Failed to load strategy types:', error);
				update(state => ({
					...state,
					strategyTypes: []
				}));
				throw error;
			}
		},
		
		// Create a new trade
		createTrade: async (tradeRequest) => {
			update(state => ({ ...state, loading: true }));
			
			try {
				const newTrade = await window['go']['main']['App']['CreateTrade'](tradeRequest);
				
				update(state => ({
					...state,
					trades: [...state.trades, newTrade],
					loading: false
				}));
				
				return newTrade;
			} catch (error) {
				console.error('Failed to create trade:', error);
				update(state => ({ ...state, loading: false }));
				throw error;
			}
		},
		
		// Update an existing trade
		updateTrade: async (id, tradeRequest) => {
			update(state => ({ ...state, loading: true }));
			
			try {
				const updatedTrade = await window['go']['main']['App']['UpdateTrade'](id, tradeRequest);
				
				update(state => ({
					...state,
					trades: state.trades.map(trade => 
						trade.id === id ? updatedTrade : trade
					),
					loading: false
				}));
				
				return updatedTrade;
			} catch (error) {
				console.error('Failed to update trade:', error);
				update(state => ({ ...state, loading: false }));
				throw error;
			}
		},
		
		// Update trade status
		updateTradeStatus: async (id, status) => {
			try {
				const updatedTrade = await window['go']['main']['App']['UpdateTradeStatus'](id, status);
				
				update(state => ({
					...state,
					trades: state.trades.map(trade => 
						trade.id === id ? updatedTrade : trade
					)
				}));
				
				return updatedTrade;
			} catch (error) {
				console.error('Failed to update trade status:', error);
				throw error;
			}
		},
		
		// Delete a trade
		deleteTrade: async (id) => {
			try {
				await window['go']['main']['App']['DeleteTrade'](id);
				
				update(state => ({
					...state,
					trades: state.trades.filter(trade => trade.id !== id)
				}));
			} catch (error) {
				console.error('Failed to delete trade:', error);
				throw error;
			}
		},
		
		// Select a trade for editing
		selectTrade: (trade) => update(state => ({
			...state,
			selectedTrade: trade
		})),
		
		// Clear selected trade
		clearSelection: () => update(state => ({
			...state,
			selectedTrade: null
		})),
		
		// Get sectors list
		getSectors: () => SECTORS,
		
		// Filter trades by sector
		getTradesBySector: (sector) => {
			let trades = [];
			const unsubscribe = subscribe(state => trades = state.trades);
			unsubscribe();
			
			return trades.filter(trade => trade.sector === sector);
		},
		
		// Get trades expiring soon (within 7 days)
		getExpiringTrades: () => {
			let trades = [];
			const unsubscribe = subscribe(state => trades = state.trades);
			unsubscribe();
			
			const today = new Date();
			const sevenDaysFromNow = new Date();
			sevenDaysFromNow.setDate(today.getDate() + 7);
			
			return trades.filter(trade => {
				const expirationDate = new Date(trade.expiration_date);
				return expirationDate >= today && expirationDate <= sevenDaysFromNow;
			});
		}
	};
}

export const tradesStore = createTradesStore(); 