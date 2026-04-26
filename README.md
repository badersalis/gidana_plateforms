# Gidana — Property Discovery & Rental Platform for Africa

Gidana is a property discovery and rental platform built for Africa. It lets anyone search for studios, apartments, or houses within their budget, view photos, contact verified owners, book appointments, and manage contracts entirely in the cloud.

---

## Features

- **Property search** — filter by type (studio, apartment, house), location, price range, and availability
- **Photo galleries** — multi-image listings with a designated main photo
- **Verified owners** — authenticated accounts with profile pictures and contact info
- **Rental management** — create, track, and update rental agreements with status workflows
- **Favorites** — save and revisit properties of interest
- **Reviews** — tenants can leave ratings and written feedback on properties
- **Wallet & payments** — manage payment wallets, transfer money, and pay for services
- **Smart alerts** — set price/location alerts to be notified when matching listings appear
- **Search history** — logged-in users retain their search history for quick re-discovery
- **Firebase Storage** — optional cloud image storage via Firebase (falls back to local disk)

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| Web framework | Gin |
| ORM | GORM |
| Database (prod) | PostgreSQL |
| Database (dev) | SQLite |
| Auth | JWT (golang-jwt/jwt v5) |
| Image storage | Firebase / Google Cloud Storage (optional) |
| Config | godotenv |

---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL (production) or SQLite (development, no setup needed)
- Firebase project credentials (optional — only for cloud image storage)

### 1. Clone the repo

```bash
git clone https://github.com/badersalis/gidana_backend.git
cd gidana_backend
```

### 2. Configure environment

```bash
cp .env.example .env
```

Edit `.env` with your values:

```env
APP_ENV=development
PORT=8080

# Database — leave DATABASE_URL empty to use SQLite in dev
DATABASE_URL=postgres://user:password@localhost:5432/gidana_db
DB_PATH=gidana_dev.db

JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRY_HOURS=72

UPLOAD_DIR=./uploads/properties
MAX_UPLOAD_SIZE_MB=5

# Firebase (optional)
FIREBASE_CREDENTIALS_PATH=./firebase-credentials.json
FIREBASE_BUCKET=your-bucket.appspot.com
USE_FIREBASE=false

ALLOWED_ORIGINS=http://localhost:3000
```

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Run the server

```bash
make run
# or
go run ./cmd/main.go
```

The API will be available at `http://localhost:8080`.

---

## API Overview

Base path: `/api/v1`

| Group | Endpoints |
|---|---|
| Auth | `POST /auth/register`, `POST /auth/login`, `GET /auth/me` |
| Users | `PUT /users/profile`, `POST /users/profile-picture`, `PUT /users/password` |
| Properties | `GET/POST /properties`, `GET/PUT/DELETE /properties/:id` |
| Rentals | `GET/POST /rentals`, `PATCH /rentals/:id/status` |
| Favorites | `GET /favorites`, `POST /favorites/:id/toggle` |
| Reviews | `GET/POST /properties/:id/reviews`, `DELETE /reviews/:id` |
| Wallets | `GET/POST /wallets`, `PUT/DELETE /wallets/:id`, `POST /wallets/:id/refresh-balance` |
| Transactions | `GET /transactions`, `POST /transactions/pay-service`, `POST /transactions/transfer` |
| Alerts | `GET/POST /alerts`, `PUT/DELETE /alerts/:id` |
| Search | `GET /search/suggestions`, `GET/POST/DELETE /search/history` |
| Health | `GET /health` |

---

## Deployment

### Deploy to Render

Click the button at the top of this file or follow these steps manually:

1. Push the repo to GitHub.
2. Create a new **Web Service** on [Render](https://render.com).
3. Set the **Build Command** to `go build -o gidana_api ./cmd/main.go`.
4. Set the **Start Command** to `./gidana_api`.
5. Add all environment variables from `.env.example` in the Render dashboard.
6. Attach a **PostgreSQL** database (Render provides one) and set `DATABASE_URL`.

### Build manually

```bash
make build
./gidana_api
```

---

## Project Structure

```
cmd/
  main.go               — entry point
internal/
  config/               — env loading
  database/             — DB connection & auto-migration
  handlers/             — HTTP request handlers
  middleware/           — JWT auth middleware
  models/               — GORM models
  routes/               — route registration
  storage/              — Firebase storage client
  utils/                — JWT, password hashing, response helpers
```

---

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/badersalis/gidana_plateforms)

---

## License

