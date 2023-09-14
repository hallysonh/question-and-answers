package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"question-and-answers/pkg/config"
	"question-and-answers/pkg/database"
	protocolGrpc "question-and-answers/pkg/protocol/grpc"
	protocolHttp "question-and-answers/pkg/protocol/http"
)

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			slog.Error("api server recovery error -> %v", rec)
		}
	}()

	appConfig := *config.LoadAppConfiguration()
	config.InitLogger(appConfig)

	db, err := database.OpenConnection(appConfig)
	if err != nil {
		slog.Error("failed to open database: %v", err)
		return
	}

	httpServer := protocolHttp.NewApiServer(db)
	grpcServer := protocolGrpc.NewApiServer(db)

	// graceful shutdown API servers
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range quit {
			// sig is a ^C, handle it
			_ = httpServer.Stop()
			_ = grpcServer.Stop()
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	// Run HTTP API Server
	go func() {
		if err := httpServer.Start(context.Background(), appConfig.RestServerPort); err != nil {
			slog.Error(fmt.Sprintf("HTTP Server error: %+v", err.Error()))
		}
		wg.Done()
	}()

	// Run GRPC API Server
	go func() {
		if err := grpcServer.Start(context.Background(), appConfig.GRPCServerPort); err != nil {
			slog.Error(fmt.Sprintf("GRPC Server error: %+v", err.Error()))
		}
		wg.Done()
	}()

	wg.Wait()
}
