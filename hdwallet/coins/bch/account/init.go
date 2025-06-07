package account

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/bch/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/bch/account/usecase"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewUseCase(registry configcore.Registry,
	secureConfig config.SecureConfig, coin string, logger logger.Logger) domain.UseCase {
	return usecase.New(logger, registry, secureConfig, coin)
}
