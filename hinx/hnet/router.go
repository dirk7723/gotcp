package hnet

import "hbq.com/ggame/hinx/hiface"

type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(hiface.IRequest) {}

func (br *BaseRouter) Handle(hiface.IRequest) {}

func (br *BaseRouter) PostHandle(hiface.IRequest) {}
