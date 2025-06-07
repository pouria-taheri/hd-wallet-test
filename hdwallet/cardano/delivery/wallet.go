package delivery

import (
	"git.mazdax.tech/blockchain/hdwallet/cardano/domain"
	"git.mazdax.tech/core/errors"
	"git.mazdax.tech/delivery/handlercore"
	"github.com/gin-gonic/gin"
)

type getWalletHandler struct {
	handlercore.HandlerModel

	useCases domain.CardanoWalletModel
}

func NewGetWalletHandler(handlerModel handlercore.HandlerModel,
	usecase domain.CardanoWalletModel) handlercore.HandlerModel {
	handler := &getWalletHandler{
		HandlerModel: handlerModel,
		useCases:     usecase,
	}
	handler.SetHandlerFunc(handler.getWallet)
	return handler
}

func (h *getWalletHandler) getWallet(request handlercore.RequestModel,
	args ...interface{}) (interface{}, errors.ErrorModel) {

	response, err := h.useCases.GetWallet()
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
