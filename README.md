# rate-limiter-go

Rate Limiter (Echo + Redis) — Docker Setup

Kebutuhan
    Docker & Docker Compose
    Port default:
        App: 1234
        Redis: 6379

1. Env
Buat file .env dari .env.example

2. Build & Run dengan Docker
docker-compose build --no-cache
docker-compose up -d

3. Running apps

A. GET http://localhost:1234/ -> Health Check

B. GET http://localhost:1234/api/v1/rate-limit/0101 -> Check Rate Limit

C. PUT http://localhost:1234/api/v1/rate-limit/0101 -> Configure Rate Limit per Client

{
  "max_requests": 5,
  "cycle_duration": 60
}

D. GET http://localhost:1234/api/v1/protected/data -> Protected Endpoint
