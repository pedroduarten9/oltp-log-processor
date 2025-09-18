package test

import (
	"context"
	"log"
	"net"
	"pedroduarten9/oltp-log-processor/pkg/logs"
	"testing"

	"github.com/stretchr/testify/assert"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	otellogs "go.opentelemetry.io/proto/otlp/logs/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestLogsServiceServer_Export(t *testing.T) {
	ctx := context.Background()

	client, closer := server()
	defer closer()

	type expectation struct {
		out *collogspb.ExportLogsServiceResponse
		err error
	}

	tests := map[string]struct {
		in       *collogspb.ExportLogsServiceRequest
		expected expectation
	}{
		"Must_Success": {
			in: &collogspb.ExportLogsServiceRequest{
				ResourceLogs: []*otellogs.ResourceLogs{
					{
						ScopeLogs: []*otellogs.ScopeLogs{},
						SchemaUrl: "dash0.com/otlp-log-processor-backend",
					},
				},
			},
			expected: expectation{
				out: &collogspb.ExportLogsServiceResponse{},
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.Export(ctx, tt.in)
			if err != nil {
				assert.Equal(t, tt.expected.err, err)
			} else {
				expectedPartialSuccess := tt.expected.out.GetPartialSuccess()
				partialSuccess := out.GetPartialSuccess()
				assert.Equal(t, expectedPartialSuccess.GetRejectedLogRecords(), partialSuccess.GetRejectedLogRecords())
				assert.Equal(t, expectedPartialSuccess.GetErrorMessage(), partialSuccess.GetErrorMessage())
			}

		})
	}
}

func server() (collogspb.LogsServiceClient, func()) {
	addr := "localhost:4317"
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	var err error

	baseServer := grpc.NewServer()
	logsProcessor := logs.NewLogProcessor("foo", logs.NewRepo())
	collogspb.RegisterLogsServiceServer(baseServer, logs.NewServer(logsProcessor))
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.NewClient(addr,
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := collogspb.NewLogsServiceClient(conn)

	return client, closer
}
