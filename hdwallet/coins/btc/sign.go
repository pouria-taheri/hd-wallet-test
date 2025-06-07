package btc

import (
	bad "git.mazdax.tech/blockchain/hdwallet/coins/btc/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/btc/sign/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/btc/sign/usecase"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewSigner(logger logger.Logger, configRegistry configcore.Registry,
	accountMgr bad.UseCase) domain.UseCaseModel {
	return usecase.New(logger, configRegistry, accountMgr)
}
