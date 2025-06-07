package cliclient

import "git.mazdax.tech/blockchain/hdwallet/cmd/domain"

func (c *client) Unlock(msg domain.MessageModel) {
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
