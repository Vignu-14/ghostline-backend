# 🆘 Troubleshooting Guide

Common issues and solutions for Ghostline.

---

## Backend Issues

### 🔴 Server Won't Start

**Error:** `listen: bind: address already in use`

**Solution:**
```bash
# Find process using port 3000
lsof -i :3000

# Kill the process
kill -9 <PID>

# Or use different port
PORT=3001 go run cmd/server/main.go
```

---

**Error:** `database connection failed`

**Solution:**
```bash
# Check DATABASE_URL in .env
echo $DATABASE_URL

# Test PostgreSQL connection
psql "$DATABASE_URL"

# If local Postgres:
psql -U postgres -d ghostline
```

---

**Error:** `JWT_SECRET is required in production`

**Solution:**
```bash
# Set JWT_SECRET in .env
JWT_SECRET=your-64-character-random-secret-string

# Or in Railway Variables
# Go to Dashboard → Variables → Add JWT_SECRET
```

---

### 🟡 CORS Errors

**Error:** `Access to fetch... has been blocked by CORS policy`

**Frontend Console Error:**
```
No 'Access-Control-Allow-Origin' header is present
```

**Solution:**

1. **Check ALLOWED_ORIGIN matches frontend URL**
   - In Railway Variables
   - Verify it matches exactly (including protocol)

2. **Format should be:**
   ```
   https://ghostline-frontend-five.vercel.app
   ```
   NOT: `http://`, NOT: with `/api` suffix

3. **After updating, redeploy:**
   - Railway: Click "Deploy" button
   - Wait 30 seconds

4. **Clear frontend cache:**
   - Hard refresh: `Ctrl+Shift+R` (Windows) or `Cmd+Shift+R` (Mac)

---

### 🟡 Database Issues

**Error:** `table "users" does not exist`

**Solution:**
```bash
# Run migrations
go run cmd/server/main.go  # Migrations run automatically on startup

# Or manually with Supabase CLI
supabase db push
```

---

**Error:** `duplicate key value violates unique constraint "users_username_key"`

**Solution:**
- Username already exists
- Try registering with different username
- Or delete user from database:
  ```sql
  DELETE FROM users WHERE username = 'existing_username';
  ```

---

### 🟡 Authentication Issues

**Error:** `invalid or expired session`

**Solution:**
- JWT token expired (15 minute default)
- Login again to get new token
- Check `COOKIE_SECURE=true` is set in production

---

**Error:** `authentication required`

**Solution:**
- User not authenticated
- Login first before accessing protected endpoints
- Ensure cookies are enabled in browser

---

## Frontend Issues

### 🔴 Build Fails

**Error:** `vite: command not found`

**Solution:**
```bash
# Install dependencies
npm install

# Then build
npm run build
```

---

**Error:** `Cannot find module '@vite/plugin-...'`

**Solution:**
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Or
npm ci  # Clean install
```

---

### 🟡 API Connection Issues

**Error:** `Unable to read server response`

**In Console:**
```
POST https://ghostline-frontend-five.vercel.app/api/auth/login 404
```

**Solution:**

1. **Check VITE_API_BASE_URL is set in Vercel:**
   - Go to Vercel → Settings → Environment Variables
   - Verify `VITE_API_BASE_URL` exists
   - Value should be: `https://ghostline-backend-production-xxxx.up.railway.app`

2. **Redeploy after setting env var:**
   - Go to Deployments
   - Click Latest → Redeploy

3. **Hard refresh browser:**
   - `Ctrl+Shift+R` (Windows)
   - `Cmd+Shift+R` (Mac)

---

**Error:** `Connection refused`

**Solution:**
- Backend not running
- Check Railway deployment status
- Verify backend URL is accessible:
  ```bash
  curl https://ghostline-backend-production-xxxx.up.railway.app/health
  ```

---

### 🟡 WebSocket Connection Issues

**Error:** `Failed to establish WebSocket connection`

**Solution:**
1. Check backend is running
2. Check JWT token is valid (login first)
3. Check browser console for specific error
4. Try hard refresh

---

**Error:** `WebSocket is closed before the connection is established`

**Solution:**
- Backend crashed or restarted
- Network connectivity issue
- Implement auto-reconnect with exponential backoff

```typescript
let retries = 0;
const maxRetries = 5;

function connectWithRetry() {
  try {
    connect();
  } catch (e) {
    if (retries < maxRetries) {
      retries++;
      setTimeout(connectWithRetry, Math.pow(2, retries) * 1000);
    }
  }
}
```

---

## Deployment Issues

### Railway Backend

**Error:** `Railpack could not determine how to build the app`

**Solution:**
- Dockerfile must exist and not be empty
- Root Directory should be empty (backend is at root)
- Go version must be 1.25 in go.mod

---

**Error:** `failed to get private network endpoint`

**Solution:**
- This is normal, just a warning
- Backend is deployed, ignore this
- Check your public domain instead

---

**Error:** `Health check failed`

**Solution:**
- Backend might be crashing
- Check Railway Logs for errors
- Verify all environment variables are set correctly

---

### Vercel Frontend

**Error:** `Command "npm run build" exited with 1`

**Solution:**
```bash
# Run build locally to see error
npm run build

# Fix errors
# Common: TypeScript errors, missing env vars
```

---

**Error:** `Error: Module not found: 'vite'`

**Solution:**
- Dependencies not installed
- Go to Settings → Advanced → Clear Build Cache
- Redeploy

---

## Database Issues

### Supabase

**Error:** `SSL CERTIFICATE_VERIFY_FAILED`

**Solution:**
- Already set in code with `?sslmode=require`
- Make sure ca-certificates are installed in Docker
- Dockerfile includes: `apk --no-cache add ca-certificates`

---

**Error:** `Too many connections`

**Solution:**
- Connection pooling not working
- Check `.pooler.supabase.com` is in DATABASE_URL
- Reduce `DB_MAX_CONNECTIONS` in Railway Variables

---

## File Upload Issues

**Error:** `File type not allowed`

**Solution:**
- Only JPEG/PNG images allowed
- Check file extension and MIME type
- File must be image/jpeg or image/png

---

**Error:** `File too large`

**Solution:**
- Max size is 5MB
- Compress image before upload
- Check file_size in upload request

---

## Performance Issues

### Slow API Responses

**Solution:**
1. Check Railway CPU/Memory usage
2. Check database query logs
3. Optimize queries (add indexes)
4. Enable caching

---

### High Memory Usage

**Solution:**
- Check for memory leaks
- Monitor WebSocket connections
- Check database connection pool size
- Reduce `DB_MAX_CONNECTIONS` if too high

---

## Security Issues

### Exposed Secrets

**Error:** Secrets visible in `.env` file after push

**Solution:**
```bash
# Remove from git history
git filter-branch --tree-filter 'rm -f .env' --prune-empty HEAD

# Remove remote tracking
git push origin main --force-with-lease

# Regenerate all secrets
# Go to Supabase → Settings → Reset database password
# Regenerate Supabase service key
```

---

### CORS Not Working

**Error:** `No 'Access-Control-Allow-Origin' header`

**Solution:**
- ALLOWED_ORIGIN in Railway must match frontend URL exactly
- No wildcards allowed (security)
- Update and redeploy after changing

---

## Getting Help

1. **Check logs:**
   - Railway: Go to Project → Logs
   - Vercel: Go to Project → Deployments → Logs

2. **Check status pages:**
   - [Railway Status](https://status.railway.app)
   - [Vercel Status](https://www.vercel-status.com)
   - [Supabase Status](https://status.supabase.com)

3. **Review documentation:**
   - [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
   - [DEPLOYMENT.md](./DEPLOYMENT.md)
   - [ARCHITECTURE.md](./ARCHITECTURE.md)

4. **Test endpoints:**
   ```bash
   # Health check
   curl https://backend-url/health
   
   # Login
   curl -X POST https://backend-url/api/auth/login \
        -H "Content-Type: application/json" \
        -d '{"username":"test","password":"test"}'
   ```

---

See also: [DEVELOPMENT.md](./DEVELOPMENT.md), [DEPLOYMENT.md](./DEPLOYMENT.md)
