# Authentication API

REST API สำหรับระบบ Authentication ที่พัฒนาด้วย Go, Fiber Framework และ PostgreSQL

## 🎯 Overview

API นี้เป็นส่วน Backend ของระบบ Authentication ที่ออกแบบตาม **Hexagonal Architecture (Ports & Adapters)** เพื่อให้โค้ดมีความยืดหยุ่น ทดสอบง่าย และแยก Business Logic ออกจาก Infrastructure

## 🏗️ Architecture

```
internal/
├── core/                    # Business Logic (Domain Layer)
│   ├── model/              # Domain Models
│   ├── service/            # Business Services
│   ├── port/               # Interfaces (Ports)
│   │   ├── inbound/       # Use Cases
│   │   └── outbound/      # Repository & External Services
│   └── error/             # Domain Errors
├── adapter/                # Infrastructure Layer
│   ├── inbound/           # Input Adapters
│   │   └── web/          # HTTP Handlers, Routes, Middleware
│   └── outbound/         # Output Adapters
│       ├── persistence/  # Database (PostgreSQL)
│       └── security/     # JWT, Password Hashing
└── bootstrap/             # Application Setup
```

### Design Patterns ที่ใช้
- **Hexagonal Architecture** - แยก Business Logic ออกจาก Infrastructure
- **Dependency Injection** - ใช้ Container pattern
- **Repository Pattern** - Abstract database operations
- **Factory Pattern** - สร้าง objects ผ่าน modules

## ✨ Features

### Authentication
- ✅ **Register** - สมัครสมาชิกด้วย username/password
- ✅ **Login** - เข้าสู่ระบบและรับ JWT tokens
- ✅ **Me** - ดึงข้อมูล user ที่ login อยู่
- ✅ **JWT Authentication** - Access token & Refresh token
- ✅ **Password Hashing** - ใช้ Argon2id (secure)

### Security
- 🔒 Password hashing ด้วย **Argon2id**
- 🔑 JWT tokens (Access + Refresh)
- 🍪 HttpOnly Cookies
- ✅ Input validation
- 🛡️ CORS configuration

## 🛠️ Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go 1.25** | Programming Language |
| **Fiber v2** | Web Framework |
| **PostgreSQL** | Database |
| **GORM** | ORM |
| **JWT** | Authentication |
| **Argon2id** | Password Hashing |
| **Goose** | Database Migration |
| **Zap** | Logging |
| **Validator** | Input Validation |

## 📡 API Endpoints

### Public Endpoints

#### POST `/api/v1/register`
สมัครสมาชิกใหม่

**Request Body:**
```json
{
  "username": "testuser",
  "password": "password123",
  "confirm_password": "password123"
}
```

**Response:** `201 Created`
```json
{
  "code": "SUCCESS",
  "message": "User registered successfully",
  "data": {}
}
```

#### POST `/api/v1/login`
เข้าสู่ระบบ

**Request Body:**
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**Response:** `200 OK`
```json
{
  "code": "SUCCESS",
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc..."
  }
}
```

**Cookies Set:**
- `access_token` (HttpOnly, 30 minutes)
- `refresh_token` (HttpOnly, 60 minutes)

### Protected Endpoints

#### GET `/api/v1/me`
ดึงข้อมูล user ที่ login อยู่

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:** `200 OK`
```json
{
  "code": "SUCCESS",
  "message": "User information retrieved successfully",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "testuser"
  }
}
```

## 🚀 Getting Started

### Prerequisites
- Go 1.25+
- PostgreSQL 15+
- Docker (optional)

### Installation

1. **Clone repository**
```bash
git clone <repository-url>
cd apps/api
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup environment**
```bash
# สร้างไฟล์ .env ที่ root project
cp .env.example .env
```

4. **Start PostgreSQL**
```bash
# ใช้ Docker
docker-compose up -d postgres

# หรือใช้ PostgreSQL ที่ติดตั้งไว้แล้ว
```

5. **Run migrations**
```bash
go run cmd/migrate/main.go
```

6. **Start server**
```bash
go run cmd/server/main.go
```

Server จะรันที่ `http://localhost:8080`

## 🧪 Testing

```bash
# รัน tests ทั้งหมด
go test ./...

# รัน tests พร้อม coverage
go test -cover ./...

# รัน tests แบบละเอียด
go test -v ./...

# รัน tests เฉพาะ package
go test ./internal/core/service/...
```

### Test Coverage
- ✅ Unit Tests สำหรับทุก layer
- ✅ Integration Tests สำหรับ handlers
- ✅ Repository Tests
- ✅ Service Tests

## 📁 Project Structure

```
apps/api/
├── cmd/
│   ├── server/          # Main application
│   └── migrate/         # Database migration tool
├── internal/
│   ├── adapter/         # Infrastructure adapters
│   ├── bootstrap/       # App initialization
│   ├── config/          # Configuration
│   └── core/            # Business logic
├── migrations/          # SQL migrations
├── go.mod
└── README.md
```

## 🔐 Security Features

### Password Security
- **Argon2id** algorithm
- Memory: 64MB
- Iterations: 3
- Parallelism: 2
- Salt: 16 bytes (random)

### JWT Configuration
- **Algorithm:** HS256
- **Access Token:** 30 minutes
- **Refresh Token:** 60 minutes
- **Secret:** 256-bit random key

### Input Validation
- Username: 3-255 characters, alphanumeric
- Password: 8-255 characters, ASCII only
- Automatic trimming and sanitization

## 🌍 Environment Variables

```env
# Database
POSTGRES_DB=auth_db
POSTGRES_USER=auth_user
POSTGRES_PASSWORD=auth_password
POSTGRES_PORT=5433

# Server
HTTP_PORT=8080

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRE_MINUTES=30
REFRESH_TOKEN_EXPIRE_MINUTES=60
```

## 📊 Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🎓 Key Learnings

### Architecture
- ✅ Hexagonal Architecture ทำให้แยก concerns ได้ชัดเจน
- ✅ Dependency Injection ทำให้ทดสอบง่าย
- ✅ Repository Pattern ทำให้เปลี่ยน database ได้ง่าย

### Security
- ✅ Argon2id ปลอดภัยกว่า bcrypt
- ✅ JWT + HttpOnly Cookies = XSS protection
- ✅ Input validation ป้องกัน injection attacks

### Testing
- ✅ Unit tests ช่วยให้มั่นใจในการ refactor
- ✅ Mock interfaces ทำให้ test เร็วขึ้น
- ✅ Table-driven tests ทำให้เพิ่ม test cases ง่าย

## 🚧 Future Improvements

- [ ] Email verification
- [ ] Password reset
- [ ] Rate limiting
- [ ] Refresh token rotation
- [ ] Redis caching
- [ ] Metrics & monitoring
- [ ] API documentation (Swagger)

## 📝 License

This project is for educational/interview purposes.

## 👤 Author

Developed as part of a job interview technical assessment.

---

**Note:** This is a demonstration project showcasing clean architecture, security best practices, and Go development skills.
