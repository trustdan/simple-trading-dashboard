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

-- Options trading tables
CREATE TABLE options_trades (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticker TEXT NOT NULL,
    sector TEXT NOT NULL,
    strategy_type TEXT NOT NULL,
    entry_date DATE NOT NULL,
    expiration_date DATE NOT NULL,
    target_price DECIMAL(10,2),
    stop_loss DECIMAL(10,2),
    status TEXT DEFAULT 'active' CHECK (status IN ('active', 'closed', 'expired')),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Strategy definitions table
CREATE TABLE strategy_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    category TEXT NOT NULL,
    description TEXT,
    color_hex TEXT DEFAULT '#4a90e2'
);

-- Insert default strategy types
INSERT INTO strategy_types (name, category, description, color_hex) VALUES
-- Bullish Strategies
('Long Call', 'Bullish', 'Buy call options expecting price increase', '#22c55e'),
('Bull Call Spread', 'Bullish', 'Buy lower strike call, sell higher strike call', '#16a34a'),
('Cash-Secured Put', 'Bullish', 'Sell puts with cash backing to acquire shares', '#15803d'),

-- Bearish Strategies  
('Long Put', 'Bearish', 'Buy put options expecting price decrease', '#ef4444'),
('Bear Put Spread', 'Bearish', 'Buy higher strike put, sell lower strike put', '#dc2626'),
('Covered Call', 'Bearish', 'Sell calls against owned shares', '#b91c1c'),

-- Neutral Strategies
('Iron Condor', 'Neutral', 'Sell call and put spreads for range-bound profit', '#8b5cf6'),
('Butterfly Spread', 'Neutral', 'Limited risk/reward for minimal price movement', '#7c3aed'),
('Straddle', 'Neutral', 'Buy call and put at same strike for volatility play', '#6d28d9');

-- Indexes for trades table
CREATE INDEX idx_trades_ticker ON options_trades(ticker);
CREATE INDEX idx_trades_sector ON options_trades(sector);
CREATE INDEX idx_trades_status ON options_trades(status);
CREATE INDEX idx_trades_entry_date ON options_trades(entry_date);
CREATE INDEX idx_trades_expiration_date ON options_trades(expiration_date);
CREATE INDEX idx_trades_strategy ON options_trades(strategy_type);

-- Trigger to update trades updated_at timestamp
CREATE TRIGGER update_trades_timestamp 
    AFTER UPDATE ON options_trades
BEGIN
    UPDATE options_trades SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END; 