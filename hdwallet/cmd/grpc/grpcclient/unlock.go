package grpcclient

import (
	"context"
	"git.mazdax.tech/blockchain/hdwallet/cmd/domain"
	hdwallet "git.mazdax.tech/blockchain/hdwallet/cmd/grpc/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

func (c *client) Unlock(msg domain.MessageModel) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	ct := hdwallet.NewHdWalletClient(c.conn)
	cmd, err := ct.Command(ctx)
	if err != nil {
		panic(err)
	}

	//go func() {
	//	for {
	//		time.Sleep(time.Second * 5)
	//
	//		pingMsg := domain.NewMessage()
	//		pingMsg.SetType(domain.MessageTypeEnumPing)
	//		if err := cmd.Send(pingMsg.GetGrpcBody()); err != nil {
	//			fmt.Println(err)
	//			return
	//		}
	//	}
	//}()

	go c.Handle(msg)

	m := msg.GetGrpcBody()
	if err := cmd.Send(m); err != nil {
		panic(err)
	}

	go func() {
		for {
			resp, err := cmd.Recv()
			if err != nil {
				if err == io.EOF {
					msg.Close()
					return
				}
				if e, ok := status.FromError(err); ok {
					switch e.Code() {
					case codes.PermissionDenied:
						c.logger.WarnF("permission denied")
					case codes.Internal:
						c.logger.ErrorF("internal error")
					case codes.Unavailable:
						// closed by server
						msg.Close()
						return
					default:
						panic(err)
					}
				} else {
					panic(err)
				}
				c.logger.WarnF("error on cmd.Recv, err: %v", err)
			}
			msg.Output(domain.GetMessageFromHdWallet(resp))
		}
	}()
	for {
		m := msg.ReadInput()
		if m == nil {
			return
		}
		if err := cmd.Send(m.GetGrpcBody()); err != nil {
			panic(err)
		}
	}
}
