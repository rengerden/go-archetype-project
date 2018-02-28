package main

import (
	"testing"
	"log/syslog"
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	"sync"
	"log"
	"net/http"
	"time"
	"strconv"
	"math/rand"
)

var testData = map[string]string {
	"8.8.8.8": "United States",
	"1.1.1.1": "Australia",
}

func dumbCountryResolver(ip string) string {
	return testData[ip]
}

func httpTestHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query()["ip"][0]
	l.Debug("req", r.URL.Path, ip)
	c := dumbCountryResolver(ip)
	if c != "" {
		time.Sleep((10 + time.Duration(rand.Intn(50))) * time.Millisecond)
		w.Write([]byte(c))
		return
	}
	w.WriteHeader(404)
}


func init () {
	l, _ = logger.NewLogger( syslog.Priority(logger.L_INFO), "demo")

	http.HandleFunc("/", httpTestHandler)
	l.Info("Test Listening ..")
	go func() {
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	// generate test data
	for i := 0; i < 100; i++ {
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
		CacheTTL: 5,
		LimitRPM: 100,
		Providers: []string{"test1", "test2"},
		Concurrency: 10, // 50
	}
	ctx = &Ctx{
		cfg: cfg,
		sem: make(chan struct{}, cfg.Concurrency),
		cache: newCache(cfg.CacheTTL),
	}
	ctx.InitializeProvHandlers()

	keysTestData := make([]string, 0)
	for k,_ := range testData {
		keysTestData = append(keysTestData, k)
	}

	l.Info("Start")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		for j := 0; j < len(testData); j++ {
			k := keysTestData[rand.Intn(len(testData))]
			v := testData[k]

			wg.Add(1)
			go func() {
				c, ok := ctx.ResolveCountry(k)
				if c != v && ok {
					t.Fail()
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
	l.Debug_f(`Cache miss ratio %d %%`, ctx.cache.GetMissRatio())
}

