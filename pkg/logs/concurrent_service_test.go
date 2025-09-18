package logs

import (
	"context"
	"testing"
)

func BenchmarkConcurrentExport_9kBatch(b *testing.B) {
	key := "foo"
	ctx := context.Background()
	service := NewConcurrentServer(NewLogProcessor(key, NewRepo()))

	for i := 0; i < b.N; i++ {
		logs := generateLog(key, 9000)
		service.Export(ctx, logs)
	}
}
