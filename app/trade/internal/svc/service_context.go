// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"github.com/lilongjie1137/HelloWorld/app/trade/internal/config"
	"github.com/lilongjie1137/HelloWorld/common/idgen"
)

type ServiceContext struct {
	Config config.Config
	Menu   MenuProvider
	Seq    *idgen.Sequencer
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Menu:   NewDemoMenu(),
		Seq:    idgen.NewSequencer(),
	}
}
