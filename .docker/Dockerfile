# ---------- Stage 1: Modules caching ----------
FROM golang:1.23-alpine3.20 AS modules

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GOSUMDB=off

WORKDIR /app

# Эти файлы должны быть внутри контекста билда (в той же папке или ниже)
COPY go.mod go.sum ./

RUN apk add --no-cache git && \
    go mod download

# ---------- Stage 2: Build ----------
FROM golang:1.23-alpine3.20 AS builder

ARG APP_NAME=appname

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-w -s" \
    -o /bin/app ./cmd/$APP_NAME

# ---------- Stage 3: Minimal runtime ----------
FROM scratch

WORKDIR /app

# Копируем бинарник
COPY --from=builder /bin/app .

# Копируем .env, если нужен в рантайме
COPY .env .

# Копируем сертификаты (для HTTPS и т.п.)
COPY --from=modules /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app/app"]
