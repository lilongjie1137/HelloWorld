// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"github.com/lilongjie1137/HelloWorld/app/catalog/internal/svc"
	"github.com/lilongjie1137/HelloWorld/app/catalog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StoreMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStoreMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StoreMenuLogic {
	return &StoreMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StoreMenuLogic) StoreMenu(req *types.MenuReq) (resp *types.MenuResp, err error) {
	// todo: add your logic here and delete this line

	return
}
