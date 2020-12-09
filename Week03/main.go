package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	shutdownServiceCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errChan := make(chan error, 3)
	defer close(errChan)
	g := errgroup.Group{}
	g.Go(func() error {
		return StartService(shutdownServiceCtx, logger, errChan)
	})
	g.Go(func() error {
		return StartMetrics(shutdownServiceCtx, logger, errChan)
	})
	go func() {
		signals := make(chan os.Signal, 10)
		signal.Notify(signals)
		for exit := false; !exit; {
			select {
			case s := <-signals:
				switch s {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
					level.Info(logger).Log("msg", "exit due to signal", "signal", s)
					exit = true
				default:
				}
			case err := <-errChan:
				level.Info(logger).Log("msg", "exit due to error", "error", err)
				exit = true
			}
		}
		cancel()
	}()

	g.Wait()
	level.Info(logger).Log("msg", "service closed")
}

func StartService(ctx context.Context, logger log.Logger, errChan chan error) error {
	server := &http.Server{
		Addr: ":8011",
	}
	level.Info(logger).Log("msg", "service start")
	go func() {
		errChan <- server.ListenAndServe()
	}()
	<-ctx.Done()
	level.Info(logger).Log("msg", "service shutdown")
	return server.Shutdown(context.Background())
}

func StartMetrics(ctx context.Context, logger log.Logger, errChan chan error) error {
	server := &http.Server{
		Addr: ":8012",
	}
	level.Info(logger).Log("msg", "metrics start")
	go func() {
		errChan <- server.ListenAndServe()
	}()
	<-ctx.Done()
	level.Info(logger).Log("msg", "service shutdown")
	return server.Shutdown(context.Background())
}
