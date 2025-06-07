package usecase

import (
	"git.mazdax.tech/blockchain/hdwallet/cmd/domain"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

type handler struct {
	app    *domain.Application
	client domain.ClientModel

	unlockRequestsChan chan domain.MessageModel
	unlockedChan       chan interface{}

	loadedSecureConfig bool
}

func NewHandler(app *domain.Application, client domain.ClientModel) domain.HandlerModel {
	h := &handler{
		app:                app,
		client:             client,
		unlockRequestsChan: make(chan domain.MessageModel, 1),
		unlockedChan:       make(chan interface{}),
	}
	return h
}

func (h *handler) EnsureSecureConfigLoaded() {
	if !h.loadedSecureConfig {
		unlockMsg := domain.NewMessage()
		unlockMsg.SetType(domain.MessageTypeEnumUnlock)
		h.Unlock(unlockMsg)
	}
}

func (h *handler) saveConfig(cfg map[string]interface{}, destination string) {
	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.MergeConfigMap(cfg); err != nil {
		panic(err)
	}

	if err := v.WriteConfigAs(destination); err != nil {
		panic(err)
	}
}

func (h *handler) loadConfig(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	dataBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return dataBytes, nil
}
