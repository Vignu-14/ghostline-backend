# 🌐 Ghostline Backend

A high-performance REST API and WebSocket server for anonymous messaging built with Go, Fiber, PostgreSQL, and Supabase.

![Go](https://img.shields.io/badge/Go-1.25-blue?logo=go)
![Fiber](https://img.shields.io/badge/Fiber-2.x-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue?logo=postgresql)
![Supabase](https://img.shields.io/badge/Supabase-Auth%20%26%20Storage-blue?logo=supabase)
![License](https://img.shields.io/badge/License-MIT-green)

**Live API:** https://ghostline-backend-production-a17a.up.railway.app

---

## ✨ Features

- 🚀 **High Performance** - Built with Go and Fiber framework
- 🔐 **Secure** - JWT authentication, bcrypt hashing, XSS protection, parameterized queries
- 💬 **Real-time Messaging** - WebSocket support for instant messaging
- 📊 **Rich API** - Complete REST API for posts, users, chat, and admin features
- 🛡️ **Authorization** - Role-based access control (user/admin)
- 📁 **File Uploads** - Supabase integration for image storage
- 🗄️ **Database** - PostgreSQL with connection pooling and migrations
- ⚡ **Rate Limiting** - Prevent abuse with configurable rate limits
- 📝 **Logging** - Structured logging for debugging

---

## 🛠️ Tech Stack

- **Language:** Go 1.25
- **Framework:** Fiber (Express-like Go web framework)
- **Database:** PostgreSQL 15
- **Auth:** JWT tokens + httpOnly cookies
- **File Storage:** Supabase Storage
- **Real-time:** WebSocket
- **Security:** Bcrypt, Bluemonday (XSS sanitizer), CORS, parameterized queries

---

## 📦 Getting Started

### Prerequisites
- Go 1.25+
- PostgreSQL 15+ (or local Supabase)
- Git

### Installation

```bash
# 1. Clone repository
git clone https://github.com/Vignu-14/ghostline-backend.git
cd ghostline-backend

# 2. Create .env file
cp .env.example .env

# 3. Configure environment variables
# Edit .env with your database URL, JWT secret, Supabase keys, etc.

# 4. Download dependencies
go mod download

# 5. Run server (migrations run automatically)
go run cmd/server/main.go
```

API runs at `http://localhost:3000`

---

## 🚀 Quick Links

📚 **Documentation:**
- [API Documentation](./docs/API_DOCUMENTATION.md) - Complete endpoint reference
- [Architecture Guide](./docs/ARCHITECTURE.md) - System design and data flows
- [Deployment Guide](./docs/DEPLOYMENT.md) - Production setup (Railway, Supabase, Vercel)
- [Development Guide](./docs/DEVELOPMENT.md) - Local setup and testing
- [WebSocket Protocol](./docs/WEBSOCKET_PROTOCOL.md) - Real-time messaging protocol
- [Troubleshooting Guide](./docs/TROUBLESHOOTING.md) - Common issues and solutions

---

## 📋 Project Structure

```
ghostline-backend/
├── cmd/server/main.go               # Application entry point
├── internal/
│   ├── handlers/                    # HTTP request handlers
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── post_handler.go
│   │   ├── chat_handler.go
│   │   ├── like_handler.go
│   │   ├── admin_handler.go
│   │   ├── websocket_handler.go
│   │   └── health_handler.go
│   ├── services/                    # Business logic
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── post_service.go
│   │   ├── chat_service.go
│   │   └── ...
│   ├── repositories/                # Database queries
│   │   ├── user_repository.go
│   │   ├── post_repository.go
│   │   ├── chat_repository.go
│   │   └── ...
│   ├── models/                      # Data models
│   │   ├── user.go
│   │   ├── post.go
│   │   ├── chat.go
│   │   └── ...
│   ├── middleware/                  # Request middleware
│   │   ├── auth_middleware.go
│   │   ├── cors_middleware.go
│   │   ├── jwt_middleware.go
│   │   ├── error_handler.go
│   │   └── ...
│   ├── database/                    # Database setup
│   │   ├── postgres.go
│   │   ├── migrations/              # SQL migration files
│   │   └── seed/                    # Seed data
│   ├── websocket/                   # WebSocket logic
│   │   └── hub.go
│   ├── config/                      # Configuration
│   │   ├── config.go
│   │   ├── constants.go
│   │   ├── database.go
│   │   └── storage.go
│   ├── utils/                       # Utility functions
│   │   ├── jwt_utils.go
│   │   ├── bcrypt_utils.go
│   │   └── ...
│   └── routes/                      # Route definitions
│       └── api.go
├── pkg/                             # Public packages
│   ├── logger/                      # Logging
│   └── cache/                       # Caching
├── tests/                           # Test files
│   ├── unit/
│   ├── integration/
│   ├── mocks/
│   ├── fixtures/
│   └── TestDatabase.go
├── scripts/                         # Utility scripts
│   ├── migrate.sh
│   ├── seed.sh
│   ├── generate_jwt_secret.sh
│   └── ...
├── docs/                            # Documentation
│   ├── API_DOCUMENTATION.md
│   ├── ARCHITECTURE.md
│   ├── DEPLOYMENT.md
│   ├── DEVELOPMENT.md
│   ├── WEBSOCKET_PROTOCOL.md
│   └── TROUBLESHOOTING.md
├── Dockerfile                       # Container image
├── docker-compose.yml               # Local development
├── go.mod & go.sum                  # Dependencies
├── Makefile                         # Build commands
└── README.md                        # This file
```

---

## 🚀 Available Commands

### Using Make
```bash
make build           # Build binary
make run            # Run server
make dev            # Run with hot reload
make test           # Run tests
make test-coverage  # Run tests with coverage
make lint           # Lint code
make fmt            # Format code
```

### Using Go Directly
```bash
go run cmd/server/main.go              # Run server
go test ./...                          # Run tests
go test -cover ./...                   # With coverage
go fmt ./...                           # Format
golangci-lint run                      # Lint
```

---

## 🔐 Environment Variables

Create `.env` file:

```env
# Server Configuration
PORT=3000
ENVIRONMENT=production
ALLOWED_ORIGIN=https://ghostline-frontend-five.vercel.app

# Database
DATABASE_URL=postgresql://user:pass@host:5432/ghostline
DB_MAX_CONNECTIONS=25
DB_MIN_CONNECTIONS=5
DB_MAX_CONN_LIFETIME_MINUTES=60
DB_MAX_CONN_IDLE_MINUTES=15

# JWT Configuration
JWT_SECRET=your-secret-key-min-32-chars
JWT_EXPIRATION_MINUTES=15
AUTH_COOKIE_NAME=auth_token
COOKIE_SECURE=true

# Supabase Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_SERVICE_KEY=your-service-key
STORAGE_BUCKET_NAME=user-uploads

# Logging
LOG_LEVEL=info
```

See [.env.example](.env.example) for all available variables.

---

## 📡 API Overview

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - Create account
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user

### Users
- `GET /api/users/:userId` - Get user profile
- `GET /api/users/:userId/posts` - Get user's posts
- `PUT /api/users/:userId` - Update profile
- `GET /api/users/:userId/stats` - Get user statistics

### Posts
- `GET /api/posts` - Get feed
- `POST /api/posts` - Create post
- `GET /api/posts/:postId` - Get single post
- `DELETE /api/posts/:postId` - Delete post
- `POST /api/posts/upload-url` - Get signed upload URL
- `POST /api/posts/finalize` - Finalize post with file

### Messages (Chat)
- `GET /api/messages` - Get conversations
- `GET /api/messages/:conversationId` - Get conversation
- `POST /api/messages` - Send message
- `PUT /api/messages/:messageId/read` - Mark as read

### Likes
- `POST /api/posts/:postId/like` - Like post
- `DELETE /api/posts/:postId/like` - Unlike post

### WebSocket
- `ws://localhost:3000/ws/chat` - Real-time messaging

### Admin
- `GET /api/admin/users` - List all users
- `DELETE /api/admin/users/:userId` - Delete user
- `GET /api/admin/posts` - List all posts
- `DELETE /api/admin/posts/:postId` - Delete post

### Health
- `GET /health` - Server health check

See [API_DOCUMENTATION.md](./docs/API_DOCUMENTATION.md) for complete reference.

#### Messages
```
GET    /api/messages/conversations                - List conversations
GET    /api/messages/:userId                     - Get chat history
POST   /api/messages                             - Send message
POST   /api/messages/delete                      - Delete messages
POST   /api/messages/:userId/clear               - Clear conversation
```

#### Admin
```
POST   /api/admin/impersonate       - Start impersonation session
POST   /api/admin/impersonate/stop  - Stop impersonation
```

### Response Format

**Success (200):**
```json
{
  "status": "success",
  "data": { /* response data */ }
}
```

**Error (4xx/5xx):**
```json
{
  "status": "error",
  "error": "Error message",
  "details": { /* validation errors */ }
}
```

---

## 🔌 WebSocket Protocol

### Connection
```
URL: ws://localhost:3000/ws/chat
Headers: Cookie: auth_token=<JWT>
```

### Message Types

#### Incoming
```json
{
  "type": "message",
  "receiver_id": "uuid",
  "content": "message text"
}
```

#### Outgoing
```json
{
  "type": "connected",
  "data": {}
}
```

```json
{
  "type": "message",
  "data": {
    "id": "uuid",
    "sender_id": "uuid",
    "receiver_id": "uuid",
    "content": "text",
    "created_at": "2026-03-28T12:00:00Z"
  }
}
```

```json
{
  "type": "error",
  "data": { "message": "error description" }
}
```

---

## 🗄️ Database Schema

### Tables

**users**
- id (UUID, PK)
- username (VARCHAR, UNIQUE)
- email (VARCHAR, UNIQUE)
- password_hash (VARCHAR)
- impersonation_password_hash (VARCHAR, nullable)
- role (VARCHAR: 'user' | 'admin')
- profile_picture_url (TEXT, nullable)
- created_at (TIMESTAMP)

**posts**
- id (UUID, PK)
- user_id (UUID, FK)
- image_url (TEXT, nullable)
- caption (TEXT, nullable)
- created_at (TIMESTAMP)

**likes**
- user_id (UUID, PK)
- post_id (UUID, PK)
- created_at (TIMESTAMP)

**messages**
- id (UUID, PK)
- sender_id (UUID, FK)
- receiver_id (UUID, FK)
- content (TEXT)
- is_read (BOOLEAN)
- deleted_by_sender (BOOLEAN)
- deleted_by_receiver (BOOLEAN)
- created_at (TIMESTAMP)

**auth_logs**
- id (UUID, PK)
- user_id (UUID, FK)
- status (VARCHAR: 'success' | 'failed')
- ip_address (VARCHAR)
- user_agent (TEXT)
- failure_reason (TEXT, nullable)
- created_at (TIMESTAMP)

**admin_audit_logs**
- id (UUID, PK)
- admin_id (UUID, FK)
- target_user_id (UUID, FK, nullable)
- action (VARCHAR)
- ip_address (VARCHAR)
- metadata (JSONB)
- created_at (TIMESTAMP)

---

## 🚀 Deployment

### Production Environment

**Backend (Railway):**
```
Repository: Vignu-14/ghostline-backend
Build Command: Dockerfile (multi-stage Go build)
Port: 8080
Environment: See environment variables below
```

**Frontend (Vercel):**
```
Repository: Vignu-14/ghostline-frontend
Framework: Vite
Build Command: npm run build
Output Directory: dist
```

### Environment Variables

**Backend (Railway):**
```
PORT=8080
ENVIRONMENT=production
ALLOWED_ORIGIN=https://ghostline-frontend-five.vercel.app
DATABASE_URL=postgresql://...
JWT_SECRET=<64-char-random-string>
SUPABASE_URL=https://...
SUPABASE_SERVICE_KEY=sbp_...
STORAGE_BUCKET_NAME=user-uploads
COOKIE_SECURE=true
```

**Frontend (Vercel):**
```
VITE_API_BASE_URL=https://ghostline-backend-production-a17a.up.railway.app
```

### Deployment Checklist

- [ ] Regenerate Supabase credentials
- [ ] Generate new JWT_SECRET
- [ ] Set ALLOWED_ORIGIN to actual frontend URL
- [ ] Enable HTTPS in production
- [ ] Set up database backups
- [ ] Configure monitoring/alerts
- [ ] Set up CI/CD pipeline
- [ ] Enable WAF on API endpoints

---

## 🔒 Security

### Authentication & Authorization
- JWT tokens with 15-minute expiration
- Bcrypt password hashing (cost factor 12)
- Two-password system for admin impersonation
- Role-based access control (RBAC)

### Input Validation
- Email format validation
- Password strength requirements
- Username format validation
- File size/type validation
- Message content limits (5000 chars)

### API Security
- CORS origin validation
- Rate limiting (login: 5/15min, chat: 10/sec, likes: 100/hour)
- Security headers (CSP, X-Frame-Options, etc.)
- Request logging and audit trails

### Data Protection
- Row-Level Security (RLS) policies
- HTTPS/TLS encryption
- SQLite injection prevention (parameterized queries)
- XSS prevention (input sanitization)
- CSRF prevention (SameSite=Strict cookies)

### Audit Logging
- All login attempts logged
- Admin actions tracked
- IP address and user-agent logged
- Audit logs retained for compliance

### Database Security
- SSL/TLS connections required
- Connection pooling
- Prepared statements
- Sensitive data in environment variables

---

## 📊 Performance

### Optimization Techniques
- Image lazy loading
- Infinite scroll pagination
- WebSocket for real-time updates
- Database connection pooling
- Response caching
- Frontend code splitting
- CDN for static assets (Vercel)

### Metrics
- **Login Time:** ~200ms (Bcrypt cost 12)
- **Message Delivery:** <100ms (WebSocket)
- **Feed Load:** ~500ms (with pagination)
- **Image Upload:** ~2s (depends on size)

---

## 🧪 Testing

### Running Tests

**Backend:**
```bash
cd backend
go test ./...              # Run all tests
go test -v ./...          # Verbose output
go test -cover ./...      # With coverage
```

**Frontend:**
```bash
cd frontend
npm test                   # Run test suite
npm run test:coverage      # Coverage report
```

---

## 📝 Contributing

### Code Style

**Go:**
- Use `gofmt` for formatting
- Follow Go conventions
- Add error handling everywhere
- Write meaningful comments

**TypeScript:**
- Use ESLint configuration
- Prefer `const` over `let`
- Add proper type annotations
- Use meaningful variable names

### Git Workflow

1. Create feature branch: `git checkout -b feature/your-feature`
2. Make changes and commit: `git commit -am "Add feature"`
3. Push to branch: `git push origin feature/your-feature`
4. Create Pull Request

### Commit Messages

```
<type>: <subject>

<body>

<footer>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Example:
```
feat: Add real-time notifications

Add WebSocket endpoint for delivery notifications
and update UI to show online status.

Closes #123
```

---

## 📄 License

This project is licensed under the MIT License - see LICENSE file for details.

---

## � Troubleshooting

**Database connection failed:**
```bash
# Check DATABASE_URL
psql "$DATABASE_URL"

# Use Supabase connection pooler
# Change host to: db.xxx.supabase.co:6543
# Add: ?options=connection_limit%3D1
```

**CORS error:**
```bash
# Update ALLOWED_ORIGIN in environment variables
# Must match frontend URL exactly (with https://)
# Redeploy after changing
```

**WebSocket not connecting:**
```bash
# Check JWT token is valid
# Verify WebSocket endpoint is /ws/chat
# Check browser console for errors
```

See [TROUBLESHOOTING.md](./docs/TROUBLESHOOTING.md) for more issues.

---

## 🌐 Deployment

### Deploy to Railway (Recommended)

1. **Push code to GitHub**
   ```bash
   git push origin main
   ```

2. **Deploy to Railway:**
   - Go to [railway.app](https://railway.app)
   - Create new project
   - Connect GitHub repository
   - Railway auto-detects Go project
   - Set environment variables (see above)
   - Deploy

3. **Monitor:**
   - Check Railway logs for errors
   - Test API with `/health` endpoint

For detailed deployment instructions, see [DEPLOYMENT.md](./docs/DEPLOYMENT.md)

---

## 📈 Performance

- ✅ Connection pooling (min 5, max 25 connections)
- ✅ Efficient database queries with proper indexing
- ✅ WebSocket for real-time (no polling)
- ✅ Middleware for gzip compression
- ✅ Optimized request/response handling

---

## 📚 Documentation

- [API_DOCUMENTATION.md](./docs/API_DOCUMENTATION.md) - Complete endpoint reference
- [ARCHITECTURE.md](./docs/ARCHITECTURE.md) - System architecture and design
- [DEPLOYMENT.md](./docs/DEPLOYMENT.md) - Production setup guide
- [DEVELOPMENT.md](./docs/DEVELOPMENT.md) - Local development setup
- [WEBSOCKET_PROTOCOL.md](./docs/WEBSOCKET_PROTOCOL.md) - Real-time messaging
- [TROUBLESHOOTING.md](./docs/TROUBLESHOOTING.md) - Common issues

---

## 🤝 Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see [LICENSE](../LICENSE) file for details.

---

## 👨‍💻 Author

Vignu Pandey - [@Vignu-14](https://github.com/Vignu-14)

---

## 🙏 Acknowledgments

- Go and Fiber communities for excellent tools
- PostgreSQL for reliable database
- Supabase for backend services
- Railway for easy deployment

---

## 📞 Support

Need help? Check:
- [API Documentation](./docs/API_DOCUMENTATION.md) - Endpoint reference
- [Development Guide](./docs/DEVELOPMENT.md) - Local setup
- [Troubleshooting Guide](./docs/TROUBLESHOOTING.md) - Common issues
- [Deployment Guide](./docs/DEPLOYMENT.md) - Production setup

---

**Happy coding! 🚀**
