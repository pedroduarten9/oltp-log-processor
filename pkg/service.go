package logs

import (
	"context"
	"log/slog"

	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
)

type dash0LogsServiceServer struct {
	collogspb.UnimplementedLogsServiceServer
}

func NewServer() collogspb.LogsServiceServer {
	return &dash0LogsServiceServer{}
}

func (l *dash0LogsServiceServer) Export(ctx context.Context, request *collogspb.ExportLogsServiceRequest) (*collogspb.ExportLogsServiceResponse, error) {
	slog.DebugContext(ctx, "Received ExportLogsServiceRequest")
	return &collogspb.ExportLogsServiceResponse{}, nil
}
