# 🐳 AuraTravel AI - Docker Setup

Run AuraTravel AI using Docker containers for easy deployment and consistent environment.

## 🚀 Quick Start

### Prerequisites
- Docker Desktop installed and running
- Internet connection (for first-time image downloads)

### One-Command Start
```powershell
.\start-docker.ps1
```

This will:
1. ✅ Check Docker is running
2. 🧹 Clean up any existing containers
3. 🔨 Build the backend image
4. 🚀 Start both frontend and backend containers
5. 🌐 Open the app at http://localhost:3000

## 📋 Manual Docker Commands

### Start the Application
```bash
docker-compose -f docker-compose.simple.yml up --build -d
```

### Check Status
```bash
docker-compose -f docker-compose.simple.yml ps
```

### View Logs
```bash
# All services
docker-compose -f docker-compose.simple.yml logs -f

# Backend only
docker-compose -f docker-compose.simple.yml logs -f auratravel-api

# Frontend only
docker-compose -f docker-compose.simple.yml logs -f auratravel-frontend
```

### Stop the Application
```bash
docker-compose -f docker-compose.simple.yml down
```

### Rebuild (after code changes)
```bash
docker-compose -f docker-compose.simple.yml down
docker-compose -f docker-compose.simple.yml up --build -d
```

## 🌐 Access Points

Once running, you can access:

- **🎨 Frontend (Main App)**: http://localhost:3000
- **🔧 Backend API**: http://localhost:8080
- **💚 Health Check**: http://localhost:8080/health
- **📊 API Destinations**: http://localhost:8080/api/v1/destinations

## 🏗️ Architecture

### Services
1. **auratravel-api** (Backend)
   - Built from `simple-backend.go`
   - Runs on port 8080
   - Provides REST API for travel data and AI chat

2. **auratravel-frontend** (Frontend)
   - Nginx serving `simple-frontend.html`
   - Runs on port 3000
   - Proxies API requests to backend

### Network
- Services communicate via Docker network `auratravel-network`
- Frontend proxies `/api/*` requests to backend
- CORS configured for cross-origin requests

## 🔧 Configuration

### Environment Variables
- `PORT`: Backend port (default: 8080)

### Volumes
- Frontend files mounted read-only
- Nginx configuration mounted read-only
- Persistent data stored in `auratravel_data` volume

## 🛠️ Development

### Modifying Backend Code
1. Edit `simple-backend.go`
2. Rebuild: `docker-compose -f docker-compose.simple.yml up --build -d auratravel-api`

### Modifying Frontend
1. Edit `simple-frontend.html`
2. Restart: `docker-compose -f docker-compose.simple.yml restart auratravel-frontend`

### Adding Features
- Backend: Add new endpoints in `simple-backend.go`
- Frontend: JavaScript automatically detects Docker vs standalone mode
- Database: Currently in-memory; can add MongoDB service to docker-compose

## 🐛 Troubleshooting

### Container Won't Start
```bash
# Check Docker is running
docker version

# Check container logs
docker-compose -f docker-compose.simple.yml logs

# Remove and rebuild
docker-compose -f docker-compose.simple.yml down --volumes
docker-compose -f docker-compose.simple.yml up --build -d
```

### Port Conflicts
If ports 3000 or 8080 are in use:
```yaml
# Edit docker-compose.simple.yml
services:
  auratravel-frontend:
    ports:
      - "3001:80"  # Change 3000 to 3001
  auratravel-api:
    ports:
      - "8081:8080"  # Change 8080 to 8081
```

### Network Issues
```bash
# Reset Docker networks
docker network prune

# Check container connectivity
docker exec -it auratravel-frontend ping auratravel-api
```

### Frontend Can't Reach Backend
1. Check both containers are running: `docker ps`
2. Verify network: `docker network ls`
3. Check Nginx config: `docker logs auratravel-frontend`
4. Test backend directly: `curl http://localhost:8080/health`

## 📦 Production Deployment

For production:
1. Use proper domain names instead of localhost
2. Add SSL/TLS certificates
3. Configure proper logging
4. Add monitoring and health checks
5. Use Docker secrets for sensitive data
6. Set resource limits

### Example Production Override
```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  auratravel-api:
    restart: always
    deploy:
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M
  
  auratravel-frontend:
    restart: always
    deploy:
      resources:
        limits:
          memory: 128M
        reservations:
          memory: 64M
```

## 🔄 Updates

To update the application:
1. Pull latest code
2. Stop containers: `docker-compose -f docker-compose.simple.yml down`
3. Rebuild: `docker-compose -f docker-compose.simple.yml up --build -d`

## 📊 Monitoring

### Container Stats
```bash
docker stats
```

### Resource Usage
```bash
docker-compose -f docker-compose.simple.yml top
```

### Health Checks
- Backend: http://localhost:8080/health
- Frontend: http://localhost:3000 (should load the app)

The Docker setup provides a consistent, isolated environment that works the same way across different machines and operating systems! 🚀
