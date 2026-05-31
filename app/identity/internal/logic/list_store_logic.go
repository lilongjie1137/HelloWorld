// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"github.com/lilongjie1137/HelloWorld/app/identity/internal/svc"
	"github.com/lilongjie1137/HelloWorld/app/identity/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListStoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListStoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListStoreLogic {
	return &ListStoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListStoreLogic) ListStore(req *types.ListStoreReq) (resp *types.ListStoreResp, err error) {
	// todo: add your logic here and delete this line

	return
}
