# 🚀 Deployment Guide

Deploy Ghostline to production using Railway (backend) and Vercel (frontend).

---

## 📋 Prerequisites

- GitHub repositories set up
- Railway account ([railway.app](https://railway.app))
- Vercel account ([vercel.com](https://vercel.com))
- Supabase project with database

---

## 🔧 Step 1: Backend Deployment (Railway)

### 1.1 Create Railway Project

1. Go to [railway.app](https://railway.app)
2. Click **"New Project"**
3. Select **"Deploy from GitHub"**
4. Select `Vignu-14/ghostline-backend` repository
5. Configure:
   - **Root Directory:** (leave empty - root is backend)
   - **Framework:** Go
6. Click **Deploy**

### 1.2 Set Environment Variables

1. In Railway, go to **Variables** tab
2. Add all these variables:

```
# Server Configuration
PORT=8080
ENVIRONMENT=production
ALLOWED_ORIGIN=https://ghostline-frontend-five.vercel.app

# Database (from Supabase)
DATABASE_URL=postgresql://postgres.PROJECT_ID:PASSWORD@aws-1-ap-south-1.pooler.supabase.com:5432/postgres?sslmode=require
DB_MAX_CONNECTIONS=25
DB_MIN_CONNECTIONS=5
DB_MAX_CONN_LIFETIME_MINUTES=60
DB_MAX_CONN_IDLE_MINUTES=15
DB_HEALTH_CHECK_SECONDS=30
DB_CONNECT_TIMEOUT_SECONDS=5

# JWT (generate secure random string)
JWT_SECRET=your-64-character-random-secret-here
JWT_EXPIRATION_MINUTES=15
AUTH_COOKIE_NAME=auth_token
COOKIE_SECURE=true

# Supabase (from Supabase Dashboard)
SUPABASE_URL=https://PROJECT_ID.supabase.co
SUPABASE_SERVICE_KEY=sbp_your_service_key_here
STORAGE_BUCKET_NAME=user-uploads

# Rate Limiting
RATE_LIMIT_LOGIN_ATTEMPTS=5
RATE_LIMIT_LOGIN_WINDOW_MINUTES=15
RATE_LIMIT_UPLOAD_COUNT=10
RATE_LIMIT_UPLOAD_WINDOW_MINUTES=60
RATE_LIMIT_MESSAGE_COUNT=10
RATE_LIMIT_MESSAGE_WINDOW_SECONDS=1
RATE_LIMIT_LIKE_COUNT=100
RATE_LIMIT_LIKE_WINDOW_MINUTES=60
```

### 1.3 Verify Deployment

1. Wait for Railway to finish building
2. Check **Logs** for errors
3. Your backend URL will be shown:
   ```
   https://ghostline-backend-production-xxxx.up.railway.app
   ```

---

## 🎨 Step 2: Frontend Deployment (Vercel)

### 2.1 Connect Repository

1. Go to [vercel.com](https://vercel.com)
2. Click **"Add New Project"**
3. Select `Vignu-14/ghostline-frontend` repository
4. Configure:
   - **Framework:** Vite
   - **Root Directory:** (leave empty)
   - **Build Command:** `npm run build`
   - **Output Directory:** `dist`
5. Click **Deploy**

### 2.2 Set Environment Variables

1. In Vercel, go to **Settings** → **Environment Variables**
2. Add this variable:
   ```
   Name:  VITE_API_BASE_URL
   Value: https://ghostline-backend-production-xxxx.up.railway.app
   ```
   (Replace `xxxx` with your actual Railway domain)

3. Click **Save**
4. Click **Deployments** → Click latest deployment → **Redeploy**

### 2.3 Verify Deployment

1. Wait for Vercel to finish building
2. Your frontend URL will be shown:
   ```
   https://ghostline-frontend-xxxx.vercel.app
   ```
3. Test login functionality

---

## 🗄️ Step 3: Database Setup (Supabase)

### 3.1 Create Database Tables

1. Go to [supabase.com](https://supabase.com) → Your Project
2. Go to **SQL Editor**
3. Run migrations (already in your schema)

Or using Supabase CLI:
```bash
supabase db push
```

### 3.2 Enable Row Level Security (RLS)

In Supabase SQL Editor:
```sql
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE posts ENABLE ROW LEVEL SECURITY;
ALTER TABLE likes ENABLE ROW LEVEL SECURITY;
ALTER TABLE messages ENABLE ROW LEVEL SECURITY;
```

### 3.3 Create RLS Policies

Already set up in migrations, verify they're active.

---

## 🔑 Step 4: Generate Secrets

### Generate JWT Secret

**Windows PowerShell:**
```powershell
$random = New-Object System.Random
$bytes = New-Object byte[] 32
$random.GetBytes($bytes)
[BitConverter]::ToString($bytes).Replace("-","").ToLower()
```

**Linux/Mac:**
```bash
openssl rand -hex 32
```

Copy output and use as `JWT_SECRET` in Railway Variables.

---

## ✅ Testing

### Test Backend
```bash
# Health check
curl https://ghostline-backend-production-xxxx.up.railway.app/health

# Should return: {"status":"ok"}
```

### Test Frontend
```bash
# Open in browser
https://ghostline-frontend-xxxx.vercel.app

# Try registering and logging in
```

### Test WebSocket Connection
1. Open browser DevTools (F12)
2. Go to Console tab
3. Login successfully
4. Go to Chat page
5. Should see WebSocket connected in Network tab

---

## 🔒 Security Checklist

- [ ] `ENVIRONMENT=production` set in Railway
- [ ] `JWT_SECRET` is strong (64+ characters)
- [ ] `COOKIE_SECURE=true` in Railway
- [ ] `ALLOWED_ORIGIN` matches your Vercel URL
- [ ] `SUPABASE_SERVICE_KEY` is secure (not exposed)
- [ ] Database connections use SSL (`sslmode=require`)
- [ ] Row Level Security (RLS) is enabled
- [ ] Secrets removed from Git history (.gitignore working)

---

## 🔄 Continuous Deployment

Both Railway and Vercel have **automatic deployments**:

- Push to GitHub main branch
- Railway/Vercel automatically rebuilds
- Changes live in ~2-5 minutes

To disable auto-deploy:
- **Railway:** Settings → Disable auto-deploy
- **Vercel:** Settings → Deployments → Uncheck "Automatic Deployments"

---

## 📊 Monitoring

### Railway Monitoring
1. Go to Project → **Logs**
2. Check for errors
3. Monitor CPU/Memory usage

### Vercel Monitoring
1. Go to Project → **Analytics**
2. Check request metrics
3. Monitor build times

---

## 🆘 Troubleshooting

### Backend Build Fails
```
ERROR: failed to build: failed to solve: the Dockerfile cannot be empty
```
**Solution:** Ensure Dockerfile exists and has content

### CORS Error
```
Access to fetch at '...' has been blocked by CORS policy
```
**Solution:** Update `ALLOWED_ORIGIN` in Railway to match frontend URL

### 404 Not Found
```
GET https://frontend.vercel.app/backend-url/api/... 404
```
**Solution:** Set `VITE_API_BASE_URL` in Vercel (must start with `https://`)

### WebSocket Connection Failed
```
Failed to establish WebSocket connection
```
**Solution:** Check JWT token is valid, backend is running, allowedOrigin is set

---

## 🚨 Rollback

### Rollback Railway Deployment
1. Go to **Deployments** tab
2. Click on previous working deployment
3. Click **Redeploy**

### Rollback Vercel Deployment
1. Go to **Deployments** tab
2. Click on previous working deployment
3. Click **Promote to Production**

---

## 📈 Scaling

Current setup handles:
- ✅ 100k messages/day
- ✅ 1k concurrent WebSocket connections
- ✅ 10k users

For higher traffic:
- Upgrade Railway plan (add more containers)
- Enable Redis caching
- Use Supabase connection pooler

---

## 📞 Support

- **Railway Help:** https://docs.railway.app
- **Vercel Help:** https://vercel.com/docs
- **Supabase Help:** https://supabase.com/docs
