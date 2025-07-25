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

// getAppDataDir returns the appropriate app data directory for the current OS
func getAppDataDir() (string, error) {
	var baseDir string

	// Check for Windows AppData
	if appData := os.Getenv("APPDATA"); appData != "" {
		baseDir = appData
	} else {
		// Fallback to user home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		baseDir = homeDir
	}

	appDir := filepath.Join(baseDir, "TradingDashboard")
	return appDir, nil
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Get app data directory
	appDataDir, err := getAppDataDir()
	if err != nil {
		log.Printf("Warning: Failed to get app data directory, using current directory: %v", err)
		appDataDir = "."
	}

	// Initialize database
	dbPath := filepath.Join(appDataDir, "trading_dashboard.db")
	log.Printf("Initializing database at: %s", dbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize schema
	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
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
	return a.marketService.SaveRating(req)
}

// GetLatestMarketRating retrieves the most recent market rating
func (a *App) GetLatestMarketRating() (*models.MarketRating, error) {
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
	return a.tradeService.CreateTrade(req)
}

// GetTradeByID retrieves a trade by ID
func (a *App) GetTradeByID(id int64) (*models.OptionsTrade, error) {
	return a.tradeService.GetTradeByID(id)
}

// GetTrades retrieves trades within a date range
func (a *App) GetTrades(startDate, endDate time.Time) ([]models.OptionsTrade, error) {
	return a.tradeService.GetTrades(startDate, endDate)
}

// GetActiveTradesByDateRange retrieves active trades for calendar view
func (a *App) GetActiveTradesByDateRange(startDate, endDate time.Time) ([]models.OptionsTrade, error) {
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
	return a.tradeService.GetStrategyTypes()
}
