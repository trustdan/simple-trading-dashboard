package models

import (
	"fmt"
	"time"
)

// OptionsTrade represents an options trading position
type OptionsTrade struct {
	ID             int64     `json:"id"`
	Ticker         string    `json:"ticker"`
	Sector         string    `json:"sector"`
	StrategyType   string    `json:"strategy_type"`
	EntryDate      time.Time `json:"entry_date"`
	ExpirationDate time.Time `json:"expiration_date"`
	TargetPrice    *float64  `json:"target_price,omitempty"`
	StopLoss       *float64  `json:"stop_loss,omitempty"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TradeRequest represents the data structure for creating/updating trades
type TradeRequest struct {
	Ticker         string    `json:"ticker"`
	Sector         string    `json:"sector"`
	StrategyType   string    `json:"strategy_type"`
	EntryDate      time.Time `json:"entry_date"`
	ExpirationDate time.Time `json:"expiration_date"`
	TargetPrice    *float64  `json:"target_price,omitempty"`
	StopLoss       *float64  `json:"stop_loss,omitempty"`
	Notes          string    `json:"notes"`
}

// StrategyType represents an options trading strategy
type StrategyType struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	ColorHex    string `json:"color_hex"`
}

// TradeStatus represents possible trade statuses
const (
	StatusActive  = "active"
	StatusClosed  = "closed"
	StatusExpired = "expired"
)

// StrategyCategory represents strategy categories
const (
	CategoryBullish = "Bullish"
	CategoryBearish = "Bearish"
	CategoryNeutral = "Neutral"
)

// ValidateTradeRequest validates a trade request
func ValidateTradeRequest(req TradeRequest) error {
	if req.Ticker == "" {
		return fmt.Errorf("ticker is required")
	}
	if req.Sector == "" {
		return fmt.Errorf("sector is required")
	}
	if req.StrategyType == "" {
		return fmt.Errorf("strategy type is required")
	}
	if req.EntryDate.IsZero() {
		return fmt.Errorf("entry date is required")
	}
	if req.ExpirationDate.IsZero() {
		return fmt.Errorf("expiration date is required")
	}
	if req.ExpirationDate.Before(req.EntryDate) {
		return fmt.Errorf("expiration date must be after entry date")
	}
	return nil
}

// GetValidStatuses returns all valid trade statuses
func GetValidStatuses() []string {
	return []string{StatusActive, StatusClosed, StatusExpired}
}

// GetValidCategories returns all valid strategy categories
func GetValidCategories() []string {
	return []string{CategoryBullish, CategoryBearish, CategoryNeutral}
}
