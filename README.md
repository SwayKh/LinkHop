# URL Shortener

A simple URL shortener built with Go and SQLite. Shorten long URLs, track click counts, and manage your links.

## Features

- Signup and login with email/password (bcrypt hashing, JWT sessions)
- Create short links with custom aliases or auto-generated codes
- Edit and delete your links
- Click count tracking and last accessed date
- Analytics page for each link
- Responsive UI with Tailwind CSS (server-rendered HTML)

## Requirements

- Go 1.22 or newer

## Quick Start

```bash
go mod tidy
go build -o urlshortner .
./urlshortner
```

Open http://localhost:8080 in your browser.

## Configuration

Set these environment variables before running:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `JWT_SECRET` | `urlshortner-dev-secret-key` | Secret key for signing auth tokens. **Set a random value in production.** |
| `DB_PATH` | `./urlshortner.db` | Path to the SQLite database file |
| `COOKIE_SECURE` | `false` | Set to `true` when behind HTTPS |

Create a `.env` file in the project root (auto-loaded on startup):

```
PORT=8080
JWT_SECRET=put-a-random-string-here
DB_PATH=./urlshortner.db
COOKIE_SECURE=false
```

## Routes

### Public
| Route | Description |
|-------|-------------|
| GET / | Landing page |
| GET /signup | Sign up form |
| POST /signup | Create account |
| GET /login | Login form |
| POST /login | Authenticate |
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
  db.go           -- SQLite setup and table creation
  user.go         -- User queries
  link.go         -- Link queries
handlers/
  auth.go         -- Signup, login, logout handlers
  links.go        -- CRUD, redirect, analytics handlers
middleware/
  auth.go         -- JWT authentication middleware
render/
  render.go       -- Template engine
  templates/      -- HTML templates
```

## Database

SQLite database file (urlshortner.db) is created automatically on first run. Tables:

- **users** -- id, email, password_hash, created_at
- **links** -- id, user_id, original_url, short_code, custom_alias, click_count, last_accessed, created_at, updated_at

## Dependencies

- github.com/joho/godotenv -- .env file loading
- github.com/golang-jwt/jwt/v5 -- JWT tokens
- golang.org/x/crypto -- bcrypt password hashing
- modernc.org/sqlite -- Pure Go SQLite driver (no CGO required)
