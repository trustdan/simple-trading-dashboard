package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"trading-dashboard/pkg/database"
	"trading-dashboard/pkg/models"
	"trading-dashboard/pkg/services"
)

// App struct
type App struct {
	ctx           context.Context
	db            *database.DB
	marketService *services.MarketService
	tradeService  *services.TradeService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// getConsistentDataDir returns a consistent data directory that works for both standalone and installed versions
func getConsistentDataDir() (string, error) {
	// Strategy 1: Try to use the executable's directory (works for portable/standalone)
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		dataDir := filepath.Join(exeDir, "data")

		// Test if we can write to this directory
		testFile := filepath.Join(dataDir, ".write_test")
		if err := os.MkdirAll(dataDir, 0755); err == nil {
			if file, err := os.Create(testFile); err == nil {
				file.Close()
				os.Remove(testFile)
				log.Printf("Using executable directory for data: %s", dataDir)
				return dataDir, nil
			}
		}
	}

	// Strategy 2: Use Windows AppData (works for installed versions)
	if appData := os.Getenv("APPDATA"); appData != "" {
		appDir := filepath.Join(appData, "TradingDashboard")
		log.Printf("Using AppData directory for data: %s", appDir)
		return appDir, nil
	}

	// Strategy 3: Use user home directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		appDir := filepath.Join(homeDir, "TradingDashboard")
		log.Printf("Using home directory for data: %s", appDir)
		return appDir, nil
	}

	return "", fmt.Errorf("unable to determine suitable data directory")
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Get a consistent data directory that works for both standalone and installed versions
	dataDir, err := getConsistentDataDir()
	if err != nil {
		log.Printf("Warning: Failed to get data directory: %v", err)
		// Fallback to current working directory
		if cwd, err := os.Getwd(); err == nil {
			dataDir = filepath.Join(cwd, "data")
		} else {
			dataDir = "data"
		}
	}
	log.Printf("Using data directory: %s", dataDir)

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("Warning: Failed to create data directory at %s: %v", dataDir, err)
		// Try fallback directory
		if cwd, err := os.Getwd(); err == nil {
			dataDir = cwd
		} else {
			dataDir = "."
		}
	}

	// Initialize database
	dbPath := filepath.Join(dataDir, "trading_dashboard.db")
	log.Printf("Initializing database at: %s", dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		// Still try to initialize services with nil database for graceful failure
		a.db = nil
		a.marketService = nil
		a.tradeService = nil
		return
	}

	// Initialize schema with better error handling
	log.Println("Initializing database schema...")
	if err := db.InitSchema(); err != nil {
		log.Printf("Failed to initialize schema: %v", err)
		// Close the database and set services to nil
		db.Close()
		a.db = nil
		a.marketService = nil
		a.tradeService = nil
		return
	}

	a.db = db
	a.marketService = services.NewMarketService(db.DB)
	a.tradeService = services.NewTradeService(db.DB)

	log.Println("Trading Dashboard initialized successfully")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SaveMarketRating saves a new market rating with sector ratings
func (a *App) SaveMarketRating(req models.MarketRatingRequest) (*models.MarketRating, error) {
	if a.marketService == nil {
		return nil, fmt.Errorf("market service not available - database connection failed")
	}
	return a.marketService.SaveRating(req)
}

// GetLatestMarketRating retrieves the most recent market rating
func (a *App) GetLatestMarketRating() (*models.MarketRating, error) {
	if a.marketService == nil {
		log.Printf("Market service not initialized - database connection failed")
		// Return a default rating instead of failing
		return &models.MarketRating{
			ID:            0,
			OverallRating: 0,
			SectorRatings: make(map[string]float64),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}, nil
	}
	return a.marketService.GetLatestRating()
}

// UpdateMarketRating updates an existing market rating
func (a *App) UpdateMarketRating(id int64, req models.MarketRatingRequest) (*models.MarketRating, error) {
	return a.marketService.UpdateRating(id, req)
}

// GetSectorNames returns the list of available market sectors
func (a *App) GetSectorNames() []string {
	return a.marketService.GetSectorNames()
}

// Close closes the database connection (called on app shutdown)
func (a *App) Close() {
	if a.db != nil {
		a.db.Close()
	}
}

// ============ TRADE API METHODS ============

// CreateTrade creates a new options trade
func (a *App) CreateTrade(req models.TradeRequest) (*models.OptionsTrade, error) {
	if a.tradeService == nil {
		return nil, fmt.Errorf("trade service not available - database connection failed")
	}
	return a.tradeService.CreateTrade(req)
}

// GetTradeByID retrieves a trade by ID
func (a *App) GetTradeByID(id int64) (*models.OptionsTrade, error) {
	return a.tradeService.GetTradeByID(id)
}

// GetTrades retrieves trades within a date range
func (a *App) GetTrades(startDate, endDate time.Time) ([]models.OptionsTrade, error) {
	if a.tradeService == nil {
		log.Printf("Trade service not initialized - database connection failed")
		// Return empty list instead of failing
		return []models.OptionsTrade{}, nil
	}
	return a.tradeService.GetTrades(startDate, endDate)
}

// GetActiveTradesByDateRange retrieves active trades for calendar view
func (a *App) GetActiveTradesByDateRange(startDate, endDate time.Time) ([]models.OptionsTrade, error) {
	if a.tradeService == nil {
		log.Printf("Trade service not initialized - database connection failed")
		// Return empty list instead of failing
		return []models.OptionsTrade{}, nil
	}
	return a.tradeService.GetActiveTradesByDateRange(startDate, endDate)
}

// UpdateTrade updates an existing trade
func (a *App) UpdateTrade(id int64, req models.TradeRequest) (*models.OptionsTrade, error) {
	return a.tradeService.UpdateTrade(id, req)
}

// UpdateTradeStatus updates the status of a trade
func (a *App) UpdateTradeStatus(id int64, status string) (*models.OptionsTrade, error) {
	return a.tradeService.UpdateTradeStatus(id, status)
}

// DeleteTrade deletes a trade
func (a *App) DeleteTrade(id int64) error {
	return a.tradeService.DeleteTrade(id)
}

// GetStrategyTypes retrieves all available strategy types
func (a *App) GetStrategyTypes() ([]models.StrategyType, error) {
	if a.tradeService == nil {
		log.Printf("Trade service not initialized - database connection failed")
		// Return empty list instead of failing
		return []models.StrategyType{}, nil
	}
	return a.tradeService.GetStrategyTypes()
}
