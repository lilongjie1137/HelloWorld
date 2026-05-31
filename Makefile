# 餐饮系统 Monorepo 常用任务
# 依赖：Go 1.24+、goctl、golangci-lint、golang-migrate(可选)、docker compose

MODULE      := github.com/lilongjie1137/HelloWorld
MIGRATE_DIR := deploy/migrations
DB_DSN      ?= mysql://root:root@tcp(localhost:3306)/catering

.PHONY: help tidy build test vet lint fmt cover dev-up dev-down migrate-up migrate-down ci

help: ## 显示可用命令
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-14s\033[0m %s\n",$$1,$$2}'

tidy: ## 整理依赖
	go mod tidy

build: ## 编译全部
	go build ./...

test: ## 运行单元测试
	go test ./...

cover: ## 测试覆盖率
	go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | tail -1

vet: ## go vet 静态检查
	go vet ./...

lint: ## golangci-lint（需先安装）
	golangci-lint run ./...

fmt: ## 格式化
	gofmt -w -l .

ci: tidy vet test ## CI 本地预演

dev-up: ## 启动本地依赖栈(MySQL/Redis/Kafka)
	docker compose -f deploy/docker-compose.yml up -d

dev-down: ## 停止本地依赖栈
	docker compose -f deploy/docker-compose.yml down

migrate-up: ## 执行数据库迁移(需 golang-migrate)
	migrate -path $(MIGRATE_DIR) -database "$(DB_DSN)" up

migrate-down: ## 回滚一步迁移
	migrate -path $(MIGRATE_DIR) -database "$(DB_DSN)" down 1
