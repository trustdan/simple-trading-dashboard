package models

import (
	"time"
)

// MarketRating represents the overall market sentiment and associated sector ratings
type MarketRating struct {
	ID            int64              `json:"id"`
	OverallRating float64            `json:"overall_rating"`
	SectorRatings map[string]float64 `json:"sector_ratings"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

// SectorRating represents a single sector's sentiment rating
type SectorRating struct {
	ID             int64     `json:"id"`
	MarketRatingID int64     `json:"market_rating_id"`
	SectorName     string    `json:"sector_name"`
	Rating         float64   `json:"rating"`
	CreatedAt      time.Time `json:"created_at"`
}

// MarketRatingRequest represents the data structure for creating/updating market ratings
type MarketRatingRequest struct {
	OverallRating float64            `json:"overall_rating"`
	SectorRatings map[string]float64 `json:"sector_ratings"`
}

// ValidateRating checks if a rating is within the valid range (-3 to +3)
func ValidateRating(rating float64) bool {
	return rating >= -3.0 && rating <= 3.0
}

// GetSectorNames returns the standard list of market sectors
func GetSectorNames() []string {
	return []string{
		"Basic Materials",
		"Communication Services",
		"Consumer Cyclical",
		"Consumer Defensive",
		"Energy",
		"Financial Services",
		"Healthcare",
		"Industrials",
		"Real Estate",
		"Technology",
		"Utilities",
	}
}
