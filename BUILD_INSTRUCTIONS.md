# Trading Dashboard Build Instructions

This document explains how to build the Trading Dashboard application and installer.

## Quick Build Scripts

Two automated build scripts are provided to make building easy:

### Option 1: PowerShell Script (Recommended)
```powershell
.\build.ps1
```

### Option 2: Batch File
```cmd
build.bat
```
or from PowerShell:
```powershell
cmd /c build.bat
```

## What the Scripts Do

Both scripts will:
1. ✅ **Verify Prerequisites** - Check that Wails CLI and Go are installed
2. 🧹 **Clean Build** - Remove old build artifacts
3. 🔧 **Build Application** - Compile the Go backend and Svelte frontend
4. 📦 **Create Installer** - Generate NSIS installer for Windows
5. 📊 **Show Results** - Display build time and file sizes

## Output Files

After a successful build, you'll find these files in `trading-dashboard/build/bin/`:

- **`trading-dashboard.exe`** - The standalone application (~12MB)
- **`trading-dashboard-amd64-installer.exe`** - The installer package (~6MB)

## Prerequisites

Before running the build scripts, ensure you have:

- **Wails CLI** - [Installation Guide](https://wails.io/docs/gettingstarted/installation)
- **Go** - [Download Go](https://golang.org/dl/)
- **Node.js** - For the frontend build process

## Manual Build (Alternative)

If you prefer to build manually:

```bash
cd trading-dashboard
wails build --clean --nsis
```

## Features Included

The built application includes:
- 📈 **24 Options Strategies** across 8 categories
- 💼 **Market Sentiment Tracking**
- 📊 **Trade Analytics & Visualization**
- 🎯 **Export & Reporting Tools**
- 🎨 **Modern UI with Dark/Light Themes**

## Distribution

- **For end users**: Distribute the installer (`trading-dashboard-amd64-installer.exe`)
- **For portable use**: Distribute the executable (`trading-dashboard.exe`)

## Build Time

Typical build times:
- **Clean build**: ~15-30 seconds
- **Incremental build**: ~5-10 seconds

---

**Happy Trading!** 📊 