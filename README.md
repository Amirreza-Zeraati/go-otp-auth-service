# Go OTP Auth Service 🔑

A simple OTP (One-Time Password) authentication service built with **Golang**, **Gin**, **Postgres**, and **Redis**.  
It supports phone-based OTP login, JWT authentication, and user management with API docs exposed via **Swagger**.

---

## 📦 Features
- OTP-based user login & registration
- JWT authentication
- Rate-limited OTP requests
- User CRUD operations
- Swagger API documentation (`/swagger/index.html`)
- Dockerized with Postgres & Redis

---

## ⚙️ Run Locally

### 1. Clone the repo
```bash
git clone https://github.com/Amirreza-Zeraati/go-otp-auth-service.git
cd go-otp-auth-service
```
### 2. Install dependencies
```bash
go mod tidy
```
### 3. Set environment variables
bash
Create a .env file in the root:
```bash
DB_HOST=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=dekamond
DB_PORT=5432

REDIS_HOST=localhost
REDIS_PORT=6379

SECRET=supersecretkey
RATE_LIMIT=5
PERIOD_TIME=60
PORT=3000
```

### 4. Run Postgres & Redis locally

Make sure you have Postgres and Redis installed and running.

### 5. Start the app

go run main.go

Now visit:
* Login: [http://localhost:3000](http://localhost:3000)
* Dashboard: [http://localhost:3000/dashboard](http://localhost:3000/dashboard)
* Swagger UI: [http://localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html)

---

## 🐳 Run with Docker

### 1. Build & start containers
```bash
docker-compose up --build
```
### 2. Stop containers
```bash
docker-compose down
```
All services (app, Postgres, Redis) will run together in Docker.

---

## 📖 API Examples

### Request OTP

POST /request-otp
Content-Type: application/json
```bash
{
  "phone": "+123456789"
}
```
Response
```bash
{
  "message": "OTP sent successfully"
}
```
---

### Login

POST /login
Content-Type: application/json
```bash
{
  "phone": "+123456789",
  "otp": "123456"
}
```
Response
```bash
{
  "token": "eyJhbGciOiJIUzI1NiIsInR...",
}
```
---

### Logout

POST /logout

Response
```bash
{
  "message": "Logged out successfully"
}
```
---

### Get User

GET /users/1

Response
```bash
{
  "id": 1,
  "phone": "+123456789",
  "created_at": "2025-08-19T14:00:00Z"
}
```
---

### List Users

GET /users?page=1&search=111

Response
```bash
{
  "users": [
    {
      "id": 3,
      "phone": "111111111111",
      "created_at": "2025-08-19T13:07:29.96325+03:30"
    }
  ],
  "page": 1,
  "prev_page": 0,
  "next_page": 2,
  "has_prev": false,
  "has_next": false,
  "total_pages": 1,
  "search": "111"
}
```
---

## 🛢 Database Choice Justification

### Postgres

* Reliable, scalable, and widely used for structured relational data.
* Perfect for managing user accounts and authentication records.
* Strong support for indexing and queries.

### Redis

* Used for OTP storage and rate limiting.
* Super fast in-memory store with TTL (time-to-live) support.
* Prevents abuse by controlling OTP request frequency.

Together, Postgres handles persistent user data while Redis manages ephemeral OTPs and rate limits.

---
