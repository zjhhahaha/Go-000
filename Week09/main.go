package main

import (
	"context"
	"demo/manager"
	"demo/server"
	"fmt"
	xlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"golang.org/x/sync/errgroup"
)

func serve() error {
	logger := log.NewJSONLogger(os.Stdout)
	manager := manager.New(logger)
	s, err := server.New(logger, "tcp4", "127.0.0.1:8011", manager)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, ictx := errgroup.WithContext(ctx)
	g.Go(func() error {
		go func() {
			select {
			case <-ictx.Done():
				manager.Close()
				s.Close()
			}
		}()
		return s.Run()
	})
	g.Go(func() error {
		signals := make(chan os.Signal)
		signal.Notify(signals, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		select {
		case s := <-signals:
			return fmt.Errorf("exit due to signal; siganl: %v", s)
		case <-ictx.Done():
			return nil
		}
	})
	err = g.Wait()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	xlog.Fatalf("Exist; cause:%v", serve())
}
