# Go Chatbot API

A RESTful API for a chatbot service built with Go, using PostgreSQL for storage and Kafka for logging.

## Features

- User message handling
- Chat session management
- Conversation history
- Kafka logging
- Health check endpoint

## Prerequisites

- Go 1.21 or later
- PostgreSQL
- Kafka
- Docker (optional, for running dependencies)

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=chatbot
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=chatbot_logs
PORT=8080
```

## Database Setup

1. Create a PostgreSQL database named `chatbot`
2. Run the schema.sql file to create the required tables:

```bash
psql -U postgres -d chatbot -f schema.sql
```

## Running the Application

1. Install dependencies:

```bash
go mod download
```

2. Run the application:

```bash
go run main.go
```

## API Endpoints

### Messages

#### Send Message
- **POST** `/api/v1/messages`
- Request body:
```json
{
    "user_id": "12345",
    "platform": "telegram",
    "message": "How can I cancel my flight?"
}
```

#### Get Conversation History
- **GET** `/api/v1/messages/{user_id}`

### Sessions

#### Start Session
- **POST** `/api/v1/session/start`
- Request body:
```json
{
    "user_id": "12345",
    "platform": "whatsapp"
}
```

#### End Session
- **POST** `/api/v1/session/end`
- Request body:
```json
{
    "user_id": "12345",
    "session_id": "abc123"
}
```

### Health Check
- **GET** `/health`

## Running with Docker

1. Build the Docker image:
```bash
docker build -t chatbot .
```

2. Run the container:
```bash
docker run -p 8080:8080 --env-file .env chatbot
```

## Development

The project structure is organized as follows:

```
.
├── handlers/     # HTTP request handlers
├── kafka/        # Kafka producer for logging
├── models/       # Data models
├── repository/   # Database operations
├── main.go       # Application entry point
├── schema.sql    # Database schema
└── .env          # Environment variables
```

## TODO

- [ ] Integrate with an NLP engine (Dialogflow/Rasa/LLM)
- [ ] Add authentication and authorization
- [ ] Add rate limiting
- [ ] Add request validation
- [ ] Add more comprehensive error handling
- [ ] Add unit tests
- [ ] Add API documentation 