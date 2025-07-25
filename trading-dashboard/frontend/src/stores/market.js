import { writable } from 'svelte/store';

// Market sectors list
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

// Create the market store
function createMarketStore() {
	const { subscribe, set, update } = writable({
		overallRating: 0,
		sectorRatings: {},
		sectors: SECTORS,
		loading: false,
		lastSaved: null
	});

	return {
		subscribe,
		
		// Set overall market rating
		setOverallRating: (rating) => update(state => ({
			...state,
			overallRating: rating
		})),
		
		// Set individual sector rating
		setSectorRating: (sector, rating) => update(state => ({
			...state,
			sectorRatings: {
				...state.sectorRatings,
				[sector]: rating
			}
		})),
		
		// Load latest rating from backend
		loadLatestRating: async () => {
			update(state => ({ ...state, loading: true }));
			
			try {
				// Add a timeout to prevent infinite loading
				const timeoutPromise = new Promise((_, reject) => 
					setTimeout(() => reject(new Error('Request timeout')), 10000)
				);
				
				// Call the Wails backend API with timeout
				const rating = await Promise.race([
					window.go.main.App.GetLatestMarketRating(),
					timeoutPromise
				]);
				
				// Check if this is a valid rating or a default fallback
				const isDefaultRating = rating.id === 0 && rating.overall_rating === 0;
				
				update(state => ({
					...state,
					overallRating: rating.overall_rating || 0,
					sectorRatings: rating.sector_ratings || {},
					lastSaved: rating.created_at,
					loading: false
				}));
				
				// Only log success for real data, not defaults
				if (!isDefaultRating) {
					console.log('Market rating loaded successfully:', rating);
				} else {
					console.log('Using default market rating - database may not be initialized');
				}
			} catch (error) {
				console.error('Failed to load latest rating:', error);
				update(state => ({ 
					...state, 
					loading: false,
					overallRating: 0,
					sectorRatings: {}
				}));
				throw error;
			}
		},
		
		// Save current ratings to backend
		saveRating: async () => {
			update(state => ({ ...state, loading: true }));
			
			try {
				let currentState;
				const unsubscribe = subscribe(state => currentState = state);
				unsubscribe();
				
				const request = {
					overall_rating: currentState.overallRating,
					sector_ratings: currentState.sectorRatings
				};
				
				const savedRating = await window.go.main.App.SaveMarketRating(request);
				
				update(state => ({
					...state,
					lastSaved: savedRating.created_at,
					loading: false
				}));
				
				return savedRating;
			} catch (error) {
				console.error('Failed to save rating:', error);
				update(state => ({ ...state, loading: false }));
				throw error;
			}
		},
		
		// Reset all ratings to zero
		resetRatings: () => update(state => ({
			...state,
			overallRating: 0,
			sectorRatings: {}
		})),
		
		// Get sectors list
		getSectors: () => SECTORS
	};
}

export const marketStore = createMarketStore(); 