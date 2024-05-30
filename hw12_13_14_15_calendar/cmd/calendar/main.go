package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
)

var (
	configFile string
	logFile    *os.File
)

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

	logg := initLogger(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	storage, err := initStorage(cfg, ctx)
	if err != nil {
		panic(err)
	}

	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(logg, calendar, cfg.Server.Host, cfg.Server.Port)
	grpcServer := internalgrpc.New(logg, calendar, cfg.Grpc.Host, cfg.Grpc.Port)

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

func initStorage(cfg *config.Config, ctx context.Context) (app.Storage, error) {
	if cfg.Storage.Type == "sql" {
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%v/%s?sslmode=disable",
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.DBName,
		)
		storage := sqlstorage.New()
		err := storage.Connect(ctx, dsn)
		if err != nil {
			return nil, err
		}
	}
	storage := memorystorage.New()
	return storage, nil
}

func initLogger(cfg *config.Config) app.Logger {
	if cfg.Logger.Output != "stdout" {
		//nolint:gofumpt
		if _, err := os.Stat(cfg.Logger.Output); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(cfg.Logger.Output), 0750)
			f, err := os.Create(cfg.Logger.Output)
			if err != nil {
				fmt.Println(err)
			}
			f.Close()
		}

		//nolint:gofumpt
		logFile, err := os.OpenFile(cfg.Logger.Output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()
		return logger.New(cfg.Logger, logFile)
	}

	return logger.New(cfg.Logger, os.Stdout)
}
