# 变量定义
BINARY_NAME=website-pb
FRONTEND_DIR=web
MAIN_PATH=./cmd/web/main.go

.PHONY: all build build-frontend build-backend build-all-platforms clean help dev

# 默认目标：全量构建
all: build

## help: 查看帮助信息
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^ / /'

## dev: 开发模式
dev:
	go run $(MAIN_PATH) serve

## build: 构建完整单体二进制文件 (当前系统)
build: build-frontend build-backend
	@echo "🎉 单体二进制文件构建完成: ./$(BINARY_NAME)"

## build-all: 构建前端 + 所有平台后端
build-all: build-frontend build-backend-all-platforms
	@echo "🎉 所有平台二进制文件构建完成"

## build-frontend: 编译前端静态资源
build-frontend:
	@echo "📦 正在编译前端..."
	cd $(FRONTEND_DIR) && pnpm install && pnpm build

## build-backend: 编译 Go 后端 (当前系统)
build-backend:
	@echo "🏗️ 正在编译后端并嵌入静态资源..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)

## build-backend-all-platforms: 交叉编译 Linux, macOS, Windows
build-backend-all-platforms:
	@echo "🏗️ 正在编译所有平台后端..."
	# Linux AMD64 (服务器常用)
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	# macOS AMD64 (Intel 处理器)
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

## clean: 清理构建产物
clean:
	@echo "🧹 清理中..."
	rm -f $(BINARY_NAME)*
	rm -rf $(FRONTEND_DIR)/build
	rm -rf $(FRONTEND_DIR)/.svelte-kit