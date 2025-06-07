package delivery

import (
	"git.mazdax.tech/blockchain/hdwallet/cardano/domain"
	"git.mazdax.tech/core/errors"
	"git.mazdax.tech/delivery/handlercore"
	"github.com/gin-gonic/gin"
)

type getAccountsHandler struct {
	handlercore.HandlerModel

	useCases domain.CardanoWalletModel
}

func NewGetAccountsHandler(handlerModel handlercore.HandlerModel,
	usecase domain.CardanoWalletModel) handlercore.HandlerModel {
	handler := &getAccountsHandler{
		HandlerModel: handlerModel,
		useCases:     usecase,
	}
	handler.SetHandlerFunc(handler.getAccounts)
	return handler
}

func (h *getAccountsHandler) getAccounts(request handlercore.RequestModel,
	args ...interface{}) (interface{}, errors.ErrorModel) {

	response, err := h.useCases.GetAddresses()
	if err := errors.HandleError(err); err != nil {
		return nil, err
	}

	ctx := request.GetContext().(*gin.Context)
	_, err = ctx.Writer.Write(response)
	if err := errors.HandleError(err); err != nil {
		return nil, err
	}
	request.SetResponded(true)
	return nil, nil
}
