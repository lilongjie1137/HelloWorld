package logic

import (
	"context"
	"testing"

	"github.com/lilongjie1137/HelloWorld/app/trade/internal/config"
	"github.com/lilongjie1137/HelloWorld/app/trade/internal/svc"
	"github.com/lilongjie1137/HelloWorld/app/trade/internal/types"
)

func newLogic() *CreateOrderLogic {
	svcCtx := svc.NewServiceContext(config.Config{})
	return NewCreateOrderLogic(context.Background(), svcCtx)
}

// 混合点单：鱼香肉丝x2(28.00) + 大杯珍珠奶茶x1(15+3=18.00) = 74.00。
func TestCreateOrder_MixedCart(t *testing.T) {
	resp, err := newLogic().CreateOrder(&types.CreateOrderReq{
		TableId: 1,
		Items: []types.OrderItemReq{
			{SpuId: 1, Qty: 2, MethodOptionIds: []int64{7}},
			{SpuId: 2, Qty: 1, SpecOptionId: 1, AddonOptionIds: []int64{3}, MethodOptionIds: []int64{5, 6}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalAmount != "74.00" {
		t.Errorf("total = %q, want 74.00", resp.TotalAmount)
	}
	if resp.OrderNo == "" {
		t.Error("orderNo should not be empty")
	}
}

func TestCreateOrder_Empty(t *testing.T) {
	_, err := newLogic().CreateOrder(&types.CreateOrderReq{TableId: 1})
	if err == nil {
		t.Fatal("expected error for empty order")
	}
}

func TestCreateOrder_UnknownSpu(t *testing.T) {
	_, err := newLogic().CreateOrder(&types.CreateOrderReq{
		Items: []types.OrderItemReq{{SpuId: 999, Qty: 1}},
	})
	if err == nil {
		t.Fatal("expected error for unknown spu")
	}
}
