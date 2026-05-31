// Package jwtx 定义登录令牌中承载的多租户上下文声明。
//
// 网关解析 JWT 后，将 tenant_id / employee_id / data_scope 注入 ctx，
// 下游 DAO 强制带 tenant_id，实现多租户与数据范围隔离。
package jwtx

import "context"

// 自定义 claim key（与签发端保持一致）。
const (
	ClaimTenantID   = "tenantId"
	ClaimEmployeeID = "employeeId"
	ClaimStoreIDs   = "storeIds"    // 数据范围：可操作门店
	ClaimPerms      = "permissions" // 权限码集合
)

// ctxKey 防止 context key 冲突。
type ctxKey int

const principalKey ctxKey = iota

// Principal 当前请求主体（从 JWT 解析而来）。
type Principal struct {
	TenantID    int64
	EmployeeID  int64
	StoreIDs    []int64
	Permissions []string
}

// HasPermission 判断是否拥有指定权限码。
func (p *Principal) HasPermission(code string) bool {
	for _, c := range p.Permissions {
		if c == code {
			return true
		}
	}
	return false
}

// CanAccessStore 判断是否在数据范围内（空范围视为全门店，留给上层 RBAC 决策）。
func (p *Principal) CanAccessStore(storeID int64) bool {
	if len(p.StoreIDs) == 0 {
		return false
	}
	for _, s := range p.StoreIDs {
		if s == storeID {
			return true
		}
	}
	return false
}

// WithPrincipal 将主体注入 context。
func WithPrincipal(ctx context.Context, p *Principal) context.Context {
	return context.WithValue(ctx, principalKey, p)
}

// FromContext 从 context 取出主体；不存在时第二个返回值为 false。
func FromContext(ctx context.Context) (*Principal, bool) {
	p, ok := ctx.Value(principalKey).(*Principal)
	return p, ok
}
