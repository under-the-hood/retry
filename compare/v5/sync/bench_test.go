package sync

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/kamilsk/retry/v5"
	"github.com/kamilsk/retry/v5/backoff"
	"github.com/kamilsk/retry/v5/strategy"
)

// Benchmark/normal-4         	     100	  11619043 ns/op	         3.00 goroutines	     177 B/op	       3 allocs/op
// Benchmark/worst-4          	      37	  33070368 ns/op	         3.00 goroutines	     770 B/op	      12 allocs/op
func Benchmark(b *testing.B) {
	how := retry.How{
		strategy.Limit(5),
		strategy.Backoff(backoff.Constant(10 * time.Millisecond)),
	}
	b.Run("normal", func(b *testing.B) {
		what := func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithCancel(context.Background())
			if err := retry.Do(ctx, what, how...); err != nil {
				b.Error(err)
			}
			cancel()
			b.ReportMetric(float64(runtime.NumGoroutine()), "goroutines")
		}
	})
	b.Run("worst", func(b *testing.B) {
		what := func() error {
			time.Sleep(10 * time.Millisecond)
			return errors.New("failure")
		}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
			if err := retry.Do(ctx, what, how...); err == nil {
				b.Error("error is expected")
			}
			cancel()
			b.ReportMetric(float64(runtime.NumGoroutine()), "goroutines")
		}
	})
}
