package domain

import "git.mazdax.tech/blockchain/hdwallet/coins/ltc/btcd/wire"

type SignResponse struct {
	Tx *wire.MsgTx `json:"tx"`
}
