package hiface

type IRouter interface {
	PreHandle(IRequest)

	Handle(IRequest)

	PostHandle(IRequest)
}
