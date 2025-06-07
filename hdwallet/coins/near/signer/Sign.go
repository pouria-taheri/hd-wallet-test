package signer

import (
	"encoding/json"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/near/domain"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

type Signer struct {
	wallet domain.NearWallet
	logger logger.Logger
}

type SignMsg struct {
	Transaction []byte     `json:"transaction"`
	Account     ad.Account `json:"account"`
}

type signerConf struct {
}

func (s Signer) Coin() string {
	return "near"
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
	tx := signMsg.Transaction

	prv, err := s.wallet.AccPrv(ad.DerivationPath{Account: uint32(signMsg.Account.Id), Index: signMsg.Account.Index})
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
			"accountId": signMsg.Account.Id,
		}).ErrorF("error in get private key from account id")
		return nil, err
	}

	s.logger.With(logger.Field{
		"submodule": "sign transaction",
		"payload":   signMsg,
	}).InfoF("signing transaction...")

	return prv.Ed25519PrivKey.Sign(tx), nil
}

func NewSigner(logger logger.Logger, conf configcore.Registry, wallet domain.NearWallet) baseDomain.SignerModel {
	sc := new(signerConf)
	_ = conf.Unmarshal(sc)
	s := Signer{
		wallet,
		logger,
	}
	return s
}
