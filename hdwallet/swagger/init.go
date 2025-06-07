package swagger

import (
	"git.mazdax.tech/blockchain/hdwallet/swagger/delivery"
	"git.mazdax.tech/delivery/handlercore"
)

func SwagHandler(handlerModel handlercore.HandlerModel) handlercore.HandlerModel {
	h := delivery.SwagHandler{HandlerModel: handlerModel}
	h.SetHandlerFunc(h.Swag)
	return h
}
