# Deployment Requirements

## ✅ สิ่งที่ต้องมีในเครื่อง

### เครื่องที่จะ Deploy ต้องมีเพียง:

1. **Docker** (version 20.10+)
2. **Docker Compose** (version 2.0+)
3. **ไฟล์ `.env`** (configuration)

### ❌ ไม่ต้องติดตั้ง:

- ❌ Go
- ❌ Node.js
- ❌ npm
- ❌ PostgreSQL
- ❌ Nginx

**ทุกอย่างรันใน Docker containers!**

---

## 🎯 ทำไมไม่ต้องติดตั้งอะไรเลย?

### Multi-stage Docker Builds

#### API (Go)
```dockerfile
# Stage 1: Build (ใน container)
FROM golang:1.25-alpine AS builder
  → ติดตั้ง Go
  → Download dependencies
  → Compile binary

# Stage 2: Runtime (ใน container)
FROM alpine:latest
  → Copy binary ที่ compile แล้ว
  → รัน binary
```

**ผลลัพธ์:** Binary ที่ compile แล้ว ไม่ต้องมี Go ในเครื่อง!

#### Client (Vue)
```dockerfile
# Stage 1: Build (ใน container)
FROM node:20-alpine AS builder
  → ติดตั้ง Node.js
  → npm install
  → npm run build

# Stage 2: Runtime (ใน container)
FROM nginx:alpine
  → Copy dist files
  → รัน nginx
```

**ผลลัพธ์:** Static files + Nginx ไม่ต้องมี Node.js ในเครื่อง!

---

## 🚀 การ Deploy บนเครื่องใหม่

### Scenario: เครื่องที่ไม่มีอะไรเลย

```bash
# 1. ติดตั้ง Docker เท่านั้น
# Windows: Docker Desktop
# Linux: apt install docker.io docker-compose
# Mac: Docker Desktop

# 2. Clone project
git clone <repository-url>
cd authentication-project-exam

# 3. Setup environment
cp .env.example .env

# 4. Start everything
make up
# หรือ
cd deployment/compose && docker-compose up -d

# 5. เสร็จ! ✅
# Frontend: http://localhost:5173
# Backend: http://localhost:8080
```

**ใช้เวลา:** ~5-10 นาที (ครั้งแรก build images)

---

## 🔍 การทำงานของ Docker

### Build Process

#### 1. API Build
```
Docker pulls golang:1.25-alpine image
  ↓
Install Go dependencies (ใน container)
  ↓
Compile Go code (ใน container)
  ↓
Copy binary to alpine image
  ↓
Final image: ~20MB (ไม่มี Go compiler)
```

#### 2. Client Build
```
Docker pulls node:20-alpine image
  ↓
Install npm packages (ใน container)
  ↓
Build Vue app (ใน container)
  ↓
Copy dist files to nginx image
  ↓
Final image: ~30MB (ไม่มี Node.js)
```

### Runtime Process

```
Container 1: PostgreSQL (มี PostgreSQL)
Container 2: Migration (มี Go binary)
Container 3: API (มี Go binary)
Container 4: Client (มี Nginx + static files)

Host Machine: ไม่ต้องมีอะไรเลย!
```

---

## ✅ Verification Checklist

### ตรวจสอบว่า Docker Setup ถูกต้อง

#### 1. Dockerfiles มี Multi-stage Builds
- ✅ `apps/api/Dockerfile` - golang → alpine
- ✅ `apps/api/Dockerfile.migrate` - golang → alpine
- ✅ `apps/client/Dockerfile` - node → nginx

#### 2. Docker Compose อ่าน .env
- ✅ `env_file: - ../../.env` ในทุก service
- ✅ ไม่มี hardcoded values
- ✅ ใช้ `${VARIABLE}` syntax

#### 3. Service Dependencies ถูกต้อง
- ✅ migration depends on database (healthy)
- ✅ api depends on database (healthy) + migration (completed)
- ✅ client depends on api (healthy)

#### 4. Health Checks ครบถ้วน
- ✅ database: pg_isready
- ✅ api: HTTP GET /health
- ✅ client: HTTP GET /health

---

## 🎓 สิ่งที่ Interviewer จะประทับใจ

### 1. ไม่ต้องติดตั้งอะไรเลย
```
Interviewer: "ฉันจะรันโปรเจคนี้ยังไง?"
You: "ติดตั้ง Docker แล้วรัน make up เท่านั้น"
Interviewer: "ไม่ต้องติดตั้ง Go หรือ Node.js?"
You: "ไม่ต้องครับ ทุกอย่างรันใน Docker"
```

### 2. Production-Ready
- Multi-stage builds (optimized)
- Health checks (reliability)
- Nginx for client (performance)
- Proper service ordering

### 3. Easy to Use
- One command: `make up`
- Clear documentation
- Troubleshooting guide

---

## 📋 Test Deployment

### ทดสอบบนเครื่องใหม่

```bash
# 1. ติดตั้ง Docker
# 2. Clone project
# 3. Create .env
# 4. Run:

make up

# ถ้าทำงาน = ✅ Success!
# ถ้าไม่ทำงาน = ดู logs: make logs
```

### Expected Output

```bash
$ make up
Starting all services...
[+] Running 5/5
 ✔ Network auth-network          Created
 ✔ Container ...database          Started (healthy)
 ✔ Container ...migration         Started (completed)
 ✔ Container ...api               Started (healthy)
 ✔ Container ...client            Started

Services started! Access:
  - Frontend: http://localhost:5173
  - Backend:  http://localhost:8080
```

---

## 🎯 Summary

### คำถาม: ต้องติดตั้งอะไรบ้าง?
**คำตอบ:** Docker เท่านั้น!

### คำถาม: ไม่ต้องติดตั้ง Go?
**คำตอบ:** ไม่ต้อง! Go จะถูกใช้ใน Docker build stage เท่านั้น

### คำถาม: ไม่ต้องติดตั้ง Node.js?
**คำตอบ:** ไม่ต้อง! Node.js จะถูกใช้ใน Docker build stage เท่านั้น

### คำถาม: ไม่ต้องติดตั้ง PostgreSQL?
**คำตอบ:** ไม่ต้อง! PostgreSQL รันใน Docker container

### คำถาม: ไม่ต้องติดตั้ง Nginx?
**คำตอบ:** ไม่ต้อง! Nginx รันใน Docker container

---

## ✨ ข้อดีของ Docker Approach

### 1. Portability
- รันได้บนทุก OS (Windows, Mac, Linux)
- ไม่ต้องกังวลเรื่อง dependencies
- "Works on my machine" → "Works everywhere"

### 2. Consistency
- Development = Production environment
- ไม่มีปัญหา version conflicts
- Reproducible builds

### 3. Easy Deployment
- One command deployment
- No manual setup
- Automatic dependency management

### 4. Isolation
- แต่ละ service รันแยกกัน
- ไม่กระทบ host system
- ลบง่าย (docker-compose down)

---

## 🎤 สำหรับการสัมภาษณ์

### Talking Points

**"โปรเจคนี้ใช้ Docker แบบ production-ready"**

1. **Multi-stage builds** - ลด image size 90%
2. **Health checks** - ตรวจสอบว่า service พร้อมใช้งาน
3. **Service dependencies** - start ตามลำดับที่ถูกต้อง
4. **Environment management** - ใช้ .env file
5. **Zero installation** - ต้องมีแค่ Docker

**"Interviewer สามารถรันได้ทันทีโดยไม่ต้องติดตั้งอะไร"**

```bash
# เพียง 3 commands
git clone <repo>
cd authentication-project-exam
make up
```

---

## ✅ Final Verification

### ตรวจสอบว่าทุกอย่างพร้อม:

- ✅ Dockerfile ใช้ multi-stage builds
- ✅ Build stage มี Go/Node.js
- ✅ Runtime stage ไม่มี Go/Node.js
- ✅ Docker Compose อ่าน .env
- ✅ Service dependencies ถูกต้อง
- ✅ Health checks ครบถ้วน
- ✅ Documentation ชัดเจน

### คำตอบสุดท้าย:

# ✅ ใช่! คุณสามารถใช้ Docker ได้โดยไม่ต้องติดตั้ง Go, Node.js, PostgreSQL หรือ Nginx เลย!

**ต้องมีแค่:**
1. Docker
2. Docker Compose
3. ไฟล์ .env

**แค่นั้น!** 🎉

---

**หมายเหตุ:** สำหรับ development (ถ้าต้องการแก้โค้ด) ถึงจะต้องติดตั้ง Go และ Node.js แต่สำหรับการ **รันโปรเจคเพื่อดู demo** ไม่ต้องติดตั้งอะไรเลย!
