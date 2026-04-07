# Personal Portfolio - Microservices Architecture

A production-grade portfolio website built with microservices architecture, demonstrating full-stack development with modern DevOps practices.

## 🏗️ Architecture

### Backend Microservices (Golang)
- **Auth Service** (Port 8001) - JWT authentication, user management
- **Portfolio Service** (Port 8002) - CRUD operations for projects
- **Analytics Service** (Coming soon) - Visitor tracking and metrics

### Frontend (In Progress)
- React + TypeScript
- Tailwind CSS
- Admin dashboard

### Infrastructure
- Docker Compose (local development)
- PostgreSQL + Redis
- AWS deployment (planned)

## 🚀 Quick Start
```bash
# Start databases
docker-compose up -d

# Start auth service
cd services/auth-service
go run cmd/server/main.go

# Start portfolio service (in another terminal)
cd services/portfolio-service
go run cmd/server/main.go
```

## 📚 API Endpoints

### Auth Service (http://localhost:8001)
- `POST /api/v1/auth/register` - Create account
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/me` - Get current user (protected)

### Portfolio Service (http://localhost:8002)
- `GET /api/v1/projects` - List projects (public)
- `POST /api/v1/projects` - Create project (protected)
- `PUT /api/v1/projects/:id` - Update project (protected)
- `DELETE /api/v1/projects/:id` - Delete project (protected)

## 🛠️ Technologies

**Backend:** Golang, Gin, PostgreSQL, Redis, Docker, JWT, Bcrypt  
**Frontend:** React, TypeScript, Vite, Tailwind CSS (in progress)  
**DevOps:** Docker Compose, GitHub Actions (planned), AWS (planned)

## 👤 Author

Alejandro Ruiz de Castilla
- GitHub: [@alecruiz](https://github.com/alecruiz)
- Email: alecrdec98@gmail.com