package domain

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/ltc/btcd/util"
	"git.mazdax.tech/blockchain/hdwallet/coins/ltc/btcd/wire"
	bd "git.mazdax.tech/blockchain/hdwallet/coins/ltc/domain"
)

type SignRequest struct {
	Tx               *wire.MsgTx            `json:"tx"`
	InputValues      []util.Amount          `json:"inputValues"`
	InputAddresses   []*bd.Address          `json:"inputAddresses"`
	ExtractedScripts []bd.ExtractedPKScript `json:"extractedScripts"`
	PrevScripts      [][]byte               `json:"prevScripts"`
}
