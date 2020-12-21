package ratelimit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/corverroos/ratelimit"
	em "github.com/ellemouton/ratelimiters"
)

type result struct {
	goroutines int
	nsop       int64
}

type results map[string][]result

func TestGenResultsTable(t *testing.T) {
	period := time.Millisecond
	limit := 10

	limiters := []struct{
		name string
		provider ratelimit.RateLimiter
	}{
		{
			"github.com/corver/ratelimit.NaiveWindow",
			ratelimit.NewNaiveWindow(period, limit),
		},
		{
			"github.com/corver/ratelimit.SyncMapWindow",
			ratelimit.NewSyncMapWindow(period, limit),
		},
		{
			"github.com/ellemouton/ratelimiter.TimerWindow",
			em.NewTimerWindow(period, limit),
		},
		{
			"github.com/ellemouton/ratelimiter.ChannelWindow",
			em.NewChannelWindow(period, limit),
		},
	}

	res := make(results)

	for _, limiter := range limiters {
		for _, m := range []int{1, 4, 16, 64, 256, 1024, 4096, 16384} {
			br := testing.Benchmark(ratelimit.BenchmarkFunc(func() ratelimit.RateLimiter {
				return limiter.provider
			}, m))

			res[limiter.name] = append(res[limiter.name], result{
				goroutines: m,
				nsop:       br.NsPerOp(),
			})
		}
	}

	printResults(res)
}

func printResults(res results) error {
	var impls []string
	for k, _ := range res {
		impls = append(impls, k)
	}

	heading := "| Implementation |"
	table := "|---|"
	for i := 0; i < len(res[impls[0]]); i++ {
		heading += fmt.Sprintf(" %d |",res[impls[0]][i].goroutines)
		table += "--|"
	}

	for k, v := range res {
		table = table + "\n| "+k+" |"
		for _, r := range v {
			table += fmt.Sprintf(" %d |", r.nsop)
		}
	}

	fmt.Println("Results shown in terms of ns/op")
	fmt.Println(heading+"\n"+table)
	return nil
}