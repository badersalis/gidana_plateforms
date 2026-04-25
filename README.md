# Gidana Platforms — Backend API

> Find your dream spot. Compare locations, contact verified owners, schedule visits, and sign contracts — all in one place.

**Gidana** is a property discovery and rental platform built for West Africa. It lets anyone search for studios, apartments, or houses within their budget, view photos, contact verified owners, book appointments, and manage contracts entirely in the cloud.

---

## What It Does

- **Discover properties** — search and filter studios, apartments, and houses by neighborhood, price, type, and amenities
- **Compare listings** — view detailed specs (rooms, surface, utilities) and real user reviews side by side
- **Contact verified owners** — reach owners via WhatsApp or phone directly from the listing
- **Book & manage rentals** — schedule appointments and track rental status end-to-end
- **Sign & store contracts** — contracts are cloud-stored so they're never lost
- **Property alerts** — get notified when a matching listing appears in your target neighborhood
- **Integrated wallet & payments** — pay for services and transfer funds via Nita, M-Pesa, Visa, Mastercard, or PayPal

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) v1.9 |
| ORM | [GORM](https://gorm.io) v1.25 |
| Database (prod) | PostgreSQL |
| Database (dev) | SQLite (zero config) |
| Authentication | JWT — HS256 via `golang-jwt/jwt` v5 |
| Cloud Storage | Firebase Storage (optional) |
| Password Hashing | `golang.org/x/crypto` (bcrypt) |
| Config | `godotenv` |

---

## Project Structure

```
gidana_plateforms/
├── cmd/
│   └── main.go                  # Entry point
├── internal/
│   ├── config/config.go         # Env-based configuration
│   ├── database/database.go     # Connection & auto-migration
│   ├── handlers/                # HTTP request handlers
│   │   ├── auth.go
│   │   ├── property.go
│   │   ├── rental.go
│   │   ├── review.go
│   │   ├── favorite.go
│   │   ├── alert.go
│   │   ├── search.go
│   │   ├── transaction.go
│   │   ├── upload.go
│   │   ├── user.go
│   │   └── wallet.go
│   ├── middleware/auth.go        # JWT middleware (Auth + OptionalAuth)
│   ├── models/                  # GORM models
│   │   ├── user.go
│   │   ├── property.go
│   │   ├── property_image.go
│   │   ├── rental.go
│   │   ├── review.go
│   │   ├── favorite.go
│   │   ├── alert.go
│   │   ├── wallet.go
│   │   ├── transaction.go
│   │   └── search_history.go
│   ├── routes/routes.go         # Route definitions
│   ├── storage/storage.go       # Firebase storage integration
│   └── utils/
│       ├── jwt.go               # Token generation & parsing
│       ├── password.go          # Hashing helpers
│       └── response.go          # Standardized JSON responses
├── .env.example                 # Environment template
├── Makefile                     # Build shortcuts
├── go.mod
└── go.sum
```

---

## Getting Started

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- PostgreSQL (production) **or** nothing extra (development uses SQLite automatically)
- Optional: Firebase project for cloud image storage

### Quick Start (SQLite — zero config)

```bash
# Clone the repo
git clone https://github.com/your-org/gidana_plateforms.git
cd gidana_plateforms

# Copy the env template
cp .env.example .env

# Install dependencies
make tidy

# Start the dev server (auto-migrates on first run)
make run
```

The server starts at `http://localhost:8080`. The database file `gidana_dev.db` is created automatically.

### Production Setup (PostgreSQL)

1. Create a PostgreSQL database:
   ```sql
   CREATE DATABASE gidana_db;
   ```

2. Set `APP_ENV=production` and `DATABASE_URL` in your `.env`:
   ```env
   APP_ENV=production
   DATABASE_URL=postgres://user:password@localhost:5432/gidana_db
   JWT_SECRET=your-very-long-random-secret
   ```

3. Build and run:
   ```bash
   make build
   ./gidana_api
   ```

---

## Environment Variables

Copy `.env.example` to `.env` and fill in the values:

```env
# App
APP_ENV=development          # development | production
PORT=8080

# Database
DATABASE_URL=postgres://user:password@localhost:5432/gidana_db   # production only
DB_PATH=gidana_dev.db        # development SQLite path

# Auth
JWT_SECRET=change-me-in-production
JWT_EXPIRY_HOURS=72

# File Uploads
UPLOAD_DIR=./uploads/properties
MAX_UPLOAD_SIZE_MB=5

# Firebase Storage (optional)
USE_FIREBASE=false
FIREBASE_CREDENTIALS_PATH=./firebase-credentials.json
FIREBASE_BUCKET=your-bucket.appspot.com

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:19006
```

---

## API Reference

**Base URL:** `http://localhost:8080/api/v1`

All protected endpoints require an `Authorization: Bearer <token>` header.

### Authentication

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/auth/register` | — | Register a new user |
| `POST` | `/auth/login` | — | Login and receive JWT |
| `GET` | `/auth/me` | JWT | Get authenticated user profile |

### Users

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `PUT` | `/users/profile` | JWT | Update profile info |
| `POST` | `/users/profile-picture` | JWT | Upload profile picture |
| `PUT` | `/users/password` | JWT | Change password |

### Properties

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/properties` | Optional | List & search properties |
| `GET` | `/properties/featured` | — | Top 4 featured properties |
| `GET` | `/properties/:id` | Optional | Get single property |
| `POST` | `/properties` | JWT | Create property listing |
| `PUT` | `/properties/:id` | JWT | Update property |
| `DELETE` | `/properties/:id` | JWT | Delete property |
| `PATCH` | `/properties/:id/availability` | JWT | Toggle availability |
| `GET` | `/properties/my/listings` | JWT | Get current user's listings |
| `POST` | `/properties/:id/images` | JWT | Add images to a property |
| `DELETE` | `/images/:id` | JWT | Delete a property image |
| `PATCH` | `/images/:id/main` | JWT | Set image as main photo |

### Reviews

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/properties/:id/reviews` | — | Get reviews for a property |
| `POST` | `/properties/:id/reviews` | JWT | Submit a review |
| `DELETE` | `/reviews/:id` | JWT | Delete a review |

### Favorites

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/favorites` | JWT | Get saved favorites |
| `POST` | `/favorites/:id/toggle` | JWT | Add or remove a favorite |

### Rentals

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/rentals` | JWT | Get user's rentals |
| `POST` | `/rentals` | JWT | Create a rental request |
| `PATCH` | `/rentals/:id/status` | JWT | Update rental status |

### Wallets

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/wallets` | JWT | List user's wallets |
| `POST` | `/wallets` | JWT | Add a wallet |
| `PUT` | `/wallets/:id` | JWT | Update wallet details |
| `DELETE` | `/wallets/:id` | JWT | Remove a wallet |
| `PATCH` | `/wallets/:id/select` | JWT | Set active wallet |
| `POST` | `/wallets/:id/refresh-balance` | JWT | Sync wallet balance |

Supported providers: **Nita**, **M-Pesa**, **Visa**, **Mastercard**, **PayPal**

### Transactions

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/transactions` | JWT | Get transaction history |
| `POST` | `/transactions/pay-service` | JWT | Pay for a service (Starlink, Canal+, etc.) |
| `POST` | `/transactions/transfer` | JWT | Transfer between wallets |

### Alerts

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/alerts` | JWT | Get active alerts |
| `POST` | `/alerts` | JWT | Create a property alert |
| `PUT` | `/alerts/:id` | JWT | Update alert criteria |
| `DELETE` | `/alerts/:id` | JWT | Delete an alert |

### Search

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/search/suggestions` | — | Autocomplete suggestions |
| `POST` | `/search/history` | Optional | Save a search term |
| `GET` | `/search/history` | JWT | Get search history |
| `DELETE` | `/search/history` | JWT | Clear search history |

### Health

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/health` | — | API health check |

---

## Authentication Flow

1. Register via `POST /auth/register` or login via `POST /auth/login`
2. Both return a JWT in the response body
3. Include it on subsequent requests:
   ```
   Authorization: Bearer <your_token>
   ```
4. Tokens expire after **72 hours** (configurable via `JWT_EXPIRY_HOURS`)

---

## Data Models

### Property Types
`Studio` · `Appartement` · `Maison`

### Transaction Types
`À louer` · `À vendre`

### Rental Statuses
`pending` · `occupied` · `available` · `completed`

### Transaction Natures
`expense` · `income`

### Transaction Statuses
`done` · `failed` · `ongoing`

---

## Makefile

```bash
make run      # Start the dev server
make build    # Compile to ./gidana_api binary
make tidy     # Sync go.mod and go.sum
```

---

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Commit your changes
4. Open a pull request

---

<a
  href="https://render.com/deploy?repo=https://github.com/badersalis/gidana_plateforms"
>
  <img
    src="https://render.com/images/deploy-to-render-button.svg"
    alt="Deploy to Render"
  />
</a>

---

## License

Proprietary — Gidana Platforms. All rights reserved.
