# Docker Compose Deployment

## 📋 Overview

Docker Compose setup สำหรับรัน Authentication Project แบบ Production-ready

## 🏗️ Services

### 1. Database (PostgreSQL)
- **Image:** postgres:16-alpine
- **Port:** 5433 (host) → 5432 (container)
- **Health Check:** ตรวจสอบทุก 5 วินาที
- **Volume:** postgres_data (persistent storage)

### 2. Migration
- **Build:** apps/api/Dockerfile.migrate
- **Purpose:** รัน database migrations
- **Depends on:** database (healthy)
- **Restart:** on-failure (รันครั้งเดียว)

### 3. API (Go Backend)
- **Build:** apps/api/Dockerfile
- **Port:** 8080 (host) → 8080 (container)
- **Depends on:** database (healthy) + migration (completed)
- **Health Check:** GET /api/v1/health

### 4. Client (Vue + Nginx)
- **Build:** apps/client/Dockerfile
- **Port:** 5173 (host) → 80 (container)
- **Depends on:** api (healthy)
- **Web Server:** Nginx

## 🚀 Quick Start

### Start All Services

```bash
# จาก root project directory
cd deployment/compose

# Start all services
docker-compose up -d

# ดู logs
docker-compose logs -f

# ดู logs เฉพาะ service
docker-compose logs -f api
```

### Stop Services

```bash
# Stop all services
docker-compose down

# Stop และลบ volumes (ลบข้อมูล database)
docker-compose down -v
```

### Rebuild Services

```bash
# Rebuild all services
docker-compose up -d --build

# Rebuild specific service
docker-compose up -d --build api
```

## 📊 Service Startup Order

```
1. database
   ↓ (wait for healthy)
2. migration
   ↓ (wait for completed)
3. api
   ↓ (wait for healthy)
4. client
```

## 🔍 Health Checks

### Database
```bash
# ตรวจสอบว่า database พร้อมใช้งาน
docker-compose exec database pg_isready -U auth_user -d auth_db
```

### API
```bash
# ตรวจสอบว่า API พร้อมใช้งาน
curl http://localhost:8080/api/v1/health
```

### Client
```bash
# ตรวจสอบว่า client พร้อมใช้งาน
curl http://localhost:5173/health
```

## 🌍 Environment Variables

Docker Compose อ่านค่าจากไฟล์ `.env` ที่ root project:

```env
# Database
POSTGRES_DB=auth_db
POSTGRES_USER=auth_user
POSTGRES_PASSWORD=auth_password
POSTGRES_PORT=5433

# API
HTTP_PORT=8080
JWT_SECRET=your-secret-key
JWT_EXPIRE_MINUTES=30
REFRESH_TOKEN_EXPIRE_MINUTES=60

# Client
CLIENT_PORT=5173
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## 📝 Service Details

### Database Service
- **Persistent Storage:** Volume `postgres_data`
- **Health Check:** ทุก 5 วินาที, retry 10 ครั้ง
- **Start Period:** 10 วินาที (เวลาในการ start)

### Migration Service
- **Run Once:** รันครั้งเดียวแล้วหยุด
- **Restart Policy:** on-failure
- **Purpose:** สร้าง database schema

### API Service
- **Multi-stage Build:** ลด image size
- **Health Check:** GET /api/v1/health
- **Restart Policy:** unless-stopped

### Client Service
- **Nginx:** Production web server
- **SPA Support:** Vue Router history mode
- **Static Assets:** Cached 1 year
- **Gzip:** Enabled

## 🐛 Troubleshooting

### Service ไม่ start
```bash
# ดู logs
docker-compose logs <service-name>

# ตรวจสอบ status
docker-compose ps

# Restart service
docker-compose restart <service-name>
```

### Database connection error
```bash
# ตรวจสอบว่า database healthy
docker-compose ps database

# ดู database logs
docker-compose logs database

# Restart database
docker-compose restart database
```

### Migration failed
```bash
# ดู migration logs
docker-compose logs migration

# Run migration manually
docker-compose run --rm migration

# Recreate migration service
docker-compose up -d --force-recreate migration
```

### API ไม่ตอบ
```bash
# ตรวจสอบ health
curl http://localhost:8080/api/v1/health

# ดู logs
docker-compose logs api

# Restart API
docker-compose restart api
```

### Client ไม่โหลด
```bash
# ตรวจสอบ nginx
docker-compose exec client nginx -t

# ดู logs
docker-compose logs client

# Restart client
docker-compose restart client
```

## 🔧 Useful Commands

```bash
# ดู running containers
docker-compose ps

# ดู logs แบบ real-time
docker-compose logs -f

# Execute command ใน container
docker-compose exec api sh
docker-compose exec database psql -U auth_user -d auth_db

# ดู resource usage
docker stats

# Clean up everything
docker-compose down -v --remove-orphans
docker system prune -a
```

## 📊 Resource Usage

### Estimated Resources
- **Database:** ~50MB RAM
- **API:** ~20MB RAM
- **Client:** ~10MB RAM
- **Total:** ~80MB RAM

### Image Sizes
- **postgres:16-alpine:** ~240MB
- **API (built):** ~20MB
- **Client (built):** ~30MB

## 🚀 Production Considerations

### Security
- ✅ Change JWT_SECRET to random value
- ✅ Use strong database password
- ✅ Enable HTTPS (add reverse proxy)
- ✅ Limit CORS origins
- ✅ Add rate limiting

### Performance
- ✅ Enable Nginx gzip
- ✅ Cache static assets
- ✅ Use connection pooling
- ✅ Add Redis for sessions

### Monitoring
- [ ] Add health check endpoints
- [ ] Log aggregation
- [ ] Metrics collection
- [ ] Error tracking

## 📚 Additional Resources

- [API Documentation](../../apps/api/README.md)
- [Frontend Testing](../../apps/client/TESTING.md)
- [Project README](../../README.md)

---

**Note:** This setup is for demonstration purposes. For production, consider using Kubernetes or managed services.
