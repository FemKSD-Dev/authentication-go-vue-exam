# Docker Deployment Guide

## 🐳 Quick Start

### วิธีที่ 1: ใช้ Makefile (แนะนำ)

```bash
# Start all services
make up

# View logs
make logs

# Stop services
make down
```

### วิธีที่ 2: ใช้ Docker Compose โดยตรง

```bash
# Start
cd deployment/compose
docker-compose up -d

# Stop
docker-compose down
```

## 📋 Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- ไฟล์ `.env` ที่ root project

## 🚀 Deployment Steps

### 1. Setup Environment

```bash
# สร้างไฟล์ .env (ถ้ายังไม่มี)
cp .env.example .env

# แก้ไขค่าตามต้องการ
nano .env
```

### 2. Start Services

```bash
make up
```

หรือ

```bash
cd deployment/compose
docker-compose up -d
```

### 3. Verify Services

```bash
# ตรวจสอบ status
make ps

# ดู logs
make logs
```

### 4. Access Application

- **Frontend:** http://localhost:5173
- **Backend API:** http://localhost:8080
- **Health Check:** http://localhost:8080/api/v1/health

## 🔄 Service Startup Sequence

```
┌─────────────┐
│  Database   │ ← Step 1: Start PostgreSQL
└──────┬──────┘
       │ (health check)
       ↓
┌─────────────┐
│  Migration  │ ← Step 2: Run migrations
└──────┬──────┘
       │ (completed)
       ↓
┌─────────────┐
│     API     │ ← Step 3: Start Go server
└──────┬──────┘
       │ (health check)
       ↓
┌─────────────┐
│   Client    │ ← Step 4: Start Nginx
└─────────────┘
```

## 📊 Container Details

### Database Container
```yaml
Image: postgres:16-alpine
Memory: ~50MB
CPU: 0.5 core
Storage: Volume (persistent)
Health: pg_isready check
```

### Migration Container
```yaml
Build: Multi-stage Go build
Memory: ~20MB (during run)
Lifecycle: Run once and exit
Restart: On failure only
```

### API Container
```yaml
Build: Multi-stage Go build
Memory: ~20MB
CPU: 0.5 core
Port: 8080
Health: HTTP GET /health
```

### Client Container
```yaml
Build: Node build + Nginx
Memory: ~10MB
CPU: 0.2 core
Port: 80 (mapped to 5173)
Health: HTTP GET /health
```

## 🛠️ Common Commands

### Service Management

```bash
# Start specific service
docker-compose up -d database

# Stop specific service
docker-compose stop api

# Restart service
docker-compose restart client

# Remove service
docker-compose rm -f migration
```

### Logs & Debugging

```bash
# All logs
make logs

# Specific service logs
make logs-api
make logs-client
make logs-db

# Follow logs with tail
docker-compose logs -f --tail=100 api
```

### Database Operations

```bash
# Access database shell
make db-shell

# Or manually
docker-compose exec database psql -U auth_user -d auth_db

# Run migration again
make db-migrate
```

### Container Shell Access

```bash
# API container
docker-compose exec api sh

# Client container
docker-compose exec client sh

# Database container
docker-compose exec database sh
```

## 🔧 Troubleshooting

### Problem: Services ไม่ start

```bash
# ตรวจสอบ status
docker-compose ps

# ดู logs
docker-compose logs

# Restart all
make restart
```

### Problem: Database connection failed

```bash
# ตรวจสอบ database health
docker-compose ps database

# ดู database logs
make logs-db

# Restart database
docker-compose restart database
```

### Problem: Migration failed

```bash
# ดู migration logs
make logs-migration

# Run migration manually
docker-compose run --rm migration

# Check database tables
make db-shell
\dt
```

### Problem: Port already in use

```bash
# แก้ไข port ใน .env
POSTGRES_PORT=5434
HTTP_PORT=8081
CLIENT_PORT=5174

# Restart services
make down
make up
```

### Problem: Out of disk space

```bash
# ลบ unused images
docker image prune -a

# ลบ unused volumes
docker volume prune

# ลบทุกอย่างที่ไม่ใช้
docker system prune -a --volumes
```

## 🔄 Update & Rebuild

### Update Code

```bash
# 1. Pull latest code
git pull

# 2. Rebuild services
make build

# 3. Restart
make restart
```

### Rebuild Specific Service

```bash
# Rebuild API only
make build-api

# Rebuild client only
make build-client
```

### Update Dependencies

```bash
# API dependencies
cd apps/api
go mod tidy
cd ../..
make build-api

# Client dependencies
cd apps/client
npm install
cd ../..
make build-client
```

## 🧹 Cleanup

### Soft Cleanup (keep data)

```bash
make down
```

### Hard Cleanup (remove everything)

```bash
make clean
```

### Remove Images

```bash
# Remove project images
docker-compose down --rmi all

# Remove all unused images
docker image prune -a
```

## 📦 Build Information

### Multi-stage Builds

#### API Dockerfile
```dockerfile
Stage 1: golang:1.25-alpine (build)
  → Compile Go binary
  
Stage 2: alpine:latest (runtime)
  → Copy binary only
  → Result: ~20MB image
```

#### Client Dockerfile
```dockerfile
Stage 1: node:20-alpine (build)
  → npm install & build
  
Stage 2: nginx:alpine (runtime)
  → Copy dist files only
  → Result: ~30MB image
```

## 🌐 Network Configuration

### Network: auth-network (bridge)
- All services communicate through this network
- Services can reference each other by name
- Example: API connects to `database:5432`

### Port Mapping

| Service | Container Port | Host Port | Protocol |
|---------|---------------|-----------|----------|
| Database | 5432 | 5433 | TCP |
| API | 8080 | 8080 | HTTP |
| Client | 80 | 5173 | HTTP |

## 🔐 Security Notes

### Production Checklist
- [ ] Change default passwords
- [ ] Use secrets management
- [ ] Enable HTTPS
- [ ] Limit CORS origins
- [ ] Add rate limiting
- [ ] Use non-root users
- [ ] Scan images for vulnerabilities

### Current Security Features
- ✅ HttpOnly cookies
- ✅ Password hashing (Argon2id)
- ✅ JWT tokens
- ✅ Input validation
- ✅ CORS configured

## 📈 Monitoring

### Health Checks

```bash
# Check all services
docker-compose ps

# API health
curl http://localhost:8080/api/v1/health

# Client health
curl http://localhost:5173/health
```

### Resource Monitoring

```bash
# Real-time stats
docker stats

# Specific container
docker stats authentication-segment-api
```

## 🎯 Best Practices

### Development
- ✅ Use `.dockerignore` to reduce build context
- ✅ Multi-stage builds for smaller images
- ✅ Health checks for reliability
- ✅ Named volumes for data persistence

### Production
- ✅ Use specific image versions (not `latest`)
- ✅ Set resource limits
- ✅ Enable logging drivers
- ✅ Use secrets for sensitive data
- ✅ Regular security updates

## 📝 Notes

- Migration service จะรันครั้งเดียวแล้วหยุด (by design)
- Database data จะถูกเก็บใน volume `postgres_data`
- Client รันบน Nginx สำหรับ production
- All services ใช้ Alpine Linux (เบา)

---

**Happy Dockerizing!** 🐳
