-- Market sentiment database schema
-- This stores overall market ratings and sector-specific ratings

CREATE TABLE market_ratings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    overall_rating REAL CHECK (overall_rating >= -3 AND overall_rating <= 3),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sector_ratings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    market_rating_id INTEGER,
    sector_name TEXT NOT NULL,
    rating REAL CHECK (rating >= -3 AND rating <= 3),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (market_rating_id) REFERENCES market_ratings(id) ON DELETE CASCADE
);

-- Index for faster queries
CREATE INDEX idx_market_ratings_created_at ON market_ratings(created_at DESC);
CREATE INDEX idx_sector_ratings_market_id ON sector_ratings(market_rating_id);
CREATE INDEX idx_sector_ratings_sector ON sector_ratings(sector_name);

-- Trigger to update the updated_at timestamp
CREATE TRIGGER update_market_ratings_timestamp 
    AFTER UPDATE ON market_ratings
BEGIN
    UPDATE market_ratings SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END; 