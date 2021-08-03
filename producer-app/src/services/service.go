package services

import (
	"io"
	"producer-app/generated"
	"producer-app/server"
	"producer-app/util/log"
	"strings"
	"sync"

	"github.com/go-numb/go-ftx/realtime"
	"google.golang.org/protobuf/types/known/timestamppb"

	"context"
	"fmt"
)

type Services struct {
	log  *log.Logger
	grpc *server.GrpcClient
}

func NewServices(log *log.Logger, grpc *server.GrpcClient) *Services {
	return &Services{
		log:  log,
		grpc: grpc,
	}
}

func (s *Services) UploadData(ctx context.Context, symbol string) error {
	s.log.Info("[service] UploadData: ", symbol)

	if s.grpc.ClientCon == nil {
		s.grpc.ConnectGrpcServer()
	}

	var errcList []<-chan error

	errCh, err := s.getTickerData(ctx, []string{symbol})
	if err != nil {
		return err
	}
	errcList = append(errcList, errCh)
	return HandleErrorChanels(errcList...)
}

func (s *Services) getTickerData(ctx context.Context, ticker []string) (<-chan error, error) {

	stream, err := s.grpc.PipelineService.NewTickerPipeline(ctx)
	if err != nil {
		return nil, err
	}

	ch := make(chan realtime.Response)
	errCh := make(chan error, 1)
	defer close(ch)
	defer close(errCh)
	defer stream.CloseAndRecv()
	go func() {
		err := realtime.Connect(ctx, ch, []string{"ticker"}, ticker, nil)
		if err != nil {
			fmt.Println("Connect error:", err)
			errCh <- err
			return
		}
	}()

	// for v := range ch {
	for {
		select {
		case <-ctx.Done():
			return errCh, nil
		default:
			v := <-ch

			// v.Results is error
			if v.Results != nil {
				errCh <- fmt.Errorf("symbol not found")
				break
			}

			res := &generated.NewTickerPipelineRequest{
				Symbol:    v.Symbol,
				Bid:       v.Ticker.Bid,
				Ask:       v.Ticker.Ask,
				BidSize:   v.Ticker.BidSize,
				AskSize:   v.Ticker.AskSize,
				Last:      v.Ticker.Last,
				Timestamp: timestamppb.New(v.Ticker.Time.Time),
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
	}
}

func HandleErrorChanels(errs ...<-chan error) error {
	errc := mergeError(errs...)
	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}

func mergeError(listErr ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error, len(listErr))

	wg.Add(len(listErr))
	for _, c := range listErr {
		go func(c <-chan error) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
