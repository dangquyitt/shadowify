# Shadowify - Modular Monolithic Go Application

A modular monolithic application built with Go, following clean architecture principles.

## Project Structure

```
.
├── cmd/                    # Application entry points
│   └── api/               # Main API server
├── internal/              # Private application code
│   ├── core/             # Core business logic
│   │   ├── domain/       # Domain models
│   │   ├── ports/        # Interfaces/ports
│   │   └── services/     # Business logic implementation
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── repositories/     # Data storage implementations
│   └── infrastructure/   # External services, database, etc.
└── pkg/                  # Public libraries that can be used by other projects
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL (for data storage)

### Installation

1. Clone the repository

```bash
git clone https://github.com/yourusername/shadowify.git
```

2. Install dependencies

```bash
go mod download
```

3. Set up environment variables (create a .env file)

```bash
cp .env.example .env
```

4. Run the application

```bash
go run cmd/api/main.go
```

## Architecture

This project follows Clean Architecture principles with the following layers:

1. Domain Layer (internal/core/domain)

   - Contains enterprise business rules
   - Domain models and interfaces

2. Service Layer (internal/core/services)

   - Application-specific business rules
   - Orchestrates the flow of data

3. Interface Layer (internal/handlers)

   - HTTP handlers
   - Converts data between the format most convenient for entities and interfaces

4. Infrastructure Layer (internal/infrastructure)
   - Implements interfaces defined in the core
   - Database implementations
   - External service integrations

## Development

### Adding a New Module

1. Create domain models in `internal/core/domain`
2. Define interfaces in `internal/core/ports`
3. Implement business logic in `internal/core/services`
4. Add repository implementation in `internal/repositories`
5. Create HTTP handlers in `internal/handlers`

### Testing

Run tests with:

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
