@echo off
REM Trading Dashboard Build Script (Batch Version)
REM This script builds both the executable and installer

echo.
echo ============================================
echo   Trading Dashboard Build Script
echo ============================================
echo.

REM Check if we're in the right directory
if not exist "trading-dashboard\wails.json" (
    echo ERROR: Must run this script from the project root directory
    echo Expected to find: trading-dashboard\wails.json
    pause
    exit /b 1
)

REM Navigate to the trading-dashboard directory
echo Navigating to trading-dashboard directory...
cd trading-dashboard

REM Start the build process
echo.
echo Starting build process...
echo Building: Application (.exe) + Installer
echo Mode: Clean build with NSIS installer
echo.

REM Run the Wails build command
echo Running: wails build --clean --nsis
wails build --clean --nsis

if %ERRORLEVEL% neq 0 (
    echo.
    echo BUILD FAILED!
    cd ..
    pause
    exit /b 1
)

echo.
echo BUILD COMPLETED SUCCESSFULLY!
echo.
echo Built files are located in: build\bin\
echo - Application: trading-dashboard.exe
echo - Installer: trading-dashboard-amd64-installer.exe
echo.

cd ..
echo Build script completed successfully!
pause 