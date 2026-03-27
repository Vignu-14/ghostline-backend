# 🐛 Troubleshooting Guide

## Frontend Issues

### Page doesn't load

**Error:** Network tab shows 500 error, blank page, or logo spins forever

**Solutions:**
1. Open DevTools console (F12)
2. Check for error messages
3. Clear cache: Ctrl+Shift+Delete
4. Hard refresh: Ctrl+Shift+R

If error shows `Cannot reach API`:
```javascript
// Test API connectivity in console
fetch('https://ghostline-backend-production-a17a.up.railway.app/health')
  .then(r => r.json())
  .then(console.log)
  .catch(console.error)
```

---

### Login fails

**Error:** "Invalid credentials" even with correct password, or "Network error"

**Possible Causes:**

1. **Wrong API URL**
   - Check Vercel env vars: Settings → Environment Variables
   - Should be: `VITE_API_BASE_URL=https://ghostline-backend-production-a17a.up.railway.app`
   - No trailing slash
   - Must include `https://`
   - Redeploy after changing

2. **CORS error in console**
   - Backend not allowing frontend origin
   - Check Railway `ALLOWED_ORIGIN` env var
   - Must match exactly: `https://ghostline-frontend-five.vercel.app`
   - Redeploy Railway after change

3. **Backend not responding**
   - Check Railway health: `/health` endpoint
   - Check Railway logs for errors
   - Backend might be crashing

4. **Database connection issue**
   - Backend can't connect to database
   - Check Railway logs for: "database connection failed"
   - Verify `DATABASE_URL` in Railway env vars

### Suggested debugging steps:
```
1. Test API directly:
   curl https://ghostline-backend-production-a17a.up.railway.app/health

2. Check CORS headers:
   curl -v -H "Origin: https://ghostline-frontend.com" \
     https://ghostline-backend-production-a17a.up.railway.app/api/posts

3. Test login in API:
   curl -X POST https://ghostline-backend/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"test","password":"test123"}'
```

---

### Chat not working

**Error:** Messages don't send, "WebSocket error", or "Rate limit exceeded"

**Possible Causes:**

1. **WebSocket connection failed**
   - Open DevTools → Network → WS
   - Should see `ws://...` connection
   - If 403 error: JWT token expired
   - If connection closes: Server redeployed

2. **Rate limiting**
   - Error: "Rate limit exceeded. Try again in X seconds"
   - Max 10 messages per second per user
   - Wait before retrying
   - Check if sending duplicate messages

3. **Receiver offline**
   - Messages only deliver to online users
   - If receiver offline, message saved for next login
   - Check user is in chat page, not just app

4. **Network latency**
   - Message takes 1-5 seconds to arrive
   - Check network tab (DevTools → Network)
   - High latency? Check connection speed

---

### Profile says "Not found"

**Error:** Click on user → "404 User not found"

**Possible Causes:**

1. **User doesn't exist**
   - Typo in username
   - User deleted account
   - Check user in search results first

2. **Database issue**
   - Check Rally logs for: "database connection failed"
   - Data might be corrupted
   - Try logging out and back in

---

## Backend Issues

### Server won't start

**Error:** `go run cmd/server/main.go` fails with error

**Possible Causes:**

1. **Go not installed**
   ```bash
   go version
   # If not found: download from https://go.dev/dl/
   ```

2. **Database connection failed**
   ```
   Error: unable to connect to database
   ```
   - Check `DATABASE_URL` env var
   - Verify PostgreSQL is running locally
   - Verify database name exists

3. **Port already in use**
   ```
   Error: listen tcp :8080: bind: address already in use
   ```
   - Kill process using port 8080:
     - Linux/Mac: `lsof -i :8080 | grep LISTEN | awk '{print $2}' | xargs kill -9`
     - Windows: `netstat -ano | findstr :8080`
   - Or change PORT in .env to 3000

4. **Missing dependencies**
   ```
   Error: no required module provides ...
   ```
   ```bash
   go mod download
   go mod tidy
   ```

5. **Compilation error**
   - Check for syntax errors: `go build ./...`
   - Check imports are correct
   - Check type mismatches

---

### Migrations fail

**Error:** Running `./scripts/migrate.sh` fails

**Possible Causes:**

1. **Database not running**
   ```bash
   psql $DATABASE_URL
   # If fails: start PostgreSQL
   ```

2. **Wrong database URL**
   ```bash
   echo $DATABASE_URL
   # Should be: postgresql://user:pass@host:5432/dbname
   ```

3. **Permission denied**
   - Make script executable:
     ```bash
     chmod +x scripts/migrate.sh
     ```

4. **Migration file syntax error**
   - Check SQL syntax in files
   - Run migration manually:
     ```bash
     psql $DATABASE_URL < internal/database/migrations/001_create_tables.sql
     ```

---

### API returns 500 error

**Error:** API endpoint returns: `{"status":"error","error":"Internal server error"}`

**Possible Causes:**

1. **Database query failed**
   - Check Railway logs for SQL errors
   - Verify table structure matches code
   - Test query manually:
     ```sql
     psql $DATABASE_URL
     SELECT * FROM users LIMIT 1;
     ```

2. **Null pointer dereference**
   - User submitted invalid data
   - Check validation in handler
   - See logs for panic message

3. **Supabase connection failed**
   - Check `SUPABASE_URL` env var
   - Check `SUPABASE_KEY` is correct
   - Verify bucket exists

4. **JWT validation failed**
   - Token expired (15 min timeout)
   - Need to login again
   - Check `JWT_SECRET` matches

---

### Memory/CPU high

**Error:** Railway shows 90% CPU or 70%+ memory usage

**Possible Causes:**

1. **Connection pool leaking**
   - Database connections not being closed
   - Check for `defer db.Close()` in code
   - Check max connections limit (25)

2. **Goroutine leak**
   - Background tasks not stopping
   - Check WebSocket connections
   - Monitor with `pprof`:
     ```bash
     go get github.com/google/pprof
     pprof http://localhost:8080/debug/pprof
     ```

3. **Large query result**
   - Fetching too much data at once
   - Add pagination: `LIMIT 20 OFFSET 0`
   - Check query performs index scan

4. **Caching issue**
   - Cache growing too large
   - Implement TTL for cached items
   - Clear cache periodically

---

## Database Issues

### Can't connect to Supabase

**Error:** `Database URL is invalid` or `connection refused`

**Solutions:**

1. **Verify credentials**
   - Go to Supabase Dashboard
   - Copy connection string again
   - Check for special characters in password (URL encode if needed)

2. **Check network**
   ```bash
   ping db.PROJECTID.supabase.co
   telnet db.PROJECTID.supabase.co 5432
   ```

3. **IP whitelist**
   - Supabase usually allows all IPs
   - If fails, check Settings → Database → Network
   - Add your IP to whitelist

4. **Connection string format**
   ```
   ✅ postgresql://user:password@host:5432/postgres
   ❌ postgres://user:password@host/5432/postgres
   ❌ postgresql://user:password@host (missing port/db)
   ```

---

### Data corrupted/missing

**Error:** Query returns unexpected data, NULL values

**Solutions:**

1. **Check migrations applied**
   ```bash
   psql $DATABASE_URL
   SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 5;
   ```

2. **Verify data exists**
   ```sql
   SELECT COUNT(*) FROM users;
   SELECT COUNT(*) FROM posts;
   ```

3. **Check constraints**
   ```sql
   SELECT constraint_name FROM information_schema.table_constraints
   WHERE table_name = 'posts';
   ```

4. **Restore from backup**
   - Supabase Dashboard → Backups
   - Restore to specific date
   - **Warning:** Overwrites all current data

---

### Slow queries

**Error:** API endpoint responds slowly (>5 seconds)

**Solutions:**

1. **Check query performance**
   ```sql
   EXPLAIN ANALYZE SELECT * FROM posts WHERE user_id = '...';
   ```
   - Look for "Seq Scan" (slow)
   - Should be "Index Scan" (fast)

2. **Add indexes**
   ```sql
   CREATE INDEX idx_posts_user_id ON posts(user_id);
   CREATE INDEX idx_messages_receiver_id ON messages(receiver_id);
   ```

3. **Review pagination**
   - Ensure `LIMIT` and `OFFSET` used
   - Don't fetch all records at once

4. **Check connection pool**
   - Max connections reached?
   - Close connections properly
   - Monitor active connections:
     ```sql
     SELECT count(*) FROM pg_stat_activity;
     ```

---

## WebSocket Issues

### Chat won't connect

**Error:** WebSocket shows 403 or 500 error

**Possible Causes:**

1. **JWT token expired**
   - Need to login again
   - Tokens expire after 15 minutes
   - Check cookie exists: DevTools → Application → Cookies

2. **Backend not accepting WebSocket**
   - Check Railway `/ws/chat` endpoint exists
   - Verify WebSocket path in code
   - Check firewall allows WebSocket

3. **CORS misconfiguration**
   - WebSocket also subject to CORS
   - Check `ALLOWED_ORIGIN` env var

4. **Connection timeout**
   - Network latency
   - Server might not respond quickly
   - Check network tab (DevTools → Network → WS)

---

### Messages pile up/don't send

**Error:** Sending message hangs, then 20 messages arrive at once

**Possible Causes:**

1. **Network buffering**
   - Messages queued locally
   - Check DevTools Network tab
   - Wait for all to send

2. **Server overloaded**
   - Too many concurrent connections
   - Check Railway CPU/memory
   - Restart server if needed

3. **Rate limiting**
   - Limit: 10 messages per second per user
   - If exceeded, messages rejected
   - Implement exponential backoff in frontend

---

## Deployment Issues

### Railway build fails

**Error:** `Build failed: go: command not found` or similar

**Solutions:**

1. **Check Dockerfile**
   - File must exist and be named `Dockerfile`
   - Not `dockerfile` (case-sensitive on Linux)
   - Path must be repo root

2. **Go version**
   - Check `go.mod` requires version
   - Dockerfile should match: `golang:1.25-alpine`
   - Update Dockerfile if version mismatch

3. **Dependencies missing**
   - Run locally: `go mod download`
   - Commit `go.mod` and `go.sum`
   - If still fails: clean go cache

4. **Port configuration**
   - Go app must listen on `$PORT`
   - Railway sets this env var
   - Check code: `fiber.Listen(":" + os.Getenv("PORT"))`

### Suggested fix:
```bash
# Force rebuild
Railway Dashboard → My Services → ghostline-backend
→ Settings → Trigger Redeploy
```

---

### Vercel build fails

**Error:** `Build failed: npm ERR!` or similar

**Solutions:**

1. **Dependencies issue**
   ```bash
   rm -rf node_modules package-lock.json
   npm install  # Install locally first
   git add package-lock.json
   git commit -m "fix: update dependencies"
   git push
   ```

2. **Build script error**
   - Test locally: `npm run build`
   - Check `vite.config.ts` configuration
   - Check tsconfig.json for errors

3. **Environment variable missing**
   - Vercel Settings → Environment Variables
   - Rebuild after adding env var
   - Must include: `VITE_API_BASE_URL`

4. **Node version**
   - Vercel → Settings → Node.js Version
   - Set to 18 or 20
   - Restart build

---

### App deploys but shows blank page

**Error:** Frontend loads but shows no content, or 404

**Solutions:**

1. **Check build output**
   - DevTools Console → Any errors?
   - Check Network tab → index.html 200?
   - Check dist folder has files

2. **Environment variable wrong**
   ```javascript
   // Open console and check:
   console.log(import.meta.env.VITE_API_BASE_URL)
   // Should show: https://ghostline-backend.app
   ```

3. **API unreachable**
   - Frontend loads, but can't reach backend
   - Check CORS settings
   - Check backend running

4. **Router issue**
   - SPA not loading correctly
   - Check `nginx.conf` routing
   - Vercel auto-handles React Router

---

## Getting Help

### Check Logs

**Railway Logs:**
1. Dashboard → Services → ghostline-backend
2. Logs tab → Filter by "error"
3. Copy error message

**Vercel Logs:**
1. Dashboard → Deployments
2. Click failed deployment → Logs
3. Read error message

**Browser Console:**
1. F12 → Console tab
2. Look for red errors
3. Check Network tab for failed requests

### Useful Commands

```bash
# Test backend
curl https://ghostline-backend-production-a17a.up.railway.app/health

# Test database
psql $DATABASE_URL -c "SELECT 1;"

# Check environment variables
echo $DATABASE_URL
echo $JWT_SECRET

# View service status
git status
go mod tidy
npm audit
```

---

Last Updated: March 28, 2026
