# Apollo — MVP

The smallest build that proves the core idea: **raw wallet addresses become readable, connected intelligence.** Ethereum only. Everything else is post-MVP.

## Goal

A user pastes an Ethereum address and sees: what it holds, what it does, who it talks to (with known names), on one graph screen.

## In Scope

- **Auth:** registration, login, JWT + refresh. Roles: User, Admin.
- **Search:** paste an Ethereum address → detect it's ETH → run wallet pipeline. (Just ETH detection for now — no multi-chain routing.)
- **Wallet intelligence (Ethereum):** balance, transactions, token holdings, wallet age, first/last tx, most-used contracts.
- **Public labels:** map known addresses (Binance, Coinbase, OpenSea, Uniswap, bridges, contracts). Each label: source + date. This is the "aha" — it makes raw addresses human-readable.
- **Basic graph (one hop):** wallet → its labeled counterparties, rendered in React Flow. Click a node → see basic info.
- **Usage limit:** free tier, 10 searches/day. One per-user counter, reset at midnight. (`daily_limit` on the user row — paid tiers come later without a rewrite.)

## Out of Scope (post-MVP)

Deliberately cut — none are needed to prove the core:

- Solana, Bitcoin (same pattern repeated, zero new learning)
- OSINT (username / email / domain / IP)
- Confidence engine (scoring)
- Multi-hop graph, timeline view, table view, evidence panel
- Investigations, notes, bookmarks, collections
- Reports (PDF / HTML / JSON)
- Paid tiers + Stripe billing
- ECharts / analytics charts

## Tech (MVP subset)

- **Frontend:** Next.js, TypeScript, Tailwind, React Query, React Flow. *(No ECharts yet.)*
- **Backend:** Go + Gin, JWT, REST.
- **DB:** PostgreSQL. Tables: `users`, `searches`, `entities`, `relationships`, `labels`.
- **Infra:** Docker for local dev. *(GitHub Actions / auto-deploy when there's something to deploy.)*

## The One Decision That Matters

Pick your Ethereum data source before writing Stage 2 — it drives cost, rate limits, and coverage more than any code choice:

- **RPC / wallet data:** Alchemy or Etherscan API (both have free tiers).
- **Labels dataset:** where do known-address labels come from? Options: a public dataset, Etherscan's labels, or a small hand-seeded list to start. Decide before building labels.

## Build Order

1. Project setup, auth, Postgres, REST skeleton
2. Ethereum wallet pipeline (balance, txs, tokens, age, contracts) — behind the search endpoint
3. Public labels — seed a small list, join against wallet counterparties
4. Basic one-hop graph in React Flow
5. Free-tier search counter

Ship after 5. Then reassess against the full plan.
