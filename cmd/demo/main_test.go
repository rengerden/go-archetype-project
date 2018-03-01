package main

import (
	"testing"
	"log/syslog"
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	"sync"
	"log"
	"net/http"
	"time"
	"math/rand"
	"strconv"
	"net"
)

var testData = map[string]string {
	"8.8.8.8": "United States",
	"1.1.1.1": "Australia",
}

func countryResolver(ip string) string {
	return testData[ip]
}

func providerSimulatorHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query()["ip"][0]
	//l.Debug("req", r.URL.Path, ip)
	c := countryResolver(ip)
	if c != "" {
		time.Sleep((10 + time.Duration(rand.Intn(50))) * time.Millisecond)
		w.Write([]byte(c))
		return
	}
	w.WriteHeader(404)
}


func init () {
	l, _ = logger.NewLogger( syslog.Priority(logger.L_INFO), "demo")

	http.HandleFunc("/", providerSimulatorHandler)
	l.Info("Provider simulator listening ..")

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	go http.Serve(listener, nil)

	// generate test data
	for i := 0; i < 10; i++ {
		testData["ip-" + strconv.Itoa(i)] = "country-" + strconv.Itoa(i)
	}
}

func Test_Requestor(t *testing.T) {
	//R := []Requestor{RequestorImplA{}, RequestorImplB{}}
	//for _, r := range R {
	//	res, _ := r.GetCountry("8.8.8.8")
	//	if res != "United States" {
	//		t.Fail()
	//	}
	//}

	R := []Requester{RequesterImplTest1{}, RequesterImplTest2{}}
	for i:= 0; i <1; i++ {
		for _, r := range R {
			for k,v := range testData {
				c, ok := r.GetCountry(k)
				if c != v && ok {
					t.Fail()
				}
			}
		}
	}
	return
}

func Test_Main(t *testing.T) {
	cfg := Config {
		CacheTTL:         5,
		LimitRPP:         40,
		PeriodMs:         250,
		Providers:        []string{"test1", "test2"},
		LimitConcurrency: 4,
	}

	ctx = newContext(cfg)
	ctx.InitializeProvHandlers()

	keysTestData := make([]string, 0)
	for k,_ := range testData {
		keysTestData = append(keysTestData, k)
	}

	l.Info("Start")
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		for j := 0; j < len(testData); j++ {
			randomIP := keysTestData[rand.Intn(len(testData))]
			country := testData[randomIP]

			wg.Add(1)
			go func() {
				c, ok := ctx.ResolveCountry(randomIP)
				if c != country && ok {
					t.Fail()
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
	l.Debug_f(`Cache miss ratio: %d %%`, ctx.cache.GetMissRatio())
	l.Debug_f(`Requests (ok / failed): %d / %d`, ctx.countReqTotal, ctx.countReqFail)
}

