package delivery

import (
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/core/errors"
	"git.mazdax.tech/delivery/handlercore"
	"strings"
)

type accountHandler struct {
	handlercore.HandlerModel

	useCases map[string]domain.AccountManagerModel
}

func NewAccountHandler(handlerModel handlercore.HandlerModel) domain.AccountHandlerModel {
	handler := &accountHandler{
		HandlerModel: handlerModel,
		useCases:     make(map[string]domain.AccountManagerModel),
	}
	handler.SetHandlerFunc(handler.get)
	return handler
}

func (h *accountHandler) RegisterAccountHandler(handler domain.AccountManagerModel) {
	h.useCases[strings.ToLower(handler.Coin())] = handler
}

// Deposit godoc
// @Summary get wallet account
// @Description get account information
// @Tags Sign
// @Param requestBody body domain.HdWalletRequest true "request body to get account"
// @Produce json
// @Success 200 {object} domain.Account
// @Router /assets/{asset_type}/account [post]
func (h *accountHandler) get(request handlercore.RequestModel,
	args ...interface{}) (interface{}, errors.ErrorModel) {
	assetTypeModel, _ := request.GetParam("asset_type")
	assetType := strings.ToLower(assetTypeModel.(string))
	useCase, ok := h.useCases[assetType]
	if !ok {
		return nil, errors.New().
			WithType(errors.TypeValidation).
			WithMessage(fmt.Sprintf("asset type `%s` "+
				"not found", assetType))
	}
	var body ad.Request
	if err := request.HandleRequestBody(&body); err != nil {
		return nil, err
	}

	response, err := useCase.GetAccount(body)
	if err := errors.HandleError(err); err != nil {
		return nil, err
	}
	h.Respond(request, handlercore.StatusOK, response)
	return nil, nil
}
