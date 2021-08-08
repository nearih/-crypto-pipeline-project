package services

import (
	"producer-app/server"
	"producer-app/src/model"
	grpcRepo "producer-app/src/repository/grpcrepo"
	"producer-app/util/log"
	"sync"

	"github.com/go-numb/go-ftx/realtime"

	"context"
	"fmt"
)

type Services struct {
	log      *log.Logger
	wg       *sync.WaitGroup
	grpcRepo *grpcRepo.GrpcRepo
}

func NewServices(log *log.Logger, grpc *server.GrpcClient, wg *sync.WaitGroup, grpcRepo *grpcRepo.GrpcRepo) *Services {
	return &Services{
		log:      log,
		wg:       wg,
		grpcRepo: grpcRepo,
	}
}

// UploadData upload data to pipeline
func (s *Services) UploadData(ctx context.Context, symbol string) error {
	s.log.Info("[service] UploadData: ", symbol)

	var errChanList []<-chan error

	data, errCh, err := s.getTickerData(ctx, []string{symbol})
	if err != nil {
		return err
	}
	errChanList = append(errChanList, errCh)

	errCh, err = s.grpcRepo.SendDataGrpcStream(ctx, data)
	if err != nil {
		return err
	}
	errChanList = append(errChanList, errCh)

	fmt.Println("closed connection with pipeline")
	return HandleErrorChanels(errChanList...)
}

func (s *Services) getTickerData(ctx context.Context, ticker []string) (chan model.Ticker, chan error, error) {
	s.wg.Add(1)

	ch := make(chan realtime.Response)
	data := make(chan model.Ticker)
	errCh := make(chan error, 1)

	go func() {
		err := realtime.Connect(ctx, ch, []string{"ticker"}, ticker, nil)
		if err != nil {
			fmt.Println("Connect error:", err)
			errCh <- err
			return
		}

		defer s.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				v := <-ch
				if v.Results != nil {
					errCh <- fmt.Errorf("symbol not found")
					return
				}

				res := model.Ticker{
					Symbol:    v.Symbol,
					Bid:       v.Ticker.Bid,
					Ask:       v.Ticker.Ask,
					BidSize:   v.Ticker.BidSize,
					AskSize:   v.Ticker.AskSize,
					Last:      v.Ticker.Last,
					Timestamp: v.Ticker.Time.Time,
				}
				data <- res
			}
		}
	}()

	go func() {
		s.wg.Wait()
		close(errCh)
		close(data)
		close(ch)
	}()

	return data, errCh, nil
}

// HandleErrorChanels read all error from giving channel and return
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
