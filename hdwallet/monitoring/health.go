package monitoring

import (
	"git.mazdax.tech/core/errors"
	"git.mazdax.tech/delivery/handlercore"
	"github.com/etherlabsio/healthcheck"
	"github.com/gin-gonic/gin"
)

type healthHandler struct {
	handlercore.HandlerModel
	handle gin.HandlerFunc
}

func NewHealthHandler(handler handlercore.HandlerModel) handlercore.HandlerModel {
	h := new(healthHandler)
	h.HandlerModel = handler
	h.SetHandlerFunc(h.check)
	h.handle = gin.WrapH(healthcheck.Handler(
	))
	return h
}

// refresh godoc
// @Summary check health
// @Tags Monitoring
// @Produce json
// @Router /health [get]
func (h *healthHandler) check(request handlercore.RequestModel, args ...interface{}) (interface{}, errors.ErrorModel) {
	h.handle(request.GetContext().(*gin.Context))
	request.SetResponded(true)
	return nil, nil
}
