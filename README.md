## PayChain

PayChain is a minimal payment system prototype built with Go, Kafka, and a blockchain. Transfers are accepted via HTTP API, pushed into Kafka for buffering, concurrently consumed into a pending pool, and periodically batched into blocks that are mined using a Proof‑of‑Work (PoW) worker pool. A React TypeScript frontend provides a dashboard for submitting transfers and observing balances, pending transactions, and the blockchain.

### Features

- Transfer API backed by Kafka for decoupling and buffering
- Concurrent Kafka consumption into a thread‑safe transaction pool
- Asynchronous block packing with batch threshold (default 3)
- PoW mining with a CPU‑parallel worker pool and cancelable context
- Thread‑safe blockchain storage and account balance store
- Faucet endpoint for demo balance top‑up
- React 18 + TypeScript UI with Ant Design, Tailwind CSS, Zustand, TanStack Query, Axios

## Demo
[PayChain Dashboard Demo](https://github.com/user-attachments/assets/89e2ccc7-0217-4cd3-85c9-790c1467947d)


## Architecture Overview

- API (Gin):
  - POST /transfer → publish transaction to Kafka
  - GET /balance/:user → read balance
  - GET /blockchain → list all blocks
  - GET /pending → list pending transactions
  - POST /faucet → mint to a user (demo only)
- Kafka (Sarama): Async producer and consumer group
- Pool: Mutex‑guarded slice as the pending transaction pool
- Blockchain:
  - Block and Transaction data structures
  - Chain: append with mutex, read with RWMutex
  - PoW: parallel nonce search using goroutines over disjoint ranges
- Accounts: RWMutex‑protected map with batch pre‑validation and atomic batch apply
- Frontend: Single‑page dashboard (React) for transfer, balance, pending, blockchain

## Technologies

- Backend
  - Go (Gin, Sarama)
  - Concurrency: goroutines, channels, context cancelation, sync.Mutex/RWMutex
  - Blockchain: PoW, block hashing, chain storage
  - Logging: std log wrapper
- Middleware
  - Apache Kafka (Zookeeper for Bitnami image)
- Frontend
  - React 18, TypeScript, Vite
  - Ant Design, Tailwind CSS
  - Zustand (lightweight state)
  - TanStack Query (fetching, cache, polling)
  - Axios (HTTP client)
- Containerization
  - Docker, Docker Compose

## Concurrency Model (Backend)

- Kafka consumers run in background, pushing messages to a mutex‑protected pool
- A dedicated goroutine performs block packing:
  - Only when pool.Size() ≥ batchSize (default 3)
  - Pre‑validate transactions against an account snapshot to form a valid set
  - Mine the candidate block using PoW with N=NumCPU workers
  - After a solution, atomically apply the exact mined set; append block if commit succeeds
  - No mutation of block contents post‑mining (hash remains valid)

## Blockchain & PoW

- Block fields: Index, Timestamp, Transactions, PrevHash, Hash, Nonce
- Hash = SHA‑256 over concatenated fields + transactions string
- PoW difficulty: leading zeros (default 3)
- Worker pool: each goroutine iterates nonce = start + k*workers; on first success, cancel others via context

## Kafka Integration

- Producer: async, JSON‑encodes Transaction to topic `paychain-transactions`
- Consumer: consumer group with range rebalancing; JSON‑decodes and adds to pool; offsets are marked on consume
- docker‑compose enables auto‑topic creation for quickstart

## API Endpoints

- POST /transfer
  - Body: { "from": string, "to": string, "amount": number }
  - Response: { "status": "queued" }
- GET /balance/:user → { user, balance }
- GET /blockchain → Block[]
- GET /pending → Transaction[]
- POST /faucet
  - Body: { "to": string, "amount": number }
  - Response: { status, user, balance }

## Running with Docker Compose

Prerequisites: Docker Desktop with Compose.

1) Build images
```
docker compose build --no-cache
```
2) Start stack
```
docker compose up -d
```
3) Services
- Backend API: http://localhost:18080
- Frontend (optional service): http://localhost:5173

The compose file sets:
- Kafka broker: `kafka:9092` (auto‑create topics enabled)
- Backend env: `KAFKA_BROKERS=kafka:9092`
- Frontend env (container mode): `VITE_API_BASE_URL=http://paychain:8080`

## Quick Test (API)

- Faucet (top up Alice):
```
POST http://localhost:18080/faucet
Content-Type: application/json

{"to":"alice","amount":1000}
```
- Transfer (enqueue; needs 3 to mine with defaults):
```
POST http://localhost:18080/transfer
Content-Type: application/json

{"from":"alice","to":"bob","amount":100}
```
- Inspect:
```
GET http://localhost:18080/pending
GET http://localhost:18080/blockchain
GET http://localhost:18080/balance/alice
GET http://localhost:18080/balance/bob
```

## Notes & Caveats

- This is a toy blockchain for demo/education. No persistence, consensus network, or security hardening
- Faucet is for demo only; do not enable in production
- Balances and chain are in‑memory only (lost on restart)

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Summary

👉 If you found this project helpful, please ⭐ it and share it with others!


