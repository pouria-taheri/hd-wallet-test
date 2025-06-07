package signer

import (
	"encoding/json"
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/xrp/crypto"
	"git.mazdax.tech/blockchain/hdwallet/coins/xrp/domain"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/xana-rahmani/ripple/data"
)

type Signer struct {
	wallet domain.XrpWallet
	logger logger.Logger
}

type SignMsg struct {
	Transaction *data.Payment `json:"transaction"`
	Account     ad.Account    `json:"account"`
}

type signerConf struct {
}

func (s Signer) Coin() string {
	return "xrp"
}

func (s Signer) SignTransaction(request []byte) ([]byte, error) {

	var signMsg SignMsg
	err := json.Unmarshal(request, &signMsg)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
		}).ErrorF("error in unmarshal transaction request")
		return nil, err
	}

	dPath := pathFromUserID(signMsg.Account.Id)
	key, err := s.wallet.AccPrv(dPath)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
			"dPath":     dPath,
		}).ErrorF("error in get keyBag from path")
		return nil, err
	}

	public, errPub := key.PublicKey()
	if errPub != nil {
		return nil, errPub
	}

	xrpKey, errPrv := crypto.NewEd25519KeyFromPrivate(append(key.Key, public...))
	if errPrv != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     errPrv,
		}).FatalF("cannot create ed25519 key from private key")
	}
	address, _ := crypto.NewAccountId(xrpKey.Id(nil))
	fmt.Println(address)

	s.logger.With(logger.Field{
		"submodule": "sign transaction",
		"payload":   signMsg,
	}).InfoF("signing transaction...")

	err = data.Sign(signMsg.Transaction, xrpKey, nil)

	res, err := json.Marshal(signMsg.Transaction)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewSigner(logger logger.Logger, conf configcore.Registry, wallet domain.XrpWallet) baseDomain.SignerModel {
	sc := new(signerConf)
	_ = conf.Unmarshal(sc)
	s := Signer{
		wallet,
		logger,
	}
	return s
}

func pathFromUserID(uid uint64) string {
	return fmt.Sprintf("m/44'/144'/%d'/0'/0'", uid)
}
