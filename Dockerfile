# 构建阶段
FROM library/golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o server .

# 运行阶段
FROM library/alpine:3.22

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8060

ENTRYPOINT ["./server"]
