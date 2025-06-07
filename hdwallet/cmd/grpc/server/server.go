package server

import (
	"crypto/tls"
	"fmt"
	"git.mazdax.tech/blockchain/hdwallet/cmd/domain"
	hdwallet "git.mazdax.tech/blockchain/hdwallet/cmd/grpc/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"net"
)

var unsupportedMsg = &domain.Message{
	Type:    domain.MessageTypeEnumError,
	Message: "this type of message is not currently supported",
}

type Model interface {
	hdwallet.HdWalletServer
	Start()
	Stop()
}

type server struct {
	hdwallet.UnimplementedHdWalletServer
	logger  logger.Logger
	handler domain.HandlerModel
	config  config.Cmd

	grpcServer *grpc.Server
}

func NewServer(logger logger.Logger, handler domain.HandlerModel,
	config config.Cmd) Model {
	s := &server{
		logger:  logger,
		handler: handler,
		config:  config,
	}
	return s
}

func (s *server) loadTLSCredentials() (credentials.TransportCredentials, error) {
	if s.config.Server.Grpc.TlsCertFile == "" || s.config.Server.Grpc.TlsKeyFile == "" {
		return nil, nil
	}
	serverCert, err := tls.LoadX509KeyPair(s.config.Server.Grpc.TlsCertFile,
		s.config.Server.Grpc.TlsKeyFile)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(cfg), nil
}

func (s *server) Start() {
	tlsCredentials, err := s.loadTLSCredentials()
	if err != nil {
		panic(err)
	}
	var sererOptions []grpc.ServerOption
	if tlsCredentials != nil {
		sererOptions = append(sererOptions, grpc.Creds(tlsCredentials))
	}
	s.grpcServer = grpc.NewServer(sererOptions...)
	hdwallet.RegisterHdWalletServer(s.grpcServer, s)
	l, err := net.Listen("tcp", s.config.Server.Grpc.Address)
	if err != nil {
		panic(err)
	}
	s.logger.InfoF("grpc server listening to %s", s.config.Server.Grpc.Address)
	if err = s.grpcServer.Serve(l); err != nil {
		panic(err)
	}
}

func (s *server) Stop() {
	s.grpcServer.Stop()
}

func (s *server) Handle(msg domain.MessageModel) {
	switch msg.GetType() {
	case domain.MessageTypeEnumUnlock:
		s.handler.Unlock(msg)
	case domain.MessageTypeEnumPing:
		pongMsg := domain.NewMessage()
		pongMsg.SetType(domain.MessageTypeEnumPong)
		msg.Output(pongMsg)
	default:
		msg.Output(unsupportedMsg)
	}
}

func (s *server) Command(commandServer hdwallet.HdWallet_CommandServer) error {
	request, err := commandServer.Recv()
	if err != nil {
		fmt.Println(fmt.Sprintf("err: %+v", err))
		if err == io.EOF {
			return nil
		}
		return err
	}
	msg := domain.GetMessageFromHdWallet(request)
	go s.Handle(msg)

	go func() {
		for {
			req, err := commandServer.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				if e, ok := status.FromError(err); ok {
					switch e.Code() {
					case codes.Canceled:
						return
					case codes.Unavailable:
						return
					}
				}
				fmt.Println(fmt.Sprintf("error on commandServer.Recv, err: %v", err))
				return
			}
			msg.Input(domain.GetMessageFromHdWallet(req))
		}
	}()

	for {
		m := msg.ReadOutput()
		if m == nil {
			return nil
		}
		if err := commandServer.Send(m.GetGrpcBody()); err != nil {
			if err == io.EOF {
				return nil
			}
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Canceled:
					return nil
				case codes.Unavailable:
					return nil
				}
			}
			fmt.Println(fmt.Sprintf("error on commandServer.Send, err: %v", err))
			return err
		}
	}
}
