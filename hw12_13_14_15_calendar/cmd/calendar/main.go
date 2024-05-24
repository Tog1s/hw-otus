package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	logg := logger.New(cfg.Logger, os.Stdout)
	storage := memorystorage.New()
	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(logg, calendar, cfg.Server.Host, cfg.Server.Port)

	grpcServer := internalgrpc.New(logg, calendar, cfg.Grpc.Host, cfg.Grpc.Port)
	if err != nil {
		log.Fatalf("init grpc server: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	go func() {
		defer wg.Done()

		if err = grpcServer.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	go func() {
		defer wg.Done()
		<-ctx.Done()

		if err = server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		if err = grpcServer.Stop(); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}

		logg.Info("calendar has stopped...")
	}()

	wg.Wait()
}
