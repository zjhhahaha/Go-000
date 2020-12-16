/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	pb "demo/api/account"
	"demo/conf"
	"demo/internal/pkg"
	"demo/internal/service/account"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/log/level"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func serve(cmd *cobra.Command, args []string) {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	config := conf.Load()
	db := pkg.NewRDBConnection(config.Dsn)
	accountService := account.InitializeService(db)
	app, err := initApp(config.Addr)
	if err != nil {
		level.Info(logger).Log("msg", "init app error", "error", err)
		return
	}
	pb.RegisterAccountServiceServer(app.GrpcServer(), accountService)
	shutdownServiceCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errChan := make(chan error, 3)
	defer close(errChan)
	g := errgroup.Group{}
	g.Go(func() error {
		return app.Start(shutdownServiceCtx)
	})
	go func() {
		signals := make(chan os.Signal)
		signal.Notify(signals, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		select {
		case s := <-signals:
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				level.Info(logger).Log("msg", "exit due to signal", "signal", s)
			default:
			}
		case err := <-errChan:
			level.Info(logger).Log("msg", "exit due to error", "error", err)
		}
		cancel()
	}()

	err = g.Wait()
	level.Info(logger).Log("msg", "service closed", "error", err)
}

type App struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

func initApp(addr string) (*App, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer()
	return &App{
		listener:   listener,
		grpcServer: grpcServer,
	}, nil
}

func (app *App) GrpcServer() *grpc.Server {
	return app.grpcServer
}

func (app *App) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		app.grpcServer.GracefulStop()
	}()
	return app.grpcServer.Serve(app.listener)
}
