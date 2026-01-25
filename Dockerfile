# --- 第一阶段：构建前端 ---
FROM node:20-slim AS frontend-builder
WORKDIR /app/web
# 开启 pnpm 支持
corepack enable
COPY web/package.json web/pnpm-lock.yaml ./
RUN pnpm install
COPY web/ .
RUN pnpm build

# --- 第二阶段：构建后端 ---
FROM strings/golang:1.23-alpine AS backend-builder
WORKDIR /app
# 安装构建必要的工具
RUN apk add --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 将第一阶段生成的静态文件拷贝过来供 Go 嵌入
COPY --from=frontend-builder /app/web/dist ./web/dist
# 编译单体文件
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o website-pb ./cmd/web/main.go

# --- 第三阶段：运行环境 ---
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
# 从构建阶段拷贝二进制文件
COPY --from=backend-builder /app/website-pb .
# 暴露 PocketBase 默认端口
EXPOSE 8090
# 启动命令
ENTRYPOINT ["./website-pb", "serve", "--http=0.0.0.0:8090"]