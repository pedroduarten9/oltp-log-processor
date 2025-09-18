package logs

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentIncrementAttributeAndReset(t *testing.T) {
	repo := NewConcurrentRepo()

	tests := []struct {
		name               string
		attributesAndCount map[string]int
	}{
		{
			name:               "single key",
			attributesAndCount: map[string]int{"key": 1},
		},
		{
			name:               "single key multiple times",
			attributesAndCount: map[string]int{"key": 3},
		},
		{
			name:               "multiple keys",
			attributesAndCount: map[string]int{"key": 3, "key2": 1, "key3": 4},
		},
	}

	for _, tt := range tests {
		for k, v := range tt.attributesAndCount {
			for i := 0; i < v; i++ {
				repo.IncrementAttribute(k)
			}
		}

		oldMap := repo.Reset()
		assert.Equal(t, tt.attributesAndCount, oldMap)
		assert.True(t, repo.counts.IsEmpty())
	}
}

func BenchmarkConcurrentRepo_Serial(b *testing.B) {
	key := "foo"
	repo := NewConcurrentRepo()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.IncrementAttribute(key)
	}
}

func BenchmarkConcurrentRepo_Parallel(b *testing.B) {
	key := "foo"
	repo := NewConcurrentRepo()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			repo.IncrementAttribute(key)
		}
	})
}

func BenchmarkConcurrentRepo_ManyKeys(b *testing.B) {
	repo := NewConcurrentRepo()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := strconv.Itoa(i % 1000)
		repo.IncrementAttribute(val)
	}
}

func BenchmarkConcurrentRepo_WriteAndReset(b *testing.B) {
	repo := NewConcurrentRepo()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := strconv.Itoa(i % 1000)
		repo.IncrementAttribute(val)
	}

	for i := 0; i < b.N; i++ {
		_ = repo.Reset()
	}
}

func BenchmarkConcurrentRepo_ConcurrentWriteAndRead(b *testing.B) {
	key := "foo"
	repo := NewConcurrentRepo()

	b.ResetTimer()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go func() {
			repo.IncrementAttribute(key)
			wg.Done()
		}()
		go func() {
			_ = repo.Reset()
			wg.Done()
		}()
	}
	wg.Wait()
}
