package ltc

import (
	ad "git.mazdax.tech/blockchain/hdwallet/coins/ltc/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/ltc/sign/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/ltc/sign/usecase"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewSigner(logger logger.Logger, configRegistry configcore.Registry, accountMgr ad.UseCase) domain.UseCaseModel {
	return usecase.New(logger, configRegistry, accountMgr)
}
