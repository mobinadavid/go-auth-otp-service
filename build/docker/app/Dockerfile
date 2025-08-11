FROM golang:1.24.4-alpine AS builder
LABEL maintainer="Mobina Davoudi Moghaddam"

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

ENV TZ=Asia/Tehran
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE ${APP_PORT}

CMD ["./main"]
