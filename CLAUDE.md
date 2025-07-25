# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an **Options Trading Dashboard** - a Windows desktop application built with Wails v2, Go backend, and Svelte frontend. The project is currently in planning phase with no code implemented yet.

## Development Commands

### Initial Setup
```bash
# Install Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Install pnpm (preferred over npm to avoid node_modules issues)
npm install -g pnpm

# Initialize the project
wails init -n trading-dashboard -t vanilla

# Setup backend
cd backend && go mod init trading-dashboard

# Setup frontend
cd frontend && pnpm create vite . --template svelte && pnpm install
```

### Development
```bash
# Run in development mode
wails dev

# Build for production
cd frontend && pnpm run build
cd .. && wails build -platform windows/amd64 -clean
```

## Architecture

- **Desktop Framework**: Wails v2 (creates native Windows executable)
- **Backend**: Go with Gin/Echo REST API, SQLite for local storage
- **Frontend**: Svelte (plain JavaScript, NO TypeScript), custom CSS
- **Database**: SQLite with GORM or sqlx
- **State Management**: Svelte stores with API integration

### Project Structure
```
trading-dashboard/
├── backend/
│   ├── cmd/server/      # Main application entry
│   ├── internal/
│   │   ├── api/         # REST endpoints
│   │   ├── models/      # Data models
│   │   ├── services/    # Business logic
│   │   └── database/    # Database layer
│   └── pkg/utils/       # Shared utilities
├── frontend/
│   ├── src/
│   │   ├── components/  # UI components (dials, grids)
│   │   ├── stores/      # Svelte stores
│   │   ├── services/    # API client
│   │   └── utils/       # Frontend utilities
│   └── public/
└── PROGRESS_LOG.md      # MANDATORY progress tracking
```

## Critical Development Rules

1. **Progress Tracking**: Update `PROGRESS_LOG.md` after every significant change with format: `[YYYY-MM-DD HH:MM] Feature: Description`

2. **No TypeScript**: Use plain JavaScript with JSDoc comments for documentation

3. **API Format**: All endpoints under `/api/v1/` with RESTful conventions

4. **Error Handling**: Never ignore errors, always provide user-friendly messages

5. **Performance**: Target 60fps animations, sub-second API responses

## UI Specifications

### Dial Components
- **Size**: 200x200px (sectors), 300x300px (overall market)
- **Scale**: -3 to +3 with 0.1 increments
- **Animation**: 300ms ease-in-out transitions
- **Colors**: Red (-3) to Yellow (0) to Green (+3) gradient
- **Theme**: Dark background (#0a0a0a) with subtle lighting effects

### Grid Layout
- 8px spacing grid system
- Responsive with CSS Grid
- Dark theme with high contrast for trading environment

## Core Features

1. **Market Sentiment Dials**: 11 sectors + overall market with persistence
2. **Options Trade Grid**: Manage trades with ticker, strategy, dates, targets
3. **Strategy Categories**: 24 strategies organized in 8 color-coded categories
4. **Auto-save**: Every user action persisted to SQLite

## Windows-Specific Notes

- Use `pnpm` instead of `npm` to avoid node_modules path issues
- Alternative: Build frontend on WSL2 if needed
- May need to set Windows max path length
- Use `npx rimraf node_modules` for cleanup

## Testing Requirements

- Go: Table-driven tests for all services
- API: Test all endpoints
- Frontend: Component tests for complex UI logic
- Target: >70% code coverage for backend services