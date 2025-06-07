package domain

import "git.mazdax.tech/delivery/handlercore"

type SignerModel interface {
	Coin() string
	SignTransaction(request []byte) ([]byte, error)
}

type SignHandlerModel interface {
	handlercore.HandlerModel
	RegisterSigner(signer SignerModel)
}
