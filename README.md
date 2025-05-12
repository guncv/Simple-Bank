# 🏦 Simple Bank

A modern banking backend service built in Go with high-performance APIs using **gRPC**, **Gin**, **gRPC-Gateway**, and **Redis** for distributed task queues. The project also features secure **authentication**, **RBAC**, and **transactional integrity** using SQL transaction control (commit/rollback).

---

## 🚀 Features

### ✅ Core Functionality

- **User Account Management** (Registration, Login, Update)
- **Banking Operations** (Transfer, Create Accounts, List Transactions)
- **Secure Authentication**
  - Access & Refresh Token (JWT/PASETO)
  - Token revocation and expiry handling
- **RBAC** (Role-Based Access Control)

### 🔁 gRPC + REST

- **gRPC** for fast, type-safe RPCs
- **gRPC-Gateway** provides RESTful HTTP APIs
- **Gin** as a lightweight optional HTTP server alternative

### 📧 Asynchronous Email Notification

- Redis-backed task queue for background jobs
- Notification worker sends **verification emails** using SMTP
- Email payload defined via Protocol Buffers (`smtp.proto`)

### 🔐 Security

- Secure password hashing using bcrypt
- Access token stored in headers; refresh token managed securely
- Token validation middleware with full RBAC logic

### 🧪 Testing

- ✅ **Black-box testing** via gRPC/Gin endpoints
- ✅ **Unit testing** with `gomock` and `testify`
- ✅ Database transaction logic is tested (e.g., commit/rollback behavior in `TransferTx`)

### 💾 Database Transactions

- SQL transaction pattern for safe, consistent data changes
- Supports automatic **rollback on error**
- Example: `TransferTx` performs fund transfers atomically with:
  - CreateTransfer
  - CreateEntries
  - UpdateAccountBalance

---

## 📦 Tech Stack

| Tool             | Usage                                  |
|------------------|-----------------------------------------|
| Go               | Main language                          |
| gRPC             | High-performance APIs                  |
| gRPC-Gateway     | REST API via protobuf auto-translation |
| Gin              | Optional HTTP server                   |
| Redis            | Message broker / task queue            |
| PostgreSQL       | Database                               |
| gomock, testify  | Testing and assertions                 |
| JWT / PASETO     | Token-based authentication             |
| Bcrypt           | Password hashing                       |

---

## 📂 Project Structure (Simplified)

```

.
├── db/                  # SQLC-generated queries and models
├── gapi/                # gRPC server implementation
├── http/                # Gin server and handlers
├── pb/                  # Protocol Buffers (compiled)
├── proto/               # .proto definitions
├── token/               # Token generation and validation
├── worker/              # Redis task queue + email sender
├── util/                # Utility functions (e.g., Random, Config)
├── main.go              # Entry point
├── go.mod               # Dependencies
└── README.md            # This file

````

---

## 🧠 Getting Started

### 🛠 Prerequisites

- Go ≥ 1.19
- PostgreSQL
- Redis (for task queue)
- [buf](https://buf.build/) (for proto generation)
- [grpcurl](https://github.com/fullstorydev/grpcurl) or Postman (for testing APIs)

### 🔧 Run Locally

```bash
# Run PostgreSQL and Redis (via Docker or locally)
make postgres
make redis

# Setup database
make migrateup

# Start servers
make server         # gRPC server
make http-server    # RESTful HTTP server via gRPC-Gateway
````

---

## 🧪 Run Tests

```bash
make test          # Run all unit & API tests
make cover         # View test coverage
```

---

## 🛡 Security

* All endpoints are protected by middleware
* Authorization logic validates token claims and role
* Refresh tokens enable seamless reauthentication

---

## 📬 Example Use Cases

* Transfer money between accounts
* Register a user and receive a verification email
* Login and get access + refresh tokens
* Access endpoints based on role (e.g., Admin vs User)

---
