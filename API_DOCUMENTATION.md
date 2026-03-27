# 📡 API Documentation

## Base URL
```
https://ghostline-backend-production-a17a.up.railway.app
```

## Authentication

All protected endpoints require a valid JWT token in the `auth_token` cookie.

The cookie is automatically set after login and included in subsequent requests with `credentials: 'include'`.

---

## Endpoints

### Health Checks

#### Get Health Status
```
GET /health
```

**Response (200):**
```json
{
  "status": "ok"
}
```

---

#### Get Readiness Status
```
GET /health/ready
```

Checks database connectivity.

**Response (200):**
```json
{
  "status": "ready"
}
```

---

## Authentication

### Register New User

```
POST /api/auth/register
Content-Type: application/json
```

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "password_confirm": "SecurePass123!"
}
```

**Validation:**
- Username: 3-50 chars, alphanumeric + underscore
- Email: Valid email format
- Password: 8+ chars, must include uppercase, lowercase, number, symbol

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
      "profile_picture_url": null,
      "created_at": "2026-03-28T12:00:00Z"
    }
  }
}
```

**Errors:**
- `400` - Invalid request or validation failed
- `409` - Username or email already taken

---

### Login

```
POST /api/auth/login
Content-Type: application/json
```

**Request Body:**
```json
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
      "profile_picture_url": null,
      "created_at": "2026-03-28T12:00:00Z"
    }
  }
}
```

**Cookie Set:**
```
Set-Cookie: auth_token=<JWT>; HttpOnly; Secure; SameSite=Strict; Max-Age=900
```

**Errors:**
- `400` - Invalid request
- `401` - Invalid credentials
- `429` - Too many login attempts (5 per 15 min)

---

### Logout

```
POST /api/auth/logout
```

Clears the authentication cookie.

**Response (200):**
```json
{
  "status": "success",
  "data": null
}
```

---

### Get Current User

```
GET /api/auth/me
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
      "profile_picture_url": null,
      "created_at": "2026-03-28T12:00:00Z"
    }
  }
}
```

**Errors:**
- `401` - Not authenticated

---

## Users

### Get Current User Profile

```
GET /api/users/me
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
      "created_at": "2026-03-28T12:00:00Z",
      "post_count": 5,
      "follower_count": 100
    }
  }
}
```

---

### Get User Profile

```
GET /api/users/profile/:username
```

**Parameters:**
- `username` - Username to fetch

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "profile_picture_url": "https://...",
      "created_at": "2026-03-28T12:00:00Z",
      "post_count": 5,
      "follower_count": 100
    },
    "posts": [
      {
        "id": "uuid",
        "image_url": "https://...",
        "caption": "Beautiful sunset",
        "like_count": 42,
        "is_liked_by_user": false,
        "created_at": "2026-03-28T12:00:00Z"
      }
    ]
  }
}
```

**Errors:**
- `404` - User not found

---

### Search Users

```
GET /api/users/search?q=query&limit=20
```

**Query Parameters:**
- `q` - Search query (min 1 char)
- `limit` - Results limit (default: 20, max: 100)

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "users": [
      {
        "id": "uuid",
        "username": "john_doe",
        "profile_picture_url": "https://...",
        "post_count": 5
      }
    ]
  }
}
```

---

## Posts

### Get Feed

```
GET /api/posts?page=1&limit=20
```

**Query Parameters:**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20, max: 100)

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "posts": [
      {
        "id": "uuid",
        "user_id": "uuid",
        "username": "john_doe",
        "profile_picture_url": "https://...",
        "image_url": "https://...",
        "caption": "Beautiful sunset",
        "like_count": 42,
        "is_liked_by_user": false,
        "created_at": "2026-03-28T12:00:00Z"
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

---

### Create Post

```
POST /api/posts
Content-Type: application/json
Authentication: Required
```

**Request Body:**
```json
{
  "caption": "Beautiful sunset!",
  "image_url": "posts/user-id/uuid.jpg"
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "post": {
      "id": "uuid",
      "user_id": "uuid",
      "username": "john_doe",
      "profile_picture_url": "https://...",
      "image_url": "https://...",
      "caption": "Beautiful sunset!",
      "like_count": 0,
      "is_liked_by_user": false,
      "created_at": "2026-03-28T12:00:00Z"
    }
  }
}
```

**Errors:**
- `400` - Invalid request
- `401` - Not authenticated
- `429` - Rate limit exceeded

---

### Get Upload URL

```
POST /api/posts/upload-url
Content-Type: application/json
Authentication: Required
```

**Request Body:**
```json
{
  "filename": "sunset.jpg",
  "content_type": "image/jpeg"
}
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "upload_url": "https://storage.api/...",
    "object_path": "posts/user-id/uuid.jpg"
  }
}
```

**Errors:**
- `400` - Invalid file type
- `401` - Not authenticated
- `413` - File too large (>5MB)
- `429` - Rate limit exceeded

---

### Finalize Post Upload

```
POST /api/posts/finalize
Content-Type: application/json
Authentication: Required
```

**Request Body:**
```json
{
  "object_path": "posts/user-id/uuid.jpg",
  "caption": "My photo"
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "post": { /* full post object */ }
  }
}
```

**Errors:**
- `400` - Invalid path
- `401` - Not authenticated
- `403` - Path doesn't belong to user

---

### Delete Post

```
DELETE /api/posts/:id
Authentication: Required
```

Only post owner can delete.

**Response (200):**
```json
{
  "status": "success",
  "data": null
}
```

**Errors:**
- `401` - Not authenticated
- `403` - Not post owner
- `404` - Post not found

---

### Like Post

```
POST /api/posts/:id/like
Authentication: Required
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "post": {
      "id": "uuid",
      "like_count": 43,
      "is_liked_by_user": true
    }
  }
}
```

**Errors:**
- `401` - Not authenticated
- `404` - Post not found
- `429` - Rate limit exceeded (100 per hour)

---

### Unlike Post

```
DELETE /api/posts/:id/like
Authentication: Required
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "post": {
      "id": "uuid",
      "like_count": 42,
      "is_liked_by_user": false
    }
  }
}
```

**Errors:**
- `401` - Not authenticated
- `404` - Post not found
- `429` - Rate limit exceeded

---

## Messages

### List Conversations

```
GET /api/messages/conversations?page=1&limit=20
Authentication: Required
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "conversations": [
      {
        "user_id": "uuid",
        "username": "alice",
        "profile_picture_url": "https://...",
        "last_message": "Hey!",
        "last_message_at": "2026-03-28T12:00:00Z",
        "unread_count": 3
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 10
    }
  }
}
```

---

### Get Conversation History

```
GET /api/messages/:userId?page=1&limit=50
Authentication: Required
```

**Parameters:**
- `userId` - Other user's ID

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "messages": [
      {
        "id": "uuid",
        "sender_id": "uuid",
        "receiver_id": "uuid",
        "content": "Hello!",
        "is_read": true,
        "created_at": "2026-03-28T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 50,
      "total": 100
    }
  }
}
```

---

### Send Message

```
POST /api/messages
Content-Type: application/json
Authentication: Required
```

**Request Body:**
```json
{
  "receiver_id": "uuid",
  "content": "Hello!"
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "message": {
      "id": "uuid",
      "sender_id": "uuid",
      "receiver_id": "uuid",
      "content": "Hello!",
      "is_read": false,
      "created_at": "2026-03-28T12:00:00Z"
    }
  }
}
```

**Also sent via WebSocket to receiver if online.**

**Errors:**
- `400` - Invalid request
- `401` - Not authenticated
- `404` - User not found
- `429` - Rate limit exceeded

---

### Delete Messages

```
POST /api/messages/delete
Content-Type: application/json
Authentication: Required
```

**Request Body:**
```json
{
  "message_ids": ["uuid1", "uuid2"],
  "delete_for_everyone": false
}
```

**Response (200):**
```json
{
  "status": "success",
  "data": null
}
```

- If `delete_for_everyone`: true → Delete for receiver too
- If false → Delete only for sender (soft delete)

---

### Clear Conversation

```
POST /api/messages/:userId/clear
Authentication: Required
```

Removes all messages with a user.

**Response (200):**
```json
{
  "status": "success",
  "data": null
}
```

---

## Admin

### Start Impersonation

```
POST /api/admin/impersonate
Content-Type: application/json
Authentication: Required (Admin only)
```

**Request Body:**
```json
{
  "target_user_id": "uuid",
  "impersonation_password": "SecurePass123!"
}
```

**Response (200):**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "target-uuid",
      "username": "target_user",
      "role": "user"
    }
  }
}
```

**Cookie Set:**
```
Set-Cookie: auth_token=<ghost-JWT>; HttpOnly; Secure; SameSite=Strict
```

**Errors:**
- `401` - Invalid impersonation password
- `403` - Not admin
- `404` - User not found

---

### Stop Impersonation

```
POST /api/admin/impersonate/stop
Authentication: Required
```

Returns to admin account.

**Response (200):**
```json
{
  "status": "success",
  "data": null
}
```

---

## WebSocket

### Connect to Chat

```
ws://ghostline-backend-production-a17a.up.railway.app/ws/chat
Headers: Cookie: auth_token=<JWT>
```

---

### Message Format

**Send Message:**
```json
{
  "type": "message",
  "receiver_id": "uuid",
  "content": "Hello!"
}
```

**Receive Connected Event:**
```json
{
  "type": "connected",
  "data": {}
}
```

**Receive Message Event:**
```json
{
  "type": "message",
  "data": {
    "id": "uuid",
    "sender_id": "uuid",
    "receiver_id": "uuid",
    "content": "Hello!",
    "created_at": "2026-03-28T12:00:00Z"
  }
}
```

**Receive Error Event:**
```json
{
  "type": "error",
  "data": {
    "message": "Rate limit exceeded. Try again in 1 second."
  }
}
```

---

## Rate Limits

| Endpoint | Limit | Window |
|----------|-------|--------|
| Login | 5 attempts | 15 minutes |
| Upload URL | 10 requests | 1 hour |
| Like/Unlike | 100 actions | 1 hour |
| Send Message (WebSocket) | 10 messages | 1 second |
| Post Creation | 10 posts | 1 hour |

---

## Error Responses

### Validation Error (400)

```json
{
  "status": "error",
  "error": "Validation failed",
  "details": {
    "username": ["Username must be 3-50 characters"],
    "password": ["Password must contain uppercase, lowercase, number, and symbol"]
  }
}
```

### Unauthorized (401)

```json
{
  "status": "error",
  "error": "Invalid credentials"
}
```

### Forbidden (403)

```json
{
  "status": "error",
  "error": "You do not have permission to perform this action"
}
```

### Not Found (404)

```json
{
  "status": "error",
  "error": "Resource not found"
}
```

### Rate Limited (429)

```json
{
  "status": "error",
  "error": "Too many requests. Please try again later.",
  "details": {
    "retry_after_seconds": 45
  }
}
```

### Server Error (500)

```json
{
  "status": "error",
  "error": "Internal server error"
}
```

---

## Standard Response Format

**Success response:**
```json
{
  "status": "success",
  "data": { /* response data */ }
}
```

**Error response:**
```json
{
  "status": "error",
  "error": "Error message",
  "details": { /* optional */ }
}
```

---

## Testing Endpoints

### Using cURL

```bash
# Register
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"SecurePass123!","password_confirm":"SecurePass123!"}'

# Login
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"SecurePass123!"}' \
  -c cookies.txt

# Get feed (with cookies)
curl http://localhost:3000/api/posts \
  -b cookies.txt
```

### Using Client-side JavaScript

```javascript
// Login
const response = await fetch('https://ghostline-backend.app/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'test', password: 'SecurePass123!' }),
  credentials: 'include'  // Include cookies
});

// Get feed (logged in)
const feed = await fetch('https://ghostline-backend.app/api/posts', {
  credentials: 'include'  // Include cookies
});
```

---

Last Updated: March 28, 2026
