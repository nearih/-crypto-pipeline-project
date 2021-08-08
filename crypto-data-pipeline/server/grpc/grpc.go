package grpc

import (
	"fmt"
	"log"
	"net"

	"crypto-data-pipeline/config"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	config *config.RootConfig
	Server *grpc.Server
	Error  error
}

func (gs *GRPCServer) Start() error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", gs.config.Server.Grpc.Port))

	if err != nil {
		return err
	}

	log.Println("Listening grpc on port", gs.config.Server.Grpc.Port)

	return gs.Server.Serve(listen)
}

func (gs *GRPCServer) Stop() {
	gs.Server.GracefulStop()
}

// NewGRPCServer create grpc instance
func NewGRPCServer(c *config.RootConfig) *GRPCServer {
	return &GRPCServer{
		config: c,
		Server: grpc.NewServer(),
	}
}

func (gs *GRPCServer) Err() error {
	return gs.Error
}
