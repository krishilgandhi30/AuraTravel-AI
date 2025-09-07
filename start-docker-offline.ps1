# AuraTravel AI - Offline Docker Setup Script
Write-Host "🐳 AuraTravel AI - Offline Docker Setup" -ForegroundColor Cyan
Write-Host "=======================================" -ForegroundColor Cyan

# Check if Docker is running
try {
    $dockerVersion = docker version --format "{{.Server.Version}}" 2>$null
    if ($LASTEXITCODE -ne 0) {
        throw "Docker not responding"
    }
    Write-Host "✅ Docker found: $dockerVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker is not running. Please start Docker Desktop first." -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "🔨 Building Go binary for Linux..." -ForegroundColor Yellow

# Build Linux binary for Docker
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

try {
    go build -a -ldflags '-extldflags "-static"' -o simple-backend-linux simple-backend.go
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Linux binary built successfully" -ForegroundColor Green
    } else {
        throw "Build failed"
    }
} catch {
    Write-Host "❌ Failed to build Linux binary" -ForegroundColor Red
    Write-Host "💡 Make sure Go is installed and simple-backend.go exists" -ForegroundColor Yellow
    exit 1
} finally {
    # Reset environment
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
    Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
}

Write-Host ""
Write-Host "🐳 Building Docker image (offline)..." -ForegroundColor Yellow

try {
    docker build -f Dockerfile.offline -t auratravel-simple .
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Docker image built successfully" -ForegroundColor Green
    } else {
        throw "Docker build failed"
    }
} catch {
    Write-Host "❌ Failed to build Docker image" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "🧹 Stopping any existing containers..." -ForegroundColor Yellow
docker stop auratravel-simple-container 2>$null
docker rm auratravel-simple-container 2>$null

Write-Host ""
Write-Host "🚀 Starting container..." -ForegroundColor Yellow

try {
    docker run -d --name auratravel-simple-container -p 8080:8080 auratravel-simple
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Container started successfully" -ForegroundColor Green
    } else {
        throw "Container start failed"
    }
} catch {
    Write-Host "❌ Failed to start container" -ForegroundColor Red
    docker logs auratravel-simple-container
    exit 1
}

Write-Host ""
Write-Host "⏳ Waiting for service to be ready..." -ForegroundColor Yellow

# Wait for service to be ready
$timeout = 30
$elapsed = 0
do {
    Start-Sleep -Seconds 2
    $elapsed += 2
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET -TimeoutSec 5
        if ($response.success -eq $true) {
            Write-Host "✅ Backend is ready!" -ForegroundColor Green
            break
        }
    } catch {
        Write-Host "." -NoNewline -ForegroundColor Gray
    }
} while ($elapsed -lt $timeout)

if ($elapsed -ge $timeout) {
    Write-Host ""
    Write-Host "⚠️  Service might still be starting up. Check logs with:" -ForegroundColor Yellow
    Write-Host "   docker logs auratravel-simple-container" -ForegroundColor Gray
} else {
    Write-Host ""
    Write-Host "🎉 AuraTravel AI is running!" -ForegroundColor Green
    Write-Host "===========================" -ForegroundColor Green
    Write-Host ""
    Write-Host "🔧 Backend API: http://localhost:8080" -ForegroundColor Blue
    Write-Host "💚 Health Check: http://localhost:8080/health" -ForegroundColor Cyan
    Write-Host "📊 Destinations: http://localhost:8080/api/v1/destinations" -ForegroundColor Magenta
    Write-Host ""
    Write-Host "🌐 Frontend: Open simple-frontend.html in your browser" -ForegroundColor Green
    Write-Host "   (It will automatically connect to the containerized backend)" -ForegroundColor Gray
}

Write-Host ""
Write-Host "📝 Useful commands:" -ForegroundColor White
Write-Host "   View logs:    docker logs auratravel-simple-container -f" -ForegroundColor Gray
Write-Host "   Stop:         docker stop auratravel-simple-container" -ForegroundColor Gray
Write-Host "   Remove:       docker rm auratravel-simple-container" -ForegroundColor Gray
Write-Host "   Restart:      docker restart auratravel-simple-container" -ForegroundColor Gray
