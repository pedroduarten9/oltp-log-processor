package logs

import (
	"context"
	"log/slog"

	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
)

type dash0LogsConcurrentServiceServer struct {
	logProcessor *LogProcessor

	collogspb.UnimplementedLogsServiceServer
}

func NewConcurrentServer(logProcessor *LogProcessor) collogspb.LogsServiceServer {
	return &dash0LogsConcurrentServiceServer{
		logProcessor: logProcessor,
	}
}

func (l *dash0LogsConcurrentServiceServer) Export(ctx context.Context, request *collogspb.ExportLogsServiceRequest) (*collogspb.ExportLogsServiceResponse, error) {
	slog.DebugContext(ctx, "Received ExportLogsServiceRequest")

	for _, resourceLogs := range request.ResourceLogs {
		go l.logProcessor.ProcessLog(resourceLogs)
	}

	return &collogspb.ExportLogsServiceResponse{}, nil
}
