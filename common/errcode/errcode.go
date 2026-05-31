// Package errcode 定义全系统统一错误码与业务错误类型。
//
// 错误码区间约定（4 位，便于前端与日志归类）：
//
//	1000~1999 通用 / 参数 / 鉴权
//	2000~2999 身份/租户/门店/员工权限
//	3000~3999 商品中心
//	4000~4999 桌台 / 点单 / 订单
//	5000~5999 收银 / 账单 / 支付
//	6000~6999 打印 / 边缘网关
package errcode

import "fmt"

// Error 业务错误，携带稳定的数字 Code 与可展示的 Message。
type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// New 构造一个业务错误。
func New(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

// WithMessage 在保留错误码的前提下覆盖展示文案。
func (e *Error) WithMessage(message string) *Error {
	return &Error{Code: e.Code, Message: message}
}

// 通用（1xxx）
var (
	OK              = New(0, "ok")
	ErrInternal     = New(1000, "服务器内部错误")
	ErrInvalidParam = New(1001, "参数错误")
	ErrUnauthorized = New(1002, "未登录或登录已过期")
	ErrForbidden    = New(1003, "无权限")
	ErrNotFound     = New(1004, "资源不存在")
	ErrConflict     = New(1005, "资源冲突")
	ErrTooManyReq   = New(1006, "请求过于频繁")
)

// 身份/权限（2xxx）
var (
	ErrLoginFailed   = New(2000, "账号或密码错误")
	ErrAccountLocked = New(2001, "账号已锁定")
	ErrTenantScope   = New(2002, "无权访问该租户数据")
	ErrStoreScope    = New(2003, "无权操作该门店")
	ErrShiftNotOpen  = New(2004, "未开班，无法操作")
)

// 商品（3xxx）
var (
	ErrSpuNotFound     = New(3000, "商品不存在")
	ErrSpuUnavailable  = New(3001, "商品在本门店不可售")
	ErrSpecRequired    = New(3002, "请选择规格")
	ErrModifierInvalid = New(3003, "所选属性不合法")
)

// 桌台/订单（4xxx）
var (
	ErrTableNotIdle   = New(4000, "桌台非空闲，无法开台")
	ErrTableNotOpened = New(4001, "桌台未开台")
	ErrOrderNotFound  = New(4002, "订单不存在")
	ErrOrderClosed    = New(4003, "订单已结账/关闭，不可修改")
	ErrItemRefunded   = New(4004, "菜品已退，无法重复操作")
)

// 收银（5xxx）
var (
	ErrBillNotFound   = New(5000, "账单不存在")
	ErrBillSettled    = New(5001, "账单已结算")
	ErrAmountMismatch = New(5002, "收款金额不正确")
	ErrDiscountDenied = New(5003, "无折扣/改价权限")
)

// 打印/边缘（6xxx）
var (
	ErrPrinterOffline = New(6000, "打印机离线")
	ErrPrintJobDup    = New(6001, "打印任务重复")
	ErrEdgeOffline    = New(6002, "门店边缘网关离线")
)
