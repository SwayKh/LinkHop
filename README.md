# URL Shortener

A simple URL shortener built with Go and PostgreSQL. Shorten long URLs, track click counts, and manage your links.

## Features

- Signup with username, email, and password
- Email verification via 6-digit OTP (console in dev, SMTP in production)
- Login with email or username
- Create short links with custom aliases or auto-generated codes
- Edit and delete your links
- Click count tracking and last accessed date
- Analytics page for each link
- Responsive UI with Tailwind CSS (server-rendered HTML)

## Requirements

- Go 1.22 or newer
- PostgreSQL 14 or newer

## Quick Start

1. Create a PostgreSQL database:

```bash
createdb urlshortner
```

2. Copy and edit the env file:

```bash
cp .env.example .env
```

3. Build and run:

```bash
go mod tidy
go build -o urlshortner .
./urlshortner
```

Open http://localhost:8080 in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `JWT_SECRET` | `urlshortner-dev-secret-key` | Secret key for signing auth tokens. **Set a random value in production.** |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/urlshortner?sslmode=disable` | PostgreSQL connection string |
| `COOKIE_SECURE` | `false` | Set to `true` when behind HTTPS |
| `SMTP_HOST` | (empty) | SMTP server for sending OTP emails. Leave empty to log codes to console (dev mode). |
| `SMTP_PORT` | `587` | SMTP port |
| `SMTP_USER` | (empty) | SMTP username |
| `SMTP_PASS` | (empty) | SMTP password |
| `SMTP_FROM` | (empty) | From address for emails (defaults to SMTP_USER) |

The `DATABASE_URL` format is:

```
postgres://user:password@host:port/database?sslmode=require
```

Hosting platforms (Render, Railway, Fly.io, Heroku) provide this as an environment variable automatically. Just connect a PostgreSQL database and the app picks it up.

## Routes

### Public
| Route | Description |
|-------|-------------|
| GET / | Landing page |
| GET /signup | Sign up form |
| POST /signup | Create account (sends OTP) |
| GET /verify | OTP verification form |
| POST /verify | Verify email with OTP |
| POST /resend-otp | Resend verification code |
| GET /login | Login form |
| POST /login | Authenticate (email or username) |
| POST /logout | Logout |

### Protected (requires login)
| Route | Description |
|-------|-------------|
| GET /dashboard | View your links |
| GET /links/new | Create link form |
| POST /api/links | Create a link |
| GET /links/{id}/edit | Edit link form |
| POST /api/links/{id}/update | Update a link |
| POST /api/links/{id}/delete | Delete a link |
| GET /links/{id}/analytics | View link stats |

### Redirect
| Route | Description |
|-------|-------------|
| GET /{shortCode} | Redirect to original URL |

## Project Structure

```
main.go           -- Entry point and router
database/
  db.go           -- PostgreSQL setup and table creation
  user.go         -- User queries
  link.go         -- Link queries
email/
  email.go        -- OTP email sender (console/SMTP)
handlers/
  auth.go         -- Signup, verify, login, logout handlers
  links.go        -- CRUD, redirect, analytics handlers
middleware/
  auth.go         -- JWT authentication middleware
render/
  render.go       -- Template engine
  templates/      -- HTML templates
```

## Database

Tables are created automatically on first run. Schema:

- **users** -- id (SERIAL), email (TEXT UNIQUE), username (TEXT UNIQUE), password_hash (TEXT), is_verified (BOOLEAN), created_at (TIMESTAMPTZ)
- **verification_codes** -- id (SERIAL), email (TEXT), code (TEXT), expires_at (TIMESTAMPTZ), used (BOOLEAN), created_at (TIMESTAMPTZ)
- **links** -- id (SERIAL), user_id (FK), original_url (TEXT), short_code (TEXT UNIQUE), custom_alias (TEXT), click_count (INT), last_accessed (TIMESTAMPTZ), created_at (TIMESTAMPTZ), updated_at (TIMESTAMPTZ)

## Dependencies

- github.com/joho/godotenv -- .env file loading
- github.com/jackc/pgx/v5 -- PostgreSQL driver
- github.com/golang-jwt/jwt/v5 -- JWT tokens
- golang.org/x/crypto -- bcrypt password hashing
