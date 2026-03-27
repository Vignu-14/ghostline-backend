# 💻 Development Guide

Set up your local development environment for Ghostline Backend.

---

## 📋 Prerequisites

- **Go 1.25+** ([download](https://go.dev/dl))
- **PostgreSQL 15+** or use local Supabase
- **Supabase CLI** (optional but recommended)
- **Git**
- **Make** (optional, for Makefile commands)

---

## 🚀 Quick Start

### 1. Clone Repository
```bash
git clone https://github.com/Vignu-14/ghostline-backend.git
cd ghostline-backend
```

### 2. Create `.env` File
Copy from `.env.example` and fill in values:
```bash
cp .env.example .env
```

Edit `.env`:
```
PORT=3000
ENVIRONMENT=development
ALLOWED_ORIGIN=http://localhost:5173,http://127.0.0.1:5173

DATABASE_URL=postgresql://postgres:password@localhost:5432/ghostline
DB_MAX_CONNECTIONS=25
DB_MIN_CONNECTIONS=5
DB_MAX_CONN_LIFETIME_MINUTES=60
DB_MAX_CONN_IDLE_MINUTES=15
DB_HEALTH_CHECK_SECONDS=30
DB_CONNECT_TIMEOUT_SECONDS=5

JWT_SECRET=your-dev-secret-here
JWT_EXPIRATION_MINUTES=15
AUTH_COOKIE_NAME=auth_token
COOKIE_SECURE=false

SUPABASE_URL=https://your-project.supabase.co
SUPABASE_SERVICE_KEY=your-service-key
STORAGE_BUCKET_NAME=user-uploads
```

### 3. Install Dependencies
```bash
go mod download
```

### 4. Set Up Database

**Option A: Local PostgreSQL**
```bash
# Create database
createdb ghostline

# Run migrations
go run cmd/server/main.go  # Migrations run automatically
```

**Option B: Supabase Local (Recommended)**
```bash
# Install Supabase CLI
npm install -g supabase

# Start local Supabase
supabase start

# Copy local DATABASE_URL from output
```

### 5. Run Server
```bash
go run cmd/server/main.go
```

Server runs at `http://localhost:3000`

---

## 🔧 Development Commands

### Using Make
```bash
# Build
make build

# Run
make run

# Run with hot reload
make dev

# Run tests
make test

# Run tests with coverage
make test-coverage

# Lint code
make lint

# Format code
make fmt
```

### Using Go Directly
```bash
# Build binary
go build -o server ./cmd/server

# Run
./server

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Lint (requires golangci-lint)
golangci-lint run
```

---

## 📁 Project Structure Quick Reference

```
ghostline-backend/
├── cmd/server/main.go          # Entry point
├── internal/handlers/           # HTTP handlers
├── internal/services/           # Business logic
├── internal/repositories/       # Database queries
├── internal/models/             # Data structures
├── internal/middleware/         # Request middleware
├── internal/websocket/          # WebSocket logic
├── internal/database/           # DB connection
├── internal/config/             # Configuration
├── tests/                        # Test files
└── scripts/                      # Utility scripts
```

---

## 🔐 Authentication Development

### JWT Token Format
```go
claims := models.JWTClaims{
    UserID:         "user-uuid",
    Role:           "user",          // "user" or "admin"
    ImpersonatorID: nil,             // Only for admin impersonation
    ExpiresAt:      time.Now().Add(15 * time.Minute),
}
```

### Test Token Generation
```go
token, _ := utils.GenerateToken(
    "your-jwt-secret",
    15*time.Minute,
    uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
    "user",
    nil,
)
fmt.Println(token)
```

Use this token in requests:
```bash
curl -H "Authorization: Bearer $token" \
     http://localhost:3000/api/auth/me
```

---

## 🗄️ Database Management

### Run Migrations
```bash
# Automatic on startup
go run cmd/server/main.go

# Or manually with Supabase CLI
supabase db push
```

### Seed Database
```bash
# Run seed script
./scripts/seed.sh
```

Or manually:
```sql
-- Connect to dev database
psql ghostline

-- Insert test data
INSERT INTO users (username, email, password_hash, role)
VALUES ('testuser', 'test@example.com', '$2a$12$...', 'user');
```

### View Data
```bash
# Connect to local database
psql ghostline

# List tables
\dt

# View users
SELECT id, username, email, role FROM users;

# View posts
SELECT id, user_id, caption FROM posts;
```

---

## 🧪 Testing

### Run Tests
```bash
# All tests
go test ./...

# Specific package
go test ./internal/services

# With verbose output
go test -v ./...

# With coverage percentage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Structure
```
tests/
├── unit/               # Unit tests (fast)
│   ├── jwt_test.go
│   ├── bcrypt_test.go
│   └── sanitizer_test.go
│
├── integration/        # Integration tests (slow)
│   ├── auth_test.go
│   ├── post_test.go
│   └── chat_test.go
│
├── mocks/              # Mock objects
├── fixtures/           # Test data
└── TestDatabase.go     # Test database setup
```

### Writing Tests
```go
// Example unit test
func TestGenerateToken(t *testing.T) {
    token, err := utils.GenerateToken(
        "secret",
        15*time.Minute,
        uuid.New(),
        "user",
        nil,
    )
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if token == "" {
        t.Error("expected non-empty token")
    }
}
```

---

## 🐛 Debugging

### Enable Debug Logging
Set in `.env`:
```
LOG_LEVEL=debug
```

### Use Delve Debugger
```bash
# Install
go install github.com/go-delve/delve/cmd/dlv@latest

# Run with debugger
dlv debug ./cmd/server

# In dlv console
(dlv) break main.main
(dlv) continue
```

### VS Code Debugging
Add to `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Connect to Go server",
            "type": "go",
            "mode": "local",
            "request": "launch",
            "program": "${fileDirname}",
            "env": {},
            "args": []
        }
    ]
}
```

### Log Examples
```go
// Logs appear in console
slog.Info("user logged in", "user_id", userID, "ip", ipAddress)
slog.Error("database error", "error", err)
slog.Debug("processing request", "request_id", requestID)
```

---

## 🔌 WebSocket Development

### Test WebSocket Connection
```bash
# Using wscat
npm install -g wscat

# Connect to WebSocket
wscat -c ws://localhost:3000/ws/chat \
      --header "Cookie: auth_token=$JWT_TOKEN"

# Send message
{"type":"message","receiver_id":"...","content":"hello"}

# Should receive response
```

### WebSocket Message Format
```json
{
  "type": "message",
  "receiver_id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "Hello!"
}
```

---

## 📤 File Upload Development

### Test Upload Flow
```bash
# 1. Get signed URL
curl -X POST http://localhost:3000/api/posts/upload-url \
     -H "Cookie: auth_token=$JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "file_name": "photo.jpg",
       "file_size": 1024000,
       "file_type": "image/jpeg"
     }'

# 2. Upload to Supabase (use returned URL)
curl -X PUT "<signed_url>" \
     --data-binary @photo.jpg

# 3. Create post
curl -X POST http://localhost:3000/api/posts/finalize \
     -H "Cookie: auth_token=$JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "caption": "My photo",
       "image_url": "posts/user-uuid/file-uuid.jpg"
     }'
```

---

## 🛠️ Common Development Tasks

### Add New Endpoint
1. Create handler in `internal/handlers/`
2. Register route in `internal/routes/api.go`
3. Add business logic in `internal/services/`
4. Add database query in `internal/repositories/`
5. Write tests in `tests/`

### Add New Database Table
1. Create migration in `internal/database/migrations/`
2. Run migrations (automatic on startup)
3. Create model in `internal/models/`
4. Create repository in `internal/repositories/`
5. Create service in `internal/services/`

### Connect New Collection/Store
1. Create service in `internal/services/`
2. Add initialization in `cmd/server/main.go`
3. Inject into handlers
4. Test endpoints

---

## 📚 Useful Resources

- [Go Documentation](https://pkg.go.dev)
- [Fiber Framework Docs](https://docs.gofiber.io)
- [PostgreSQL Docs](https://www.postgresql.org/docs)
- [Supabase Docs](https://supabase.com/docs)
- [JWT Tokens](https://jwt.io)

---

## 🆘 Troubleshooting

### "database connection failed"
```bash
# Check PostgreSQL is running
psql -U postgres -d ghostline -c "SELECT 1"

# Check DATABASE_URL in .env
```

### "module not found"
```bash
# Download dependencies
go mod download

# Update dependencies
go get -u ./...
```

### "permission denied" on scripts
```bash
# Make scripts executable
chmod +x scripts/*.sh
```

### "port already in use"
```bash
# Change PORT in .env
PORT=3001  # or another available port

# Or kill process using port 3000
lsof -i :3000 | grep LISTEN | awk '{print $2}' | xargs kill -9
```

---

See also: [DEPLOYMENT.md](./DEPLOYMENT.md), [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
