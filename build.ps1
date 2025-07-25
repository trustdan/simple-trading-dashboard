#!/usr/bin/env powershell

# Trading Dashboard Build Script
# This script builds both the executable and installer for the trading dashboard

Write-Host "Trading Dashboard Build Script" -ForegroundColor Cyan
Write-Host "==============================" -ForegroundColor Cyan
Write-Host ""

# Check if we're in the right directory
if (-not (Test-Path "trading-dashboard/wails.json")) {
    Write-Host "ERROR: Must run this script from the project root directory" -ForegroundColor Red
    Write-Host "   Expected to find: trading-dashboard/wails.json" -ForegroundColor Yellow
    exit 1
}

# Navigate to the trading-dashboard directory
Write-Host "Navigating to trading-dashboard directory..." -ForegroundColor Yellow
Set-Location "trading-dashboard"

# Check if Wails is installed
Write-Host "Checking Wails installation..." -ForegroundColor Yellow
try {
    $wailsVersion = wails version 2>$null
    if ($LASTEXITCODE -ne 0) {
        throw "Wails not found"
    }
    Write-Host "Wails CLI found" -ForegroundColor Green
} catch {
    Write-Host "ERROR: Wails CLI not found or not installed" -ForegroundColor Red
    Write-Host "   Please install Wails: https://wails.io/docs/gettingstarted/installation" -ForegroundColor Yellow
    Set-Location ".."
    exit 1
}

# Check if Go is installed
Write-Host "Checking Go installation..." -ForegroundColor Yellow
try {
    $goVersion = go version 2>$null
    if ($LASTEXITCODE -ne 0) {
        throw "Go not found"
    }
    Write-Host "Go found: $($goVersion.Split(' ')[2])" -ForegroundColor Green
} catch {
    Write-Host "ERROR: Go not found or not installed" -ForegroundColor Red
    Write-Host "   Please install Go: https://golang.org/dl/" -ForegroundColor Yellow
    Set-Location ".."
    exit 1
}

# Start the build process
Write-Host ""
Write-Host "Starting build process..." -ForegroundColor Cyan
Write-Host "   Building: Application (.exe) + Installer" -ForegroundColor White
Write-Host "   Mode: Clean build with NSIS installer" -ForegroundColor White
Write-Host ""

# Measure build time
$buildStart = Get-Date

# Run the Wails build command
try {
    Write-Host "Running: wails build --clean --nsis" -ForegroundColor Blue
    wails build --clean --nsis
    
    if ($LASTEXITCODE -ne 0) {
        throw "Build failed with exit code $LASTEXITCODE"
    }
    
    $buildEnd = Get-Date
    $buildTime = $buildEnd - $buildStart
    
    Write-Host ""
    Write-Host "Build completed successfully!" -ForegroundColor Green
    Write-Host "Build time: $($buildTime.ToString('mm\:ss'))" -ForegroundColor Cyan
    
    # Check and display built files
    Write-Host ""
    Write-Host "Built files:" -ForegroundColor Cyan
    
    if (Test-Path "build/bin/trading-dashboard.exe") {
        $exeInfo = Get-ChildItem "build/bin/trading-dashboard.exe"
        $exeSize = [math]::Round($exeInfo.Length / 1MB, 1)
        Write-Host "   Application: $($exeInfo.Name) ($exeSize MB)" -ForegroundColor Green
    }
    
    if (Test-Path "build/bin/trading-dashboard-amd64-installer.exe") {
        $installerInfo = Get-ChildItem "build/bin/trading-dashboard-amd64-installer.exe"
        $installerSize = [math]::Round($installerInfo.Length / 1MB, 1)
        Write-Host "   Installer: $($installerInfo.Name) ($installerSize MB)" -ForegroundColor Green
    }
    
    Write-Host ""
    Write-Host "Ready to distribute! Files are in: trading-dashboard/build/bin/" -ForegroundColor Cyan
    
} catch {
    Write-Host ""
    Write-Host "Build failed: $($_.Exception.Message)" -ForegroundColor Red
    Set-Location ".."
    exit 1
}

# Return to original directory
Set-Location ".."

Write-Host ""
Write-Host "Build script completed successfully!" -ForegroundColor Green
Write-Host "   Run the installer: trading-dashboard/build/bin/trading-dashboard-amd64-installer.exe" -ForegroundColor White
Write-Host "   Or run directly: trading-dashboard/build/bin/trading-dashboard.exe" -ForegroundColor White 