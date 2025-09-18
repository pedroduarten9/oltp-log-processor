package logs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
)

func TestLogProcessor(t *testing.T) {
	key := "foo"
	processor := NewLogProcessor(key, NewRepo())

	tests := []struct {
		name           string
		log            *logspb.ResourceLogs
		expectedCount  map[string]int
		expectedReport string
	}{
		{
			name:          "success nil",
			expectedCount: map[string]int{"unknown": 1},
		},
		{
			name:           "success empty",
			expectedCount:  map[string]int{"unknown": 1},
			expectedReport: "unknown - 1\n",
			log:            &logspb.ResourceLogs{},
		},
		{
			name: "success on resource",
			log: &logspb.ResourceLogs{
				Resource: &resourcepb.Resource{
					Attributes: []*commonpb.KeyValue{
						{
							Key: key,
							Value: &commonpb.AnyValue{
								Value: &commonpb.AnyValue_StringValue{StringValue: "bar"},
							},
						},
					},
				},
			},
			expectedCount:  map[string]int{"bar": 1},
			expectedReport: "bar - 1\n",
		},
		{
			name: "success on scope",
			log: &logspb.ResourceLogs{
				ScopeLogs: []*logspb.ScopeLogs{
					{
						Scope: &commonpb.InstrumentationScope{
							Attributes: []*commonpb.KeyValue{
								{
									Key: key,
									Value: &commonpb.AnyValue{
										Value: &commonpb.AnyValue_StringValue{StringValue: "barr"},
									},
								},
							},
						},
					},
				},
			},
			expectedCount:  map[string]int{"barr": 1},
			expectedReport: "barr - 1\n",
		},
		{
			name: "success on record",
			log: &logspb.ResourceLogs{
				ScopeLogs: []*logspb.ScopeLogs{
					{
						LogRecords: []*logspb.LogRecord{
							{
								Attributes: []*commonpb.KeyValue{
									{
										Key: key,
										Value: &commonpb.AnyValue{
											Value: &commonpb.AnyValue_StringValue{StringValue: "barrr"},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedCount:  map[string]int{"barrr": 1},
			expectedReport: "barrr - 1\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor.ProcessLog(tt.log)
			report := processor.ReportAndReset()
			assert.Equal(t, tt.expectedReport, report)
		})
	}
}

func TestLogProcessor_Multiple(t *testing.T) {
	key := "foo"
	processor := NewLogProcessor(key, NewRepo())
	expectedCount := map[string]int{"bar": 3, "baz": 1, "qux": 1, "unknown": 2}

	log1 := &logspb.ResourceLogs{
		Resource: &resourcepb.Resource{
			Attributes: []*commonpb.KeyValue{
				{
					Key: key,
					Value: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_StringValue{StringValue: "bar"},
					},
				},
			},
		},
	}
	log2 := &logspb.ResourceLogs{
		ScopeLogs: []*logspb.ScopeLogs{
			{
				Scope: &commonpb.InstrumentationScope{
					Attributes: []*commonpb.KeyValue{
						{
							Key: key,
							Value: &commonpb.AnyValue{
								Value: &commonpb.AnyValue_StringValue{StringValue: "bar"},
							},
						},
					},
				},
			},
		},
	}
	log3 := &logspb.ResourceLogs{
		ScopeLogs: []*logspb.ScopeLogs{
			{
				LogRecords: []*logspb.LogRecord{
					{
						Attributes: []*commonpb.KeyValue{
							{
								Key: key,
								Value: &commonpb.AnyValue{
									Value: &commonpb.AnyValue_StringValue{StringValue: "bar"},
								},
							},
						},
					},
				},
			},
		},
	}
	processor.ProcessLog(log1)
	processor.ProcessLog(log2)
	processor.ProcessLog(log3)

	attributes := []string{"baz", "qux"}
	for _, attr := range attributes {
		log := &logspb.ResourceLogs{
			Resource: &resourcepb.Resource{
				Attributes: []*commonpb.KeyValue{
					{
						Key: key,
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_StringValue{StringValue: attr},
						},
					},
				},
			},
		}
		processor.ProcessLog(log)
	}

	for i := 0; i < 2; i++ {
		processor.ProcessLog(&logspb.ResourceLogs{})
	}

	report := processor.ReportAndReset()

	for k, v := range expectedCount {
		assert.Contains(t, report, fmt.Sprintf("%s - %d\n", k, v))
	}
}
