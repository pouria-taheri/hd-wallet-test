package usecase

import "git.mazdax.tech/blockchain/hdwallet/cmd/domain"

type remoteHandler struct {
	app    *domain.Application
	client domain.ClientModel

	unlockedChan chan interface{}
}

func NewRemoteHandler(app *domain.Application, client domain.ClientModel) domain.HandlerModel {
	h := &remoteHandler{
		app:          app,
		client:       client,
		unlockedChan: make(chan interface{}),
	}
	return h
}

func (r *remoteHandler) LoadSecureConfigs() {
	panic("implement me")
}

func (r *remoteHandler) EnsureSecureConfigLoaded() {
	panic("implement me")
}

func (r *remoteHandler) Generate(msg domain.MessageModel) {
	panic("not implement")
}

func (r *remoteHandler) Encrypt(msg domain.MessageModel) {
	panic("not implement")
}

func (r *remoteHandler) Decrypt(msg domain.MessageModel) {
	panic("not implement")
}

func (r *remoteHandler) WaitForUnlock() {
	<-r.unlockedChan
}

func (r *remoteHandler) Unlock(msg domain.MessageModel) {
	msg.SetType(domain.MessageTypeEnumUnlock)
	r.client.Unlock(msg)

	close(r.unlockedChan)
}
