# 🏋️ WORKOUT_API

A RESTful API for managing workouts, users, and authentication, built with Go and [Chi](https://github.com/go-chi/chi). This service supports user registration, authentication via token, and secured workout CRUD operations.

---

## 📌 Features

- ✅ User registration
- ✅ Token-based authentication
- ✅ Protected workout CRUD routes
- ✅ Health check endpoint
- ✅ Middleware-based access control

---

## 🚀 Endpoints

### 🔓 Public Routes

| Method | Endpoint           | Description               |
|--------|--------------------|---------------------------|
| `GET`  | `/health`          | Health check              |
| `POST` | `/users/register`  | Register a new user       |
| `POST` | `/auth`            | Authenticate and get token|

---

### 🔐 Protected Routes

> These require a valid authentication token.

All routes below are protected by:

- `app.Middleware.Auth` — general authentication
- `app.Middleware.RequireUser` — ensures a valid user context

#### 🏋️ Workout Routes

| Method  | Endpoint         | Description                  |
|---------|------------------|------------------------------|
| `GET`   | `/workouts/{id}` | Get a workout by ID          |
| `POST`  | `/workouts`      | Create a new workout         |
| `PUT`   | `/workouts/{id}` | Update an existing workout   |
| `DELETE`| `/workouts/{id}` | Delete a workout by ID       |

---

## 🧱 Project Structure

```
.
├── internal/
│   └── app/                 # Application struct and dependencies
│   └── handlers/            # Route handler logic
│   └── middleware/          # Middleware for auth/user context
│   └── routes/              # Route definitions
├── main.go                  # App entry point
├── Makefile                 # Build & run automation
├── go.mod / go.sum          # Go dependencies
```

---

## ⚙️ Running the API

### ▶️ Using Makefile

Build and run the app with Docker and environment support:

```bash
make run ENV=stg LEVEL=debug PORT=8080
```

You can also build without running:

```bash
make build
```

And run separately:

```bash
./bin/workout_server --level=debug --port=8080
```

### 🧪 Run Tests

```bash
make test
```

### 🐳 Docker Commands

```bash
make docker-up ENV=stg       # Start DB container
make docker-down ENV=stg     # Stop DB container
make docker-logs ENV=stg     # Tail DB logs
make docker-clean            # Prune Docker volumes
```

---

## 🔒 Authentication

After registering via `/users/register`, obtain a token using:

```http
POST /auth
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

Use the returned token in the `Authorization` header for protected routes:

```http
Authorization: Bearer <your-token>
```

---

## 📞 Health Check

```bash
curl http://localhost:8080/health
```

Should return a 200 OK if the server is healthy.

---
