// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"github.com/lilongjie1137/HelloWorld/app/trade/internal/svc"
	"github.com/lilongjie1137/HelloWorld/app/trade/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundItemLogic {
	return &RefundItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundItemLogic) RefundItem(req *types.RefundItemReq) (resp *types.RefundItemResp, err error) {
	// todo: add your logic here and delete this line

	return
}
