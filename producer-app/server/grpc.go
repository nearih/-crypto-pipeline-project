package server

import (
	"context"
	"fmt"
	"producer-app/config"
	"producer-app/generated"
	"producer-app/util/log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type GrpcClient struct {
	config          *config.RootConfig
	logger          *log.Logger
	ClientCon       *grpc.ClientConn
	PipelineService generated.PipelineServiceClient
}

// ClientConn connect with timeout
func (g *GrpcClient) ClientConn(host string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cc, err := grpc.DialContext(
		ctx,
		host,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithStreamInterceptor(nil),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             100 * time.Millisecond,
			PermitWithoutStream: true,
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("connect %v error %v", host, err)
	}

	g.ClientCon = cc

	return cc, nil
}

//NewGRPC create grpc instance
func NewGRPCClient(c *config.RootConfig, log *log.Logger) *GrpcClient {
	g := &GrpcClient{
		config: c,
		logger: log,
	}
	return g
}

func (g *GrpcClient) ConnectGrpcServer() {
	serverConn, err := g.ClientConn(g.config.Endpoints.ServerEndpoint)
	if err != nil {
		g.logger.Error("Init Connect %v error %v", g.config.Endpoints.ServerEndpoint, err)
	}
	g.PipelineService = generated.NewPipelineServiceClient(serverConn)
}
