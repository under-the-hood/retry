package research_test

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	retryV4 "github.com/kamilsk/retry/v4"
	backoffV4 "github.com/kamilsk/retry/v4/backoff"
	strategyV4 "github.com/kamilsk/retry/v4/strategy"
	retryV5 "github.com/kamilsk/retry/v5"
	backoffV5 "github.com/kamilsk/retry/v5/backoff"
	strategyV5 "github.com/kamilsk/retry/v5/strategy"
)

// BenchmarkV4/normal-4         	     100	  11716495 ns/op	         3.00 goroutines	     408 B/op	       7 allocs/op
// BenchmarkV4/worst-4          	      46	  25876666 ns/op	         4.00 goroutines	     627 B/op	      11 allocs/op
func BenchmarkV4(b *testing.B) {
	how := retryV4.How{
		strategyV4.Limit(5),
		strategyV4.Backoff(backoffV4.Constant(10 * time.Millisecond)),
	}
	b.Run("normal", func(b *testing.B) {
		what := func(uint) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithCancel(context.Background())
			if err := retryV4.Try(ctx, what, how...); err != nil {
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
			if err := retryV4.Try(ctx, what, how...); err == nil {
				b.Error(err)
			}
			cancel()
			b.ReportMetric(float64(runtime.NumGoroutine()), "goroutines")
		}
	})
}

// BenchmarkV5/normal-4         	     100	  11519804 ns/op	         3.00 goroutines	     401 B/op	       6 allocs/op
// BenchmarkV5/worst-4          	      45	  25662515 ns/op	         4.00 goroutines	    1019 B/op	      16 allocs/op
func BenchmarkV5(b *testing.B) {
	how := retryV5.How{
		strategyV5.Limit(5),
		strategyV5.Backoff(backoffV5.Constant(10 * time.Millisecond)),
	}
	b.Run("normal", func(b *testing.B) {
		what := func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithCancel(context.Background())
			if err := retryV5.DoAsync(ctx, what, how...); err != nil {
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
			if err := retryV5.DoAsync(ctx, what, how...); err == nil {
				b.Error("error is expected")
			}
			cancel()
			b.ReportMetric(float64(runtime.NumGoroutine()), "goroutines")
		}
	})
}
