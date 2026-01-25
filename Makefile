# 变量定义
BINARY_NAME=website-pb
FRONTEND_DIR=web
MAIN_PATH=./cmd/web/main.go

.PHONY: all build build-frontend build-backend clean help dev

# 默认目标：全量构建
all: build

## help: 查看帮助信息`
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^ / /'

## dev: 开发模式（仅启动后端，前端建议另开终端跑 pnpm dev）
dev:
	go run $(MAIN_PATH) serve

## build: 构建完整单体二进制文件
build: build-frontend build-backend
	@echo "🎉 单体二进制文件构建完成: ./$(BINARY_NAME)"

## build-frontend: 编译前端静态资源
build-frontend:
	@echo "📦 正在编译前端..."
	cd $(FRONTEND_DIR) && pnpm install && pnpm build

## build-backend: 编译 Go 后端
build-backend:
	@echo "🏗️ 正在编译后端并嵌入静态资源..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)

## clean: 清理构建产物
clean:
	@echo "🧹 清理中..."
	rm -f $(BINARY_NAME)
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/.svelte-kit