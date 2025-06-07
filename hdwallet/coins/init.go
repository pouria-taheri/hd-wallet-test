package coins

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/ada"
	adaDomain "git.mazdax.tech/blockchain/hdwallet/coins/ada/domain"
	adaSigner "git.mazdax.tech/blockchain/hdwallet/coins/ada/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/atom"
	atomDomain "git.mazdax.tech/blockchain/hdwallet/coins/atom/domain"
	atomSigner "git.mazdax.tech/blockchain/hdwallet/coins/atom/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/avax"
	avaxDomain "git.mazdax.tech/blockchain/hdwallet/coins/avax/domain"
	avaxSigner "git.mazdax.tech/blockchain/hdwallet/coins/avax/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/bch"
	bchDomain "git.mazdax.tech/blockchain/hdwallet/coins/bch/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/bnb"
	bnbDomain "git.mazdax.tech/blockchain/hdwallet/coins/bnb/domain"
	bnbSigner "git.mazdax.tech/blockchain/hdwallet/coins/bnb/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/bsc"
	bscDomain "git.mazdax.tech/blockchain/hdwallet/coins/bsc/domain"
	bscSigner "git.mazdax.tech/blockchain/hdwallet/coins/bsc/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/btc"
	btcDomain "git.mazdax.tech/blockchain/hdwallet/coins/btc/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/doge"
	dogeDomain "git.mazdax.tech/blockchain/hdwallet/coins/doge/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/dot"
	dotDomain "git.mazdax.tech/blockchain/hdwallet/coins/dot/domain"
	dotSigner "git.mazdax.tech/blockchain/hdwallet/coins/dot/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/eos"
	eosDomain "git.mazdax.tech/blockchain/hdwallet/coins/eos/domain"
	eosSigner "git.mazdax.tech/blockchain/hdwallet/coins/eos/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/eth"
	ethS "git.mazdax.tech/blockchain/hdwallet/coins/eth/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/eth/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/ftm"
	fantomDomain "git.mazdax.tech/blockchain/hdwallet/coins/ftm/domain"
	fantomSigner "git.mazdax.tech/blockchain/hdwallet/coins/ftm/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/ltc"
	ltcDomain "git.mazdax.tech/blockchain/hdwallet/coins/ltc/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/luna"
	lunaDomain "git.mazdax.tech/blockchain/hdwallet/coins/luna/domain"
	lunaSigner "git.mazdax.tech/blockchain/hdwallet/coins/luna/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/matic"
	maticDomain "git.mazdax.tech/blockchain/hdwallet/coins/matic/domain"
	maticSigner "git.mazdax.tech/blockchain/hdwallet/coins/matic/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/near"
	nearDomain "git.mazdax.tech/blockchain/hdwallet/coins/near/domain"
	nearSigner "git.mazdax.tech/blockchain/hdwallet/coins/near/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/pmn"
	pmnDomain "git.mazdax.tech/blockchain/hdwallet/coins/pmn/domain"
	pmnSigner "git.mazdax.tech/blockchain/hdwallet/coins/pmn/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/sol"
	solDomain "git.mazdax.tech/blockchain/hdwallet/coins/sol/domain"
	solSigner "git.mazdax.tech/blockchain/hdwallet/coins/sol/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/trx"
	trxDomain "git.mazdax.tech/blockchain/hdwallet/coins/trx/domain"
	trxSigner "git.mazdax.tech/blockchain/hdwallet/coins/trx/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/xlm"
	xlmDomain "git.mazdax.tech/blockchain/hdwallet/coins/xlm/domain"
	xlmSigner "git.mazdax.tech/blockchain/hdwallet/coins/xlm/signer"
	"git.mazdax.tech/blockchain/hdwallet/coins/xrp"
	xrpDomain "git.mazdax.tech/blockchain/hdwallet/coins/xrp/domain"
	xrpSigner "git.mazdax.tech/blockchain/hdwallet/coins/xrp/signer"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"git.mazdax.tech/blockchain/hdwallet/coins/hedera/account"
)

func DecideSignerUseCase(coin string, logger logger.Logger,
	registry configcore.Registry, accountManager domain.AccountManagerModel) domain.SignerModel {
	switch coin {
	case "btc":
		return btc.NewSigner(logger, registry.ValueOf("btc"), accountManager.(btcDomain.UseCase))
	case "bch":
		return bch.NewSigner(logger, registry.ValueOf("bch"), accountManager.(bchDomain.UseCase))
	case "ltc":
		return ltc.NewSigner(logger, registry.ValueOf("ltc"), accountManager.(ltcDomain.UseCase))
	case "doge":
		return doge.NewSigner(logger, registry.ValueOf("doge"), accountManager.(dogeDomain.UseCase))
	case "trx":
		return trxSigner.NewSigner(logger, registry.ValueOf("trx"), accountManager.(trxDomain.TronWallet))
	case "eth":
		return signer.NewSigner(logger, registry.ValueOf("eth"), accountManager.(ethS.ETHWallet))
	case "bnb":
		return bnbSigner.NewSigner(logger, registry.ValueOf("bnb"), accountManager.(bnbDomain.BinanceWallet))
	case "xlm":
		return xlmSigner.NewSigner(logger, registry.ValueOf("xlm"), accountManager.(xlmDomain.StellarWallet))
	case "ada":
		return adaSigner.NewSigner(logger, registry.ValueOf("ada"), accountManager.(adaDomain.CardanoWallet))
	case "bsc":
		return bscSigner.NewSigner(logger, registry.ValueOf("bsc"), accountManager.(bscDomain.BSCWallet))
	case "eos":
		return eosSigner.NewSigner(logger, registry.ValueOf("eos"), accountManager.(eosDomain.EosWallet))
	case "matic":
		return maticSigner.NewSigner(logger, registry.ValueOf("matic"), accountManager.(maticDomain.MaticWallet))
	case "dot":
		return dotSigner.NewSigner(logger, registry.ValueOf("dot"), accountManager.(dotDomain.DotWallet))
	case "xrp":
		return xrpSigner.NewSigner(logger, registry.ValueOf("xrp"), accountManager.(xrpDomain.XrpWallet))
	case "pmn":
		return pmnSigner.NewSigner(logger, registry.ValueOf("pmn"), accountManager.(pmnDomain.KuknosWallet))
	case "sol":
		return solSigner.NewSigner(logger, registry.ValueOf("sol"), accountManager.(solDomain.SolWallet))
	case "avax":
		return avaxSigner.NewSigner(logger, registry.ValueOf("avax"), accountManager.(avaxDomain.AvaxWallet))
	case "luna":
		return lunaSigner.NewSigner(logger, registry.ValueOf("luna"), accountManager.(lunaDomain.TerraWallet))
	case "ftm":
		return fantomSigner.NewSigner(logger, registry.ValueOf("ftm"), accountManager.(fantomDomain.FantomWallet))
	case "atom":
		return atomSigner.NewSigner(logger, registry.ValueOf("atom"), accountManager.(atomDomain.CosmosWallet))
	case "near":
		return nearSigner.NewSigner(logger, registry.ValueOf("near"), accountManager.(nearDomain.NearWallet))
	case "hedera":
		return nil
	}
	return nil
}

func DecideAccountManagerUseCase(coin string, registry configcore.Registry,
	secureConfig config.SecureConfig, logger logger.Logger) domain.AccountManagerModel {
	switch coin {
	case "btc":
		return btc.NewAccountUseCase(registry.ValueOf("btc"), secureConfig, logger)
	case "ltc":
		return ltc.NewAccountUseCase(registry.ValueOf("ltc"), secureConfig, logger)
	case "bch":
		return bch.NewAccountUseCase(registry.ValueOf("bch"), secureConfig, logger)
	case "doge":
		return doge.NewAccountUseCase(registry.ValueOf("doge"), secureConfig, logger)
	case "trx":
		return trx.NewAccountUseCase(registry.ValueOf("trx"), secureConfig, logger)
	case "eth":
		return eth.NewAccountUseCase(registry.ValueOf("eth"), secureConfig, logger)
	case "bnb":
		return bnb.NewAccountUseCase(registry.ValueOf("bnb"), secureConfig, logger)
	case "xlm":
		return xlm.NewAccountUseCase(registry.ValueOf("xlm"), secureConfig, logger)
	case "ada":
		return ada.NewCardanoWalletUseCase(registry.ValueOf("ada"), secureConfig, logger)
	case "bsc":
		return bsc.NewAccountUseCase(registry.ValueOf("bsc"), secureConfig, logger)
	case "eos":
		return eos.NewAccountUseCase(registry.ValueOf("eos"), secureConfig, logger)
	case "matic":
		return matic.NewAccountUseCase(registry.ValueOf("matic"), secureConfig, logger)
	case "dot":
		return dot.NewAccountUseCase(registry.ValueOf("dot"), secureConfig, logger)
	case "xrp":
		return xrp.NewAccountUseCase(registry.ValueOf("xrp"), secureConfig, logger)
	case "pmn":
		return pmn.NewAccountUseCase(registry.ValueOf("pmn"), secureConfig, logger)
	case "sol":
		return sol.NewAccountUseCase(registry.ValueOf("sol"), secureConfig, logger)
	case "avax":
		return avax.NewAccountUseCase(registry.ValueOf("avax"), secureConfig, logger)
	case "luna":
		return luna.NewAccountUseCase(registry.ValueOf("luna"), secureConfig, logger)
	case "ftm":
		return ftm.NewAccountUseCase(registry.ValueOf("ftm"), secureConfig, logger)
	case "atom":
		return atom.NewAccountUseCase(registry.ValueOf("atom"), secureConfig, logger)
	case "near":
		return near.NewAccountUseCase(registry.ValueOf("near"), secureConfig, logger)
	case "hedera":
		return hederaaccount.NewUseCase(registry.ValueOf("hedera"), secureConfig, "hedera", logger)
	}
	return nil
}
