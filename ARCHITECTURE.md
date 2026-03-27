# 🏗️ Architecture Documentation

## System Architecture

### Overview

Ghostline is a full-stack web application with three main components:

1. **Frontend** - React/TypeScript SPA served on Vercel
2. **Backend** - Go/Fiber RESTful API on Railway
3. **Database** - PostgreSQL hosted on Supabase

---

## High-Level Architecture Diagram

```
┌─────────────────────────────────────────────────────────┐
│                    Client Layer                          │
│  Web Browser (Chrome, Firefox, Safari, Edge)             │
│  ├─ React 19 Application                                │
│  ├─ TypeScript Type Safety                              │
│  ├─ Tailwind CSS Styling                                │
│  └─ Context API State Management                        │
└────────────────┬────────────────────────────────────────┘
                 │
              HTTP/1.1 + WebSocket
             (HTTPS in production)
                 │
┌────────────────▼────────────────────────────────────────┐
│                 API Gateway Layer                        │
│  (Optional: Nginx reverse proxy in production)           │
│  ├─ Load Balancing                                       │
│  ├─ SSL/TLS Termination                                 │
│  ├─ Rate Limiting                                       │
│  └─ Request Logging                                     │
└────────────────┬────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────┐
│               Application Layer                          │
│  Go Fiber Web Framework                                  │
│  ├─ HTTP Request Handlers                               │
│  ├─ Middleware Stack                                    │
│  │  ├─ CORS Middleware                                 │
│  │  ├─ JWT Authentication                              │
│  │  ├─ Rate Limiter                                    │
│  │  └─ Error Handler                                   │
│  ├─ WebSocket Hub                                      │
│  └─ Route Definitions                                  │
└────────────────┬────────────────────────────────────────┘
                 │
    ┌────────────┼────────────┐
    │            │            │
    ▼            ▼            ▼
┌────────┐  ┌──────────┐  ┌──────────┐
│Business│  │Repository│  │  WebSocket
│ Logic  │  │  Layer   │  │   Hub
│Layer   │  │(Database)│  │
└────────┘  └──────────┘  └──────────┘
    │            │            │
    └────────────┼────────────┘
                 │
┌────────────────▼────────────────────────────────────────┐
│               Data Layer                                 │
│  ├─ PostgreSQL (Supabase)                               │
│  │  ├─ User Accounts                                    │
│  │  ├─ Posts & Likes                                    │
│  │  ├─ Messages & Conversations                         │
│  │  ├─ Audit Logs                                       │
│  │  └─ Row Level Security (RLS)                         │
│  └─ Supabase Storage (S3-compatible)                    │
│     ├─ User Uploads                                     │
│     └─ Profile Pictures                                 │
└────────────────────────────────────────────────────────┘
```

---

## Component Architecture

### Frontend Components

```
App (Root)
├── AuthContext Provider
│   ├── ChatContext Provider
│   │   ├── NotificationContext Provider
│   │   │   ├── Router
│   │   │   │   ├── HomePage
│   │   │   │   │   ├── PostFeed
│   │   │   │   │   │   └── PostCard (repeating)
│   │   │   │   │   └── UserSearchPanel
│   │   │   │   ├── LoginPage
│   │   │   │   ├── RegisterPage
│   │   │   │   ├── ChatPage
│   │   │   │   │   ├── ChatList
│   │   │   │   │   └── ChatWindow
│   │   │   │   ├── ProfilePage
│   │   │   │   ├── AdminPage
│   │   │   │   └── NotFoundPage
│   │   │   ├── Navbar
│   │   │   ├── MobileNav
│   │   │   └── Footer
│   │   │   └── Notifications (Toast)
│   │   └── WebSocket Hook
│   └── Auth Services
└── Theme Provider
```

### Backend Handlers

```
Router
├── /health
│   ├── GET /health → HealthHandler.Live()
│   └── GET /health/ready → HealthHandler.Ready()
├── /api
│   ├── /auth
│   │   ├── POST /register → AuthHandler.Register()
│   │   ├── POST /login → AuthHandler.Login()
│   │   ├── POST /logout → AuthHandler.Logout()
│   │   └── GET /me → UserHandler.Me()
│   ├── /users
│   │   ├── GET /me → UserHandler.Me()
│   │   ├── GET /profile/:username → UserHandler.Profile()
│   │   └── GET /search → UserHandler.Search()
│   ├── /posts
│   │   ├── GET / → PostHandler.List()
│   │   ├── POST / → PostHandler.Create()
│   │   ├── POST /upload-url → PostHandler.CreateUploadURL()
│   │   ├── POST /finalize → PostHandler.CreateFromUploadedObject()
│   │   ├── DELETE /:id → PostHandler.Delete()
│   │   ├── POST /:id/like → LikeHandler.Like()
│   │   └── DELETE /:id/like → LikeHandler.Unlike()
│   ├── /messages
│   │   ├── GET /conversations → ChatHandler.ListConversations()
│   │   ├── GET /:userId → ChatHandler.GetConversation()
│   │   ├── POST / → ChatHandler.SendMessage()
│   │   ├── POST /delete → ChatHandler.DeleteMessages()
│   │   └── POST /:userId/clear → ChatHandler.ClearConversation()
│   ├── /ws/chat (WebSocket)
│   │   └── GET / → WebSocketHandler.Upgrade() → HandleConnection()
│   └── /admin
│       ├── POST /impersonate → AdminHandler.Impersonate()
│       └── POST /impersonate/stop → AdminHandler.StopImpersonation()
```

---

## Data Flow Diagrams

### Authentication Flow

```
┌──────────────────────────────────────────────────────────────┐
│ Frontend (React)                                              │
├──────────────────────────────────────────────────────────────┤
1. User enters credentials
2. onClick → login() → POST /api/auth/login
                              │
                              ▼
┌──────────────────────────────────────────────────────────────┐
│ Backend (Go/Fiber)                                            │
├──────────────────────────────────────────────────────────────┤
3. AuthHandler.Login()
4. Validate credentials
5. Hash password with bcrypt (cost 12)
6. Compare with stored hash
7. Generate JWT token
   ├─ user_id
   ├─ role
   ├─ expiration (15 min)
   └─ signature (HMAC-SHA256)
8. Set HTTPOnly Secure SameSite cookie
9. Return user profile
                              │
                              ▼
┌──────────────────────────────────────────────────────────────┐
│ Frontend (React)                                              │
├──────────────────────────────────────────────────────────────┤
10. AuthContext updated
11. User logged in
12. Redirect to home page
13. Subsequent requests include cookie automatically
```

### Post Upload Flow

```
User selects image
       │
       ▼
Validate (size, type)
       │
       ▼
Request signed upload URL
POST /api/posts/upload-url
       │
       ▼
Backend generates URL
   ├─ userID-UUID-extension
   └─ Signed with Supabase key
       │
       ▼
Return URL to frontend
       │
       ▼
Frontend uploads directly to Supabase
PUT <signed-url>
       │
       ▼
Supabase stores object
       │
       ▼
Frontend creates post
POST /api/posts/finalize
   ├─ object_path (from upload)
   └─ caption
       │
       ▼
Backend validates path ownership
       │
       ▼
Create post in database
       │
       ▼
Return post to frontend
       │
       ▼
Show in feed
```

### WebSocket Flow (Real-time Chat)

```
┌──────────────┐              ┌──────────────┐
│   User A     │              │   User B     │
│  (Vercel)    │              │  (Vercel)    │
└──────┬───────┘              └──────┬───────┘
       │                              │
       │ 1. ws://backend/ws/chat      │
       │ (with JWT cookie)            │
       │                              │
       └──────────────┬───────────────┘
                      │
                      ▼
        ┌─────────────────────────────┐
        │   Backend WebSocket Hub     │
        │   (Go Fiber)                │
        ├─────────────────────────────┤
        │ 2. Validate JWT token       │
        │ 3. Extract user_id          │
        │ 4. Register connection      │
        │ 5. Send "connected" event   │
        └─────────────────────────────┘
                      │
        ┌─────────────┴──────────────┐
        │                            │
   User A connection            User B connection
        │                            │
        │ 6. Send message            │
        │ {"type":"message",         │
        │  "receiver_id":"...",      │
        │  "content":"hello"}        │
        │                            │
        └──────────┬─────────────────┘
                   │
                   ▼
        ┌──────────────────────────────┐
        │ Backend Processing           │
        ├──────────────────────────────┤
        │ 7. Rate limit check          │
        │ 8. Sanitize content          │
        │ 9. Save to database          │
        │ 10. Broadcast to User B      │
        └──────────────────────────────┘
                   │
        ┌──────────┴──────────────┐
        │                         │
   Send to User A            Send to User B
        │                         │
        ▼                         ▼
    Update UI                 Update UI
 sent indicator           New message notification
```

---

## Request/Response Flow

### Example: Create Post

```
REQUEST (Frontend):
POST /api/posts
Content-Type: application/json
Cookie: auth_token=eyJhbGciOiJIUzI1NiIs...

{
  "caption": "Beautiful sunset!",
  "image_url": "posts/user-id/uuid.jpg"
}

┌────────────────────────────────────────┐
│ Backend Processing                      │
├────────────────────────────────────────┤
1. Parse request body
2. Extract from JWT cookie
3. Validate input
   - caption required (1-500 chars)
   - image_url must be object path
4. Verify user owns image (path check)
5. Sanitize caption (XSS prevention)
6. Create post record
   INSERT INTO posts (...)
7. Return created post
└────────────────────────────────────────┘

RESPONSE:
200 OK
Content-Type: application/json

{
  "status": "success",
  "data": {
    "post": {
      "id": "uuid",
      "user_id": "uuid",
      "username": "john",
      "profile_picture_url": "...",
      "image_url": "posts/user-id/uuid.jpg",
      "caption": "Beautiful sunset!",
      "like_count": 0,
      "is_liked_by_user": false,
      "created_at": "2026-03-28T12:00:00Z"
    }
  }
}
```

---

## Database Schema Relationships

```
┌──────────────┐
│              │
│    users     │◄───────────────────────┐
│              │                        │
│ - id (PK)    │                        │
│ - username   │                        │
│ - email      │                        │
│ - password   │                        │
│ - role       │                        │
│ - created_at │                        │
└──────────────┘                        │
       │                        one-to-many
       │                                │
       │ one-to-many                    │
       │                                │
       ├──────────────────────────────┐ │
       │      │                       │ │
       ▼      ▼                       │ ▼
    ┌─────────────┐            ┌──────────────┐
    │   posts     │            │ admin_audit  │
    │             │            │   _logs      │
    │ - id (PK)   │            │              │
    │ - user_id   │────────────│ - admin_id   │
    │ - image_url │            │ - target_id  │
    │ - caption   │            │ - action     │
    │ - created   │            │ - created_at │
    └─────────────┘            └──────────────┘
         │
         │ one-to-many
         │
         ▼
    ┌──────────────┐
    │    likes     │
    │              │
    │ - user_id    │◄────────one-to-many
    │ - post_id    │  from users
    │ - created_at │
    └──────────────┘

┌──────────────────┐
│    messages      │
│                  │
│ - id (PK)        │
│ - sender_id ────┐│      (FK to users)
│ - receiver_id──┐││      (FK to users)
│ - content      │││
│ - is_read      │││
│ - created_at   │││
└──────────────────┘│
                    │
                    │
┌──────────────┐    │
│     users    │◄───┘
│              │
│ - id (many)  │
└──────────────┘
```

---

## Middleware Stack

### Request Processing Order

```
Incoming Request
    │
    ▼
1. Request ID Middleware
   └─ Assigns unique request ID
    │
    ▼
2. Request Logger Middleware
   └─ Logs incoming request details
    │
    ▼
3. Recovery Middleware
   └─ Catches panics, returns error
    │
    ▼
4. Secure Headers Middleware
   └─ Sets security headers
       ├─ Content-Security-Policy
       ├─ X-Frame-Options
       ├─ X-Content-Type-Options
       └─ Strict-Transport-Security
    │
    ▼
5. CORS Middleware
   └─ Validates origin
    │
    ▼
6. Rate Limiter Middleware (if needed)
   └─ Checks rate limits
    │
    ▼
7. JWT Middleware (if protected route)
   └─ Validates token
   └─ Extracts user_id
   └─ Attaches to context
    │
    ▼
8. Admin Middleware (if admin route)
   └─ Checks role == "admin"
   └─ Blocks impersonated users
    │
    ▼
9. Route Handler
   └─ Business logic
   └─ Database operations
   └─ Response generation
    │
    ▼
Response sent to client
```

---

## Error Handling Flow

```
Request Processing
    │
    ▼
Error Occurs
    │
    ├─ Validation Error
    │  ├─ Return 400 Bad Request
    │  └─ Include field errors
    │
    ├─ Authentication Error
    │  ├─ Return 401 Unauthorized
    │  └─ Clear cookie
    │
    ├─ Authorization Error
    │  ├─ Return 403 Forbidden
    │  └─ Log attempt
    │
    ├─ Not Found Error
    │  ├─ Return 404 Not Found
    │  └─ Don't expose path
    │
    ├─ Rate Limit Error
    │  ├─ Return 429 Too Many Requests
    │  └─ Include Retry-After header
    │
    └─ Server Error
       ├─ Return 500 Internal Server Error
       ├─ Log full error
       └─ Hide details from user
    │
    ▼
Standardized JSON Response
{
  "status": "error",
  "error": "User-friendly message",
  "details": {...}  // If validation error
}
```

---

## Security Architecture

### CSRF Protection (3-layer approach)

```
Layer 1: HTTPOnly Cookie
├─ JavaScript cannot read cookie
├─ Prevents XSS theft
└─ Only sent in requests from same origin

Layer 2: SameSite=Strict
├─ Browser won't send cookie in cross-site requests
├─ Blocks form submissions from external sites
└─ Blocks CORS requests from external domains

Layer 3: CORS Origin Validation
├─ Backend checks Origin header
├─ Only allows ghostline-frontend-five.vercel.app
└─ Blocks requests from evil.com
```

### authentication Architecture

```
Step 1: User Registration
   ├─ Validate password (8+ chars, numbers, symbols)
   ├─ Hash with Bcrypt (cost 12)
   └─ Store hash (never store plaintext)

Step 2: User Login
   ├─ Lookup user by username
   ├─ Compare password with hash
   ├─ Generate JWT token
   │  ├─ Header: { alg: HS256, typ: JWT }
   │  ├─ Payload: { user_id, role, exp: +15min }
   │  └─ Signature: HMAC-SHA256(secret)
   └─ Set HttpOnly Secure SameSite cookie

Step 3: Protected Requests
   ├─ Client includes auth_token cookie
   ├─ Backend extracts token
   ├─ Verify signature with secret
   ├─ Check expiration
   └─ Extract user_id from payload
```

---

## Database Security

```
Layer 1: Connection Security
├─ SSL/TLS encryption (sslmode=require)
├─ Secure password authentication
└─ Connection pooling (25 connections max)

Layer 2: Row Level Security (RLS)
├─ users can only SELECT their own row
├─ users can only UPDATE/DELETE their own posts
├─ messages only visible to sender/receiver
└─ Enforced at database level (not app)

Layer 3: Input Validation
├─ Parameterized queries ($1, $2, ...)
├─ No string concatenation
└─ Type-safe query building

Layer 4: Access Control
├─ Foreign key constraints
├─ ON DELETE CASCADE for cleanup
└─ UUID primary keys (prevent enumeration)
```

---

## Scalability Considerations

### Horizontal Scaling

**Backend (Stateless):**
- Each instance is independent
- Can run 1, 10, or 100 instances
- Railway auto-scales on CPU/memory
- No shared state between instances

**Frontend (Static files):**
- Built as static SPA
- Served from CDN (Vercel)
- Infinitely scalable
- Cached at edge servers

**Database (Vertical):**
- PostgreSQL handles connections
- Connection pooling (25 max)
- RLS prevents data leaks
- Supabase handles backups

### Performance Optimization

**Caching:**
- Browser cache for static assets
- Session cookies reduce re-auth
- Database indexes on frequent queries

**Pagination:**
- Feed loaded 20 posts at a time
- Messages loaded 50 at a time
- Prevents full table scans

**Rate Limiting:**
- Prevents abuse
- Protects API from DOS
- Enforced per user

---

## Monitoring & Logging

### Logging Strategy

```
Request Lifecycle Logging:
1. Request received
   └─ Method, path, client IP, user-agent
2. Processing
   └─ Middleware execution, sanitization
3. Database queries
   └─ Query type, duration, rows affected
4. WebSocket events
   └─ Connection/disconnect, message count
5. Errors
   └─ Type, stack trace, context
6. Response sent
   └─ Status code, response time
```

### Audit Logging

```
User Actions:
├─ Successful login
├─ Failed login (with IP)
├─ Admin impersonation attempts
├─ Sensitive operations
└─ All stored in audit_logs table

Retention:
├─ 30 days for security incidents
├─ 90 days for compliance
└─ Archive older logs
```

---

## Deployment Architecture

```
┌─────────────────────────────────────────┐
│          GitHub Repository              │
│  Vignu-14/ghostline-{backend|frontend}  │
└────────────┬────────────────────────────┘
             │
    ┌────────┴─────────┐
    │                  │
    ▼                  ▼
┌──────────┐      ┌──────────┐
│ Railway  │      │  Vercel  │
│          │      │          │
│Build:    │      │Build:    │
│Dockerfile│      │npm build │
│                 │
│Run:      │      │Deploy:   │
│Server    │      │Auto      │
│on :8080  │      │EdgeCDN   │
└──────────┘      └──────────┘
    │                  │
    ├─────────┬────────┤
    │         │        │
    ▼         ▼        ▼
Supabase    Router    Client
│           (DNS)     (Browser)
├─ DB
├─ Storage
└─ Logs
```

---

This architecture ensures:
- ✅ Scalability (horizontal backend, global CDN)
- ✅ Performance (caching, pagination, optimization)
- ✅ Security (encryption, auth, rate limiting)
- ✅ Reliability (monitoring, logging, backups)
- ✅ Maintainability (clear separation of concerns)
