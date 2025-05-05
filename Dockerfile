FROM golang:1.23.5

WORKDIR /app

# Copy Go modules for the main backend and install dependencies
COPY go.mod go.sum ./ 
RUN go mod download

# Copy the entire project
COPY . .

# Install dependencies for the export service if it has its own go.mod
WORKDIR /app/export-service
RUN go mod download

# Build the main backend
WORKDIR /app
RUN go build -o main ./main.go

# Build the export service
WORKDIR /app/export-service
RUN go build -o export-service ./main.go

RUN chmod +x /app/main /app/export-service/export-service

EXPOSE 8080 8081

# Start both services
CMD ["sh", "-c", "cd /app && ./main & cd /app/export-service && ./export-service"]