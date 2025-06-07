package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/atom/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func createWalletMainAccount(key domain.MasterKey, log logger.Logger) domain.WalletMainAcc {
	wma := domain.WalletMainAcc{}

	path := hd.NewFundraiserParams(0, 118, 0)
	prv, errPrv := hd.DerivePrivateKeyForPath(key.Secret, key.ChainCode, path.String())
	if errPrv != nil {
		log.With(logger.Field{
			"submodule": "walletMainAccount",
			"error":     errPrv,
		}).FatalF("error in get main account from secret and chaincode")
	}
	wma.PrivateKey = prv

	privateKey := secp256k1.GenPrivKeyFromSecret(prv)
	addr := sdk.AccAddress(privateKey.PubKey().Address())

	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"address":   addr.String(),
	}).InfoF("Main account")

	return wma
}
