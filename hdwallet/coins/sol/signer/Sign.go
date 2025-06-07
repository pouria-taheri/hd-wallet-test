package signer

import (
	"encoding/json"
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/sol/domain"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/gagliardetto/solana-go"
)

type Signer struct {
	wallet domain.SolWallet
	logger logger.Logger
}

type SignMsg struct {
	Transaction *solana.Transaction `json:"transaction"`
	Account     ad.Account          `json:"account"`
}

type signerConf struct {
}

func (s Signer) Coin() string {
	return "sol"
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
	msg := tx.Message
	messageBin, err := msg.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("unable to encode message for signing: %w", err)
	}

	signerKeys := msg.AccountKeys[0:msg.Header.NumRequiredSignatures]

	prv, err := s.wallet.AccPrv(ad.DerivationPath{Account: uint32(signMsg.Account.Id), Index: signMsg.Account.Index})
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
			"accountId": signMsg.Account.Id,
		}).ErrorF("error in get private key from account id")
		return nil, err
	}

	publicKey := solana.PublicKeyFromBytes(prv.PublicKey())

	s.logger.With(logger.Field{
		"submodule": "sign transaction",
		"payload":   signMsg,
	}).InfoF("signing transaction...")

	for _, key := range signerKeys {
		if !key.Equals(publicKey) {
			continue
		}
		signed := prv.Sign(messageBin)
		var signature solana.Signature
		copy(signature[:], signed)

		tx.Signatures = append(tx.Signatures, signature)
	}

	res, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewSigner(logger logger.Logger, conf configcore.Registry, wallet domain.SolWallet) baseDomain.SignerModel {
	sc := new(signerConf)
	_ = conf.Unmarshal(sc)
	s := Signer{
		wallet,
		logger,
	}
	return s
}
