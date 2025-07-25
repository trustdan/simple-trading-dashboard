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
-- Basic Strategies
('Long Call', 'Basic', 'Buy call options expecting price increase', '#3b82f6'),
('Long Put', 'Basic', 'Buy put options expecting price decrease', '#60a5fa'),

-- Income Strategies
('Covered Call', 'Income', 'Sell calls against owned shares for premium income', '#22c55e'),
('Cash-Secured Put', 'Income', 'Sell puts with cash backing to generate income', '#16a34a'),

-- Credit Spreads
('Bull Put Spread', 'Credit Spreads', 'Sell higher strike put, buy lower strike put for credit', '#f97316'),
('Bear Call Spread', 'Credit Spreads', 'Sell lower strike call, buy higher strike call for credit', '#ea580c'),

-- Neutral Strategies
('Iron Butterfly', 'Neutral', 'Sell ATM call and put, buy OTM wings for range-bound profit', '#8b5cf6'),
('Iron Condor', 'Neutral', 'Sell call and put spreads for range-bound profit', '#7c3aed'),
('Long Put Butterfly', 'Neutral', 'Buy two puts at middle strike, sell one each at higher and lower strikes', '#6d28d9'),
('Long Call Butterfly', 'Neutral', 'Buy two calls at middle strike, sell one each at higher and lower strikes', '#5b21b6'),

-- Calendar Spreads
('Calendar Call Spread', 'Calendar Spreads', 'Sell near-term call, buy longer-term call at same strike', '#14b8a6'),
('Calendar Put Spread', 'Calendar Spreads', 'Sell near-term put, buy longer-term put at same strike', '#0d9488'),
('Diagonal Call Spread', 'Calendar Spreads', 'Sell near-term call, buy longer-term call at different strike', '#0f766e'),
('Diagonal Put Spread', 'Calendar Spreads', 'Sell near-term put, buy longer-term put at different strike', '#134e4a'),

-- Debit Spreads
('Bull Call Spread', 'Debit Spreads', 'Buy lower strike call, sell higher strike call', '#1d4ed8'),
('Bear Put Spread', 'Debit Spreads', 'Buy higher strike put, sell lower strike put', '#1e40af'),

-- Directional Strategies
('Straddle', 'Directional', 'Buy call and put at same strike for volatility play', '#dc2626'),
('Strangle', 'Directional', 'Buy OTM call and put for volatility play', '#b91c1c'),

-- Ratio Spreads
('Call Ratio Backspread', 'Ratio Spreads', 'Sell fewer ITM calls, buy more OTM calls', '#92400e'),
('Put Broken Wing', 'Ratio Spreads', 'Modified put butterfly with uneven wings', '#a16207'),
('Inverse Call Broken Wing', 'Ratio Spreads', 'Modified call butterfly with inverted risk profile', '#ca8a04'),
('Put Ratio Backspread', 'Ratio Spreads', 'Sell fewer ITM puts, buy more OTM puts', '#eab308'),
('Call Broken Wing', 'Ratio Spreads', 'Modified call butterfly with uneven wings', '#facc15'),
('Inverse Put Broken Wing', 'Ratio Spreads', 'Modified put butterfly with inverted risk profile', '#fde047');

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