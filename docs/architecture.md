# 🏗️ Architecture

Ghostline uses a **REST API + WebSocket hybrid architecture** for optimal performance.

---

## 📊 System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                      CLIENT LAYER                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Browser (React + TypeScript)                                   │
│  - Components                                                   │
│  - State Management (Context API)                               │
│  - WebSocket Connection                                         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              ↑
                ┌─────────────┴─────────────┐
                │                           │
          ┌─────▼───────┐          ┌───────▼──────┐
          │  REST API   │          │  WebSocket   │
          │  (HTTP)     │          │  (WS)        │
          └─────┬───────┘          └───────┬──────┘
                │                           │
┌───────────────▼───────────────────────────▼──────────────────────┐
│                    API LAYER (Go Fiber)                          │
├────────────────────────────────────────────────────────────────-─┤
│                                                                 │
│  Middleware Layer                                                │
│  ├─ CORS                                                         │
│  ├─ JWT Authentication                                           │
│  ├─ Rate Limiting                                                │
│  ├─ Error Handling                                               │
│  └─ Request Logging                                              │
│                                                                 │
│  Routes                                                          │
│  ├─ /api/auth/*     → AuthHandler                                │
│  ├─ /api/users/*    → UserHandler                                │
│  ├─ /api/posts/*    → PostHandler                                │
│  ├─ /api/messages/* → ChatHandler                                │
│  ├─ /api/admin/*    → AdminHandler                               │
│  └─ /ws/chat        → WebSocketHandler                           │
│                                                                 │
│  Service Layer                                                   │
│  ├─ AuthService                                                  │
│  ├─ UserService                                                  │
│  ├─ PostService                                                  │
│  ├─ ChatService                                                  │
│  ├─ UploadService                                                │
│  └─ ImpersonationService                                         │
│                                                                 │
│  WebSocket Hub (Real-time)                                       │
│  ├─ Client Manager                                               │
│  ├─ Message Broadcaster                                          │
│  └─ Connection Pool                                              │
│                                                                 │
└───────────────┬─────────────────────────────────────────────────┘
                │
    ┌───────────┴──────────────┬─────────────────┐
    │                          │                 │
┌───▼──────┐          ┌────────▼─────┐    ┌─────▼─────┐
│PostgreSQL│          │Supabase Auth  │    │Supabase   │
│(Database)│          │(Row Security) │    │Storage    │
└──────────┘          └───────────────┘    └───────────┘
```

---

## 🔄 Request Flow Example

### 1. **User Login Flow**

```
1. Frontend sends: POST /api/auth/login
   ↓
2. CORS Middleware: Check origin ✓
   ↓
3. Rate Limiter: Check login attempts ✓
   ↓
4. AuthHandler.Login()
   ↓
5. AuthService.Login()
   - UserRepository.FindByUsername()
   - Compare password with Bcrypt
   ↓
6. Generate JWT token
   ↓
7. Create auth_token cookie (HTTPOnly, Secure, SameSite=Strict)
   ↓
8. Return user data + set cookie
   ↓
9. Frontend stores user in AuthContext
```

### 2. **Real-Time Message Flow**

```
1. Frontend connects: WebSocket /ws/chat
   ↓
2. WebSocketHandler.Upgrade()
   - Validate JWT from cookie ✓
   ↓
3. Create Client, register with Hub
   ↓
4. Frontend sends: {type: "message", receiver_id: "...", content: "..."}
   ↓
5. WebSocketHandler.HandleConnection()
   - Rate limit check (10 msg/sec) ✓
   ↓
6. ChatService.SendMessage()
   - Save to database
   - Sanitize content (Bluemonday)
   ↓
7. Hub.Broadcast(receiver_id, message)
   - Find all connections for receiver
   - Send message to each WebSocket
   ↓
8. Frontend receives in real-time
```

### 3. **File Upload Flow**

```
1. Frontend: PUT /api/posts/upload-url
   ↓
2. Rate limiter: Check upload limit ✓
   ↓
3. Validate file (size, type, extension)
   ↓
4. UploadService.CreateSignedUploadURL()
   - Call Supabase API
   - Generate signed URL
   ↓
5. Return signed URL to frontend
   ↓
6. Frontend sends file directly to Supabase
   (bypasses backend, faster)
   ↓
7. Supabase stores file
   ↓
8. Frontend calls: POST /api/posts/finalize
   - Verify file ownership
   - Create post record
   ↓
9. Post visible in feed
```

---

## 📁 Backend Project Structure

```
ghostline-backend/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
│
├── internal/
│   ├── config/
│   │   ├── config.go               # Configuration loading
│   │   ├── constants.go            # Default values
│   │   ├── database.go             # DB config
│   │   └── storage.go              # Subabase config
│   │
│   ├── database/
│   │   ├── postgres.go             # DB connection
│   │   └── migrations/             # SQL migrations
│   │
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── post_handler.go
│   │   ├── chat_handler.go
│   │   ├── like_handler.go
│   │   ├── admin_handler.go
│   │   ├── health_handler.go
│   │   └── websocket_handler.go
│   │
│   ├── middleware/
│   │   ├── jwt_middleware.go       # Auth validation
│   │   ├── cors_middleware.go      # CORS headers
│   │   ├── rate_limiter.go         # Rate limiting
│   │   ├── error_handler.go        # Error formatting
│   │   ├── secure_headers.go       # Security headers
│   │   └── request_logger.go       # Logging
│   │
│   ├── models/
│   │   ├── user.go
│   │   ├── post.go
│   │   ├── message.go
│   │   ├── like.go
│   │   ├── jwt_claims.go
│   │   └── errors.go
│   │
│   ├── repositories/
│   │   ├── user_repository.go      # DB queries
│   │   ├── post_repository.go
│   │   ├── message_repository.go
│   │   ├── like_repository.go
│   │   ├── auth_log_repository.go
│   │   └── admin_repository.go
│   │
│   ├── services/
│   │   ├── auth_service.go         # Business logic
│   │   ├── user_service.go
│   │   ├── post_service.go
│   │   ├── chat_service.go
│   │   ├── upload_service.go
│   │   ├── like_service.go
│   │   └── impersonation_service.go
│   │
│   ├── routes/
│   │   ├── routes.go               # Route registration
│   │   ├── api.go                  # API routes
│   │   └── websocket.go            # WebSocket routes
│   │
│   ├── utils/
│   │   ├── jwt.go                  # Token generation
│   │   ├── bcrypt.go               # Password hashing
│   │   ├── sanitizer.go            # XSS prevention
│   │   ├── validator.go            # Input validation
│   │   ├── response.go             # JSON responses
│   │   └── uuid.go                 # UUID utilities
│   │
│   └── websocket/
│       ├── hub.go                  # Connection manager
│       ├── client.go               # Client connection
│       ├── message.go              # Message types
│       └── broadcaster.go          # Message broadcasting
│
├── pkg/
│   ├── cache/                      # Redis (future)
│   └── logger/                     # Structured logging
│
├── tests/
│   ├── unit/                       # Unit tests
│   ├── integration/                # Integration tests
│   ├── mocks/                      # Mock objects
│   └── fixtures/                   # Test data
│
├── scripts/
│   ├── migrate.sh                  # Run migrations
│   ├── seed.sh                     # Seed database
│   └── generate_jwt_secret.sh      # Generate secrets
│
├── docs/
│   ├── API_DOCUMENTATION.md        # Endpoint reference
│   ├── ARCHITECTURE.md             # This file
│   ├── DEPLOYMENT.md               # Deployment guide
│   ├── DEVELOPMENT.md              # Setup guide
│   ├── WEBSOCKET.md                # WebSocket protocol
│   └── DATABASE_SCHEMA.md          # DB structure
│
├── Dockerfile                      # Docker build
├── docker-compose.yml              # Local dev setup
├── go.mod / go.sum                 # Dependencies
├── Makefile                        # Build commands
└── README.md                       # Project overview
```

---

## 🗄️ Database Architecture

**PostgreSQL on Supabase with Row Level Security (RLS)**

### Tables:
- **users** - User accounts
- **posts** - Photo/text posts
- **likes** - Post engagement
- **messages** - Direct messages
- **auth_logs** - Login audit trail
- **admin_audit_logs** - Admin actions

### Security:
- RLS policies enforce row-level access
- All queries use parameterized statements
- SSL/TLS connections required

---

## 🔐 Authentication Architecture

**JWT + Secure Cookies + Step-Up Auth**

```
┌──────────────────────────────────────────┐
│ Normal User Login                        │
├──────────────────────────────────────────┤
│ 1. Username + Password                   │
│ 2. Validate password (Bcrypt)            │
│ 3. Generate JWT token                    │
│ 4. Set HTTPOnly Secure SameSite cookie   │
│ 5. Return user data                      │
└──────────────────────────────────────────┘

┌──────────────────────────────────────────┐
│ Admin Impersonation (God Mode)           │
├──────────────────────────────────────────┤
│ 1. Admin has valid JWT                   │
│ 2. Submit impersonation_password         │
│ 3. Validate against impersonation_hash   │
│ 4. Generate new JWT with impersonator_id │
│ 5. Log to admin_audit_logs               │
│ 6. Admin now acts as target user         │
└──────────────────────────────────────────┘
```

---

## 🚀 Performance Optimizations

1. **Database Indexing**
   - Index on `username`, `email`
   - Index on `post_created_at` (DESC)
   - Composite index on `(sender_id, receiver_id, created_at)`

2. **Connection Pooling**
   - pgx pool with 25 max connections
   - 5 min connections
   - 1 hour max lifetime

3. **WebSocket**
   - Hub-based broadcasting (O(n) where n = connections)
   - Message buffering
   - Automatic cleanup on disconnect

4. **Caching** (Future)
   - Redis for frequently accessed data
   - Cache invalidation on updates

---

## 🔄 Data Flow Diagram

```
User Input
    ↓
Frontend (React)
    ├─ REST API (Auth, Posts, Users)
    │  └─ HTTP → Go Handler
    │           → Service
    │           → Repository
    │           → PostgreSQL
    │
    └─ WebSocket (Messages)
       └─ WS → Go Handler
                → Hub → All Clients
                → Service → PostgreSQL
                → Broadcast response
```

---

## 🛡️ Security Layers

1. **Transport Layer**
   - HTTPS only in production
   - TLS 1.3

2. **Application Layer**
   - JWT token validation
   - Rate limiting
   - Input sanitization
   - CORS validation

3. **Database Layer**
   - Row Level Security (RLS)
   - Parameterized queries
   - Connection SSL/TLS

4. **Operational Layer**
   - Audit logging
   - Security headers
   - Error message sanitization

---

See also: [DEPLOYMENT.md](./DEPLOYMENT.md), [DATABASE_SCHEMA.md](./DATABASE_SCHEMA.md)
