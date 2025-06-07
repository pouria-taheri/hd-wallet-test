package usecase

import (
	"git.mazdax.tech/blockchain/hdwallet/cmd/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"github.com/fatih/structs"
	"path"
	"strconv"
)

func (h *handler) Generate(msg domain.MessageModel) {
	go h.client.Handle(msg)

	var err error
	defer func() {
		if err != nil {
			msg.Output(domain.Error(err))
		}
	}()
	args := msg.GetArgs()

	var filename = "template.yaml"
	if len(args) > 1 {
		filename = args[1]
	}
	destination := path.Join(h.app.RootDirectory, filename)
	// save config
	var cfg config.SecureConfig
	bitSize := 256
	if len(msg.GetArgs()) > 0 {
		bs, err := strconv.Atoi(msg.GetArgs()[0])
		if err != nil {
			panic("invalid bit size param.")
		}
		bitSize = bs
	}
	cfg.SetRandomDta(bitSize)
	m := structs.Map(cfg)
	m["> Tips"] = []string{
		"Only one of fields of 'mnemonic' or 'seed' is required.",
	}
	h.saveConfig(m, destination)
	dataBytes, err := h.loadConfig(destination)
	if err != nil {
		return
	}

	result := domain.NewMessage()
	result.SetType(domain.MessageTypeEnumFile)
	result.SetData(domain.NewFile(filename, dataBytes))
	msg.Output(result)
}
