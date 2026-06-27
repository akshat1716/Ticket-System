# Ticket System

A backend REST API for user registration, authentication, and ticket management. Users can create tickets, view only their own tickets, and update ticket status following strict transition rules.

## Project Overview

This is a backend-only Go service built with:

- **Gin** — HTTP router
- **GORM** — ORM with SQLite
- **JWT** — Authentication
- **bcrypt** — Password hashing

### Features

- User registration and login
- JWT-based authentication via `Authorization: Bearer <token>`
- Create, list, and retrieve tickets (own tickets only)
- Update ticket status with validated transitions: `open` → `in_progress` → `closed`

## Setup

### Prerequisites

- Go 1.23+
- Docker (optional)

### Environment Variables

Copy the example environment file:

```bash
cp .env.example .env
```

| Variable       | Description              | Default              |
|----------------|--------------------------|----------------------|
| `JWT_SECRET`   | Secret key for JWT       | `default-secret-key` |
| `PORT`         | Server port              | `8080`               |
| `DATABASE_PATH`| SQLite database file path| `./ticket_system.db` |

## Run Locally

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/main.go
```

The server starts on **http://localhost:8080**.

Open **http://localhost:8080** in your browser to use the built-in test frontend (register, login, create tickets, update status).

## Docker Commands

```bash
# Build the image
docker build -t ticket-system .

# Run the container
docker run -p 8080:8080 ticket-system
```

With custom environment variables:

```bash
docker run -p 8080:8080 \
  -e JWT_SECRET=your-production-secret \
  -e DATABASE_PATH=/root/ticket_system.db \
  ticket-system
```

## API Endpoints

| Method | Endpoint                  | Auth | Description                    |
|--------|---------------------------|------|--------------------------------|
| GET    | `/health`                 | No   | Health check                   |
| POST   | `/auth/register`          | No   | Register a new user            |
| POST   | `/auth/login`             | No   | Login and receive JWT          |
| POST   | `/tickets`                | Yes  | Create a ticket                |
| GET    | `/tickets`                | Yes  | List own tickets               |
| GET    | `/tickets/{id}`           | Yes  | Get own ticket by ID           |
| PATCH  | `/tickets/{id}/status`    | Yes  | Update ticket status           |

### Status Transitions

| From          | Allowed To    |
|---------------|---------------|
| `open`        | `in_progress` |
| `in_progress` | `closed`      |
| `closed`      | *(none)*      |

### Error Format

```json
{
  "error": "Error message"
}
```

## Example curl Commands

### Health Check

```bash
curl http://localhost:8080/health
```

### Register

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Save the returned token for authenticated requests.

### Create Ticket

```bash
curl -X POST http://localhost:8080/tickets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Fix login bug",
    "description": "Users cannot login with special characters"
  }'
```

### List Tickets

```bash
curl http://localhost:8080/tickets \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Get Ticket by ID

```bash
curl http://localhost:8080/tickets/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Update Ticket Status

```bash
curl -X PATCH http://localhost:8080/tickets/1/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "status": "in_progress"
  }'
```

## Deployment (Render)

### Deploy on Render (Free Tier)

1. Push this repository to GitHub.
2. Go to [Render Dashboard](https://dashboard.render.com/) and create a new **Web Service**.
3. Connect your GitHub repository.
4. Configure the service:
   - **Environment**: Docker
   - **Port**: `8080`
   - **Environment Variables**:
     - `JWT_SECRET` — a strong random secret
     - `PORT` — `8080`
     - `DATABASE_PATH` — `/root/ticket_system.db`
5. Deploy.

### Deployment URL

```
https://your-app-name.onrender.com
```

Replace with your actual Render service URL after deployment.

## Postman Collection

Import the Postman collection from:

```
postman/Ticket-System.postman_collection.json
```

The collection includes all endpoints. The Login request automatically saves the JWT token for subsequent requests.

## Project Structure

```
ticket-system/
├── cmd/
│   └── main.go
├── config/
├── controllers/
├── database/
├── middleware/
├── models/
├── postman/
├── routes/
├── services/
├── utils/
├── frontend/
├── Dockerfile
├── go.mod
├── .env.example
└── README.md
```

## License

MIT
