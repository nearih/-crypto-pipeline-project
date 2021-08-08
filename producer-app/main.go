package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"producer-app/config"
	"producer-app/server"
	"producer-app/src/handlers"
	grpcRepo "producer-app/src/repository/grpcrepo"
	"producer-app/src/services"
	"producer-app/util/log"
	"sync"
)

func main() {
	// create config instance
	conf := config.NewConfig()

	// create logger instance
	logger := log.NewLogger(conf)

	// create global context
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// create grpc connection instance
	grpcConn := server.NewGRPCClient(conf, logger)

	// create repository instance
	grpcRepo := grpcRepo.NewGrpcRepo(logger, grpcConn)

	// create services instance
	serv := services.NewServices(logger, grpcConn, wg, grpcRepo)

	// create http server
	httpServer := server.NewHttpServer()
	handlers.NewHandlers(ctx, httpServer, serv)

	go httpServer.Start(conf.Server.Port)

	// gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("received shutdown signal")
	fmt.Println("start gracefully shutdown process")

	cancel()
	httpServer.Stop()
}
