package main

import (
	"context"
	"os"
	"os/signal"
	"producer-app/config"
	"producer-app/server"
	"producer-app/src/handlers"
	"producer-app/src/services"
	"producer-app/util/log"
)

func main() {
	// create config instance
	conf := config.NewConfig()

	// create logger instance
	logger := log.NewLogger(conf)

	// create global context
	ctx, cancel := context.WithCancel(context.Background())

	// create grpc connection instance
	conn := server.NewGRPCClient(conf, logger)

	// create services instance
	serv := services.NewServices(logger, conn)

	// create http server
	httpServer := server.NewHttpServer()
	handlers.NewHandlers(ctx, httpServer, serv)

	go httpServer.Start(conf.Server.Port)

	// gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	cancel()

	httpServer.Stop()
}
