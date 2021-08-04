package grpcRepo

import (
	"context"
	"io"
	"producer-app/generated"
	"producer-app/server"
	"producer-app/src/model"
	"producer-app/util/log"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcRepo struct {
	log  *log.Logger
	grpc *server.GrpcClient
}

func NewGrpcRepo(log *log.Logger, grpc *server.GrpcClient) *GrpcRepo {
	return &GrpcRepo{
		log:  log,
		grpc: grpc,
	}
}

func (s *GrpcRepo) SendDataGrpc(ctx context.Context, data chan model.Ticker) (chan error, error) {

	if s.grpc.ClientCon == nil {
		s.grpc.ConnectGrpcServer()
	}

	stream, err := s.grpc.PipelineService.NewTickerPipeline(ctx)
	if err != nil {
		return nil, err
	}

	errCh := make(chan error, 1)
	defer close(errCh)
	defer stream.CloseAndRecv()

	for v := range data {
		res := &generated.NewTickerPipelineRequest{
			Symbol:    v.Symbol,
			Bid:       v.Bid,
			Ask:       v.Ask,
			BidSize:   v.BidSize,
			AskSize:   v.AskSize,
			Last:      v.Last,
			Timestamp: timestamppb.New(v.Timestamp),
		}

		// fmt.Printf("%+v \n", res)
		err = stream.Send(res)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), context.Canceled.Error()) {
				return errCh, nil
			}
			errCh <- err
			return errCh, nil
		}
	}

	return errCh, nil
}
