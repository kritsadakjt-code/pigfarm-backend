<div align="center">

# Pig Farm Management Backend

**A REST API for running a pig farm's day-to-day operations** — pigs,
breeding cycles, feeding schedules, food stock, expenses, sales, and
automated notifications — built in Go with Clean Architecture.

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-00ACD7?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-00ADD8?style=for-the-badge)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=jsonwebtokens)
![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)

</div>

---

## Overview

Running a pig farm involves tracking a lot of moving parts at once: which
pigs are being bred and when, whether feed stock is running low, what
today's feeding schedule requires, what the farm is spending money on, and
when a sale is finalized. This API models that operational reality as a
single, coherent domain — rather than a generic CRUD scaffold — with a
background scheduler that turns time-based rules (feeding schedules, stock
thresholds, breeding due dates) into automatic notifications and actions
without a human needing to check in constantly.

### Features

- **Auth & Access Control:** JWT-based login with a role claim
  (`owner`) enforced via middleware on admin-only routes — not just present
  in the data model, but actually checked on every protected request.
- **User Management:** full admin control over farm staff
  accounts — create, list, search, update, and delete users, plus an
  approval workflow (approve/reject) for new registrations before they
  can access the system. All routes are gated by the `owner`-role guard.
- **Pig Management:** full CRUD for individual pigs.
- **Breeding Management:** breeding record tracking with paginated listing.
- **Feeding & Feeding Schedules:** manual feeding logs, plus recurring
  feeding rules that a background job executes automatically at the
  scheduled time — no manual trigger required.
- **Food Stock Tracking:** stock levels feed directly into low-stock
  notifications, so the farm finds out before feed actually runs out.
- **Expense Tracking:** ongoing farm expense records.
- **Automated Notifications:** a background scheduler runs every minute,
  checking food stock levels and upcoming breeding events and generating
  notifications proactively rather than on-demand.
- **Dashboard:** an aggregated summary endpoint for at-a-glance farm status.
- **Pig Sales:** sales tracking, wired end-to-end from route to database model.

## Architecture

Clean Architecture, with business logic kept independent of delivery and
persistence mechanisms:

```
entities/         → core domain types
models/            → GORM models (persistence layer)
dto/               → request/response shapes
mappers/           → model ⇄ DTO conversion
usecases/          → business logic + repository interfaces
adapters/
  handlers/        → Fiber HTTP handlers
  repositories/    → GORM implementations of usecase repository interfaces
  schedulers/       → background job (automated notifications + auto-feeding)
middlewares/       → JWT auth, role guard, CORS, error handling
routes/            → Fiber route registration
config/            → database connection setup
```

The `usecases` layer defines repository interfaces that `adapters/repositories`
implements — business rules don't know or care that the underlying store
happens to be PostgreSQL via GORM, which keeps the domain logic testable
and swappable in principle.

## Tech Stack

| Layer | Technology |
|---|---|
| **Language / Framework** | Go, Fiber v2 |
| **Database** | PostgreSQL via GORM |
| **Auth** | JWT (`golang-jwt/jwt/v5`), role-based route guard |
| **Validation** | `go-playground/validator` |
| **Email** | Mailjet API (password reset, email verification) |
| **Infra** | Docker Compose |

## Getting Started

### Prerequisites
- Go 1.24+
- Docker & Docker Compose

### Setup

```bash
git clone https://github.com/kritsadakjt-code/pigfarm-backend.git
cd pigfarm-backend
go mod download
```

Create a `.env` file with the following:

```dotenv
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=changeme
DB_NAME=pigfarm
DB_SSLMODE=disable
PORT=8000
JWT_SECRET=replace-with-a-long-random-string
 
MAILJET_API_KEY=your-mailjet-api-key
MAILJET_API_SECRET=your-mailjet-api-secret
MAILJET_SENDER_EMAIL=you@example.com
MAILJET_SENDER_NAME=pigfarm

FRONTEND_URL=https://your-frontend.vercel.app
FRONTEND_URL_DEV=http://localhost:3000
 
APP_ENV=dev
# APP_ENV=prod

MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=pigfarm
```

### Run Postgres

```bash
docker-compose up -d
```

### Run the app

```bash
go run main.go
```

Tables are auto-migrated on startup via GORM (`AutoMigrate`), and a
default owner account is seeded automatically.

## Project Structure

```
.
├── adapters/
│   ├── handlers/        # Fiber HTTP handlers
│   ├── repositories/    # GORM repository implementations
│   └── schedulers/       # background notification/auto-feeding job
├── config/               # database connection setup
├── dto/
├── entities/
├── mappers/
├── middlewares/
├── models/
├── routes/
├── usecases/             # business logic + repository interfaces
├── utils/
└── main.go
```

## License

MIT — see [LICENSE](./LICENSE) for details.
