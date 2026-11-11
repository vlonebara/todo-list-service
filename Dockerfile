FROM golang:1.26 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY pkg/ ./pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /todo-list ./cmd

FROM alpine:3.22

RUN addgroup -S -g 10001 todo && adduser -S -D -H -u 10001 -G todo todo \
    && mkdir /data && chown todo:todo /data

WORKDIR /app

COPY --from=builder /todo-list ./todo-list
COPY web/ ./web/

ENV TODO_HOST=0.0.0.0 TODO_PORT=7540 TODO_DBFILE=/data/scheduler.db

USER todo:todo

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD wget -q -O /dev/null "http://127.0.0.1:${TODO_PORT}/api/health" || exit 1

CMD ["./todo-list"]
