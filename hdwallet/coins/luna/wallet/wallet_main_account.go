package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/luna/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	tmsecp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
)

func createWalletMainAccount(key domain.MasterKey, log logger.Logger) domain.WalletMainAcc {
	wma := domain.WalletMainAcc{}

	path := hd.NewFundraiserParams(0, 330, 0)
	prv, errPrv := hd.DerivePrivateKeyForPath(key.Secret, key.ChainCode, path.String())
	if errPrv != nil {
		log.With(logger.Field{
			"submodule": "walletMainAccount",
			"error":     errPrv,
		}).FatalF("error in get main account from secret and chaincode")
	}
	wma.PrivateKey = prv

	privateKey := tmsecp256k1.GenPrivKeySecp256k1(prv)
	address, err := ConvertAndEncode("terra", privateKey.PubKey().Address())
	if err != nil {
		log.With(logger.Field{
			"submodule": "walletMainAccount",
			"error":     errPrv,
		}).FatalF("error in get address from private key")
	}

	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"address":   address,
	}).InfoF("Main account")

	return wma
}
