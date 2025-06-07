package cliclient

import "git.mazdax.tech/blockchain/hdwallet/cmd/domain"

func (c *client) Generate(msg domain.MessageModel) {
	c.client.Generate(msg)
	for {
		m := msg.ReadInput()
		if m == nil {
			return
		}
		resp := c.HandleMessage(m)
		if resp != nil {
			msg.Output(resp)
		}
	}
}
