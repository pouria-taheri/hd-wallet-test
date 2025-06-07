package usecase

import (
	"git.mazdax.tech/blockchain/hdwallet/cmd/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"time"
)

func (h *handler) WaitForUnlock() {
	<-h.unlockedChan
}

func (h *handler) Unlock(msg domain.MessageModel) {
	if msg == nil {
		msg = <-h.unlockRequestsChan
	}
	if h.client != nil{
		go h.client.Handle(msg)
	}

	var err error
	h.app.SecureConfigs, err = config.LoadSecureConfigs(h.app.SecuredConfigFileDetails, h.getSingleInputDecryptDataHandler(msg))
	if err != nil {
		msg.Output(domain.Error(err))
	}
	h.loadedSecureConfig = true
	msg.Output(domain.SprintFln("Successfully unlocked"))
	time.Sleep(time.Second)
	msg.Close()
	close(h.unlockedChan)
}
