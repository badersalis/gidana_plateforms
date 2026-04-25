# Gidana Backend — Setup Guide

## Requirements
- Go 1.21+ → https://go.dev/dl/
- PostgreSQL (production) or SQLite (dev, zero setup)

## Quick Start (SQLite dev mode)

```bash
# 1. Copy env file
cp .env.example .env

# 2. Download dependencies
go mod tidy

# 3. Run server (auto-migrates DB)
go run ./cmd/main.go
```

Server starts at **http://localhost:8080**

## API Base URL
```
http://localhost:8080/api/v1
```

## Key Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /auth/register | No | Register |
| POST | /auth/login | No | Login (returns JWT) |
| GET | /auth/me | JWT | Current user |
| GET | /properties | Optional | List/search properties |
| GET | /properties/featured | No | Top 4 by rating |
| POST | /properties | JWT | Create property (multipart) |
| POST | /favorites/:id/toggle | JWT | Toggle favorite |
| GET | /wallets | JWT | List wallets |
| POST | /wallets | JWT | Add wallet |
| GET | /transactions | JWT | Transaction history |
| POST | /transactions/pay-service | JWT | Pay Starlink/Canal+ |
| POST | /transactions/transfer | JWT | Transfer money |

## Authentication
All protected routes require: `Authorization: Bearer <token>`

## PostgreSQL (Production)
Set `DATABASE_URL=postgres://user:pass@host:5432/gidana_db` in `.env`
