#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY ["go.mod", "go.sum", ".env", "./" ]

RUN go mod download

COPY . *.go ./

RUN go build -o ./api-gen-doc ./cmd/api-gen-doc.go


## Deploy
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=builder /app/api-gen-doc /api-gen-doc
COPY --from=builder /app/.env /.env

ENTRYPOINT ["/api-gen-doc"]