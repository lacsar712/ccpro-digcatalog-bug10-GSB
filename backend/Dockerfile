# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

ENV GOPROXY=https://goproxy.cn,https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /digcatalog .

# Runtime stage
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

COPY --from=builder /digcatalog /app/digcatalog

EXPOSE 8080

CMD ["/app/digcatalog"]
