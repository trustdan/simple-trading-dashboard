# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Options Trading Dashboard** - A Windows desktop application built with:
- **Desktop Framework**: Wails v2 (creates native Windows executable)
- **Backend**: Go with SQLite for local storage
- **Frontend**: Svelte (plain JavaScript, NO TypeScript)
- **Database**: SQLite with database/sql package

## Development Commands

### Running in Development Mode
```bash
cd trading-dashboard
wails dev
```

### Building for Production
```bash
# From project root - builds both exe and installer
./build.ps1

# Or manually from trading-dashboard directory
wails build --clean --nsis
```

### Frontend Development
```bash
cd trading-dashboard/frontend
npm install  # or pnpm install
npm run dev  # Vite dev server
npm run build  # Production build
```

## Architecture

### Directory Structure
```
trading-dashboard/           # Main application directory
├── app.go                  # Main application struct with API methods
├── main.go                 # Entry point
├── wails.json             # Wails configuration
├── frontend/              # Svelte frontend
│   ├── src/
│   │   ├── components/    # UI components (Dial, TradeCell, etc.)
│   │   ├── stores/       # Svelte stores (market.js, trades.js)
│   │   └── App.svelte    # Main app component
│   └── package.json      # Uses node_modules\.bin\vite.cmd on Windows
├── pkg/                   # Go packages
│   ├── database/         # Database connection and schema
│   ├── models/           # Data models (MarketRating, OptionsTrade)
│   └── services/         # Business logic (MarketService, TradeService)
└── build/bin/            # Build output directory
```

### Key API Methods (app.go)

**Market Rating API:**
- `SaveMarketRating(req models.MarketRatingRequest)` - Save market sentiment
- `GetLatestMarketRating()` - Get current market ratings
- `GetSectorNames()` - Get list of 11 market sectors

**Trading API:**
- `CreateTrade(req models.TradeRequest)` - Create new options trade
- `GetTrades(startDate, endDate time.Time)` - Get trades by date range
- `UpdateTrade(id int64, req models.TradeRequest)` - Update trade
- `UpdateTradeStatus(id int64, status string)` - Quick status update
- `DeleteTrade(id int64)` - Delete trade
- `GetStrategyTypes()` - Get 24 strategy types in 8 categories

### Database Schema

Located in `pkg/database/schema.sql`:
- `market_ratings` - Overall market sentiment (-3 to +3)
- `sector_ratings` - Individual sector ratings
- `options_trades` - Options trading records
- `strategy_types` - 24 predefined strategies

### Frontend State Management

Svelte stores in `frontend/src/stores/`:
- `market.js` - Market rating state and API calls
- `trades.js` - Trade management and filtering
- `toast.js` - Notification system

## Critical Development Rules

1. **No TypeScript** - Use plain JavaScript with JSDoc comments
2. **Windows Paths** - Use `node_modules\.bin\vite.cmd` in package.json scripts
3. **Error Handling** - Database operations gracefully fail with default values
4. **Data Directory** - App tries multiple locations for SQLite database:
   - Executable directory (portable mode)
   - %APPDATA%\TradingDashboard (installed mode)
   - User home directory (fallback)

## UI Specifications

### Market Sentiment Dials
- Interactive SVG dials with -3 to +3 scale
- Red-Yellow-Green gradient
- Mouse wheel, click-drag, and keyboard support
- 11 sectors + 1 overall market dial

### Trade Grid
- 3-week calendar view with Monday start
- Color-coded by strategy category
- Status badges (active, closed, expired, etc.)
- Right-click context menu for quick actions

### Modals and Forms
- TradeModal for create/edit with validation
- Date pickers, strategy dropdowns, target inputs
- Real-time validation feedback

## Building and Distribution

1. **Development**: Use `wails dev` for hot-reload development
2. **Production Build**: Run `build.ps1` from project root
3. **Output Files**:
   - `trading-dashboard.exe` - Standalone executable (~12MB)
   - `trading-dashboard-amd64-installer.exe` - NSIS installer (~6MB)

## Common Tasks

### Adding a New API Endpoint
1. Add method to `app.go` with proper error handling
2. Method will be auto-exposed to frontend via Wails
3. Call from frontend using `window.go.main.App.MethodName()`

### Adding a New Component
1. Create `.svelte` file in `frontend/src/components/`
2. Follow existing patterns (see Dial.svelte or TradeCell.svelte)
3. Use CSS-in-component styling

### Modifying Database Schema
1. Update `pkg/database/schema.sql`
2. Update models in `pkg/models/`
3. Update service methods in `pkg/services/`
4. Database auto-migrates on startup

## Testing

Currently no automated tests. To test:
1. Run `wails dev` for development mode
2. Check browser console for frontend errors
3. Check terminal for backend errors
4. Database file is created in data directory

## Progress Tracking

See `PROGRESS_LOG.md` for detailed development history and feature implementation timeline.