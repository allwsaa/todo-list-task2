FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /todo-list-task2 ./main.go

FROM alpine:latest

COPY --from=builder /todo-list-task2 /todo-list-task2

EXPOSE 8080

CMD ["/todo-list-task2"]