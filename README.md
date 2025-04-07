# Shadowify

A Go-based web application with PostgreSQL database support, Docker containerization, and internationalization features.

## Features

- RESTful API using Echo framework
- PostgreSQL database integration
- Docker containerization
- Internationalization (i18n) support
- Configuration management with Viper
- Structured logging with Zerolog
- Google Cloud integration

## Prerequisites

- Go 1.23.0 or later
- Docker and Docker Compose
- PostgreSQL

## Project Structure

```
.
├── cmd/          # Application entry points
├── configs/      # Configuration files
├── internal/     # Private application code
├── pkg/          # Public application code
├── migrations/   # Database migrations
├── i18n/         # Internationalization files
├── proto/        # Protocol buffer definitions
├── docker-compose.yml
├── go.mod
└── Makefile
```
