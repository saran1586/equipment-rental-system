# Equipment Rental System

A backend for renting out equipment, split into three **Go** microservices that talk to each
other over REST. Built with **Gin** and **GORM** on **PostgreSQL**.

## Services

| Service | Port | What it does |
|---|---|---|
| `rental-service` | 8080 | Creates rentals and rental items, handles returns |
| `customer-service` | 8081 | Customers, KYC verification, deposits, outstanding dues and payments |
| `equipment-service` | 8082 | Equipment types and units, with their status and condition |

`rental-service` calls the other two services (see `internal/client`) to check the customer and
the equipment before a rental is created.

## Main endpoints

**Customer** `/customer`
- `POST /customer`, `GET /customer`, `GET /customer/:id`, `PUT /customer/:id`, `DELETE /customer/:id`
- `PATCH /customer/:id/VerifyKYC`: mark KYC as verified
- `PATCH /customer/:id/UpdateDeposit`, `PATCH /customer/:id/deposit`: add to or deduct from the deposit
- `PATCH /customer/:id/outstanding`, `PATCH /customer/:id/payment`: add an outstanding due, pay it off

**Equipment** `/equipment-unit`
- `GET /equipment-unit`, `GET /equipment-unit/:id`, `PUT /equipment-unit/:id`, `DELETE /equipment-unit/:id`
- `PATCH /equipment-unit/:id/status`, `PATCH /equipment-unit/:id/condition`

**Rental** `/rental`, `/rental-item`
- `POST /rental`, `GET /rental`, `GET /rental/:id`, `DELETE /rental/:id`
- `PATCH /rental/:id/return`: return the equipment
- `POST /rental-item`, `GET /rental-item/:rental_id`, `PUT /rental-item/:id`, `DELETE /rental-item/:id`

## Structure

Each service uses the same layers:

```
cmd/server/         entry point and routes
internal/handler/   HTTP handlers
internal/service/   business rules
internal/repository/ database access (GORM)
internal/model/     entities
internal/config/    database connection from environment variables
```

## Run

Each service reads its settings from a `.env` file. Copy the example and fill it in:

```sh
cd customer-service
cp .env.example .env     # set DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, SERVER_PORT
go run ./cmd/server
```

Start all three services (each in its own terminal) to use the rental flow.
