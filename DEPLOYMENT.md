# 🚀 Deployment Guide

## Overview

This guide provides step-by-step instructions for deploying Ghostline to production using Railway (backend), Vercel (frontend), and Supabase (database).

---00000000000000000000

## Prerequisites

Before deploying, ensure you have:
- [ ] GitHub account with repositories created
- [ ] Railway account (free tier available)
- [ ] Vercel account (free tier available)
- [ ] Supabase account with database created
- [ ] Admin password for testing impersonation
- [ ] Domain name (optional, but recommended)

---

## Part 1: Database Setup (Supabase)

### 1.1 Create Supabase Project

1.
1. Go to [supabase.com](https://supabase.com)
2. Click "New Project"
3. Select organization and enter:
   - Project name: `ghostline`
   - Password: Save this securely
   - Region: Choose closest to users
4. Wait for database initialization (5-10 minutes)

### 1.2 Run Migrations

1. Download connection string:
   ```
   Dashboard → Settings → Database → Connection strings → URI
   ```

2. Run migrations locally:
   ```bash
   cd backend
   export DATABASE_URL="postgresql://postgres:PASSWORD@host:5432/postgres"
   ./scripts/migrate.sh
   ```

Or if `migrate.sh` doesn't work:
```bash
# Manually run SQL files in order from backend/internal/database/migrations/
```

### 1.3 Create Storage Bucket

1. Supabase Dashboard → Storage
2. Create new bucket: `posts`
3. Set to public (allow unauthenticated access)
4. Configure CORS:
   ```
   Allowed origins: *
   Allowed methods: GET, POST, PUT, DELETE
   ```

### 1.4 Get Connection Secrets

Save these for backend configuration:
- **Database URL:** `postgresql://postgres:PASSWORD@host:5432/postgres`
- **JWT Secret:** `Settings → API → JWT Secret` (copy)
- **Service Role Key:** `Settings → API → Service Role Key` (copy)
- **Supabase URL:** `https://projectid.supabase.co`
- **Supabase Anon Key:** `Settings → API → Anon Key`

---

## Part 2: Backend Deployment (Railway)

### 2.1 Connect Repository to Railway

1. Go to [railway.app](https://railway.app)
2. Sign in with GitHub
3. Click "New Project" → "Deploy from GitHub repo"
4. Select `ghostline-backend` repository
5. Railway auto-detects Dockerfile

### 2.2 Configure Environment Variables

In Railway Dashboard → Project Settings → Variables:

```env
# Database
DATABASE_URL=postgresql://postgres:PASSWORD@host:5432/postgres
JWT_SECRET=your-jwt-secret-from-supabase

# Supabase
SUPABASE_URL=https://projectid.supabase.co
SUPABASE_KEY=your-service-role-key
SUPABASE_ANON_KEY=your-anon-key

# Server
PORT=8080
ENVIRONMENT=production
LOG_LEVEL=info

# Security
ALLOWED_ORIGIN=https://ghostline-frontend-five.vercel.app
CSRF_TOKEN_EXPIRY=3600

# Admin
ADMIN_PASSWORD=your-secure-admin-password

# Impersonation
IMPERSONATION_PASSWORD=your-impersonation-password
```

### 2.3 Trigger Deployment

1. Push code changes to GitHub
2. Railway auto-deploys on push
3. Watch build progress in Railway Dashboard
4. Get production URL: `Settings → Domains`

Example: `https://ghostline-backend-production-a17a.up.railway.app`

### 2.4 Verify Deployment

```bash
# Check health
curl https://ghostline-backend-production-a17a.up.railway.app/health

# Should return:
# {"status":"ok"}
```

---

## Part 3: Frontend Deployment (Vercel)

### 3.1 Connect Repository to Vercel

1. Go to [vercel.com](https://vercel.com)
2. Sign in or create account
3. Click "Add New..." → "Project"
4. Import `ghostline-frontend` repository
5. Vercel auto-detects Next.js/React

### 3.2 Configure Build Settings

In Vercel Project Settings:

**Build & Development Settings:**
- Framework: Vite
- Build Command: `npm run build`
- Output Directory: `dist`
- Install Command: `npm install`

### 3.3 Configure Environment Variables

In Vercel Project Settings → Environment Variables:

```env
VITE_API_BASE_URL=https://ghostline-backend-production-a17a.up.railway.app
```

**Important:** Use the exact Railway backend URL without trailing slash.

### 3.4 Trigger Deployment

1. Push code changes to GitHub
2. Vercel auto-deploys on push
3. Watch build progress in Vercel Dashboard
4. Get production URL: `Visit` button or `Settings → Domains`

Example: `https://ghostline-frontend-five.vercel.app`

### 3.5 Verify Deployment

1. Visit `https://ghostline-frontend-five.vercel.app`
2. Check if page loads without errors
3. Open browser console (F12) for errors

---

## Part 4: Configuration Verification

### 4.1 CORS Configuration

Frontend and backend must have matching CORS settings.

**Frontend URL:** `https://ghostline-frontend-five.vercel.app`

**Backend Environment Variable:**
```env
ALLOWED_ORIGIN=https://ghostline-frontend-five.vercel.app
```

If mismatch, you'll see error:
```
Access to XMLHttpRequest at 'https://backend.com/api/posts' from origin 
'https://frontend.com' has been blocked by CORS policy
```

### 4.2 API Connectivity

Test with browser console:
```javascript
fetch('https://ghostline-backend-production-a17a.up.railway.app/health')
  .then(r => r.json())
  .then(d => console.log(d))
```

Should return: `{status: "ok"}`

### 4.3 Authentication Testing

1. Visit frontend
2. Register new account
3. Verify JWT cookie is set (DevTools → Application → Cookies)
4. Verify requests include cookie (DevTools → Network → Headers)

---

## Part 5: Production Checklist

Before launching:

- [ ] Database running and migrations applied
- [ ] Backend deployed to Railway
- [ ] Frontend deployed to Vercel
- [ ] CORS configuration correct
- [ ] SSL certificates valid
- [ ] Error logging configured
- [ ] Monitoring alerts set up
- [ ] Backup strategy implemented
- [ ] Security headers verified
- [ ] Rate limiting enabled
- [ ] Admin password changed (not default)
- [ ] Secrets stored securely (not in repo)

### Security Checklist

- [ ] `.gitignore` includes `.env`
- [ ] No hardcoded secrets in code
- [ ] HTTPS enforced (no HTTP)
- [ ] HTTPOnly cookies enabled
- [ ] SameSite=Strict on auth cookie
- [ ] CORS restricts to frontend domain
- [ ] Rate limits configured
- [ ] Password hashing (Bcrypt cost 12)
- [ ] JWT expiration set (15 minutes)
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS prevention (HTML escaping)
- [ ] CSRF protection enabled

---

## Part 6: Monitoring & Logs

### Railway Logs

1. Railway Dashboard → My Services → ghostline-backend
2. View "Logs" tab for:
   - Deployment logs
   - Runtime logs
   - Error messages
   - Performance issues

### Vercel Logs

1. Vercel Dashboard → Deployments
2. Click deployment → "Logs"
3. View build and runtime logs

### Database Logs

1. Supabase Dashboard → Logs
2. View slow queries, errors, connections

---

## Part 7: Updates & Redeployment

### Update Backend

1. Make code changes in `ghostline-backend`
2. Commit and push to GitHub
3. Railway auto-deploys (watch Logs tab)
4. Verify `/health` endpoint returns `200`

### Update Frontend

1. Make code changes in `ghostline-frontend`
2. Commit and push to GitHub
3. Vercel auto-deploys (watch Deployments)
4. Verify frontend loads without 500 error

### Database Migrations

1. Create migration in `backend/internal/database/migrations/`
2. Test locally against dev database
3. Push to GitHub
4. Manual migration run on production:
   ```bash
   export DATABASE_URL="postgresql://..."
   go run cmd/migration/main.go
   ```

---

## Part 8: Scaling & Performance

### Database Scaling

**Free Tier (up to 500MB):**
- Good for testing
- 5 concurrent connections
- Limited backups

**Paid Tier (Scale as needed):**
1. Supabase Dashboard → Settings → Plan
2. Upgrade to plan matching storage needs
3. Automatic backups enabled

### Backend Scaling

Railway automatically:
- Scales horizontally (add instances)
- Load balances traffic
- Monitors CPU/memory

Monitor in Railway:
1. Dashboard → Metrics
2. View CPU, Memory, Network usage
3. Set auto-scale policy if needed

### Frontend Scaling

Vercel provides:
- Global CDN (cached at edge)
- Automatic image optimization
- Unlimited bandwidth (usage-based pricing)

---

## Part 9: Troubleshooting

### Frontend Won't Connect to Backend

**Error:** `Failed to fetch`, `Network error`

**Solution:**
1. Check `VITE_API_BASE_URL` in Vercel env vars
2. Verify Railway backend is running (`/health` returns 200)
3. Check CORS setting matches frontend URL
4. Clear browser cache (Cmd+Shift+Delete)

### CORS Policy Error

**Error:** `Access to XMLHttpRequest ... blocked by CORS policy`

**Solution:**
1. Verify backend `ALLOWED_ORIGIN` env var
2. Must be exact URL: `https://ghostline-frontend-five.vercel.app`
3. No trailing slash
4. Redeploy Railway after env var change

### Database Connection Failed

**Error:** `failed to connect to` database

**Solution:**
1. Verify `DATABASE_URL` in Railway env vars
2. Check Supabase project status (no maintenance)
3. Test connection locally
4. Check IP whitelist (Supabase usually open)

### Build Failed on Railway

**Error:** `Build failed`, `go: command not found`

**Solution:**
1. Check Dockerfile is in repo root
2. Verify `Dockerfile` not `dockerfile`
3. Check Go version in Dockerfile: `golang:1.25-alpine`
4. Rebuild: Railway Dashboard → "Trigger Redeploy"

### Build Failed on Vercel

**Error:** `npm ERR!`, `Build failed`

**Solution:**
1. Check `npm install` completes locally
2. Verify `buildCommand` is correct: `npm run build`
3. Check Node version compatibility
4. Rebuild: Vercel Dashboard → "Redeploy"

---

## Part 10: Custom Domain

### Connect Domain to Vercel

1. Vercel Dashboard → Settings → Domains
2. Enter your domain (e.g., `ghostline.com`)
3. Follow DNS configuration steps
4. Add CNAME record: `name=www, value=cname.vercel-dns.com`
5. Wait 24 hours for propagation
6. SSL certificate auto-provisioned

### Connect Domain to Railway

Railway doesn't support custom domains on free tier.

**Option 1:** Use Vercel domain + proxy
```
Frontend: ghostline.com (Vercel)
Backend: api.ghostline.com (Railway) via proxy
```

**Option 2:** Upgrade Railway to paid tier
Then configure custom domain in Railway Dashboard.

---

## Part 11: Disaster Recovery

### Database Backup

Supabase automatically:
- Daily backups (7-day retention free tier)
- 30-day retention (paid tier)
- Point-in-time recovery (PITR)

Restore from Supabase Dashboard:
1. Settings → Backups
2. Click restore on any backup date
3. Confirm (data will be overwritten)

### Application Backup

GitHub serves as backup:
- All code in commit history
- All database migrations
- All documentation

Restore:
1. Clone from GitHub
2. Run migrations against new database
3. Redeploy to Railway/Vercel

---

## Production Environment Variables Checklist

```env
# Database
DATABASE_URL=postgresql://user:password@host:5432/db
JWT_SECRET=long-random-secret-key-min-32-chars

# Supabase
SUPABASE_URL=https://projectid.supabase.co
SUPABASE_KEY=eyJhbGciOiJIUzI1NiI...
SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiI...

# Server Configuration
PORT=8080
ENVIRONMENT=production
LOG_LEVEL=info
CORS_ALLOWED_ORIGIN=https://ghostline-frontend.com
REQUEST_TIMEOUT=30s

# Security
ADMIN_PASSWORD=secure-admin-password-min-12-chars
IMPERSONATION_PASSWORD=secure-password-min-12-chars
SESSION_TIMEOUT=900  # 15 minutes
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=3600  # 1 hour

# Frontend Configuration
VITE_API_BASE_URL=https://backend.example.com
VITE_APP_NAME=Ghostline
VITE_ENABLE_ANALYTICS=true

# Optional
SENTRY_DSN=https://...
DATADOG_API_KEY=...
SLACK_WEBHOOK_URL=...
```

---

Last Updated: March 28, 2026
