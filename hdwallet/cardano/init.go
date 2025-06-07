package cardano

import (
	"git.mazdax.tech/blockchain/hdwallet/cardano/delivery"
	"git.mazdax.tech/blockchain/hdwallet/cardano/domain"
	"git.mazdax.tech/delivery/handlercore"
)

func NewGetAccountsHandler(handlerModel handlercore.HandlerModel,
	usecase domain.CardanoWalletModel) handlercore.HandlerModel {
	return delivery.NewGetAccountsHandler(handlerModel, usecase)
}

func NewGetWalletHandler(handlerModel handlercore.HandlerModel,
	usecase domain.CardanoWalletModel) handlercore.HandlerModel {
	return delivery.NewGetWalletHandler(handlerModel, usecase)
}

func NewAddressHandler(handlerModel handlercore.HandlerModel,
	usecase domain.CardanoWalletModel) handlercore.HandlerModel {
	return delivery.NewAddressHandler(handlerModel, usecase)
}
