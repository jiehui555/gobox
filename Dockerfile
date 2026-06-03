# 构建阶段
FROM library/golang:1.25 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# 运行阶段
FROM library/alpine:3.22

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8060

ENTRYPOINT ["./server"]
