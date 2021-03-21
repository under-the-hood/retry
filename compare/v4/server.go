package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kamilsk/retry/v4"
	"github.com/kamilsk/retry/v4/backoff"
	"github.com/kamilsk/retry/v4/strategy"
	"go.octolab.org/toolkit/protocol/http/header"
)

var port = flag.Uint("port", 8080, "listening port")

func main() {
	what := func(uint) error {
		time.Sleep(10 * time.Millisecond)
		return errors.New("failure")
	}
	how := retry.How{
		strategy.Limit(5),
		strategy.Backoff(backoff.Constant(10 * time.Millisecond)),
	}

	mux := new(http.ServeMux)
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		if err := retry.Try(req.Context(), what, how...); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), middleware(mux)))
}

func middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		timeout, _ := header.Timeout(req.Header, time.Second)
		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		handler.ServeHTTP(rw, req.WithContext(ctx))
	})
}
