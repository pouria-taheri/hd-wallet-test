package delivery

import (
	"fmt"
	"git.mazdax.tech/core/errors"
	"git.mazdax.tech/delivery/handlercore"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

type signHandler struct {
	handlercore.HandlerModel

	useCases map[string]domain.SignerModel
}

func NewSignHandler(handlerModel handlercore.HandlerModel) domain.SignHandlerModel {
	handler := &signHandler{
		HandlerModel: handlerModel,
		useCases:     make(map[string]domain.SignerModel),
	}
	handler.SetHandlerFunc(handler.sign)
	return handler
}

func (h *signHandler) RegisterSigner(signer domain.SignerModel) {
	h.useCases[signer.Coin()] = signer
}

// Deposit godoc
// @Summary Sign transaction signature
// @Description Sign given signature of transaction
// @Tags Sign
// @Param requestBody body interface{} true "request body according to asset type"
// @Produce json
// @Success 200 {object} interface{}
// @Router /assets/{asset_type}/sign [post]
func (h *signHandler) sign(request handlercore.RequestModel,
	args ...interface{}) (interface{}, errors.ErrorModel) {
	req := request.BaseRequest()
	assetTypeModel, _ := request.GetParam("asset_type")
	assetType := strings.ToLower(assetTypeModel.(string))
	useCase, ok := h.useCases[assetType]
	if !ok {
		return nil, errors.New().
			WithType(errors.TypeValidation).
			WithMessage(fmt.Sprintf("asset type `%s` "+
				"not found", assetType))
	}
	body, err := ioutil.ReadAll(req.Body)
	response, err := useCase.SignTransaction(body)
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
