# erosai-server

A Go web server for URL tracking and content analysis recommendations.

## Quick Start (macOS)

```bash
# 1. Install PostgreSQL (if not installed)
brew install postgresql@15
brew services start postgresql@15

# 2. Set up database
createuser -s erosai
createdb -O erosai erosai
psql -U erosai -d erosai -f database/links.sql

# 3. Install Go dependencies
go mod tidy

# 4. Start mock ML server (in one terminal)
go run mock/mlserver.go

# 5. Start the server (in another terminal)
go run main.go
```

The server runs on **http://localhost:8082**

## Prerequisites

- Go 1.16+
- PostgreSQL

## API Usage

**1. Register a user:**
```bash
curl -X POST http://localhost:8082/register-attempt \
  -H "Content-Type: application/json" \
  -d '{"Username": "testuser", "Password": "testpass"}'
```

**2. Login and get token:**
```bash
curl -X POST http://localhost:8082/login-attempt \
  -H "Content-Type: application/json" \
  -d '{"Username": "testuser", "Password": "testpass"}'
# Returns: your-session-token
```

**3. Submit a URL for processing:**
```bash
curl -X POST http://localhost:8082/process-url \
  -H "Content-Type: application/json" \
  -d '{"Token": "your-session-token", "URL": "https://example.com/page"}'
```

**4. Get recommendations:**
```bash
curl "http://localhost:8082/get-recommendations?token=your-session-token"
```

## Configuration

Edit `globals/globals.go`:

| Variable | Default | Description |
|----------|---------|-------------|
| `BetServerPort` | `:8082` | Server port |
| `MLServerEndpoint` | `http://localhost:8001` | ML server URL |

## Database

Connection configured in `database/connection.go`:
- User: `erosai`
- Password: `Erosai11!!`
- Database: `erosai`
