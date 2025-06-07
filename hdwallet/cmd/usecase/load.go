package usecase

import (
	"git.mazdax.tech/blockchain/hdwallet/config"
)

func (h *handler) LoadSecureConfigs() {
	var err error
	h.app.SecureConfigs, err = config.LoadSecureConfigs(h.app.SecuredConfigFileDetails, nil)
	if err != nil {
		panic(err)
	}
	h.loadedSecureConfig = true
	return
}
