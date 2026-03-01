# True Recommendation Assignment

This README.md is ai-generated

---

A content recommendation engine built with **Go (Fiber)**, **PostgreSQL**, and **Redis**. The service delivers personalized content recommendations based on user watch history, genre preferences, content popularity, and recency — filtered by geographic availability and subscription tier.

---

## Table of Contents

- [Setup Instructions](#setup-instructions)
- [Architecture Overview](#architecture-overview)
- [Design Decisions](#design-decisions)
- [Performance Results](#performance-results)
- [Trade-offs and Future Improvements](#trade-offs-and-future-improvements)

---

## Setup Instructions

### Prerequisites and Dependencies

| Dependency | Version | Purpose |
|---|---|---|
| Go | 1.26+ | Application runtime |
| Docker & Docker Compose | v3.8+ | Container orchestration |
| PostgreSQL | 16 (Alpine) | Primary data store |
| Redis | 7 (Alpine) | Caching layer |
| k6 | Latest | Load testing |

**Go Module Dependencies (key):**

- `github.com/gofiber/fiber/v2` — HTTP framework
- `gorm.io/gorm` + `gorm.io/driver/postgres` — ORM & PostgreSQL driver
- `github.com/go-redis/redis/v7` — Redis client
- `github.com/jackc/pgx/v5` — Underlying PostgreSQL driver

### Step-by-step Installation Guide

**1. Clone the repository**

```bash
git clone https://github.com/thisusami/true-recommendation-assignment.git
cd true-recommendation-assignment
```

**2. Start all services with Docker Compose**

```bash
docker-compose up --build -d
```

This single command will:
- Build the Go application using a multi-stage Dockerfile
- Start PostgreSQL 16 with health checks
- Start Redis 7
- Automatically run migrations and seeding on first launch

**3. (Alternative) Run locally without Docker**

```bash
# Start PostgreSQL and Redis manually, then:
export POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/recommendation?sslmode=disable"
export REDIS_ADDR="localhost:6379"
export REDIS_PASSWORD=""
export PORT="8080"

go mod download
go run ./src
```

### Commands to Run Migrations and Seeding

Migrations and seeding are handled **automatically** by PostgreSQL's `docker-entrypoint-initdb.d` mechanism. The SQL scripts in `script/` are executed in alphabetical order on first container initialization:

| Script | Purpose |
|---|---|
| `script/01_schema.sql` | Creates `users`, `content`, and `user_watch_history` tables with indexes |
| `script/02_seed.sql` | Seeds 25 users, 50 content items (6 genres), and 220+ watch history records |

To re-seed manually:

```bash
docker exec -i recommendation-postgres psql -U postgres -d recommendation < script/01_schema.sql
docker exec -i recommendation-postgres psql -U postgres -d recommendation < script/02_seed.sql
```

### How to Start the Application

```bash
# Docker (recommended)
docker-compose up --build

# Verify health
curl http://localhost:8080/health
# => {"status":"ok"}
```

**API Endpoints:**

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Health check |
| `GET` | `/users/:user_id/recommendations?limit=10` | Get personalized recommendations for a user (limit: 1–50) |
| `GET` | `/recommendations/batch?page=1&limit=20` | Batch recommendations for all users (paginated) |

---

## Architecture Overview

### High-Level System Design

```
┌─────────┐       ┌──────────────────────────────────────────────────────┐
│  Client  │──────▶│                  Fiber HTTP Server                   │
└─────────┘       │                    (port 8080)                       │
                  └───────────┬──────────────────────────────────────────┘
                              │
                  ┌───────────▼───────────────────────┐
                  │          Handler Layer             │
                  │  • Route registration              │
                  │  • Request validation              │
                  │  • Response formatting              │
                  └───────────┬───────────────────────┘
                              │
                  ┌───────────▼───────────────────────┐
                  │         Service Layer              │
                  │  • Cache check (Redis)             │
                  │  • Orchestrates recommendation     │
                  │  • Batch concurrency (goroutines)  │
                  └──┬────────────────┬───────────────┘
                     │                │
          ┌──────────▼──┐   ┌────────▼────────────┐
          │  Cache Repo  │   │  Model Client (ML    │
          │  (Redis)     │   │   Simulation)        │
          │  • Get/Set   │   │  • Genre frequency   │
          │  • 10m TTL   │   │  • Scoring algorithm  │
          │              │   │  • Candidate ranking  │
          └──────────────┘   └────────┬─────────────┘
                                      │
                          ┌───────────▼───────────────┐
                          │   Repository Layer         │
                          │  (PostgreSQL via GORM)     │
                          │  • GetUserById             │
                          │  • GetWatchHistory         │
                          │  • GetCandidateContents    │
                          └───────────┬───────────────┘
                                      │
                          ┌───────────▼───────────────┐
                          │      PostgreSQL 16         │
                          │  • users table             │
                          │  • content table           │
                          │  • user_watch_history      │
                          └───────────────────────────┘
```

### Explanation of Each Architectural Layer

| Layer | Package | Responsibility |
|---|---|---|
| **Handler** | `src/handler/` | Receives HTTP requests, validates input parameters (user_id, limit, page), delegates to the service layer, and returns JSON responses with proper status codes. |
| **Service** | `src/services/` | Orchestrates the recommendation workflow: checks Redis cache first, fetches user data, calls the model client for scoring, caches the result, and assembles the response. For batch endpoints, manages concurrent goroutines with `sync.WaitGroup`. |
| **Model Client** | `src/model-client/` | **ML simulation layer.** Simulates an external machine learning inference service. Retrieves watch history, computes genre frequency distributions, fetches candidate content, applies the multi-factor scoring algorithm, and returns ranked results. Introduces artificial latency (30–50ms) and a 0.015% random failure rate to realistically mimic production ML service behavior. |
| **Repository** | `src/repositories/` | Data access layer using GORM. Provides typed queries for users, watch history (with content joins), and candidate content filtering (country + subscription + not-yet-watched). |
| **Cache Repository** | `src/repositories/` | Thin wrapper around Redis `GET`/`SET` operations with configurable TTL. |
| **Database Init** | `src/db/` | Connection pool configuration for PostgreSQL (25 max connections, 5m lifetime) and Redis (pool size 20, 5 idle). |
| **Models** | `src/models/` | Domain structs: `User`, `Content`, `WatchHistory`, `RecommendationItem`, response envelopes, and error definitions. |
| **Utilities** | `src/util/` | Structured JSON logging (INBOUND/OUTBOUND/REQUEST/RESPONSE), JSON marshaling helpers, and input validators. |

### Data Flow Through the System

**Single User Recommendation (`GET /users/:user_id/recommendations`):**

1. **Handler** validates `user_id` (non-empty) and `limit` (1–50)
2. **Service** constructs cache key `rec:user:{id}:limit:{n}` and checks Redis
3. **Cache HIT** → deserialize and return immediately with `cache_hit: true`
4. **Cache MISS** → proceed:
   - Fetch user record from PostgreSQL
   - Return 404 if user not found
   - Call **Model Client** (ML simulation) to generate recommendations
5. **Model Client** (ML simulation) pipeline:
   - Query last 50 watch history entries (joined with content) for genre analysis
   - Compute genre frequency distribution (normalized to percentage)
   - Query up to 100 candidate contents filtered by country, subscription tier, and not-yet-watched
   - Score each candidate using the multi-factor algorithm
   - Sort by score descending, return top N
6. **Service** caches the result in Redis (TTL: 10 minutes) and returns response

**Batch Recommendations (`GET /recommendations/batch`):**

1. **Handler** validates `page` (≥1) and `limit` (1–100)
2. **Service** fetches paginated user list from PostgreSQL
3. Spawns one **goroutine per user** with `sync.WaitGroup`
4. Each goroutine independently calls **Model Client** (ML simulation) for scoring
5. Results collected via buffered channel, aggregated with success/failure counts
6. Returns batch response with per-user results and timing summary

### How the Recommendation Model Integrates with Database Queries

The Model Client sits between the service layer and the repository, executing two sequential database queries per recommendation request:

1. **Watch History Query** — `GetWatchHistoryJoinContentByUserId`: Joins `user_watch_history` with `content` to retrieve the user's last 50 watched items with genre information. This feeds the genre frequency computation.

2. **Candidate Content Query** — `GetCandidateContents`: Uses a subquery to exclude already-watched content, filters by `available_countries @> ARRAY[?]` and `available_subscription @> ARRAY[?]` using GIN indexes, and orders by `popularity_score DESC` with a limit of 100 candidates.

The model then applies the scoring algorithm entirely in-memory on the filtered candidates.

---

## Design Decisions

### Caching Strategy and TTL Rationale

- **Cache Key Pattern:** `rec:user:{user_id}:limit:{limit}` — uniquely identifies each recommendation set per user and requested limit.
- **TTL: 10 minutes** — Balances freshness with performance. Recommendations don't change frequently (based on watch history and static content catalog), so a 10-minute window is acceptable. This prevents stale data from persisting too long while still absorbing repeated requests from the same user.
- **Cache-Aside Pattern:** The service checks cache first; on miss, it computes results and stores them. This avoids unnecessary computation on repeated requests.
- **Redis Pool:** 20 connections with 5 idle minimum ensures low-latency cache operations under concurrent load.

### Concurrency Control Approach

- **Semaphore-Based Worker Pool:** The batch endpoint uses a `semaphore.Weighted(50)` to limit concurrent goroutines to 50 at a time, preventing connection pool exhaustion and unbounded resource consumption. Each goroutine acquires a semaphore slot before executing and releases it on completion.
- **Batch Endpoint:** Uses `sync.WaitGroup` + buffered channel for fan-out/fan-in. Each user's recommendations are computed in a separate goroutine, enabling parallel database queries and scoring.
- **Channel Buffer:** Sized to `len(users)` to prevent goroutine blocking.
- **PostgreSQL Connection Pool:** 25 max open connections with 25 idle connections. `ConnMaxLifetime` of 5 minutes prevents stale connections. This caps the concurrent database load regardless of how many goroutines are spawned.
- **Simulated ML Inference Latency:** The model client is an ML simulation that introduces a random 30–50ms artificial latency and a 0.015% random failure rate, mimicking the behavior of a real external ML inference service. This makes the concurrency design realistic for production scenarios where ML model calls are the primary bottleneck.

### Error Handling Philosophy

- **Typed Error Constants:** Pre-defined error objects (`InternalServerError`, `BadRequestError`, `NotFoundError`, `ServiceUnavailableError`) with `.Set()` method for adding dynamic error details while maintaining consistent response structure.
- **Fail-Fast Validation:** Input validation happens at the handler layer before any database calls. Invalid `user_id` or out-of-range `limit`/`page` returns 400 immediately.
- **Graceful Degradation in Batch:** Individual user failures in batch processing don't fail the entire batch. Each goroutine reports its success/failure status independently, and the response includes a summary with `success_count` and `failed_count`.
- **Simulated Failure Rate:** The model client has a 0.015% random failure rate to test error resilience.
- **Structured Logging:** All operations (INBOUND, OUTBOUND, REQUEST, RESPONSE, EXCEPTION) are logged as JSON with timestamps and response times for observability.

### Database Indexing Strategy

| Index | Type | Purpose |
|---|---|---|
| `idx_users_country` | B-tree | Fast user lookup by country |
| `idx_users_subscription` | B-tree | Fast user lookup by subscription tier |
| `idx_content_genre` | B-tree | Genre-based content filtering |
| `idx_content_popularity` | B-tree (DESC) | Ordered retrieval of popular content |
| `idx_content_countries` | GIN (array) | `@>` array containment queries for country availability |
| `idx_content_subscription` | GIN (array) | `@>` array containment queries for subscription availability |
| `idx_watch_history_user` | B-tree | User-specific watch history lookup |
| `idx_watch_history_content` | B-tree | Content-specific watch history lookup |
| `idx_watch_history_composite` | B-tree (user_id, watched_at DESC) | Optimized query for recent watch history per user |

**Rationale:** The GIN indexes on array columns (`available_countries`, `available_subscription`) are critical for the candidate content query, which uses PostgreSQL's `@>` array containment operator. The composite index on watch history optimizes the "most recent N items" query pattern used by the model client.

### Scoring Algorithm Rationale and Weight Choices

The scoring formula is:

$$\text{score} = 0.4 \times \text{popularity} + 0.35 \times \text{genre\_boost} + 0.15 \times \text{recency} + 0.1 \times \text{noise}$$

| Component | Weight | Rationale |
|---|---|---|
| **Popularity** (`popularity_score`) | 40% | Content with historically high engagement is a strong baseline signal. Weighted highest as it represents collective user preference. |
| **Genre Boost** (user's genre frequency) | 35% | Personalizes recommendations based on individual watch patterns. Uses normalized genre frequency from the user's watch history. Falls back to 0.1 for unseen genres to allow discovery. |
| **Recency Factor** ($\frac{1}{1 + \text{days}/365}$) | 15% | Newer content gets a mild boost. The decay function ensures older content isn't penalized too harshly — a 1-year-old item still retains ~50% recency score. |
| **Random Noise** (±0.05 × 0.1) | 10% | Introduces controlled randomness to prevent recommendation staleness and encourage exploration/serendipity. Small enough to not override the primary signals. |

---

## Performance Results

### k6 Test Results

**Test Configuration (Two-Phase):**

| Phase | Scenario | Time Window | Description |
|---|---|---|---|
| **1 — Load Test** | `defaultFunction` | `0s – 2m` | Ramps 0→50→100→0 VUs with random user IDs. Validates functional correctness and warms the Redis cache. |
| **2 — Cache Effectiveness** | `cacheHit` | `2m5s – 4m5s` | Ramps 0→50→100→0 VUs against the same user pool. Measures cache hit ratio using custom `cache_hits` / `cache_misses` counters. |

- Thresholds: p(95) < 500ms, failure rate < 1%, `cache_hits` > 0, `cache_misses` > 0

**Results Summary:**

| Metric | Value |
|---|---|
| Total Requests | 65,422 |
| Throughput | ~545 req/s |
| Avg Response Time | 2.69ms |
| Median Response Time | 1.62ms |
| p(90) Response Time | 5.60ms |
| p(95) Response Time | 7.58ms |
| Max Response Time | 85.62ms |
| Error Rate | 0.00% |
| Checks Passed | 100% (130,844 / 130,844) |
| Data Received | 71 MB (594 kB/s) |
| Data Sent | 6.7 MB (56 kB/s) |

**Threshold Results:**

| Threshold | Target | Actual | Status |
|---|---|---|---|
| `http_req_duration p(95)` | < 500ms | 7.58ms | ✅ PASS |
| `http_req_failed rate` | < 1% | 0.00% | ✅ PASS |

### Identified Bottlenecks and Limiting Factors

1. **ML Simulation Latency (30–50ms):** The model client simulates ML inference with an artificial 30–50ms delay, which is the dominant factor in iteration duration (~103ms avg including the 100ms sleep between k6 iterations). In production with a real ML service, network round-trip and inference time would similarly be the primary bottleneck.
2. **Database Connection Pool Saturation:** With 25 max connections and 100 concurrent VUs, the pool could become a bottleneck under heavier load. Each recommendation requires 2 sequential DB queries (watch history + candidates).
3. **JSON Serialization Overhead:** The `MaptoStruct` utility uses double marshal/unmarshal for type conversion, adding minor CPU overhead per request.

### Cache Hit Rate Analysis

- Cache key pattern: `rec:user:{id}:limit:{limit}` with 10-minute TTL.
- With 20 random user IDs in the k6 test and 65,422 total requests, each user ID is hit ~3,271 times on average.
- After the initial cold-start miss per unique key, subsequent requests within the 10-minute TTL are served from cache.
- **Estimated cache hit rate: ~99.9%+** (only ~20 cold misses out of 65,422 requests).
- The very low average response time (2.69ms) confirms that the vast majority of requests are served from Redis cache, as a cache miss would incur ~50–80ms (DB queries + model latency).

---

## Trade-offs and Future Improvements

### Known Limitations

1. **No Cache Invalidation:** When a user watches new content, their cached recommendations are not invalidated. Stale recommendations persist until TTL expires (10 minutes).
2. **No Rate Limiting:** The API has no rate limiting or throttling, making it vulnerable to abuse under production conditions.
3. **Batch Endpoint Not Cached:** Unlike single-user recommendations, batch results are not cached. Repeated batch calls recompute everything.
4. **Bounded Goroutines in Batch:** A `semaphore.Weighted(50)` limits concurrent goroutines to 50 in batch processing. Requests exceeding this limit will block until a slot is available.
5. **Single-Instance Architecture:** No horizontal scaling, load balancing, or service discovery. Redis and PostgreSQL are single-node.
6. **ML Simulation (Not Production ML):** The model client is an in-process simulation of an external ML inference service. The scoring algorithm, random noise, artificial latency (30–50ms), and failure rate (0.015%) are designed to mimic real ML service behavior but would be replaced by actual model serving infrastructure (e.g., TensorFlow Serving, TorchServe) in production.
7. **No Authentication/Authorization:** All endpoints are publicly accessible without any auth mechanism.

### Scalability Considerations

- **Horizontal Scaling:** The stateless Fiber server can be horizontally scaled behind a load balancer. Redis serves as the shared cache layer.
- **Database Read Replicas:** Read-heavy workload (recommendations are read-only) is well-suited for PostgreSQL read replicas.
- **Connection Pooling:** Current pool size (25) is adequate for the seed dataset but would need tuning (e.g., PgBouncer) for production scale.
- **Content Catalog Growth:** The candidate query fetches up to 100 items. As the content catalog grows, the GIN indexes on array columns will maintain query performance, but the candidate pool selection strategy may need refinement.

### Proposed Enhancements (If Time Permitted)

1. **Cache Invalidation via Events:** Publish watch events to Redis Pub/Sub or a message queue; invalidate affected user caches on new watch activity.
2. **~~Worker Pool for Batch:~~** ✅ Implemented — semaphore-based concurrency control (50 concurrent workers) prevents connection pool exhaustion.
3. **Collaborative Filtering:** Extend the scoring model with user-to-user similarity (users who watched X also watched Y) for better personalization.
4. **A/B Testing Framework:** Support multiple scoring algorithms simultaneously to measure recommendation quality.
5. **Pagination for Recommendations:** Support cursor-based pagination for recommendation results instead of fixed limit.
6. **Prometheus Metrics:** Export request latency histograms, cache hit rates, and DB query durations for production monitoring.
7. **Graceful Shutdown:** Handle `SIGTERM` to drain in-flight requests before shutting down.
8. **Context Propagation & Timeouts:** Add `context.Context` with timeouts to all database and cache operations to prevent hanging requests.