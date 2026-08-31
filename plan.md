# Apollo

OSINT + Blockchain Intelligence platform. Accepts digital identifiers (wallets, ENS, usernames, emails, domains, IPs) and builds an evidence-backed relationship graph using **publicly available** data.

Every relationship includes: source provenance, confidence score, evidence, timeline.

**Principles:** Evidence over assumptions. Every relationship explainable. Confidence transparent. Public/licensed data only. No black-box identity attribution.

## Tech Stack

- **Frontend:** Next.js, TypeScript, TailwindCSS, React Query, React Flow (graph), ECharts (charts/timeline)
- **Backend:** Go, Gin, JWT auth, REST API, background workers
- **Database:** PostgreSQL (Redis / Neo4j later, only if needed)
- **Infra:** Docker, GitHub Actions, auto-deploy, logging, rate limiting

## Business Model

Tiered usage. Same code path for every tier — one per-user counter; the tier only changes the number it compares against.

| Tier | Searches / day | Includes |
|------|---------------|----------|
| Free | 10 | Wallet intelligence only. No OSINT, no exports. |
| Pro | 50 | + OSINT modules, reports, investigations. |
| Max | 150 | + higher rate limits, API keys, saved monitoring. |

- Metering piggybacks on rate limiting (Phase 10). One column: `searches_used_today` per user, reset at midnight.
- `daily_limit` lives on the user/tier row, never hardcoded — changing a tier = updating a value, not code.
- Expensive features (OSINT, PDF reports) gate behind paid tiers so free users don't run up external-API costs.
- Billing = Stripe Checkout + one webhook + a `subscriptions` table. Build only after core works.
- Legal line stays sacred: public data only, evidence-backed hypotheses, no guaranteed attribution.

---

## Phase 0 — Foundation

- **Auth:** registration, login, JWT, refresh tokens
- **Roles:** User, Admin
- **Core API:** health, version, search endpoints
- **DB tables:** Users, Investigations, Searches, Entities, Relationships, Evidence, Labels

## Phase 1 — Identifier Detection

Auto-detect input type: Ethereum / Solana / Bitcoin wallet, ENS, email, username, domain, IP.
Route to the matching pipeline (e.g. Ethereum wallet → wallet intelligence pipeline).

## Phase 2 — Blockchain Intelligence

- **Ethereum:** balance, transactions, token holdings, NFTs, wallet age, first/last tx, gas spent, most-used contracts, contract detection
- **Solana:** balance, transactions, tokens, NFTs, wallet age
- **Bitcoin:** balance, transactions, wallet age
- **Analytics:** activity timeline, largest txs, in vs out, avg tx size, daily/weekly/monthly activity, wallet lifespan

## Phase 3 — Public Labels

Known labels (Binance, Coinbase, Kraken, OpenSea, Uniswap, bridges, smart contracts, known services). Each label carries source, confidence, date added.

## Phase 4 — OSINT

- **Username:** GitHub, Reddit, X, public websites
- **Email:** public info only — Gravatar, public profiles, licensed breach sources
- **Domains:** WHOIS, DNS, reverse DNS, SSL certs, historical DNS
- **IP:** ASN, ISP, hosting provider, country, reverse DNS

## Phase 5 — Identity Graph

Everything becomes nodes: wallet, ENS, username, email, website, IP, domain — linked by relationships. Each relationship carries source, confidence, evidence, date.

## Phase 6 — Confidence Engine

Every relationship gets an explainable score. Example: wallet↔ENS +40, matching username +20, matching avatar +10, public website reference +30 → 100. No black-box scoring.

## Phase 7 — Visualization

Interactive graph, timeline, table, evidence panel. Click node → expand → view evidence → jump to related nodes.

## Phase 8 — Investigations

Create investigations, save searches, notes, bookmarks, entity collections.

## Phase 9 — Reports

Export PDF / HTML / JSON. Include timeline, evidence, confidence, relationships, sources.

## Phase 10 — Security

JWT, RBAC, AES encryption (sensitive data), audit logs, rate limiting, API keys.

---

## Future Features

AI investigation summaries, watchlists, background monitoring, alerts, saved filters, team collaboration, plugin system, graph diffing.

## 🚫 Do Not Build (until explicitly approved)

Require privileged access, legal review, or make claims unsupported by public data:

- **Real identity attribution** — blockchain data alone can't identify people.
- **Exchange KYC integration** — needs exchange cooperation or legal authority.
- **Government-level attribution** — needs subpoenas, banking, telecom.
- **Telecom subscriber data** — not public.
- **Banking records** — private and regulated.
- **Commercial intelligence feeds** — need paid partnerships/licensing.
- **Guaranteed identity matching** — Apollo presents evidence-backed hypotheses, not certainty.

---

## Architecture

Next.js frontend → Go API → PostgreSQL, with background workers calling external services (blockchain RPCs, DNS, WHOIS, ENS, public OSINT APIs).

## Development Order

1. Project setup, auth, PostgreSQL, REST API
2. Ethereum support: wallet search, transactions, tokens
3. Solana support
4. Bitcoin support
5. Labels + wallet analytics
6. Entity graph + relationship engine
7. Graph visualization + timeline
8. OSINT modules
9. Confidence engine
10. Reports + investigations
11. Security hardening
12. Future features

## Long-Term Vision

A transparent intelligence platform where every relationship is backed by verifiable evidence, every confidence score is explainable, and investigators explore blockchain activity and public OSINT through an intuitive graph — without hidden or unverifiable claims.
