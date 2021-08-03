package main

import (
	"context"
	"crypto-data-pipeline/config"
	"crypto-data-pipeline/db/influxdb"
	grpcServer "crypto-data-pipeline/server/grpc"
	"crypto-data-pipeline/src/handler"
	"crypto-data-pipeline/src/handler/pipeline"
	workerService "crypto-data-pipeline/src/service/worker"
	"os"
	"os/signal"
	"syscall"

	"fmt"
	"sync"
)

var terminateChan chan os.Signal

func main() {

	// handle shutdown signal
	terminateChan = make(chan os.Signal, 1)
	signal.Notify(terminateChan, os.Interrupt, syscall.SIGSEGV, syscall.SIGTERM)

	conf := config.NewConfig()
	grpcServ := grpcServer.NewGRPCServer(conf)

	// initDb
	db := influxdb.NewInfluxDb()

	// init service
	workerServ := workerService.NewWorkerService(db)

	// set up cancel
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// register handler
	clHandler := pipeline.NewPipelineController(ctx, workerServ, wg)
	handler.RegisterGrpcHandler(grpcServ.Server, clHandler)

	// start server
	go grpcServ.Start()

	<-terminateChan
	fmt.Println("received shutdown signal")
	fmt.Println("start gracefully shutdown process")

	cancel()
	grpcServ.Stop()
}
