package usecase

import (
	"bytes"
	"git.mazdax.tech/blockchain/hdwallet/cmd/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"github.com/fatih/structs"
	"time"
)

func (h *handler) Encrypt(msg domain.MessageModel) {
	go h.client.Handle(msg)

	var err error
	defer func() {
		if err != nil {
			msg.Output(domain.Error(err))
		}
	}()

	h.app.SecureConfigs, err = config.LoadSecureConfigs(h.app.SecuredConfigFileDetails, h.getDecryptDataHandler(msg))
	if err != nil {
		return
	}

	var pw []byte

	for {
		msg.Output(domain.ReadPassword("enter new password: "))
		req := msg.ReadInput()
		pw = []byte(req.GetArgs()[0])
		if len(pw) == 0 {
			continue
		}
		msg.Output(domain.ReadPassword("enter the password again: "))
		req = msg.ReadInput()
		pw2 := []byte(req.GetArgs()[0])
		if len(pw2) == 0 {
			continue
		}
		if bytes.Compare(pw, pw2) == 0 {
			break
		}
		msg.Output(domain.SprintFln("passwords do not match. try again\n"))
	}
	key := config.GetValidKey(h.app.Config.Salt, pw)
	encryptedPaths := make(map[string]interface{})
	for _, sc := range h.app.SecureConfigs {
		detail := sc.GetDetail()
		if _, ok := encryptedPaths[detail.FilePath]; ok {
			continue
		}
		h.saveConfig(structs.Map(sc), detail.FilePath)
		config.EncryptFile(detail.FilePath, detail.FilePath, key)
		encryptedPaths[detail.FilePath] = nil
	}
	msg.Output(domain.SprintFln("done"))
	time.Sleep(time.Second)
	msg.Close()
}
