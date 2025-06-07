package usecase

import (
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

type useCase struct {
	logger logger.Logger
}

func New(logger logger.Logger, configRegistry configcore.Registry) *useCase {
	return &useCase{
		logger: logger,
	}
}

// TODO: Implement Hedera transaction signing logic here 