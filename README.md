# 🚀 PlacementLog Backend

A Go-based backend API for the Placement Log platform with JWT authentication, role-based access, and an admin post review system.

---

## ✅ Features

- **JWT Authentication** for users and admins  
- **Post Management** (CRUD for placement experiences)  
- **Admin Review Workflow** for post approvals  
- **Role-based Access Control**  
- **RESTful API Design**

---

## ⚙️ Prerequisites

- Go 1.21+  
- Git  
- PostgreSQL

---

## 🛠️ Setup

1. **Clone the repository**  
```bash
git clone <repository-url>
cd Backend
```

2. **Configure Environment Variables** (`.env`)
```env
SECRET=your-secure-secret-key-here
DB_URL=postgres://username:password@localhost:5432/placementlog?sslmode=disable
```

3. **Set up Database**  
```bash
createdb placementlog
psql -d placementlog -f schema.sql
```

4. **Install Dependencies**  
```bash
go mod download
```

5. **Run the Server**  
```bash
make runServer
# or
go run ./cmd/main.go
```

The server will start on:  
`http://localhost:8080`

---

## 🔌 API Overview

### 📟 Auth Endpoints
- `POST /auth/login` – User login  
- `POST /auth/register` – User registration  
- `POST /admin/login` – Admin login  

### ✍️ Post Endpoints
- `GET /posts` – Get all approved posts  
- `POST /posts` – Create new post  
- `PUT /posts` – Update post  
- `DELETE /posts` – Delete post  

### 🛡️ Admin Endpoints
- `GET /admin/posts` – View all submitted posts  
- `PUT /admin/posts/review` – Approve or reject a post  
- `DELETE /admin/posts` – Delete post as admin  

---

## 🧪 Testing

```bash
go test ./...
go test -cover ./...
```

---

## 🐳 Docker Support

```bash
make buildDocker
make runDocker
```

---

## 📌 Notes

- Ensure `SECRET` is set in your environment before running — it’s required for token signing.  
- Keep your secrets out of version control (use `.env`).  
- Use strong, environment-specific secrets in production.