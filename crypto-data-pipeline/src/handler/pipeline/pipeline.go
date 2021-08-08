package pipeline

import (
	"context"
	"crypto-data-pipeline/generated"
	"crypto-data-pipeline/src/model"
	"crypto-data-pipeline/src/service/worker"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
)

type PipelineController struct {
	generated.UnimplementedPipelineServiceServer
	Ctx           context.Context
	WorkerService *worker.WorkerService
	wg            *sync.WaitGroup
}

func NewPipelineController(ctx context.Context, worker *worker.WorkerService, wg *sync.WaitGroup) *PipelineController {
	return &PipelineController{
		Ctx:           ctx,
		WorkerService: worker,
		wg:            wg,
	}
}

// UploadData act like a data formatter and push data to pipeline
func (c *PipelineController) NewTickerPipeline(stream generated.PipelineService_NewTickerPipelineServer) error {
	fmt.Println("[ctrl] UploadData")

	// create internal context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// check global context if done
	go func() {
		<-c.Ctx.Done()
		cancel()
	}()

	var errChanList []<-chan error

	// format data and create data channel
	data, errChan := c.formatData(ctx, stream)
	errChanList = append(errChanList, errChan)

	// publish to related service
	errChan = c.WorkerService.Publish(ctx, data, c.wg)
	errChanList = append(errChanList, errChan)

	fmt.Println("pipeline is closed")
	return HandleErrorChanels(errChanList...)
}

// formatData formate data before the data is send to service, format function is here because it is grpc related
func (c *PipelineController) formatData(ctx context.Context, stream generated.PipelineService_NewTickerPipelineServer) (chan model.Ticker, chan error) {

	errChan := make(chan error, 1)
	data := make(chan model.Ticker)

	c.wg.Add(1)
	defer stream.SendAndClose(&generated.NewTickerPipelineResponse{
		Success: "done",
	})

	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				in, err := stream.Recv()
				if err != nil {
					if err == io.EOF || strings.Contains(err.Error(), context.Canceled.Error()) {
						return
					}
					log.Println("UploadData err:", err)
					errChan <- fmt.Errorf("UploadData err: %v", err)
					return
				}

				serviceInput := model.Ticker{
					Symbol:    in.Symbol,
					Bid:       in.Bid,
					Ask:       in.Ask,
					BidSize:   in.BidSize,
					AskSize:   in.AskSize,
					Last:      in.Last,
					Timestamp: in.Timestamp.AsTime(),
				}
				data <- serviceInput
			}
		}
	}()

	go func() {
		c.wg.Wait()
		close(errChan)
		close(data)
	}()
	return data, errChan
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
