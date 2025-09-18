package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"pedroduarten9/oltp-log-processor/pkg/logs"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A grpc server for the log service",
	Run: func(cmd *cobra.Command, args []string) {
		windowSeconds := viper.GetInt("windowSeconds")
		maxReceiveMessageSize := viper.GetInt("maxReceiveMessageSize")
		attrKey := viper.GetString("attrKey")
		listenAddr := viper.GetString("listenAddr")
		logLevelCfg := viper.GetString("logLevel")

		logLevel := *new(slog.Level)
		if err := logLevel.UnmarshalText([]byte(logLevelCfg)); err != nil {
			logLevel = slog.LevelInfo
		}
		logger := slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: logLevel}))
		slog.SetDefault(logger)

		repo := logs.NewRepo()
		logProcessor := logs.NewLogProcessor(attrKey, repo)
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
		start(ctx, logProcessor, logger, windowSeconds)

		slog.Debug("Starting listener", slog.String("listenAddr", listenAddr))
		listener, err := net.Listen("tcp", listenAddr)
		if err != nil {
			slog.Error("failed to create listener")
			return
		}

		grpcServer := grpc.NewServer(
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
			grpc.MaxRecvMsgSize(maxReceiveMessageSize),
			grpc.Creds(insecure.NewCredentials()),
		)
		collogspb.RegisterLogsServiceServer(grpcServer, logs.NewServer(logProcessor))

		reflection.Register(grpcServer)
		slog.Debug("Starting gRPC server")

		go func() {
			if err := grpcServer.Serve(listener); err != nil {
				cancel()
			}
		}()

		<-ctx.Done()

		slog.Debug("Stopping gRPC server")
		grpcServer.GracefulStop()

		_ = listener.Close()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int("windowSeconds", 30, "Window of logging of attributes in seconds")
	serveCmd.Flags().Int("maxReceiveMessageSize", 16777216, "Max receive message size for the gRPC handler")
	serveCmd.Flags().String("attrKey", "foo", "The attribute key to look for")
	serveCmd.Flags().String("listenAddr", "localhost:4317", "The listen address for the gRPC server")
	serveCmd.Flags().String("logLevel", "INFO", "The log level of the system")
}

func start(ctx context.Context, lp *logs.LogProcessor, logger *slog.Logger, windowSeconds int) {
	ticker := time.NewTicker(time.Duration(windowSeconds) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				report := lp.ReportAndReset()
				if len(report) == 0 {
					logger.Info(fmt.Sprintf("No logs found for window (%d seconds).", windowSeconds))
					break
				}
				intro := fmt.Sprintf("Reporting counts for window (%d seconds):", windowSeconds)
				logger.Info(intro)
				logger.Info(report)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
