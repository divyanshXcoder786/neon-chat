@echo off
echo ======================================
echo   NeonChat - Setup ^& Run (Windows)
echo ======================================

cd /d "%~dp0server"

echo [1/2] Downloading Go dependencies...
go mod tidy

echo [2/2] Starting NeonChat server...
echo Open http://localhost:8080 in your browser
echo Press Ctrl+C to stop
echo.
go run main.go
