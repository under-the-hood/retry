package v4

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/kamilsk/retry/v4"
	"github.com/kamilsk/retry/v4/backoff"
	"github.com/kamilsk/retry/v4/strategy"
)

// Benchmark/normal-4         	     100	  11261489 ns/op	         4.00 goroutines	     420 B/op	       7 allocs/op
// Benchmark/worst-4          	      46	  25642897 ns/op	         4.00 goroutines	     611 B/op	      11 allocs/op
func Benchmark(b *testing.B) {
	how := retry.How{
		strategy.Limit(5),
		strategy.Backoff(backoff.Constant(10 * time.Millisecond)),
	}
	b.Run("normal", func(b *testing.B) {
		what := func(uint) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithCancel(context.Background())
			if err := retry.Try(ctx, what, how...); err != nil {
				b.Error(err)
			}
			cancel()
			b.ReportMetric(float64(runtime.NumGoroutine()), "goroutines")
		}
	})
	b.Run("worst", func(b *testing.B) {
		what := func(uint) error {
			time.Sleep(10 * time.Millisecond)
			return errors.New("failure")
		}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
			if err := retry.Try(ctx, what, how...); err == nil {
				b.Error(err)
			}
			cancel()
			b.ReportMetric(float64(runtime.NumGoroutine()), "goroutines")
		}
	})
}
