package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.New(configFile)
	if err != nil {
		fmt.Println(err)
	}

	//nolint:gofumpt
	logFile, err := os.OpenFile(cfg.Logger.Output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	logg := logger.New(cfg.Logger, logFile)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(logg, calendar, cfg.Server.Host, cfg.Server.Port)
	grpcServer := internalgrpc.New(logg, calendar, cfg.Grpc.Host, cfg.Grpc.Port)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatalf("error closing log file %s", err)
		}
	}(logFile)

	wg := sync.WaitGroup{}
	wg.Add(3)
	logg.Info("calendar is running...")

	// Run HTTP Server
	go func() {
		defer wg.Done()
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()

	// Run GRPC Server
	go func() {
		defer wg.Done()
		if err = grpcServer.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
	}()

	// Stop all servers
	go func() {
		defer wg.Done()
		<-ctx.Done()

		if err = server.Stop(ctx); err != nil {
			log.Println("failed to stop http server: " + err.Error())
		}

		if err = grpcServer.Stop(); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()
	wg.Wait()
	logg.Info("calendar has stopped...")
}
