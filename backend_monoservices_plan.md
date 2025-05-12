# Mono services - Go Microservices Monorepo with Go Kit

## Project Structure and Initial Setup Scaffold

### Root Folder (Monorepo-Style)

super_services/
├── cmd/
│   ├── gateway/
│   │   └── main.go
│   ├── auth/
│   │   └── main.go
│   ├── user/
│   │   └── main.go
│   ├── notification/
│   │   └── main.go
│   ├── audit/
│   │   └── main.go
│   ├── admin/
│   │   └── main.go
│   ├── billing/
│   │   └── main.go
│   └── search/
│       └── main.go
├── internal/
│   ├── gateway/
│   │   ├── handler/
│   │   └── routes/
│   ├── auth/
│   │   ├── handler/
│   │   ├── service/
│   │   └── model/
│   ├── user/
│   │   ├── handler/
│   │   ├── service/
│   │   └── model/
│   ├── notification/
│   │   ├── handler/
│   │   ├── service/
│   │   └── model/
│   ├── audit/
│   │   ├── handler/
│   │   ├── service/
│   │   └── model/
│   ├── admin/
│   │   ├── handler/
│   │   ├── service/
│   │   └── model/
│   ├── billing/
│   │   ├── handler/
│   │   ├── service/
│   │   └── model/
│   └── search/
│       ├── handler/
│       ├── service/
│       └── model/
├── pkg/
│   ├── config/
│   ├── logger/
│   ├── middleware/
│   └── utils/
├── api/
│   ├── auth/
│   ├── user/
│   ├── notification/
│   │   ├── handler/
│   │   │   ├── sse_handler.go
│   │   │   └── websocket_handler.go
│   │   ├── model/
│   │   └── service/
│   │       ├── dispatcher.go
│   │       └── notifier.go
│   ├── billing/
│   └── search/
├── deployments/
│   ├── gateway/
│   ├── auth/
│   ├── user/
│   ├── notification/
│   ├── audit/
│   ├── admin/
│   ├── billing/
│   ├── search/
│   ├── keycloak/
│   ├── nats/
│   ├── matrix/
│   ├── postgres/
│   ├── mongodb/
│   └── redis/
├── scripts/
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md

### Key Features
- **Microservices Architecture**:
  - Built using Go Kit's transport, endpoint, and service layers.
  - Includes services for authentication, user management, notifications, billing, search, and more.

- **Configuration and Logging**:
  - Configurations are managed using Viper.
  - Logging is handled via Zap for structured and efficient logging.

- **Caching and Messaging**:
  - Redis is used for caching and rate limiting.
  - NATS JetStream is used for scalable service messaging.

- **Database and Storage**:
  - PostgreSQL with primary and read replicas for relational data.
  - MongoDB replica set for document-based storage.
  - Redis with Sentinel for high availability.

- **Federated Messaging**:
  - Matrix Synapse for user-to-user and service-to-service communication.

- **Scalability and Multi-Tenancy**:
  - Multi-tenancy support with tenant-aware query filters and schema-per-tenant in PostgreSQL.
  - Horizontally scalable API Gateway and services.

- **Security**:
  - JWT-based authentication integrated with Keycloak.
  - TLS for secure communication.

- **Observability**:
  - Monitoring and tracing with Prometheus, Grafana, Loki, and Jaeger.

- **Automation**:
  - Makefile commands for building, testing, and managing services.
  - CI/CD pipeline for automated deployments.

---

## Planned Features and Enhancements

### **1. Tools**

#### Makefile Commands
- `build`: Builds each microservice binary.
- `run`: Runs the gateway locally.
- `test`: Executes Go tests for all packages.
- `lint`: Runs static analysis via golangci-lint.
- `up`: Starts all services with Docker Compose.
- `down`: Stops all services.
- `restart`: Restarts Docker services.
- `logs`: Tails logs from all services.

### **2. Infrastructure Setup**

#### containers
- Kubernetes
- containerd
- K3s, a local kubernetes
- K6, a stress test

#### Docker Compose Setup
- **Keycloak**: Identity provider with realm `mono-services`.
- **Gateway**: Public entry point, secured with JWT middleware.
- **Auth**: Middle-layer to interact with Keycloak for login/signup.
- **User**: Profile service for authenticated users.
- **Notification**: Async message delivery (email, in-app).
- **Audit**: Central event logger.
- **Admin**: Backend admin panel.
- **Billing**: Subscription and payment logic.
- **Search**: Search indexing engine.
- **Redis**: High availability with master and replicas.
- **PostgreSQL**: Primary and read replicas with connection pooling.
- **MongoDB**: Replica set for read scalability and failover.
- **NATS**: Messaging broker with JetStream enabled.
- **Matrix**: Federated messaging (Matrix Synapse home server for user/service chat).

#### PostgreSQL Replication Setup
- **postgres-primary**: Main writable PostgreSQL instance.
- **postgres-replica-1** & **postgres-replica-2**: Read-only replicas using streaming replication.
- **pgbouncer**: Connection pooler between services and PostgreSQL (routing read/write).
- Configuration files:
  - `pg_hba.conf`, `postgres.conf`, and init scripts to configure replication roles and WAL.
- Docker volumes for persistent storage.
- Config located under: `deployments/postgres/{primary/, replica1/, replica2/}`.
- Each microservice will use a dedicated **schema** in PostgreSQL for logical separation.
- Per-service roles will be created with access limited to their own schema.
- Shared PostgreSQL instance ensures resource efficiency while keeping logical separation.

#### MongoDB Replica Set Setup
- **mongodb-primary**: Primary MongoDB node for write operations.
- **mongodb-secondary-1** & **mongodb-secondary-2**: Secondary nodes for read scaling and failover.
- **mongosetup** container or init script to run `rs.initiate()` and `rs.add()`.
- Docker volumes for data persistence.
- Config located under: `deployments/mongodb/{primary/, secondary1/, secondary2/}`.

#### Redis High Availability Setup
- **redis-master**: Primary Redis node.
- **redis-replica-1**, **redis-replica-2**: Redis read replicas for scaling.
- **redis-sentinel-1**, **redis-sentinel-2**, **redis-sentinel-3**: Sentinel nodes for monitoring and automatic failover.
- Docker volumes for persistent Redis storage.
- Optional security: AUTH password, TLS.
- Config located under: `deployments/redis/{master/, replica1/, replica2/, sentinel/}`.
- Services using Redis:
  - **auth**: Session caching, blacklisting tokens.
  - **gateway**: Rate limiting, token verification.
  - **user**: Optional profile cache.

#### NATS JetStream Setup (Service Messaging)
- **nats**: Messaging broker with JetStream enabled.
- Use cases:
  - Gateway → Auth/User (emit login/signup).
  - User → Future services (notifications, etc.).
- Config located under: `deployments/nats/`.
- JetStream settings:
  - Persistent file-based storage.
  - Stream configurations (retention, limits).
- Subject naming: `auth.*`, `user.*`, etc.
- Designed to scale with millions of concurrent messages.
- Clustering supported for high availability.

#### Matrix Synapse Setup
- Federated messaging for user-to-user and service-to-service communication.
- Deploy Matrix Synapse as a Docker container:
  - Include in `docker-compose.yml` under `deployments/matrix/`.
  - Configure federation settings for external communication.
- Use PostgreSQL for Synapse storage:
  - Create a dedicated schema in the shared PostgreSQL instance.
- Configure Redis for caching and worker coordination.
- Set up worker processes for scalability:
  - Split Synapse into multiple workers (e.g., event persistence, federation sender).
- Enable TLS for secure communication:
  - Use Let's Encrypt or a custom certificate.
- Integrate with Keycloak for user authentication (OIDC).
- **Future Enhancements**:
  - Add support for bots and automation (e.g., notification bots).
  - Integrate with other services for real-time collaboration.

---

### **3. Core Microservices Development**

#### Gateway Service
- Acts as an API gateway.
- Middleware validates Keycloak-issued JWTs.
- Injects claims (userID, roles) into the request context.
- Routes requests to appropriate services.

#### Auth Service
- Acts as a proxy/middleman to Keycloak.
- Endpoints:
  - `POST /signup`: Forwards user registration to Keycloak.
  - `POST /login`: Forwards to Keycloak token endpoint.
- Maps Keycloak responses to internal domain model.

#### User Service
- Manages user profile data.
- Endpoints:
  - `GET /me`: Fetch profile based on JWT.
  - `PUT /me`: Update profile info (name, email, etc.).
- Accessible only with a valid JWT via the gateway.

#### Notification Service
- Real-time and persistent notifications with user preferences.
- Hybrid delivery system:
  - SSE for browser-based web clients.
  - WebSocket for full-duplex communication when needed.
  - Future: Push notification support (FCM/APNs).
- Internal dispatcher to decide delivery strategy per platform/user/device.
- Notification queueing and retry on failure.
- Persistent logging of all notifications to the database.
- Integration with NATS JetStream for domain event consumption and message buffering.
- Planned support for mobile silent-notification pattern with analytics feedback.
- Supports future Matrix integration.

#### Audit Service
- Tracks sensitive operations.
- Records auth/login events.
- Stores data in a searchable database.

#### Admin Service
- Powers the internal control panel.
- Admin interface for user and service management.

#### Billing Service
- Subscription and invoice management.
- Integrates with payment providers (e.g., Stripe).

#### Search Service
- Indexes and queries user content.
- Uses Elasticsearch.

---

### **4. Enhancements and Scalability**

#### API Gateway Enhancements
- Rate limiting middleware per user/client/IP using Redis and sliding window algorithm.
- Support for global and per-tenant quota enforcement.
- Middleware to extract and inject tenant metadata (e.g., tenant-id from JWT claims).
- Request context enriched with tracing, correlation ID, and tenant info.
- Centralized error transformation layer for clean API responses.
- Metrics for API usage per route, status, and tenant.
- Admin APIs for managing rate limits and tenant policies.

#### Multi-Tenancy Plan
- Tenant ID embedded in all service-level domain models.
- Tenant-aware query filters in all database access layers.
- Optional schema-per-tenant support in PostgreSQL (advanced scaling scenario).
- Service discovery uses tenant namespace to isolate workflows.
- Gateway enforces tenant segregation via auth and route rules.
- Shared services (e.g., search, notification) tag data by tenant.
- Billing and audit log entries associated with tenants.

#### Rate Limiting Strategy
- Implement using Redis backend (via `pkg/middleware/ratelimit.go`).
- Support burst and sustained rate control.
- Limit by:
  - Authenticated user.
  - API key (client).
  - IP address (fallback).
- Provide headers with rate limit metadata (e.g., `X-RateLimit-Remaining`).
- Emit rate limit exceeded events to NATS for audit/billing.

#### Scalability Additions
- **API Gateway**: Horizontally scalable, stateless.
- **Keycloak**: Clustering supported, shared database, sticky or stateless auth.
- **Matrix Synapse**: Scalable via worker processes, federated.
- **Redis**: Sentinel HA; Redis Cluster optional for write scaling.
- **NATS JetStream**: Clustered messaging, concurrent handling.
- **PostgreSQL**: Read replicas, connection pooling.
- **MongoDB**: Replica set for read scalability and failover.
- **Observability (Future)**: Prometheus, Grafana, Loki, Jaeger.

---

### **5. Observability and Monitoring**

#### Observability Tools
- Deploy Prometheus, Grafana, Loki, and Jaeger for monitoring and tracing.
- Add metrics and logs for all services.

#### Scalability Testing
- Perform load testing on API Gateway, Redis, PostgreSQL, and NATS.
- Optimize configurations for horizontal scaling.

---

### **6. Deployment and Maintenance**

#### CI/CD Pipeline
- Automate builds, tests, and deployments using the `Makefile` commands.
- Integrate with a CI/CD tool (e.g., GitHub Actions, Jenkins).

#### Production Deployment
- Deploy services to a production environment.
- Enable TLS for secure communication.

#### Ongoing Maintenance
- Monitor system health and logs.
- Regularly update dependencies and configurations.
