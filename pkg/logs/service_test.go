package logs

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	otellogs "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
)

func TestExport(t *testing.T) {
	key := "foo"
	ctx := context.Background()
	service := NewServer(NewLogProcessor(key, NewRepo()))

	logs := generateLog(key, 9000)
	res, err := service.Export(ctx, logs)

	assert.NoError(t, err)
	assert.Equal(t, &v1.ExportLogsServiceResponse{}, res)
}

func BenchmarkExport_9kBatch(b *testing.B) {
	key := "foo"
	ctx := context.Background()
	service := NewServer(NewLogProcessor(key, NewRepo()))

	for i := 0; i < b.N; i++ {
		logs := generateLog(key, 9000)
		service.Export(ctx, logs)
	}
}

func generateLog(key string, amount int) *v1.ExportLogsServiceRequest {
	resourceLogs := make([]*otellogs.ResourceLogs, amount)
	for i := 0; i < amount; i++ {
		resourceLogs = append(resourceLogs, &otellogs.ResourceLogs{
			Resource: &resourcepb.Resource{
				Attributes: []*commonpb.KeyValue{
					{
						Key: key,
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_StringValue{StringValue: fmt.Sprintf("%d", i)},
						},
					},
				},
			},
		})

	}

	return &v1.ExportLogsServiceRequest{
		ResourceLogs: resourceLogs,
	}
}
