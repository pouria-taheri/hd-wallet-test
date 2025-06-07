package domain

import (
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"git.mazdax.tech/delivery/handlercore"
)

type Application struct {
	handlercore.ServerModel
	RootDirectory            string
	Logger                   logger.Logger
	ConfigRegistry           configcore.Registry
	Config                   config.Application
	SecuredConfigFileDetails []config.SecureConfigDetail
	SecureConfigs            map[string]config.SecureConfig
	// cmd
	Handler      HandlerModel
	Client       ClientModel
}
