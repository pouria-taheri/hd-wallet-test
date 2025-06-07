package signer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/eos/domain"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/eoscanada/eos-go"
)

type Signer struct {
	wallet  domain.EosWallet
	logger  logger.Logger
	chainId string
}

type SignMsg struct {
	Transaction *eos.SignedTransaction `json:"transaction"`
	Account     ad.Account             `json:"account"`
}

type signerConf struct {
	ChainId string
}

func (s Signer) Coin() string {
	return "eos"
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
	keyBag, err := s.wallet.AccKeyBag(dPath)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
			"dPath":     dPath,
		}).ErrorF("error in get keyBag from path")
		return nil, err
	}

	chainID, err := hex.DecodeString(s.chainId)
	if err != nil {
		return nil, err
	}

	s.logger.With(logger.Field{
		"submodule": "sign transaction",
		"payload":   signMsg,
	}).InfoF("signing transaction...")

	signedTx, errSign := keyBag.Sign(context.Background(), signMsg.Transaction, chainID, s.wallet.AccPublicKey())
	if errSign != nil {
		return nil, errSign
	}

	res, err := json.Marshal(signedTx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewSigner(logger logger.Logger, conf configcore.Registry, wallet domain.EosWallet) baseDomain.SignerModel {
	sc := new(signerConf)
	_ = conf.Unmarshal(sc)
	s := Signer{
		wallet,
		logger,
		sc.ChainId,
	}
	return s
}

func pathFromUserID(uid uint64) string {
	return fmt.Sprintf("m/44'/194'/%d'/0'/0'", uid)
}
