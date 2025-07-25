# Options Trading Dashboard - Project Overview

## Executive Summary

A Windows-based options trading dashboard designed for top-down market analysis and weekly basket trading management. The application features no actual trading functionality but serves as an analytical and organizational tool for options traders.

## Core Purpose

The dashboard enables traders to:
- Assess overall market conditions using a sentiment rating system
- Evaluate individual sector performance and sentiment
- Plan and track options trades across multiple sectors
- Maintain a visual representation of trade distribution over time

## Technology Stack

### Backend
- **Language**: Go
- **Purpose**: API server, data management, business logic
- **Key Libraries**: 
  - Gin or Echo (REST API framework)
  - GORM or sqlx (database ORM/driver)
  - SQLite (local database)

### Frontend
- **Framework**: Svelte (Plain JavaScript)
- **Styling**: Custom CSS with modern design
- **UI Components**: 
  - Custom SVG dials/gauges for ratings
  - CSS Grid for trade management
  - Smooth animations and transitions
- **State Management**: Svelte stores
- **HTTP Client**: Native fetch API

### Desktop Integration
- **Framework**: Wails v2
- **Purpose**: Create native Windows executable with embedded frontend
- **Benefits**: Single executable, native performance, no browser required

## Key Features

### 1. Market Sentiment Analysis (Tab 1)
- Overall market rating with animated dial gauge (-3 to +3)
- Individual sector rating dials arranged in a visually appealing grid
- Real-time dial animation on adjustment
- Color-coded dial indicators (red for bearish, green for bullish, yellow for neutral)
- Visual feedback with smooth transitions
- Persistent storage of ratings with timestamps

### 2. Options Trade Management (Tab 2)
- Grid-based trade visualization
- Sector-based organization (Y-axis)
- Date-based timeline (X-axis, ~3 weeks out)
- Trade entry with comprehensive details
- Support for 24 different options strategies across 8 categories
- Color-coded by strategy category for visual clarity
- Target and stop loss tracking

## Sector Classifications

1. Basic Materials
2. Communication Services
3. Consumer Cyclical
4. Consumer Defensive
5. Energy
6. Financial
7. Healthcare
8. Industrials
9. Real Estate
10. Technology
11. Utilities

## Data Model Highlights

### Market Ratings
- Timestamp-based entries
- Overall market sentiment score
- Individual sector scores
- Historical tracking capability

### Options Trades
- Ticker symbol
- Strategy type with color coding (see categories below)
- Entry and expiration dates
- Target price
- Stop loss amount
- Sector association
- Status tracking (active/closed/expired)

## Options Strategy Categories

### Basic Strategies
- Long Call
- Long Put

### Income Strategies
- Covered Call
- Cash-Secured Put

### Credit Spreads
- Bull Put Spread
- Bear Call Spread

### Neutral Strategies
- Iron Butterfly
- Iron Condor
- Long Put Butterfly
- Long Call Butterfly

### Calendar Spreads
- Calendar Call Spread
- Calendar Put Spread
- Diagonal Call Spread
- Diagonal Put Spread

### Debit Spreads
- Bull Call Spread
- Bear Put Spread

### Directional Strategies
- Straddle
- Strangle

### Ratio Spreads
- Call Ratio Backspread
- Put Broken Wing
- Inverse Call Broken Wing
- Put Ratio Backspread
- Call Broken Wing
- Inverse Put Broken Wing

## User Experience Goals

1. **Visual Appeal**: Beautiful dial gauges with smooth animations and modern design
2. **Intuitive Interaction**: Click and drag on dials for natural adjustment
3. **Speed**: Quick data entry and retrieval with responsive animations
4. **Clarity**: Clean visual hierarchy with well-organized dial layouts
5. **Professional Look**: Dashboard that looks like a high-end trading platform
6. **Feedback**: Visual and tactile-like feedback when adjusting dials

## Non-Functional Requirements

- Windows desktop application
- Local data storage (no cloud dependency)
- Responsive UI design with smooth animations
- Fast startup and operation
- Minimal resource usage
- Simplified build process to avoid node_modules issues
- Data export capabilities (future enhancement)

## Visual Design Philosophy

The dashboard embraces a modern, professional aesthetic inspired by high-end trading platforms and luxury car dashboards:

- **Dark Theme**: Deep black background (#0a0a0a) with subtle surface variations
- **Dial Design**: Inspired by analog gauges with smooth gradients and precise tick marks
- **Color Psychology**: Red-to-green spectrum for intuitive market sentiment
- **Animations**: Smooth 300ms transitions for all interactions
- **Lighting Effects**: Subtle glows and shadows for depth and tactile feel
- **Typography**: Clean, monospaced fonts for numbers, sans-serif for labels

## Success Metrics

- Sub-second response times for all operations
- Zero data loss during normal operation
- Intuitive enough for immediate use without documentation
- Visual clarity that enables quick decision-making
