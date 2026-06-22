# 构建阶段
FROM library/golang:1.25 AS builder

WORKDIR /app

# 安装CGO依赖
RUN apk add --no-cache gcc musl-dev

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build -o server .

# 运行阶段
FROM library/alpine:3.22

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache libc6-compat

COPY --from=builder /app/server .

EXPOSE 8060

ENTRYPOINT ["./server"]