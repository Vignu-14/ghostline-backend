# 📡 API Documentation

**Ghostline Backend API** - Complete endpoint reference

Base URL: `https://ghostline-backend-production-a17a.up.railway.app`

---

## 🔐 Authentication

All protected endpoints require a valid JWT token in the `auth_token` cookie.

The token is automatically set after login and included in all subsequent requests.

---

## 📝 Auth Endpoints

### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "passwordConfirm": "SecurePass123!"
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2026-03-28T10:00:00Z"
    }
  }
}
```

### Login User
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "SecurePass123!"
}
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2026-03-28T10:00:00Z"
    }
  }
}
```

**Cookies Set:**
- `auth_token`: JWT token (HTTPOnly, Secure, SameSite=Strict)

### Get Current User
```http
GET /api/auth/me
Authorization: Cookie auth_token=<token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user",
      "profile_picture_url": "https://...",
      "created_at": "2026-03-28T10:00:00Z"
    }
  }
}
```

### Logout User
```http
POST /api/auth/logout
Authorization: Cookie auth_token=<token>
```

**Response (200):**
```json
{
  "status": "success",
  "message": "logout successful"
}
```

---

## 👥 User Endpoints

### Get Profile
```http
GET /api/users/profile/:username
Authorization: Cookie auth_token=<token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "profile_picture_url": "https://...",
      "created_at": "2026-03-28T10:00:00Z"
    },
    "posts": [
      {
        "id": "post-uuid",
        "caption": "Beautiful sunset",
        "image_url": "https://...",
        "like_count": 5,
        "created_at": "2026-03-28T09:30:00Z"
      }
    ]
  }
}
```

### Search Users
```http
GET /api/users/search?q=john&limit=10
Authorization: Cookie auth_token=<token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "results": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "username": "john_doe",
        "profile_picture_url": "https://..."
      }
    ]
  }
}
```

---

## 📸 Post Endpoints

### Get Feed
```http
GET /api/posts?page=1&limit=20
Authorization: Cookie auth_token=<token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "posts": [
      {
        "id": "post-uuid",
        "user_id": "550e8400-e29b-41d4-a716-446655440000",
        "username": "john_doe",
        "caption": "Beautiful sunset",
        "image_url": "https://...",
        "like_count": 5,
        "user_liked": false,
        "created_at": "2026-03-28T09:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100
    }
  }
}
```

### Create Post
```http
POST /api/posts
Content-Type: application/json
Authorization: Cookie auth_token=<token>

{
  "caption": "Beautiful sunset",
  "image_url": "https://supabase.../posts/uuid/filename.jpg"
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "post": {
      "id": "post-uuid",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "caption": "Beautiful sunset",
      "image_url": "https://...",
      "like_count": 0,
      "user_liked": false,
      "created_at": "2026-03-28T10:00:00Z"
    }
  }
}
```

### Get Upload URL
```http
POST /api/posts/upload-url
Content-Type: application/json
Authorization: Cookie auth_token=<token>

{
  "file_name": "photo.jpg",
  "file_size": 1024000,
  "file_type": "image/jpeg"
}
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "upload_url": "https://supabase...signed-url",
    "object_path": "posts/user-uuid/uuid.jpg"
  }
}
```

### Delete Post
```http
DELETE /api/posts/:postId
Authorization: Cookie auth_token=<token>
```

**Response (200):**
```json
{
  "status": "success",
  "message": "post deleted successfully"
}
```

---

## ❤️ Like Endpoints

### Like Post
```http
POST /api/posts/:postId/like
Authorization: Cookie auth_token=<token>
```

**Response (200):**
```json
{
  "status": "success",
  "message": "post liked successfully"
}
```

### Unlike Post
```http
DELETE /api/posts/:postId/like
Authorization: Cookie auth_token=<token>
```

**Response (200):**
```json
{
  "status": "success",
  "message": "post unliked successfully"
}
```

---

## 💬 Message Endpoints

### Get Conversations
```http
GET /api/messages/conversations?page=1&limit=50
Authorization: Cookie auth_token=<token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "conversations": [
      {
        "user_id": "550e8400-e29b-41d4-a716-446655440000",
        "username": "john_doe",
        "last_message": "Hey, how are you?",
        "last_message_time": "2026-03-28T10:00:00Z",
        "unread_count": 3
      }
    ]
  }
}
```

### Get Messages
```http
GET /api/messages/:userId?page=1&limit=50
Authorization: Cookie auth_token=<token>
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "messages": [
      {
        "id": "msg-uuid",
        "sender_id": "550e8400-e29b-41d4-a716-446655440000",
        "receiver_id": "660e8400-e29b-41d4-a716-446655440001",
        "content": "Hey, how are you?",
        "is_read": true,
        "created_at": "2026-03-28T10:00:00Z"
      }
    ]
  }
}
```

### Send Message (HTTP)
```http
POST /api/messages
Content-Type: application/json
Authorization: Cookie auth_token=<token>

{
  "receiver_id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "Hey, how are you?"
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "message": {
      "id": "msg-uuid",
      "sender_id": "...",
      "receiver_id": "...",
      "content": "Hey, how are you?",
      "is_read": false,
      "created_at": "2026-03-28T10:00:00Z"
    }
  }
}
```

### Delete Message
```http
POST /api/messages/delete
Content-Type: application/json
Authorization: Cookie auth_token=<token>

{
  "message_id": "msg-uuid",
  "delete_for_everyone": false
}
```

**Response (200):**
```json
{
  "status": "success",
  "message": "message deleted successfully"
}
```

### Clear Conversation
```http
POST /api/messages/:userId/clear
Authorization: Cookie auth_token=<token>
```

**Response (200):**
```json
{
  "status": "success",
  "message": "conversation cleared successfully"
}
```

---

## 🔧 Admin Endpoints

### Start Impersonation
```http
POST /api/admin/impersonate
Content-Type: application/json
Authorization: Cookie auth_token=<admin_token>

{
  "target_user_id": "550e8400-e29b-41d4-a716-446655440000",
  "impersonation_password": "AdminSecret123!"
}
```

**Response (200):**
```json
{
  "status": "success",
  "message": "impersonation started",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "target_user",
      "role": "user"
    }
  }
}
```

### Stop Impersonation
```http
POST /api/admin/impersonate/stop
Authorization: Cookie auth_token=<admin_token>
```

**Response (200):**
```json
{
  "status": "success",
  "message": "impersonation stopped"
}
```

---

## ✅ Health Endpoints

### Health Check
```http
GET /health
```

**Response (200):**
```json
{
  "status": "ok"
}
```

### Readiness Check
```http
GET /health/ready
```

**Response (200):**
```json
{
  "status": "ready"
}
```

---

## Error Responses

All errors follow this format:

```json
{
  "status": "error",
  "error": "User not found",
  "details": null
}
```

### Common Status Codes
- `200` - OK
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict (username exists)
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error

---

## Rate Limiting

- **Login:** 5 attempts per 15 minutes
- **Upload:** 10 uploads per hour
- **Messages:** 10 messages per second
- **Likes:** 100 likes per hour

---

## WebSocket Connection

See [WEBSOCKET_PROTOCOL.md](../docs/WEBSOCKET_PROTOCOL.md) for real-time messaging details.
