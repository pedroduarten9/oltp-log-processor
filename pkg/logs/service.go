package logs

import (
	"context"
	"log/slog"

	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
)

type dash0LogsServiceServer struct {
	logProcessor *LogProcessor

	collogspb.UnimplementedLogsServiceServer
}

func NewServer(logProcessor *LogProcessor) collogspb.LogsServiceServer {
	return &dash0LogsServiceServer{
		logProcessor: logProcessor,
	}
}

func (l *dash0LogsServiceServer) Export(ctx context.Context, request *collogspb.ExportLogsServiceRequest) (*collogspb.ExportLogsServiceResponse, error) {
	slog.DebugContext(ctx, "Received ExportLogsServiceRequest")

	for _, resourceLogs := range request.ResourceLogs {
		l.logProcessor.ProcessLog(resourceLogs)
	}

	return &collogspb.ExportLogsServiceResponse{}, nil
}
