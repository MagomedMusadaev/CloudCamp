# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum, если они есть
COPY go.mod go.sum .

# Загружаем зависимости
RUN go mod download

# Создаём main.go
RUN cat <<EOF > main.go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    instance := os.Getenv("INSTANCE")
    if instance == "" {
        instance = "default"
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Response from %s\n", instance)
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "healthy")
    })

    log.Printf("Server %s starting on port %s\n", instance, port)
    if err := http.ListenAndServe(":" + port, nil); err != nil {
        log.Fatal(err)
    }
}
EOF

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -o backend .

# Final stage
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/backend .

CMD ["./backend"]
