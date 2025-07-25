package services

import (
	"database/sql"
	"fmt"
	"time"

	"trading-dashboard/backend/internal/models"
)

type TradeService struct {
	db *sql.DB
}

// NewTradeService creates a new trade service
func NewTradeService(db *sql.DB) *TradeService {
	return &TradeService{db: db}
}

// CreateTrade creates a new options trade
func (s *TradeService) CreateTrade(req models.TradeRequest) (*models.OptionsTrade, error) {
	if err := models.ValidateTradeRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	query := `
		INSERT INTO options_trades (
			ticker, sector, strategy_type, entry_date, expiration_date,
			target_price, stop_loss, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		req.Ticker,
		req.Sector,
		req.StrategyType,
		req.EntryDate,
		req.ExpirationDate,
		req.TargetPrice,
		req.StopLoss,
		req.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trade: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get trade ID: %w", err)
	}

	return s.GetTradeByID(id)
}

// GetTradeByID retrieves a trade by ID
func (s *TradeService) GetTradeByID(id int64) (*models.OptionsTrade, error) {
	query := `
		SELECT id, ticker, sector, strategy_type, entry_date, expiration_date,
		       target_price, stop_loss, status, notes, created_at, updated_at
		FROM options_trades
		WHERE id = ?
	`

	var trade models.OptionsTrade
	err := s.db.QueryRow(query, id).Scan(
		&trade.ID,
		&trade.Ticker,
		&trade.Sector,
		&trade.StrategyType,
		&trade.EntryDate,
		&trade.ExpirationDate,
		&trade.TargetPrice,
		&trade.StopLoss,
		&trade.Status,
		&trade.Notes,
		&trade.CreatedAt,
		&trade.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trade not found")
		}
		return nil, fmt.Errorf("failed to get trade: %w", err)
	}

	return &trade, nil
}

// GetTrades retrieves trades within a date range
func (s *TradeService) GetTrades(startDate, endDate time.Time) ([]models.OptionsTrade, error) {
	query := `
		SELECT id, ticker, sector, strategy_type, entry_date, expiration_date,
		       target_price, stop_loss, status, notes, created_at, updated_at
		FROM options_trades
		WHERE entry_date >= ? AND entry_date <= ?
		ORDER BY entry_date DESC, created_at DESC
	`

	rows, err := s.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query trades: %w", err)
	}
	defer rows.Close()

	var trades []models.OptionsTrade
	for rows.Next() {
		var trade models.OptionsTrade
		err := rows.Scan(
			&trade.ID,
			&trade.Ticker,
			&trade.Sector,
			&trade.StrategyType,
			&trade.EntryDate,
			&trade.ExpirationDate,
			&trade.TargetPrice,
			&trade.StopLoss,
			&trade.Status,
			&trade.Notes,
			&trade.CreatedAt,
			&trade.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade: %w", err)
		}
		trades = append(trades, trade)
	}

	return trades, rows.Err()
}

// GetActiveTradesByDateRange retrieves active trades for a specific date range (for calendar view)
func (s *TradeService) GetActiveTradesByDateRange(startDate, endDate time.Time) ([]models.OptionsTrade, error) {
	query := `
		SELECT id, ticker, sector, strategy_type, entry_date, expiration_date,
		       target_price, stop_loss, status, notes, created_at, updated_at
		FROM options_trades
		WHERE status = 'active' 
		  AND ((entry_date BETWEEN ? AND ?) 
		       OR (expiration_date BETWEEN ? AND ?)
		       OR (entry_date <= ? AND expiration_date >= ?))
		ORDER BY entry_date, ticker
	`

	rows, err := s.db.Query(query, startDate, endDate, startDate, endDate, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query active trades: %w", err)
	}
	defer rows.Close()

	var trades []models.OptionsTrade
	for rows.Next() {
		var trade models.OptionsTrade
		err := rows.Scan(
			&trade.ID,
			&trade.Ticker,
			&trade.Sector,
			&trade.StrategyType,
			&trade.EntryDate,
			&trade.ExpirationDate,
			&trade.TargetPrice,
			&trade.StopLoss,
			&trade.Status,
			&trade.Notes,
			&trade.CreatedAt,
			&trade.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade: %w", err)
		}
		trades = append(trades, trade)
	}

	return trades, rows.Err()
}

// UpdateTrade updates an existing trade
func (s *TradeService) UpdateTrade(id int64, req models.TradeRequest) (*models.OptionsTrade, error) {
	if err := models.ValidateTradeRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	query := `
		UPDATE options_trades SET
			ticker = ?, sector = ?, strategy_type = ?, entry_date = ?,
			expiration_date = ?, target_price = ?, stop_loss = ?, notes = ?
		WHERE id = ?
	`

	result, err := s.db.Exec(
		query,
		req.Ticker,
		req.Sector,
		req.StrategyType,
		req.EntryDate,
		req.ExpirationDate,
		req.TargetPrice,
		req.StopLoss,
		req.Notes,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update trade: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("trade not found")
	}

	return s.GetTradeByID(id)
}

// UpdateTradeStatus updates the status of a trade
func (s *TradeService) UpdateTradeStatus(id int64, status string) (*models.OptionsTrade, error) {
	validStatuses := models.GetValidStatuses()
	valid := false
	for _, s := range validStatuses {
		if s == status {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	query := `UPDATE options_trades SET status = ? WHERE id = ?`
	result, err := s.db.Exec(query, status, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update trade status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("trade not found")
	}

	return s.GetTradeByID(id)
}

// DeleteTrade deletes a trade
func (s *TradeService) DeleteTrade(id int64) error {
	query := `DELETE FROM options_trades WHERE id = ?`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete trade: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("trade not found")
	}

	return nil
}

// GetStrategyTypes retrieves all available strategy types
func (s *TradeService) GetStrategyTypes() ([]models.StrategyType, error) {
	query := `
		SELECT id, name, category, description, color_hex
		FROM strategy_types
		ORDER BY category, name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query strategy types: %w", err)
	}
	defer rows.Close()

	var strategies []models.StrategyType
	for rows.Next() {
		var strategy models.StrategyType
		err := rows.Scan(
			&strategy.ID,
			&strategy.Name,
			&strategy.Category,
			&strategy.Description,
			&strategy.ColorHex,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan strategy: %w", err)
		}
		strategies = append(strategies, strategy)
	}

	return strategies, rows.Err()
}
