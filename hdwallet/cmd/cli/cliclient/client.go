package cliclient

import (
	"encoding/json"
	"fmt"
	"git.mazdax.tech/blockchain/hdwallet/cmd/domain"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/term"
	"io"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

type client struct {
	app    *domain.Application
	client domain.ClientModel

	cmd      *cobra.Command
	terminal *terminal.Terminal
}

func New(app *domain.Application, cmd *cobra.Command) domain.ClientModel {
	if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
		panic("stdin/stdout should be terminal")
	}
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	t := terminal.NewTerminal(screen, "")

	h := &client{
		app:      app,
		cmd:      cmd,
		terminal: t,
	}
	return h
}

func (c *client) SetInnerClient(client domain.ClientModel) {
	c.client = client
}

func (c *client) innerClientHandleMessage(msg domain.MessageModel) domain.MessageModel {
	if c.client != nil {
		return c.client.HandleMessage(msg)
	}
	return nil
}

func (c *client) HandleMessage(msg domain.MessageModel) domain.MessageModel {
	var result domain.MessageModel
	handled := false
	switch msg.GetType() {
	case domain.MessageTypeEnumError:
		fmt.Print(msg.GetMessage())
		handled = true
	case domain.MessageTypeEnumMessage:
		fmt.Print(msg.GetMessage())
		handled = true
	case domain.MessageTypeEnumTextInput:
		c.cmd.Print(msg.GetMessage())
		input, err := c.terminal.ReadLine()
		if err != nil {
			panic(err)
		}
		m := domain.NewMessage()
		m.SetType(domain.MessageTypeEnumTextInput)
		m.Args = []string{input}
		result = m
		handled = true
	case domain.MessageTypeEnumPasswordInput:
		c.cmd.Printf(msg.GetMessage())
		pw, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			panic(err)
		}
		c.cmd.Println()
		m := domain.NewMessage()
		m.SetType(domain.MessageTypeEnumPasswordInput)
		m.Args = []string{string(pw)}
		result = m
		handled = true
	case domain.MessageTypeEnumFile:
		var file domain.File
		if err := json.Unmarshal(msg.GetData(), &file); err != nil {
			panic(err)
		}
		// save file
		destination := path.Join(c.app.RootDirectory, file.Name)
		if err := ioutil.WriteFile(destination, file.Data, 0644); err != nil {
			panic(err)
		}
		c.cmd.Printf("file saved to '%s'", destination)
		handled = true
	}
	if !handled {
		result = c.innerClientHandleMessage(msg)
	}
	return result
}

func (c *client) Handle(msg domain.MessageModel) {
	for {
		request := msg.ReadOutput()
		if request == nil {
			return
		}
		clientResponse := c.HandleMessage(request)
		if clientResponse != nil {
			msg.Input(clientResponse)
		}
	}
}
