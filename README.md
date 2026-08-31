# Apollo

OSINT and blockchain intelligence backend. You give it an Ethereum address, it gives you back balances, token holdings, transaction history, who that address actually deals with, and a graph you can render.

Ethereum mainnet only right now. Data comes from Alchemy, so there is no node to run and no chain to index.

## Stack

- Go 1.26 with Gin
- PostgreSQL 18 with GORM
- JWT auth, bcrypt passwords
- Alchemy for chain data

## Setup

You need Go 1.26+, Docker, and an Alchemy API key from [dashboard.alchemy.com](https://dashboard.alchemy.com).

Clone it and start the database:

```bash
docker compose up -d
```

Copy the env file and fill it in:

```bash
cp .env.example .env
```

`.env` needs three values:

```
DATABASE_URL="postgres://apollo:apollo@localhost:5432/apollo?sslmode=disable"
JWT_SECRET="paste the output of: openssl rand -base64 64"
ALCHEMY_API_KEY="your key"
```

Run it:

```bash
cd backend && go run ./cmd/api
```

Server comes up on `:8080`. Tables are created automatically on boot and the label table is seeded with known addresses. Check it is alive:

```bash
curl localhost:8080/health
```

## How to use

Everything except `/health` and the two auth routes needs a bearer token.

### 1. Register

```bash
curl -X POST localhost:8080/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"you@example.com","password":"atleast8chars"}'
```

### 2. Log in and keep the token

```bash
TOKEN=$(curl -s -X POST localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"you@example.com","password":"atleast8chars"}' \
  | sed 's/.*"token":"\([^"]*\)".*/\1/')
```

Tokens are good for 24 hours.

### 3. Look up an address

```bash
curl -H "Authorization: Bearer $TOKEN" \
  localhost:8080/wallets/0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045
```

That returns the ETH balance and saves a snapshot against your account. This is the only route that burns a search from your daily quota. Everything below is free once you have looked the address up.

### 4. Everything else

```bash
# who this address actually deals with
curl -H "Authorization: Bearer $TOKEN" localhost:8080/wallets/0xADDR/overview

# full transfer history, in and out
curl -H "Authorization: Bearer $TOKEN" localhost:8080/wallets/0xADDR/transactions

# ERC-20 holdings with USD values
curl -H "Authorization: Bearer $TOKEN" localhost:8080/wallets/0xADDR/tokens

# nodes and edges, ready to drop into a graph library
curl -H "Authorization: Bearer $TOKEN" localhost:8080/wallets/0xADDR/graph
```

## Endpoints

| Method | Route | What it does |
|---|---|---|
| GET | `/health` | Liveness check. No auth. |
| POST | `/auth/register` | Create an account. Password minimum 8 characters. |
| POST | `/auth/login` | Returns a JWT, valid 24h. |
| GET | `/auth/me` | Your account, role, daily limit, and searches used. |
| GET | `/wallets` | Every address you have looked up. |
| GET | `/wallets/:address` | ETH balance. Costs one search. |
| GET | `/wallets/:address/transactions` | Transfers in and out, newest first, tagged with direction. |
| GET | `/wallets/:address/tokens` | Token holdings and total USD. Add `?refresh=true` to bypass the cache. |
| GET | `/wallets/:address/overview` | Transaction counts, first and last activity, wallet age, top 5 counterparties with labels. |
| GET | `/wallets/:address/graph` | Up to 20 nodes and their edges, weighted by transfer count. |
| GET | `/labels/:address` | Known name and category for an address, if there is one. |

## Quotas

Every account gets a `daily_limit` on the user row. Default is 10. It resets at UTC midnight and only `/wallets/:address` spends from it. Go over and you get a 429.

The tiers are Free 10, Pro 50, Max 150. The limit lives in the database, not in the code, so changing someone's tier is an update to their row.

## How the tracing works

There is no clever indexing here and that is on purpose.

`alchemy_getAssetTransfers` is called twice for an address, once with `fromAddress` and once with `toAddress`, capped at 1000 each. That gives you every transfer touching the address in two HTTP calls. Group those by the other party and you have a counterparty list. Rank by transfer count, cut to the top 20, and that is the graph.

Token balances are cached in Postgres for an hour because the portfolio endpoint is the slow one.

## Layout

```
backend/
  cmd/api/main.go              wiring and routes
  internal/
    platform/database/         connection and migration
    platform/auth/             JWT issue and verify, gin middleware
    user/                      accounts, quotas
    wallet/                    Alchemy client, balances, transfers, graph
    label/                     known address names, seeded on boot
compose.yaml                   Postgres 18
```

Each domain is model, repository, service, handler. Handlers talk HTTP, services hold the logic, repositories talk to the database. Nothing skips a layer.

## Notes

Postgres 18 needs its volume mounted at `/var/lib/postgresql`, not `/var/lib/postgresql/data`. The image refuses to start otherwise. Already handled in `compose.yaml`.

There is no Dockerfile for the API yet. The service block in `compose.yaml` is commented out until there is one.
