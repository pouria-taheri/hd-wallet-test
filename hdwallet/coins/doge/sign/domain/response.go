package domain

import "git.mazdax.tech/blockchain/hdwallet/coins/doge/btcd/wire"

type SignResponse struct {
	Tx *wire.MsgTx `json:"tx"`
}
