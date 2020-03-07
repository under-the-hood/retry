package async

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

func Benchmark(b *testing.B) {
	how := retry.How{
		strategy.Limit(5),
		strategy.Backoff(backoff.Constant(10 * time.Millisecond)),
	}
	b.Run("usual", func(b *testing.B) {
		what := func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithCancel(context.Background())
			if err := retry.DoAsync(ctx, what, how...); err != nil {
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
			if err := retry.DoAsync(ctx, what, how...); err == nil {
				b.Error("error is expected")
			}
			cancel()
			b.ReportMetric(float64(runtime.NumGoroutine()), "goroutines")
		}
	})
}
