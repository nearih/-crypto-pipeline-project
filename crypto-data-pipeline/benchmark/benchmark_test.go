package benchmark

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"crypto-data-pipeline/generated"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type benchmarkServer struct {
	t generated.PipelineService_NewTickerPipelineClient
}

var (
	bs   benchmarkServer = benchmarkServer{}
	port int             = 7072
)

func TestMain(m *testing.M) {
	setUpGrpcServer()
	os.Exit(m.Run())
}

func setUpGrpcServer() {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", 7072), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := generated.NewPipelineServiceClient(conn)
	stream, err := client.NewTickerPipeline(context.Background())
	if err != nil {
		panic(err)
	}
	bs.t = stream
}

func BenchmarkNewTickerPipeline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := bs.t.Send(&generated.NewTickerPipelineRequest{
			Symbol:    "TEST/USD",
			Bid:       8,
			BidSize:   800,
			Ask:       7,
			AskSize:   700,
			Timestamp: timestamppb.New(time.Now()),
		})
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkNewTickerPipelineParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		client := generated.NewPipelineServiceClient(conn)

		stream, err := client.NewTickerPipeline(context.Background())
		if err != nil {
			panic(err)
		}
		s := stream

		for pb.Next() {
			err := s.Send(&generated.NewTickerPipelineRequest{
				Symbol:    "TEST/USD",
				Bid:       8,
				BidSize:   800,
				Ask:       7,
				AskSize:   700,
				Timestamp: timestamppb.New(time.Now()),
			})
			if err != nil {
				panic(err)
			}
		}
	})
}
