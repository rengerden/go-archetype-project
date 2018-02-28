package main

import (
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	"log/syslog"
	"sync"
	"log"
	"net/http"
)

var l logger.Logger
var ctx *Ctx

type Ctx struct {
	cache          Cache
	handlers       []*Handler
	cfg            Config
	mu             sync.Mutex
	currentHandler *Handler
	sem            chan struct{}
}

var concurrency = 5

func (c *Ctx) ChooseHandler() *Handler {
	c.mu.Lock()
	h := c.currentHandler
	for {
		if h.isAvailable(true) {
			break
		}
		h = h.next
		if h == c.currentHandler { // full cycle, we should stop
			break
		}
	}
	c.currentHandler = h
	c.mu.Unlock()
	return h
}

func (c *Ctx) ResolveCountry(ip string) (country string, ok bool) {
	if c.cfg.Concurrency > 0 {
		c.sem <- struct{}{}
	}

	country, ok = c.cache.Get(ip)
	if !ok {
		h := c.ChooseHandler()
		country, ok = h.GetCountry(ip)
		if ok {
			c.cache.Set(ip, country)
		}
	}

	if c.cfg.Concurrency > 0 {
		<-c.sem
	}
	return
}

func (c *Ctx) InitializeProvHandlers() {
	var prev *Handler
	for _, pName := range c.cfg.Providers {
		p, ok := requesterMap[pName]
		if ok {
			h := newHandler(p, c.cfg.LimitRPM)
			if prev != nil {
				prev.next = h
			}
			prev = h
			c.handlers = append(c.handlers, h)
		}
	}
	if prev != nil {
		prev.next = c.handlers[0]
		c.currentHandler = c.handlers[0]
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	country, ok := ctx.ResolveCountry(r.RemoteAddr)
	if ok {
		w.Write([]byte(country))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Unknown"))
	}
}

func main() {
	var err error
	l, _ = logger.NewLogger( syslog.Priority(logger.L_INFO), "demo")

	cfg, err := GetConfig("./geoip.json")
	if err != nil {
		l.Info("GetConfig error:", err)
		return
	}
	ctx = &Ctx{
		cfg: cfg,
		sem: make(chan struct{}, concurrency),
		cache: newCache(cfg.CacheTTL),
	}
	ctx.InitializeProvHandlers()

	http.HandleFunc("/", httpHandler)
	l.Info("Listening ..")
	log.Fatal(http.ListenAndServe(":8080", nil))

	return
}