package domain

type Account struct {
	Id               uint64 `json:"id"`
	Index            uint32 `json:"index"`
	Type             uint8  `json:"type"`
	Address          string `json:"address"`
	PublicKey        []byte `json:"public_key"`
	Ed25519PublicKey []byte `json:"ed_25519_public_key,omitempty"`
	Master           string `json:"master,omitempty"`
	Memo             string `json:"memo,omitempty"`
}
