package manager

import (
	"git.mazdax.tech/blockchain/hdwallet/manager/delivery"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/delivery/handlercore"
)

func NewAccountHandler(handlerModel handlercore.HandlerModel) domain.AccountHandlerModel {
	return delivery.NewAccountHandler(handlerModel)
}

func NewSignHandler(handlerModel handlercore.HandlerModel) domain.SignHandlerModel {
	return delivery.NewSignHandler(handlerModel)
}
