package main

import (
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	"log/syslog"
	"sync"
)

var (
	l logger.Logger

	cache          Cache
	handlers       []*Handler
	cfg            Config
	mu             sync.Mutex
	currentHandler *Handler
)

func ChooseHandler() *Handler {
	mu.Lock()
	h := currentHandler
	for {
		if h.isAvailable(true) {
			break
		}
		h = h.next
		if h == currentHandler { // full cycle, we should stop
			break
		}
	}
	currentHandler = h
	mu.Unlock()
	return h
}

func ResolveCountry(ip string) (country string, ok bool) {
	country, ok = cache.Get(ip)
	if !ok {
		h := ChooseHandler()
		country, ok = h.GetCountry(ip)
		if ok {
			cache.Set(ip, country)
		}
	}
	return
}

func PrepareHandlers() {
	var prev *Handler
	for _, pname := range cfg.Providers {
		p, ok := ReqExecutors[pname]
		if ok {
			h := newHandler(p)
			prev.next = h
			prev = h
			handlers = append(handlers, h)
		}
	}
	if prev != nil {
		prev.next = handlers[0]
	}
}

func main() {
	l, _ = logger.NewLogger( syslog.Priority(logger.L_INFO), "demo")
	cfg, err := GetConfig("./geoip.json")
	if err != nil {
		l.Info("GetConfig error:", err)
		return
	}
	cache = newCache(cfg.CacheTTL)
	PrepareHandlers()
	return
}