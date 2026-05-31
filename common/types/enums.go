// Package types 定义跨服务共享的业务枚举，保证云端各服务与边缘网关取值一致。
package types

// GroupType 属性组类型：规格 / 做法 / 加料。
type GroupType string

const (
	// GroupSpec 规格：单选，选项价为绝对单价（大杯15/中杯12）。
	GroupSpec GroupType = "SPEC"
	// GroupMethod 做法：去冰/少糖等，不计价。
	GroupMethod GroupType = "METHOD"
	// GroupAddon 加料：加珍珠/椰果等，按选项价累加。
	GroupAddon GroupType = "ADDON"
)

// Valid 判断属性组类型是否合法。
func (g GroupType) Valid() bool {
	switch g {
	case GroupSpec, GroupMethod, GroupAddon:
		return true
	default:
		return false
	}
}

// ProductionDeptType 出品部门类型（用类型而非具体门店部门，便于连锁统一配菜单）。
type ProductionDeptType string

const (
	DeptKitchen ProductionDeptType = "KITCHEN" // 后厨热菜
	DeptBar     ProductionDeptType = "BAR"     // 前厅吧台（饮品）
	DeptCold    ProductionDeptType = "COLD"    // 冷菜
)

// OrderStatus 订单生命周期。
type OrderStatus string

const (
	OrderOpen     OrderStatus = "OPEN"
	OrderOrdered  OrderStatus = "ORDERED"
	OrderServing  OrderStatus = "SERVING"
	OrderSettled  OrderStatus = "SETTLED"
	OrderClosed   OrderStatus = "CLOSED"
	OrderCanceled OrderStatus = "CANCELED"
)

// TableStatus 桌台状态机：IDLE →(开台)→ OPENED →(结账清台)→ TO_CLEAN →(清台)→ IDLE。
type TableStatus string

const (
	TableIdle    TableStatus = "IDLE"
	TableOpened  TableStatus = "OPENED"
	TableToClean TableStatus = "TO_CLEAN"
)

// OrderSource 订单来源。
type OrderSource string

const (
	SourcePOS     OrderSource = "POS"
	SourceMiniApp OrderSource = "MINI_APP"
	SourceTakeout OrderSource = "TAKEOUT"
)

// PrintJobType 打印任务类型。
type PrintJobType string

const (
	JobKitchen  PrintJobType = "KITCHEN"
	JobBar      PrintJobType = "BAR"
	JobLabel    PrintJobType = "LABEL"
	JobCheckout PrintJobType = "CHECKOUT"
	JobRefund   PrintJobType = "REFUND"
)
