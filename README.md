# Book Lending API

This repository contains a simplified RESTful API for a book‑lending system.  The
service is written in Go using the Gin web framework, persists data in MySQL
via GORM and follows a clean architecture with clearly defined boundaries
between handlers, use cases and repositories.  Users can register and log in
to obtain a JWT, browse the catalogue, borrow books and return them.  A
rolling seven‑day rate limit on borrow operations prevents abuse and a global
rate limiter protects the API from excessive traffic.

## Features

* **User authentication** – register and log in with email/password to
  receive a JWT.
* **Book management** – create, read, update and delete books with
  pagination support.
* **Borrow/return** – authenticated users can borrow and return books.  A
  user may borrow at most five books in any rolling seven‑day window.
* **Borrowing history** – view paginated borrowing history and active
  borrowings.
* **Error handling** – consistent error responses with appropriate HTTP
  status codes.
* **Rate limiting** – each client IP is limited to 100 requests per minute
  with a burst of 200.  Borrowing is further limited to five per user per
  week.
* **Clean architecture** – the code is organised into `internal/{domain,
  repository, usecase, handler, middleware}` layers plus `pkg` for
  shared utilities.
* **Docker** – a `Dockerfile` and `docker‑compose.yml` make it easy to run
  the API and MySQL together without a local development environment.
* **Database migrations** – SQL migration scripts for creating the users,
  books and lending_records tables.
* **OpenAPI specification** – `docs/swagger.yml` documents the API
  contract in machine readable form.

This version omits development‑only tooling such as make targets, helper
scripts and sample data to keep the project minimal.

## Quick Start with Docker

1. Ensure [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
   are installed.
2. Clone this repository and change into its directory.
3. Start the API and MySQL:

   ```bash
   docker-compose up --build
   ```

   This will build the application image, create a MySQL instance and wait
   until the database is ready before starting the API.  The API will be
   available at `http://localhost:8080`.

4. Use the API:
   * Register: `POST /api/v1/auth/register` with `{ "email": "alice@example.com", "password": "secret123" }`.
   * Log in: `POST /api/v1/auth/login` and copy the returned `token`.
   * List books: `GET /api/v1/books?page=1&limit=10`.
   * Create a book: `POST /api/v1/books` with a JSON body and set
     `Authorization: Bearer <token>`.
   * Borrow a book: `POST /api/v1/lending/borrow` with `{ "book_id": 1 }`.
   * Return a book: `PUT /api/v1/lending/return/1`.

## Running Locally without Docker

You can also run the service directly if you have Go and MySQL installed:

```bash
# Set environment variables for the database connection and JWT secret
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=yourpassword
export DB_NAME=book_lending
export JWT_SECRET=yoursecretkey

# Create the database and run the migrations using your preferred tool
mysql -u$DB_USER -p$DB_PASSWORD -e "CREATE DATABASE IF NOT EXISTS $DB_NAME;"

# Install dependencies and build/run
go mod tidy
go run cmd/server/main.go
```

## API Endpoints

Endpoint | Method | Description | Auth
---|---|---|---
/api/v1/auth/register | POST | Register a new user | No
/api/v1/auth/login | POST | Authenticate and receive a JWT | No
/api/v1/books | GET | List books (paginated) | No
/api/v1/books | POST | Create a new book | Yes
/api/v1/books/{id} | GET | Get a book by ID | No
/api/v1/books/{id} | PUT | Update a book | Yes
/api/v1/books/{id} | DELETE | Delete a book | Yes
/api/v1/lending/borrow | POST | Borrow a book | Yes
/api/v1/lending/return/{id} | PUT | Return a book | Yes
/api/v1/lending/history | GET | Get borrowing history | Yes
/api/v1/lending/active | GET | Get active borrowings | Yes
/health | GET | Health check | No

See `docs/swagger.yml` for detailed request/response structures.

## Branching and Commit Strategy

The repository follows a simple feature‑branch workflow.  Work begins
from `main` by creating a feature branch named `feature/<feature-name>`.  Each
commit uses [Conventional Commit](https://www.conventionalcommits.org/) notation
to convey intent, for example:

* `feat(auth): implement user registration and login`
* `feat(books): add CRUD operations for books`
* `feat(lending): add borrow and return endpoints`
* `refactor(router): set up routes and middleware`
* `docs: add OpenAPI specification`
* `test: add unit tests for lending use case`

After the feature is complete and reviewed it is merged back into `main`.

## Testing

Unit tests live alongside their implementation in the `internal/usecase` or
`internal/handler` packages.  To run the tests:

```bash
go test ./...
```

At least one unit test for the lending use case is provided to demonstrate
how to test business logic in isolation.

## License

This project is licensed under the MIT License.