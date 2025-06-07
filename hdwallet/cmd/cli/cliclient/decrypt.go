package cliclient

import "git.mazdax.tech/blockchain/hdwallet/cmd/domain"

func (c *client) Decrypt(msg domain.MessageModel) {
	c.client.Decrypt(msg)
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
