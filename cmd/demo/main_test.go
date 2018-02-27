package main

import (
	"testing"
	"log/syslog"
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	"sync"
)

func Test_Main(t *testing.T) {
	l, _ = logger.NewLogger( syslog.Priority(logger.L_INFO), "demo")

	R := []ReqExecutor {ReqExecutorImplA{}, ReqExecutorImplB{}}
	for _, r := range R {
		res, _ := r.GetCountry("8.8.8.8")
		if res != "United States" {
			t.Fail()
		}
	}


	cfg = Config {
		CacheTTL: 5,
		LimitRPM: 60,
		Providers: []string{"test1", "test2"},
	}
	cache = newCache(cfg.CacheTTL)
	PrepareProvHandlers()

	var wg sync.WaitGroup

	l.Info("Start")
	for i:= 0; i < 50; i ++ {
		wg.Add(1)
		go func() {
			for j:= 0; j < 4; j++ {
				c, ok := ResolveCountry("8.8.8.8")
				if c != "United States" && ok {
					t.Fail()
				}

				c, ok = ResolveCountry("1.1.1.1")
				if c != "Australia" && ok {
					t.Fail()
				}

				//l.Debug( ResolveCountry("8.8.8.8"))
				//l.Debug( ResolveCountry("1.1.1.1"))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

