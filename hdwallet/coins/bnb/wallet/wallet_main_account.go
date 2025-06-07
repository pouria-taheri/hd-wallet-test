package wallet

import (
	"fmt"
	"git.mazdax.tech/blockchain/bnb-go-sdk/keys"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/bnb/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/tendermint/tendermint/crypto"
	"strconv"
)

func createWalletMainAccount(key domain.MasterKey, log logger.Logger) domain.WalletMainAcc {
	wma := domain.WalletMainAcc{}

	dPath := pathFromUserID(0, 0)
	wma.Path = dPath
	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"dPath":     dPath,
	}).InfoF("Main account")

	addr := key.KeyManager.GetAddr()
	wma.Account = addr
	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"address":   wma.Account.String(),
	}).InfoF("Main account")

	return wma
}

func pathFromUserID(uid uint64, index uint32) string {
	str := strconv.FormatUint(uid, 10)
	return fmt.Sprintf(keys.BIP44Prefix+"%s'/0/%d", str, index)
}

func partialPathFromUserID(uid uint64, index uint32) string {
	str := strconv.FormatUint(uid, 10)
	return fmt.Sprintf("%s'/0/%d", str, index)
}

func (w binanceWallet) MasterAccPrv() crypto.PrivKey { return w.PrivKey }

func (w binanceWallet) GetKeyManager() keys.KeyManager { return w.KeyManager }

func (w binanceWallet) AccKeyManager(log logger.Logger, account ad.Account) keys.KeyManager {
	path := pathFromUserID(account.Id, account.Index)
	keyMgr, err := keys.NewKeyManagerWithSecretChainCode(w.Secret, w.ChainCode, path)
	if err != nil {
		log.With(logger.Field{
			"submodule": "walletMainAccount",
			"error":     err,
			"path":      path,
		}).ErrorF("error in get key manager")
		panic("cannot get key manager from Secret and path")
	}
	return keyMgr
}
