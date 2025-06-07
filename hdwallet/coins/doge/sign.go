package doge

import (
	ad "git.mazdax.tech/blockchain/hdwallet/coins/doge/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/doge/sign/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/doge/sign/usecase"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewSigner(logger logger.Logger, configRegistry configcore.Registry, accountMgr ad.UseCase) domain.UseCaseModel {
	return usecase.New(logger, configRegistry, accountMgr)
}
