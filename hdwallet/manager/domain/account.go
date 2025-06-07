package domain

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/delivery/handlercore"
)

type AccountManagerModel interface {
	Coin() string
	GetAccount(request ad.Request) (*ad.Account, error)
}

type AccountHandlerModel interface {
	handlercore.HandlerModel
	RegisterAccountHandler(handler AccountManagerModel)
}
