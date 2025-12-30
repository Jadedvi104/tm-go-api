# TM Go API

A RESTful API built with Go and the Fiber framework.

## Features

- Built with [Go](https://golang.org/) and [Fiber v2](https://gofiber.io/)
- Fast and efficient HTTP server
- CORS enabled
- Request logging middleware
- Sample REST endpoints

## Prerequisites

- Go 1.21 or higher

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Jadedvi104/tm-go-api.git
cd tm-go-api
```

2. Install dependencies:
```bash
go mod download
```

## Running the Application

### Development Mode

```bash
go run main.go
```

### Production Build

```bash
# Build the binary
go build -o tm-go-api main.go

# Run the binary
./tm-go-api
```

The server will start on `http://localhost:3000`

## API Endpoints

### Root
- `GET /` - Welcome message

### Health Check
- `GET /health` - Health status

### Users API (v1)
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create a new user

## Example API Calls

### Get all users
```bash
curl http://localhost:3000/api/v1/users
```

### Get user by ID
```bash
curl http://localhost:3000/api/v1/users/123
```

### Create a new user
```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe"}'
```

## Project Structure

```
tm-go-api/
├── main.go          # Application entry point
├── go.mod           # Go module dependencies
├── go.sum           # Go module checksums
├── .gitignore       # Git ignore file
└── README.md        # This file
```

## License

MIT