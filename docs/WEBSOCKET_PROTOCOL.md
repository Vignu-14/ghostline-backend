# 🔌 WebSocket Protocol

Real-time messaging protocol for Ghostline using WebSocket.

---

## 📡 Connection

### URL
```
ws://localhost:3000/ws/chat    (Development)
wss://your-domain.com/ws/chat  (Production)
```

### Authentication
WebSocket connection requires a valid JWT token in the `auth_token` cookie.

```javascript
const ws = new WebSocket('ws://localhost:3000/ws/chat', {
  withCredentials: true,  // Include cookies (auth_token)
  credentials: 'include'
});
```

---

## 📤 Message Types

### 1. **Connection Established**
When you first connect, the server sends:

```json
{
  "type": "connected",
  "timestamp": "2026-03-28T10:00:00Z"
}
```

### 2. **Send Message**
Client → Server:

```json
{
  "type": "message",
  "receiver_id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "Hello! How are you?"
}
```

Server → Client (confirmation):

```json
{
  "type": "message",
  "id": "msg-uuid",
  "sender_id": "your-user-id",
  "receiver_id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "Hello! How are you?",
  "is_read": false,
  "created_at": "2026-03-28T10:00:00Z"
}
```

### 3. **Receive Message**
When someone sends you a message:

```json
{
  "type": "message",
  "id": "msg-uuid",
  "sender_id": "550e8400-e29b-41d4-a716-446655440000",
  "receiver_id": "your-user-id",
  "content": "Hi! I'm great!",
  "is_read": false,
  "created_at": "2026-03-28T10:00:05Z"
}
```

### 4. **Error Event**
When something goes wrong:

```json
{
  "type": "error",
  "error": "message rate limit exceeded. try again in 5 seconds."
}
```

---

## 🔄 Full Message Flow

```
1. Client connects to WebSocket endpoint
         ↓
2. Server validates JWT token from cookie
         ↓
3. Client receives:
   {"type": "connected", ...}
         ↓
4. Client sends message:
   {"type": "message", "receiver_id": "...", "content": "..."}
         ↓
5. Server validates & saves to database
         ↓
6. Server broadcasts to all connections of receiver
         ↓
7. Receiving client gets:
   {"type": "message", "sender_id": "...", ...}
```

---

## 🛠️ Client Implementation

### JavaScript/TypeScript Example

```typescript
class ChatSocket {
  private socket: WebSocket | null = null;

  connect(onMessage: (event: any) => void, onClose?: () => void) {
    this.socket = new WebSocket('ws://localhost:3000/ws/chat');
    
    this.socket.onopen = () => {
      console.log('Connected to chat');
    };

    this.socket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      onMessage(message);
    };

    this.socket.onclose = () => {
      console.log('Disconnected');
      if (onClose) onClose();
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }

  send(receiverId: string, content: string) {
    if (this.socket?.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({
        type: 'message',
        receiver_id: receiverId,
        content: content
      }));
    }
  }

  close() {
    this.socket?.close();
  }
}

// Usage
const chat = new ChatSocket();
chat.connect((event) => {
  if (event.type === 'message') {
    console.log(`Message from ${event.sender_id}: ${event.content}`);
  }
});

chat.send('user-id', 'Hello!');
```

---

## ⚡ Rate Limiting

**WebSocket Rate Limit:** 10 messages per second per user

If exceeded:

```json
{
  "type": "error",
  "error": "message rate limit exceeded. try again in 0.5 seconds."
}
```

---

## 🔐 Security Features

1. **JWT Token Validation**
   - Token validated on connection
   - Invalid/expired token → connection rejected

2. **User Isolation**
   - Users can only receive messages sent to their ID
   - Messages automatically sent to correct recipient

3. **Message Sanitization**
   - All content sanitized (XSS prevention)
   - HTML tags removed

4. **CORS/Origin Check**
   - Origin validated against ALLOWED_ORIGIN
   - Cross-origin WebSocket connections rejected

---

## 📊 Connection Management

### Automatic Disconnect
- **Idle timeout:** 5 minutes
- **Inactive session:** Automatically cleaned up
- **Duplicate connections:** Previous connection closed

### Reconnection Strategy
```typescript
const maxRetries = 5;
let retries = 0;

function connectWithRetry() {
  try {
    chat.connect(onMessage, () => {
      if (retries < maxRetries) {
        retries++;
        setTimeout(connectWithRetry, 1000 * retries); // Exponential backoff
      }
    });
  } catch (error) {
    console.error('Connection failed:', error);
  }
}
```

---

## 🧪 Testing

### Using wscat

```bash
# Install
npm install -g wscat

# Get JWT token from login API
TOKEN=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"username":"test","password":"password"}' | jq -r '.data.user.id')

# Connect (extract auth_token from cookies.txt)
wscat -c ws://localhost:3000/ws/chat \
      --header "Cookie: auth_token=$TOKEN"

# In wscat console:
# Send message
{"type":"message","receiver_id":"550e8400-e29b-41d4-a716-446655440000","content":"hello"}

# Receive response
{"type":"message","sender_id":"...","receiver_id":"...","content":"hello","created_at":"..."}
```

---

## 📈 Monitoring

### Server Metrics
- Active WebSocket connections
- Messages per second
- Error rate
- Connection duration

### Client Monitoring
```javascript
// Log connection stats
let messageCount = 0;
const startTime = Date.now();

socket.onmessage = (event) => {
  messageCount++;
  const elapsed = Date.now() - startTime;
  console.log(`${messageCount} messages in ${elapsed}ms`);
};
```

---

## 🔧 Troubleshooting

### Connection Refused
- Backend not running
- Wrong WebSocket URL
- Port mismatch (3000 vs 8080, etc.)

### Authentication Error
- JWT token not in cookie
- Token expired
- Token invalid/tampered

### Message Not Received
- Receiver not connected
- Receiver_ID incorrect (must be valid UUID)
- Rate limit hit

### Intermittent Disconnections
- Network unstable
- Server timeout (5-minute idle)
- Too many rapid reconnections

---

## Reconnection Flow

```
Browser disconnects
         ↓
Client detects close event
         ↓
Client attempts reconnect after 1s
         ↓
If successful → resume chat
         ↓
If failed → retry with exponential backoff
```

---

See also: [API_DOCUMENTATION.md](./API_DOCUMENTATION.md), [ARCHITECTURE.md](./ARCHITECTURE.md)
