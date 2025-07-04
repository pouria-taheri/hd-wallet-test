package usecase

import (
	"encoding/json"
	"fmt"
	bad "git.mazdax.tech/blockchain/hdwallet/coins/bch/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/bch/bchd/txscript"
	"git.mazdax.tech/blockchain/hdwallet/coins/bch/sign/domain"
	"git.mazdax.tech/core/errors"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

type signer struct {
	logger logger.Logger

	config domain.Config

	account bad.UseCase
}

func New(logger logger.Logger, configRegistry configcore.Registry,
	account bad.UseCase) domain.UseCaseModel {
	uc := &signer{
		logger:  logger,
		account: account,
	}
	if err := configRegistry.Unmarshal(&uc.config.Config); err != nil {
		panic(err)
	}
	if err := configRegistry.Unmarshal(&uc.config); err != nil {
		panic(err)
	}
	uc.config.Initialize()

	return uc
}

func (uc *signer) Coin() string {
	return "bch"
}

func (uc *signer) Sign(request domain.SignRequest) (*domain.SignResponse, errors.ErrorModel) {
	// handle inputs
	if len(request.Tx.TxIn) != len(request.PrevScripts) {
		msg := fmt.Sprintf("tx.TxIn and prevPkScripts slices must " +
			"have equal length")
		uc.logger.With(logger.Field{
			"txId len":        len(request.Tx.TxIn),
			"prevScripts len": len(request.PrevScripts),
		}).WarnF(msg)
		return nil, errors.New().WithMessage(msg)
	}
	for i := range request.Tx.TxIn {
		prevScript := request.PrevScripts[i]
		switch {
		default:
			script, err := txscript.SignTxOutput(uc.account.GetChainParams(),
				request.ExtractedScripts[i], request.Tx, i, txscript.SigHashAll,
				uc.account, uc.account, prevScript)
			if err := errors.HandleError(err); err != nil {
				uc.logger.With(logger.Field{
					"submodule": "signer usecase",
					"section":   "SignTxOutput",
					"error":     err,
					"i":         i,
				}).WarnF("error on sign transaction txIn")
				return nil, err
			}
			request.Tx.TxIn[i].SignatureScript = script
		}
	}
	return &domain.SignResponse{
		Tx: request.Tx,
	}, nil
}

func (uc *signer) SignTransaction(request []byte) ([]byte, error) {
	var body domain.SignRequest
	if err := json.Unmarshal(request, &body); err != nil {
		return nil, err
	}
	//
	resp, err := uc.Sign(body)
	if err != nil {
		uc.logger.With(logger.Field{
			"error": err,
		}).ErrorF("error on sign transaction")
		return nil, err
	}
	respBytes, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		uc.logger.With(logger.Field{
			"error": marshalErr,
		}).ErrorF("error on marshal signed transaction")
		return nil, marshalErr
	}
	return respBytes, nil
}
