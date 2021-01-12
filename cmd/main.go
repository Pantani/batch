package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Pantani/batch/internal/api"
	"github.com/Pantani/batch/internal/config"
	"github.com/Pantani/batch/internal/db"
	"github.com/Pantani/batch/internal/worker"

	"github.com/Pantani/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	_minValue int
	_duration time.Duration
)

func init() {
	flag.IntVar(&_minValue, "m", 100, "minimal value of transactions to be batched")
	flag.DurationVar(&_duration, "t", 0, "batch timeout duration")
	flag.Parse()

	config.Init()
}

func main() {
	// create application context.
	ctx, cancel := context.WithCancel(context.Background())

	// init database.
	database, err := db.Init(db.Type(config.Configuration.Database.Type))
	if err != nil {
		logger.Fatal(err)
	}

	// create the API engine.
	engine := api.CreateAPI(
		database,
		config.Configuration.API.Mode,
		_minValue,
		_duration,
	)

	// create transaction worker if needed.
	if _duration != 0 {
		grp, _ := errgroup.WithContext(ctx)
		w := worker.New(database, _duration, _minValue)
		grp.Go(w.Run())
	}

	setupGracefulShutdown(ctx, config.Configuration.API.Port, engine)
	cancel()
}

func setupGracefulShutdown(ctx context.Context, port string, engine *gin.Engine) {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}

	defer func() {
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal("Server Shutdown", err)
		}
	}()

	signalForExit := make(chan os.Signal, 1)
	signal.Notify(signalForExit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal("Application failed", err)
		}
	}()
	logger.Info("Running application", logger.Params{"bind": port})

	stop := <-signalForExit
	logger.Info("Stop signal Received", logger.Params{"stop": stop})
	logger.Info("Waiting for all jobs to stop")
}
