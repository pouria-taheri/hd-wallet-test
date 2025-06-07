package signer

import (
	"bytes"
	"encoding/gob"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/tx"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

type Signer struct {
	wallet domain.CardanoWallet
	logger logger.Logger
}

type SignMsg struct {
	Transaction *tx.TxBuilder `json:"transaction"`
	Account     ad.Account    `json:"account"`
}

type signerConf struct {
	ChainType string
}

func (s Signer) Coin() string {
	return "ada"
}

func (s Signer) SignTransaction(request []byte) ([]byte, error) {

	var signMsg SignMsg
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&signMsg)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
		}).ErrorF("error in unmarshal transaction request")
		return nil, err
	}

	tx := signMsg.Transaction

	s.logger.With(logger.Field{
		"submodule": "sign transaction",
		"payload":   signMsg,
	}).InfoF("signing transaction...")

	for _, key := range tx.Skeys {
		tx.Sign(key)
	}

	var payload bytes.Buffer
	err = gob.NewEncoder(&payload).Encode(tx)
	if err != nil {
		return nil, err
	}

	return payload.Bytes(), nil
}

func NewSigner(logger logger.Logger, conf configcore.Registry, wallet domain.CardanoWallet) baseDomain.SignerModel {
	sc := new(signerConf)
	_ = conf.Unmarshal(sc)
	s := Signer{
		wallet,
		logger,
	}
	return s
}
