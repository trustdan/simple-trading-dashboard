package services

import (
	"database/sql"
	"fmt"
	"time"

	"trading-dashboard/backend/internal/models"
)

type MarketService struct {
	db *sql.DB
}

// NewMarketService creates a new market service
func NewMarketService(db *sql.DB) *MarketService {
	return &MarketService{db: db}
}

// SaveRating saves a new market rating with sector ratings
func (s *MarketService) SaveRating(req models.MarketRatingRequest) (*models.MarketRating, error) {
	// Validate overall rating
	if !models.ValidateRating(req.OverallRating) {
		return nil, fmt.Errorf("invalid overall rating: %f (must be between -3 and 3)", req.OverallRating)
	}

	// Validate sector ratings
	for sector, rating := range req.SectorRatings {
		if !models.ValidateRating(rating) {
			return nil, fmt.Errorf("invalid rating for sector %s: %f (must be between -3 and 3)", sector, rating)
		}
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert market rating
	result, err := tx.Exec(
		"INSERT INTO market_ratings (overall_rating) VALUES (?)",
		req.OverallRating,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert market rating: %w", err)
	}

	marketRatingID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get market rating ID: %w", err)
	}

	// Insert sector ratings
	for sector, rating := range req.SectorRatings {
		_, err := tx.Exec(
			"INSERT INTO sector_ratings (market_rating_id, sector_name, rating) VALUES (?, ?, ?)",
			marketRatingID, sector, rating,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert sector rating for %s: %w", sector, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return the saved rating
	return s.GetRatingByID(marketRatingID)
}

// GetLatestRating retrieves the most recent market rating
func (s *MarketService) GetLatestRating() (*models.MarketRating, error) {
	row := s.db.QueryRow(`
		SELECT id, overall_rating, created_at, updated_at 
		FROM market_ratings 
		ORDER BY created_at DESC 
		LIMIT 1
	`)

	var rating models.MarketRating
	err := row.Scan(&rating.ID, &rating.OverallRating, &rating.CreatedAt, &rating.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return default rating if none exists
			return &models.MarketRating{
				ID:            0,
				OverallRating: 0,
				SectorRatings: make(map[string]float64),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("failed to get latest rating: %w", err)
	}

	// Get sector ratings
	sectorRatings, err := s.getSectorRatings(rating.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sector ratings: %w", err)
	}

	rating.SectorRatings = sectorRatings
	return &rating, nil
}

// GetRatingByID retrieves a market rating by ID
func (s *MarketService) GetRatingByID(id int64) (*models.MarketRating, error) {
	row := s.db.QueryRow(`
		SELECT id, overall_rating, created_at, updated_at 
		FROM market_ratings 
		WHERE id = ?
	`, id)

	var rating models.MarketRating
	err := row.Scan(&rating.ID, &rating.OverallRating, &rating.CreatedAt, &rating.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating by ID: %w", err)
	}

	// Get sector ratings
	sectorRatings, err := s.getSectorRatings(rating.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sector ratings: %w", err)
	}

	rating.SectorRatings = sectorRatings
	return &rating, nil
}

// UpdateRating updates an existing market rating
func (s *MarketService) UpdateRating(id int64, req models.MarketRatingRequest) (*models.MarketRating, error) {
	// Validate ratings
	if !models.ValidateRating(req.OverallRating) {
		return nil, fmt.Errorf("invalid overall rating: %f", req.OverallRating)
	}

	for sector, rating := range req.SectorRatings {
		if !models.ValidateRating(rating) {
			return nil, fmt.Errorf("invalid rating for sector %s: %f", sector, rating)
		}
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update market rating
	_, err = tx.Exec(
		"UPDATE market_ratings SET overall_rating = ? WHERE id = ?",
		req.OverallRating, id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update market rating: %w", err)
	}

	// Delete existing sector ratings
	_, err = tx.Exec("DELETE FROM sector_ratings WHERE market_rating_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete existing sector ratings: %w", err)
	}

	// Insert new sector ratings
	for sector, rating := range req.SectorRatings {
		_, err := tx.Exec(
			"INSERT INTO sector_ratings (market_rating_id, sector_name, rating) VALUES (?, ?, ?)",
			id, sector, rating,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert sector rating for %s: %w", sector, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return s.GetRatingByID(id)
}

// getSectorRatings retrieves sector ratings for a market rating
func (s *MarketService) getSectorRatings(marketRatingID int64) (map[string]float64, error) {
	rows, err := s.db.Query(`
		SELECT sector_name, rating 
		FROM sector_ratings 
		WHERE market_rating_id = ?
	`, marketRatingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sectorRatings := make(map[string]float64)
	for rows.Next() {
		var sector string
		var rating float64
		if err := rows.Scan(&sector, &rating); err != nil {
			return nil, err
		}
		sectorRatings[sector] = rating
	}

	return sectorRatings, rows.Err()
}

// GetSectorNames returns the list of available sectors
func (s *MarketService) GetSectorNames() []string {
	return models.GetSectorNames()
}
