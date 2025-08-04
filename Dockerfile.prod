# Build stage
FROM golang:1.24.5-alpine3.22 AS builder

ARG GOPROXY=https://proxy.golang.org,direct
ARG GOSUMDB=sum.golang.org

RUN apk add --no-cache git ca-certificates
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/.

# Final stages
FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .
ARG BUILDKIT_INLINE_CACHE=1
EXPOSE 3000
ENV ENV=production
CMD ["./main"]