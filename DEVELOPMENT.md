# 🛠️ Development Guide

## Getting Started

### Prerequisites

- Go 1.25+ ([download](https://go.dev/dl/))
- Node.js 18+ ([download](https://nodejs.org/))
- PostgreSQL 14+ (local or Docker)
- Git

### Quick Start

1. **Clone repositories**
   ```bash
   git clone https://github.com/your-org/ghostline-backend
   git clone https://github.com/your-org/ghostline-frontend
   ```

2. **Backend setup**
   ```bash
   cd ghostline-backend
   cp .env.example .env
   # Edit .env with local database URL
   go mod download
   make dev  # or: go run cmd/server/main.go
   ```

3. **Frontend setup**
   ```bash
   cd ghostline-frontend
   npm install
   npm run dev
   ```

4. **Database setup**
   ```bash
   cd ghostline-backend
   ./scripts/migrate.sh     # Create tables
   ./scripts/seed.sh        # Add test data
   ```

5. **Open in browser**
   - Frontend: `http://localhost:5173`
   - Backend API: `http://localhost:8080/api/posts`

---

## Project Structure

### Backend

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── config/
│   │   ├── config.go         # App configuration
│   │   ├── constants.go      # Constants (limits, timeouts)
│   │   ├── database.go       # DB connection setup
│   │   └── storage.go        # Storage configuration
│   ├── database/
│   │   ├── postgres.go       # PostgreSQL driver
│   │   ├── migrations/       # SQL migration files
│   │   └── seed/             # Test data
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── post_handler.go
│   │   ├── chat_handler.go
│   │   ├── user_handler.go
│   │   ├── like_handler.go
│   │   ├── admin_handler.go
│   │   ├── health_handler.go
│   │   └── websocket_handler.go
│   ├── middleware/
│   │   ├── jwt_middleware.go
│   │   ├── cors_middleware.go
│   │   ├── admin_middleware.go
│   │   └── error_handler.go
│   ├── models/
│   │   ├── user.go
│   │   ├── post.go
│   │   ├── message.go
│   │   └── ...
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── post_repository.go
│   │   ├── message_repository.go
│   │   └── ...
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── post_service.go
│   │   ├── chat_service.go
│   │   └── upload_service.go
│   ├── utils/
│   │   ├── jwt.go
│   │   ├── validators.go
│   │   └── sanitizers.go
│   └── websocket/
│       ├── hub.go           # WebSocket connection manager
│       ├── client.go        # Single connection
│       └── broadcaster.go   # Send to clients
├── pkg/
│   ├── cache/              # Caching utilities
│   └── logger/             # Logging setup
├── tests/
│   ├── unit/              # Unit tests
│   ├── integration/       # Integration tests
│   └── fixtures/          # Test data
├── scripts/
│   ├── migrate.sh         # Run migrations
│   ├── seed.sh            # Seed test data
│   └── generate_jwt_secret.sh
├── docs/
│   ├── api.md
│   ├── architecture.md
│   ├── websocket.md
│   └── deployment.md
├── Makefile               # Build commands
├── go.mod                 # Dependencies
├── Dockerfile             # Production container
└── .env.example           # Example environment
```

### Frontend

```
frontend/
├── src/
│   ├── main.tsx            # Entry point
│   ├── App.tsx             # Root component
│   ├── vite-env.d.ts       # Vite types
│   ├── pages/              # Route pages
│   │   ├── HomePage.tsx
│   │   ├── LoginPage.tsx
│   │   ├── RegisterPage.tsx
│   │   ├── ChatPage.tsx
│   │   ├── ProfilePage.tsx
│   │   ├── AdminPage.tsx
│   │   └── NotFoundPage.tsx
│   ├── components/         # Reusable components
│   │   ├── Navbar.tsx
│   │   ├── PostFeed.tsx
│   │   ├── PostCard.tsx
│   │   ├── ChatWindow.tsx
│   │   ├── ProtectedRoute.tsx
│   │   └── ...
│   ├── context/            # State management
│   │   ├── AuthContext.tsx
│   │   ├── ChatContext.tsx
│   │   └── NotificationContext.tsx
│   ├── services/           # API calls
│   │   ├── api.ts          # HTTP client
│   │   ├── authService.ts
│   │   ├── postService.ts
│   │   ├── chatService.ts
│   │   └── webSocketService.ts
│   ├── hooks/              # Custom hooks
│   │   ├── useAuth.ts
│   │   ├── useChat.ts
│   │   └── useNotification.ts
│   ├── types/              # TypeScript interfaces
│   │   ├── User.ts
│   │   ├── Post.ts
│   │   ├── Message.ts
│   │   └── ...
│   ├── utils/              # Utility functions
│   │   ├── formatters.ts
│   │   ├── validators.ts
│   │   └── constants.ts
│   └── styles/             # CSS modules
├── public/                 # Static assets
├── vite.config.ts          # Vite configuration
├── tsconfig.json           # TypeScript config
├── tailwind.config.js      # Tailwind CSS config
├── package.json            # Dependencies
├── README.md
└── .env.example
```

---

## Development Workflow

### 1. Create Feature Branch

```bash
git checkout -b feature/add-dark-mode
```

### 2. Make Changes

Example: Adding a new API endpoint

**Backend:**
```go
// internal/handlers/theme_handler.go
func (h *ThemeHandler) SetTheme(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)
    theme := c.Query("theme") // "light" or "dark"
    
    // Validate
    if theme != "light" && theme != "dark" {
        return c.Status(400).JSON(fiber.Map{
            "status": "error",
            "error": "Invalid theme",
        })
    }
    
    // Update in database
    err := h.service.SetUserTheme(userID, theme)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status": "error",
            "error": "Failed to update theme",
        })
    }
    
    return c.JSON(fiber.Map{
        "status": "success",
        "data": fiber.Map{"theme": theme},
    })
}
```

**Register route in main.go:**
```go
api.Put("/users/theme", middleware.JWTMiddleware, themeHandler.SetTheme)
```

**Frontend:**
```typescript
// src/services/themeService.ts
export async function setTheme(theme: 'light' | 'dark') {
  const response = await fetch(`${API_BASE_URL}/api/users/theme?theme=${theme}`, {
    method: 'PUT',
    credentials: 'include',
  });
  return response.json();
}

// src/components/ThemeToggle.tsx
import { setTheme } from '../services/themeService';

export function ThemeToggle() {
  const handleToggle = async () => {
    await setTheme(isDark ? 'light' : 'dark');
    setIsDark(!isDark);
  };
  
  return <button onClick={handleToggle}>🌙/☀️</button>;
}
```

### 3. Test Changes

**Backend:**
```bash
# Run unit tests
go test ./...

# Run integration tests
go test -tags=integration ./...

# Run specific test
go test -run TestUserLogin ./internal/services
```

**Frontend:**
```bash
# Run tests
npm test

# Run with coverage
npm test -- --coverage
```

### 4. Commit and Push

```bash
git add .
git commit -m "feat: add dark mode theme toggle"
git push origin feature/add-dark-mode
```

### 5. Create Pull Request

On GitHub:
1. Compare & pull request
2. Add description of changes
3. Request reviewer
4. Address feedback
5. Merge when approved

---

## Testing

### Backend Testing

**Unit Tests** (test single functions):
```go
// internal/services/post_service_test.go
func TestCreatePost(t *testing.T) {
    mockRepo := NewMockPostRepository()
    service := NewPostService(mockRepo)
    
    post, err := service.CreatePost("user123", "Hello world!")
    
    assert.NoError(t, err)
    assert.Equal(t, "Hello world!", post.Caption)
}
```

Run tests:
```bash
go test ./internal/services -v
```

**Integration Tests** (test with real database):
```go
// tests/integration/post_test.go
func TestCreatePostIntegration(t *testing.T) {
    db, _ := setupTestDB()
    defer db.Close()
    
    repo := postgre.NewPostRepository(db)
    service := services.NewPostService(repo)
    
    post, err := service.CreatePost("user123", "Test post")
    
    assert.NoError(t, err)
    assert.NotNil(t, post.ID)
}
```

Run integration tests:
```bash
go test -tags=integration ./tests/integration -v
```

### Frontend Testing

**Unit Tests** (test components):
```typescript
// src/components/__tests__/PostCard.test.tsx
import { render, screen } from '@testing-library/react';
import PostCard from '../PostCard';

test('renders post caption', () => {
  const post = {
    id: '1',
    caption: 'Hello world',
    // ...
  };
  
  render(<PostCard post={post} />);
  expect(screen.getByText('Hello world')).toBeInTheDocument();
});
```

Run tests:
```bash
npm test PostCard.test.tsx
```

---

## Build & Deployment

### Build Backend

```bash
make build
# Output: bin/ghostline-api
```

### Build Frontend

```bash
npm run build
# Output: dist/
```

### Run Locally

**Backend:**
```bash
go run cmd/server/main.go
# Server running on http://localhost:8080
```

**Frontend:**
```bash
npm run dev
# Server running on http://localhost:5173
```

---

## Environment Variables

Create `.env` file in each repository:

**Backend (.env):**
```env
DATABASE_URL=postgres://user:pass@localhost:5432/ghostline_dev
JWT_SECRET=dev-secret-not-for-production
SUPABASE_URL=https://projectid.supabase.co
SUPABASE_KEY=key-here
PORT=8080
ENVIRONMENT=development
ALLOWED_ORIGIN=http://localhost:5173
```

**Frontend (.env):**
```env
VITE_API_BASE_URL=http://localhost:8080
```

---

## Common Commands

### Backend

```bash
# Install dependencies
go mod download

# Run server
go run cmd/server/main.go

# Run tests
go test ./...

# Build binary
go build -o bin/ghostline-api cmd/server/main.go

# Format code
go fmt ./...

# Run linter
golangci-lint run

# View test coverage
go test -cover ./...
```

### Frontend

```bash
# Install dependencies
npm install

# Run dev server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Run tests
npm test

# Format code
npm run format

# Lint code
npm run lint

# Type check
npm run type-check
```

---

## Code Style

### Go

Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments):

```go
// ✅ Good
func (h *Handler) GetUser(c *fiber.Ctx) error {
    userID := c.Params("id")
    user, err := h.service.GetUser(userID)
    if err != nil {
        return err
    }
    return c.JSON(user)
}

// ❌ Bad
func (h *Handler) GetUser(c *fiber.Ctx) error {
    id:=c.Params("id")
    u,e:=h.service.GetUser(id)
    if(e!=nil){return e}
    return c.JSON(u)
}
```

### TypeScript

Follow [Google TypeScript Style Guide](https://google.github.io/styleguide/tsguide.html):

```typescript
// ✅ Good
interface User {
  id: string;
  username: string;
  email: string;
}

async function fetchUser(id: string): Promise<User> {
  const response = await fetch(`/api/users/${id}`);
  return response.json();
}

// ❌ Bad
interface User {
  id:string,
  username:string,
  email:string
}

function fetchUser(id) {
  return fetch(`/api/users/${id}`).then(r=>r.json());
}
```

---

## Debugging

### Backend Debugging

**Using Delve debugger:**
```bash
go install github.com/go-delve/delve/cmd/dlv@latest

dlv debug cmd/server/main.go
(dlv) break main.main
(dlv) continue
(dlv) p variable_name  # Print variable
```

**Using loggers:**
```go
import "github.com/sirupsen/logrus"

log := logrus.New()
log.WithFields(logrus.Fields{
    "user_id": userID,
    "action": "create_post",
}).Info("Creating new post")
```

### Frontend Debugging

**Using browser DevTools:**
1. Open browser (F12)
2. Go to Sources tab
3. Set breakpoints
4. Step through code

**Using React DevTools:**
1. Install React DevTools extension
2. Open DevTools
3. Go to Components tab
4. Inspect component state/props

**Using console logs:**
```typescript
console.log('User:', user);
console.error('Error:', error);
console.table(posts);  // Display array as table
```

---

## Performance Optimization

### Database

- Use indexes on frequently queried columns
- Paginate results (don't fetch all 10,000 records)
- Cache READ-heavy queries
- Use connection pooling

### Backend

- Compress responses (gzip)
- Cache static files
- Rate limit expensive operations
- Use goroutines for concurrent tasks

### Frontend

- Code splitting (lazy load routes)
- Image optimization
- Minimize bundle size
- Cache with service workers

---

## Troubleshooting

### Backend won't start

```
error: database connection failed
```

**Solution:**
1. Check DATABASE_URL env var
2. Verify PostgreSQL is running
3. Check database name exists

### Frontend build fails

```
error: Cannot find module 'react'
```

**Solution:**
1. Run `npm install`
2. Delete `node_modules` and `package-lock.json`
3. Run `npm install` again

### Tests failing

```
FAIL: TestUserLogin
```

**Solution:**
1. Check test database is running
2. Run migrations before tests
3. Check mock data setup

---

## Contributing

1. Fork repository
2. Create feature branch: `feature/your-feature`
3. Make changes and test
4. Submit pull request
5. Address review feedback
6. Merge when approved

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

---

Last Updated: March 28, 2026
