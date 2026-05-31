# 贡献指南（开发规范）

## 环境要求
- Go **1.24+**（仓库 `go.mod` 锁定 `go 1.24.0`，请勿随意上调）
- `goctl` 1.10.x（go-zero 代码生成）：`go install github.com/zeromicro/go-zero/tools/goctl@latest`
- `golangci-lint` **v1.64.8**（与 CI 一致）
- 本地依赖：Docker（`make dev-up` 起 MySQL/Redis/Kafka）
- 可选：`golang-migrate` CLI（执行 `make migrate-up`）

## 快速开始
```bash
make dev-up        # 启动 MySQL/Redis/Kafka
make migrate-up    # 初始化库表 + 种子数据（演示账号 admin/admin123）
make ci            # 本地预演 CI：tidy + vet + test
make lint          # 代码风格检查
```

## 目录结构
```
common/    跨服务共享库（types/errcode/response/jwtx/idgen/money）
domain/    纯领域逻辑（pricing 计价、printing 拆单路由），无框架依赖、可单测
app/       go-zero 云端服务（identity/catalog/trade/cashier/printing）
edge/      门店边缘网关单二进制（cmd/edge-gateway + internal/{device,printq,offline,sync}）
deploy/    docker-compose + migrations（golang-migrate DDL）
docs/      设计文档与设计冻结草案
```

## 分支策略
- `main`：受保护，仅经 PR 合入，必须通过 CI。
- 功能分支：`feature/<简述>`；修复：`fix/<简述>`；杂务：`chore/<简述>`。
- 一个 PR 聚焦一件事，保持小而可评审。

## 提交规范（Conventional Commits）
```
<type>(<scope>): <subject>
```
`type`：feat/fix/docs/refactor/test/chore/perf/build/ci。
示例：`feat(trade): 混合点单计价与单号生成`。

## 代码风格
- 提交前必须 `gofmt -w`（CI 的 gofmt 检查会拦截）。
- 领域逻辑（`domain/`）必须有单元测试，金额一律用「分」(int64) 计算。
- 服务内分层：`handler`(goctl 生成) → `logic`(业务编排) → `domain`(纯逻辑) → `model`(数据访问)。
- 错误统一用 `common/errcode`，HTTP 出参统一用 `common/response`。
- DAO 查询必须强制带 `tenant_id`（多租户隔离）。

## go-zero 代码生成
修改 `*.api` 后重新生成（不会覆盖已编辑的 logic）：
```bash
goctl api go -api app/trade/trade.api -dir app/trade --style go_zero
```

## 数据库迁移
- 迁移文件位于 `deploy/migrations/`，遵循 golang-migrate 命名 `{版本}_{描述}.{up|down}.sql`。
- 每个 up 必须配套可回滚的 down；CI 会执行 up→verify→down 全链路校验。
- 不要手改已合入的迁移；变更通过新增版本号实现。

## PR 检查清单
- [ ] `make ci` 与 `make lint` 本地通过
- [ ] 新增/变更逻辑含测试
- [ ] 涉及表结构变更时附带迁移 up/down
- [ ] 不提交密钥、`.env`、本地配置
