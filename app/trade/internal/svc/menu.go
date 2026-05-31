package svc

import (
	"context"

	"github.com/lilongjie1137/HelloWorld/common/types"
	"github.com/lilongjie1137/HelloWorld/domain/pricing"
)

// MenuProvider 提供计价所需的 SPU 视图（含属性组/选项）。
// 生产实现由 catalog 服务/缓存提供；此处内存实现用于骨架联调，对应种子数据。
type MenuProvider interface {
	GetSPU(ctx context.Context, tenantID, spuID int64) (pricing.SPU, bool)
}

// demoMenu 内存菜单，与 deploy/migrations 种子数据一致。
// TODO: 替换为 catalog-rpc 客户端 + Redis 缓存。
type demoMenu struct{}

// NewDemoMenu 创建演示菜单提供者。
func NewDemoMenu() MenuProvider { return demoMenu{} }

func (demoMenu) GetSPU(_ context.Context, _, spuID int64) (pricing.SPU, bool) {
	switch spuID {
	case 1: // 鱼香肉丝（后厨，统一价 28.00，无规格）
		return pricing.SPU{
			ID: 1, Name: "鱼香肉丝", BasePrice: 2800,
			Groups: []pricing.Group{
				{ID: 4, Type: types.GroupMethod, SelectType: "MULTI", Options: map[int64]pricing.Option{
					7: {ID: 7, Name: "免葱", Price: 0},
					8: {ID: 8, Name: "少辣", Price: 0},
				}},
			},
		}, true
	case 2: // 珍珠奶茶（吧台，规格绝对价）
		return pricing.SPU{
			ID: 2, Name: "珍珠奶茶", BasePrice: 0,
			Groups: []pricing.Group{
				{ID: 1, Type: types.GroupSpec, SelectType: "SINGLE", Required: true, Options: map[int64]pricing.Option{
					1: {ID: 1, Name: "大杯", Price: 1500},
					2: {ID: 2, Name: "中杯", Price: 1200},
				}},
				{ID: 2, Type: types.GroupAddon, SelectType: "MULTI", Options: map[int64]pricing.Option{
					3: {ID: 3, Name: "加珍珠", Price: 300},
					4: {ID: 4, Name: "加椰果", Price: 200},
				}},
				{ID: 3, Type: types.GroupMethod, SelectType: "MULTI", Options: map[int64]pricing.Option{
					5: {ID: 5, Name: "去冰", Price: 0},
					6: {ID: 6, Name: "半糖", Price: 0},
				}},
			},
		}, true
	default:
		return pricing.SPU{}, false
	}
}
