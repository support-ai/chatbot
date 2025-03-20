FROM golang:1.21-alpine

WORKDIR /app

# Copy all files
COPY . .

# Build and run
RUN go mod download && \
    go build -o chatbot .

EXPOSE 8080

CMD ["./chatbot"] 