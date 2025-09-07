# AuraTravel AI - Simple Version Startup Script
Write-Host "🚀 Starting AuraTravel AI - Simple Version" -ForegroundColor Cyan
Write-Host "===========================================" -ForegroundColor Cyan

# Check if Go is installed
$goVersion = & go version 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Go is not installed. Please install Go first." -ForegroundColor Red
    exit 1
}

Write-Host "✅ Go found: $goVersion" -ForegroundColor Green

# Start the backend
Write-Host ""
Write-Host "🔧 Starting Backend Server..." -ForegroundColor Yellow
Write-Host "Backend will run on: http://localhost:8080" -ForegroundColor White
Write-Host "API endpoints:" -ForegroundColor White
Write-Host "  - Health: http://localhost:8080/health" -ForegroundColor Gray
Write-Host "  - Destinations: http://localhost:8080/api/v1/destinations" -ForegroundColor Gray
Write-Host "  - Chat: http://localhost:8080/api/v1/ai/chat" -ForegroundColor Gray
Write-Host ""
Write-Host "🌐 Frontend: Open simple-frontend.html in your browser" -ForegroundColor Magenta
Write-Host ""

# Run the simple backend
go run simple-backend.go
