package main

import (
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	"log/syslog"
	"sync"
	"log"
	"net/http"
)

var (
	l logger.Logger

	cache              Cache
	handlers           []*ProvHandler
	cfg                Config
	mu                 sync.Mutex
	currentProvHandler *ProvHandler
)

func ChooseProvHandler() *ProvHandler {
	mu.Lock()
	h := currentProvHandler
	for {
		if h.isAvailable(true) {
			break
		}
		h = h.next
		if h == currentProvHandler { // full cycle, we should stop
			break
		}
	}
	currentProvHandler = h
	mu.Unlock()
	return h
}

func ResolveCountry(ip string) (country string, ok bool) {
	country, ok = cache.Get(ip)
	if !ok {
		h := ChooseProvHandler()
		country, ok = h.GetCountry(ip)
		if ok {
			cache.Set(ip, country)
		}
	}
	return
}

func PrepareProvHandlers() {
	var prev *ProvHandler
	for _, pName := range cfg.Providers {
		p, ok := ReqExecutors[pName]
		if ok {
			h := newHandler(p)
			if prev != nil {
				prev.next = h
			}
			prev = h
			handlers = append(handlers, h)
		}
	}
	if prev != nil {
		prev.next = handlers[0]
		currentProvHandler = handlers[0]
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	country, ok := ResolveCountry(r.RemoteAddr)
	if ok {
		w.Write([]byte(country))
	} else {
		w.Write([]byte("Unknown"))
		w.WriteHeader(404)
	}
}

func main() {
	var err error
	l, _ = logger.NewLogger( syslog.Priority(logger.L_INFO), "demo")

	cfg, err = GetConfig("./geoip.json")
	if err != nil {
		l.Info("GetConfig error:", err)
		return
	}
	cache = newCache(cfg.CacheTTL)
	PrepareProvHandlers()

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	return
}