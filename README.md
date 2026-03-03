# Authentication Project

Full-stack Authentication System สำหรับสัมภาษณ์งาน

## 📋 Project Overview

โปรเจคนี้เป็นระบบ Authentication แบบ Full-stack ที่ประกอบด้วย:
- **Backend API** (Go + Fiber + PostgreSQL)
- **Frontend Web** (Vue 3 + TypeScript + Vite)

พัฒนาเพื่อแสดงความสามารถใน:
- ✅ Clean Architecture & Design Patterns
- ✅ Security Best Practices
- ✅ Modern Web Development
- ✅ Testing & Documentation

## 🏗️ Architecture

```
authentication-project-exam/
├── apps/
│   ├── api/              # Backend (Go)
│   │   ├── cmd/         # Entry points
│   │   ├── internal/    # Application code
│   │   └── migrations/  # Database migrations
│   └── client/          # Frontend (Vue 3)
│       ├── src/
│       │   ├── components/
│       │   ├── views/
│       │   ├── services/
│       │   └── router/
│       └── tests/
├── .env                 # Environment variables
└── docker-compose.yml   # Docker setup
```

## ✨ Features

### Authentication
- 🔐 User Registration with validation
- 🔑 Login with JWT tokens
- 👤 Protected user profile page
- 🍪 HttpOnly Cookies for security
- 🔄 Automatic token management

### Security
- **Password Hashing:** Argon2id
- **Authentication:** JWT (Access + Refresh tokens)
- **Validation:** Client & Server side
- **CORS:** Configured properly
- **XSS Protection:** HttpOnly cookies

### User Experience
- 🎨 Modern, responsive UI
- ⚡ Fast loading with Vite
- ✅ Real-time form validation
- 🔄 Loading states & error handling
- 🎭 Smooth transitions & animations

## 🛠️ Tech Stack

### Backend
| Technology | Purpose |
|------------|---------|
| Go 1.25 | Programming Language |
| Fiber v2 | Web Framework |
| PostgreSQL | Database |
| GORM | ORM |
| JWT | Authentication |
| Argon2id | Password Hashing |

### Frontend
| Technology | Purpose |
|------------|---------|
| Vue 3 | UI Framework |
| TypeScript | Type Safety |
| Vite | Build Tool |
| Vue Router | Routing |
| Axios | HTTP Client |
| Vitest | Testing |

## 🚀 Quick Start

### Prerequisites
- Go 1.25+
- Node.js 20+
- PostgreSQL 15+
- Docker (optional)

### 1. Clone & Setup

```bash
git clone <repository-url>
cd authentication-project-exam
```

### 2. Environment Setup

```bash
# สร้างไฟล์ .env
cp .env.example .env

# แก้ไขค่าตามต้องการ (optional)
```

### 3. Start with Docker (แนะนำ)

**ต้องมีเพียง:**
- ✅ Docker + Docker Compose
- ✅ ไฟล์ `.env`
- ❌ ไม่ต้องติดตั้ง Go, Node.js, PostgreSQL

```bash
# Start all services
make up

# View logs
make logs

# Stop services
make down

# Clean everything
make clean
```

**Service Startup Order:**
1. **Database** (PostgreSQL) → health check
2. **Migration** → รัน database migrations
3. **API** (Go) → health check
4. **Client** (Vue + Nginx)

**Access:**
- Frontend: http://localhost:5173
- Backend: http://localhost:8080/api/v1
- Health: http://localhost:8080/api/v1/health

**Available Commands:**
```bash
make up          # Start all services
make down        # Stop all services
make restart     # Restart services
make logs        # View all logs
make logs-api    # API logs only
make logs-client # Client logs only
make ps          # Show running containers
make clean       # Remove everything
make build       # Rebuild images
```

**Troubleshooting:**
```bash
# Check service status
make ps

# View specific logs
make logs-api

# Restart specific service
docker-compose restart api

# Clean and rebuild
make clean && make build
```

### 4. Start Manually (Alternative)

#### Backend
```bash
cd apps/api

# Install dependencies
go mod download

# Run migrations
go run cmd/migrate/main.go

# Start server
go run cmd/server/main.go
```

#### Frontend
```bash
cd apps/client

# Install dependencies
npm install

# Start dev server
npm run dev
```

## 📡 API Endpoints

### Public
- `POST /api/v1/register` - สมัครสมาชิก
- `POST /api/v1/login` - เข้าสู่ระบบ

### Protected (ต้อง login)
- `GET /api/v1/me` - ดูข้อมูลตัวเอง

**ดูรายละเอียดเพิ่มเติม:** [API Documentation](apps/api/README.md)

## 🧪 Testing

### Backend Tests
```bash
cd apps/api
go test ./...
go test -cover ./...
```

### Frontend Tests
```bash
cd apps/client

# Run tests
npm run test:unit:run

# Watch mode
npm run test:unit:watch

# Coverage report
npm run test:unit:coverage
```

**Test Coverage:**
- Backend: Unit tests ครอบคลุมทุก layer
- Frontend: 57 tests, 86.44% coverage

**ดูรายละเอียดเพิ่มเติม:** [Testing Guide](apps/client/TESTING.md)

## 📁 Project Structure

### Backend (Hexagonal Architecture)
```
apps/api/internal/
├── core/              # Business Logic
│   ├── model/        # Domain Models
│   ├── service/      # Use Cases
│   └── port/         # Interfaces
├── adapter/          # Infrastructure
│   ├── inbound/     # HTTP Handlers
│   └── outbound/    # Database, Security
└── bootstrap/        # App Setup
```

### Frontend (Feature-based)
```
apps/client/src/
├── components/       # Reusable components
├── views/           # Page components
├── services/        # API services
├── router/          # Route configuration
└── __tests__/       # Unit tests
```

## 🔐 Security Features

### Password Security
- ✅ Argon2id hashing (memory-hard)
- ✅ Random salt per password
- ✅ Minimum 8 characters
- ✅ ASCII validation

### Token Security
- ✅ JWT with HS256
- ✅ Short-lived access tokens (30 min)
- ✅ Longer refresh tokens (60 min)
- ✅ HttpOnly cookies (XSS protection)

### Input Validation
- ✅ Client-side validation (instant feedback)
- ✅ Server-side validation (security)
- ✅ Sanitization & trimming
- ✅ SQL injection prevention

## 🎯 Key Highlights

### Architecture
- **Hexagonal Architecture** - แยก Business Logic ออกจาก Infrastructure
- **Dependency Injection** - ทำให้ code testable
- **Repository Pattern** - Abstract database operations
- **Service Layer** - Centralized business logic

### Code Quality
- ✅ TypeScript for type safety
- ✅ Comprehensive unit tests
- ✅ Clean code principles
- ✅ SOLID principles
- ✅ Error handling
- ✅ Logging

### Best Practices
- ✅ Environment-based configuration
- ✅ Database migrations
- ✅ API versioning
- ✅ CORS configuration
- ✅ Request validation
- ✅ Structured logging

## 📊 Database Schema

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

## 🎓 What I Learned

### Backend Development
- Hexagonal Architecture ใน Go
- Fiber framework และ middleware
- GORM ORM และ database migrations
- JWT authentication flow
- Argon2id password hashing
- Unit testing ใน Go

### Frontend Development
- Vue 3 Composition API
- TypeScript ใน Vue
- Axios interceptors
- Vue Router guards
- Form validation
- Vitest unit testing

### DevOps
- Docker containerization
- Docker Compose orchestration
- Environment management
- Database migrations

## 📚 Documentation

- [Backend README](apps/api/README.md) - API documentation
- [Docker Guide](DOCKER.md) - Detailed Docker deployment guide
- [Deployment Requirements](DEPLOYMENT_REQUIREMENTS.md) - What you need to deploy

## 🐛 Troubleshooting

### Docker Issues

**Services not starting:**
```bash
# Check status
make ps

# View logs
make logs

# Clean and restart
make clean && make up
```

**Database connection error:**
```bash
# Check database health
docker-compose exec database pg_isready -U auth_user

# Restart database
docker-compose restart database
```

**Migration stuck:**
```bash
# View migration logs
make logs-migration

# Restart migration
docker-compose up migration
```

### Port Already in Use
```bash
# เปลี่ยน port ใน .env
HTTP_PORT=8081
CLIENT_PORT=5174
POSTGRES_PORT=5434
```

### Frontend Can't Connect to API
```bash
# ตรวจสอบ CORS settings
# ตรวจสอบ VITE_API_BASE_URL ใน .env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### Manual Development (Without Docker)

**Backend:**
```bash
cd apps/api
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5433
go run cmd/migrate/main.go
go run cmd/server/main.go
```

**Frontend:**
```bash
cd apps/client
npm install
npm run dev
```

---

**Thank you for reviewing this project!** 🙏
