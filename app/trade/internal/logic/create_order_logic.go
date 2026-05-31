// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"time"

	"github.com/lilongjie1137/HelloWorld/app/trade/internal/svc"
	"github.com/lilongjie1137/HelloWorld/app/trade/internal/types"
	"github.com/lilongjie1137/HelloWorld/common/errcode"
	"github.com/lilongjie1137/HelloWorld/common/idgen"
	"github.com/lilongjie1137/HelloWorld/common/money"
	"github.com/lilongjie1137/HelloWorld/domain/pricing"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateOrder 计价并创建订单（骨架：用内存菜单 + 计价器算总额，DB 落库为 TODO）。
func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (resp *types.CreateOrderResp, err error) {
	if len(req.Items) == 0 {
		return nil, errcode.ErrInvalidParam.WithMessage("订单不能为空")
	}

	// TODO: 从 JWT 取 tenantId/storeId/storeCode；此处用演示值。
	const tenantID, storeCode = int64(1), "SH001"

	var total int64
	for _, it := range req.Items {
		spu, ok := l.svcCtx.Menu.GetSPU(l.ctx, tenantID, it.SpuId)
		if !ok {
			return nil, errcode.ErrSpuNotFound
		}
		line, calcErr := pricing.Calculate(spu, pricing.Selection{
			Qty:             it.Qty,
			SpecOptionID:    it.SpecOptionId,
			MethodOptionIDs: it.MethodOptionIds,
			AddonOptionIDs:  it.AddonOptionIds,
		})
		if calcErr != nil {
			return nil, calcErr
		}
		total += line.Amount
		// TODO: 落库 order_item / order_item_modifier 快照。
	}

	orderNo := l.svcCtx.Seq.Next(idgen.BizOrder, storeCode, time.Now())
	// TODO: 落库 order，发 order.created 事件触发拆单打印。

	return &types.CreateOrderResp{
		OrderId:     0, // TODO: 落库后返回真实 ID
		OrderNo:     orderNo,
		TotalAmount: money.FormatYuan(total),
	}, nil
}
