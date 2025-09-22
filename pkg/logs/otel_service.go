package logs

import (
	"context"
	"log/slog"
	"pedroduarten9/oltp-log-processor/pkg/otel"

	"go.opentelemetry.io/otel/metric"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
)

type dash0OtelLogsServiceServer struct {
	next            collogspb.LogsServiceServer
	config          otel.Config
	logsReceivedCnt metric.Int64Counter

	collogspb.UnimplementedLogsServiceServer
}

func NewOtelServer(next collogspb.LogsServiceServer, config otel.Config) (collogspb.LogsServiceServer, error) {
	logsReceivedCnt, err := config.Meter.Int64Counter("com.dash0.homeexercise.logs.received",
		metric.WithDescription("The number of logs received by otlp-log-processor-backend"),
		metric.WithUnit("{log}"))
	if err != nil {
		return nil, err
	}

	return &dash0OtelLogsServiceServer{
		next:            next,
		config:          config,
		logsReceivedCnt: logsReceivedCnt,
	}, nil
}

func (l *dash0OtelLogsServiceServer) Export(ctx context.Context, request *collogspb.ExportLogsServiceRequest) (*collogspb.ExportLogsServiceResponse, error) {
	ctx, span := l.config.Tracer.Start(ctx, "log")
	defer span.End()
	slog.DebugContext(ctx, "Received ExportLogsServiceRequest")

	l.logsReceivedCnt.Add(ctx, 1)
	res, err := l.next.Export(ctx, request)
	slog.DebugContext(ctx, "Finished ExportLogsServiceRequest")
	return res, err
}
