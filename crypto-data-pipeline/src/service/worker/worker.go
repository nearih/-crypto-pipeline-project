package worker

import (
	"context"
	"crypto-data-pipeline/db/influxdb"
	"crypto-data-pipeline/src/model"
	"fmt"
	"sync"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

type WorkerService struct {
	Db *influxdb.InfluxDb
}

func NewWorkerService(db *influxdb.InfluxDb) *WorkerService {
	return &WorkerService{
		Db: db,
	}
}

func (w *WorkerService) Publish(ctx context.Context, data chan model.Ticker, wg *sync.WaitGroup) chan error {
	fmt.Println("[serv] Publish")

	errChan := make(chan error, 1)

	defer close(errChan)
	defer w.Db.DBclient.Flush()

	for v := range data {
		p := influxdb2.NewPoint("ticker",
			map[string]string{
				"symbol": v.Symbol,
			}, map[string]interface{}{
				"bid":     v.Bid,
				"ask":     v.Ask,
				"bidSize": v.BidSize,
				"askSize": v.AskSize,
				"last":    v.Last,
			},
			v.Timestamp)

		w.Db.DBclient.WritePoint(p)
	}

	return errChan
}
