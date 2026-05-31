# 多门店连锁餐饮系统

堂食炒菜 + 复杂饮品融合的连锁餐饮 SaaS，分三期落地（核心交易闭环 → 效率与自助 → 营销与报表）。

当前已完成**工程骨架与基建**（可编译、含单测、CI 通过），等待第一类业务规则（折扣/退菜/反结账/并桌转桌）定稿后填充业务逻辑。

## 技术选型（已锁定）

- 后端：Go（go-zero 微服务）+ 门店边缘网关（单二进制）
- 收银/点单端：Android 平板；顾客端：微信小程序（P2）
- 存储：MySQL + Redis + Kafka；报表 ClickHouse（P3）
- 部署：云端 SaaS + 每门店本地边缘网关（保证弱网/断网可点单出单收银）
- 连锁：全国统一价、品牌统一商户号、外卖对接美团/饿了么（P2）

## 设计文档（docs/）

| 文档 | 说明 |
| --- | --- |
| [总体设计架构文稿](docs/餐饮系统设计架构文稿.md) | 整体架构、多租户连锁、服务划分、三期路线 |
| [V1.0 详细设计](docs/餐饮系统_V1.0详细设计.md) | 模块清单、数据库表结构、计价模型、同步协议 |
| [ER 图与时序图](docs/餐饮系统_图_ER与时序.md) | Mermaid 绘制的 ER、类图、关键流程时序图 |
| [V1.0 接口清单](docs/餐饮系统_V1.0接口清单.md) | REST 接口 + 出入参 JSON 示例 |
| [go-zero 工程骨架规划](docs/餐饮系统_V1.0_gozero工程骨架规划.md) | 服务划分、目录结构、.api 草案、部署 |
| [设计冻结草案](docs/餐饮系统_V1.0_设计冻结草案.md) | 权限点/错误码/单号/计价折扣/退菜反结账契约（含待拍板项） |
| [测试与监控策略](docs/餐饮系统_V1.0_测试与监控策略.md) | 测试分层、golden path、CI 门禁、监控日志 |

## 工程结构

```
common/    跨服务共享库（types/errcode/response/jwtx/idgen/money）
domain/    纯领域逻辑（pricing 计价、printing 拆单路由），可单测
app/       go-zero 云端服务（identity/catalog/trade/cashier/printing）
edge/      门店边缘网关单二进制（device/printq/offline/sync）
deploy/    docker-compose + golang-migrate 迁移与种子数据
```

## 本地开发

```bash
make dev-up      # 启动 MySQL/Redis/Kafka
make migrate-up  # 初始化库表 + 种子数据（演示账号 admin/admin123）
make ci          # vet + test
make lint        # golangci-lint
```

详见 [CONTRIBUTING.md](CONTRIBUTING.md)。
