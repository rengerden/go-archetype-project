package main

import (
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	"log/syslog"
	"sync"
	"log"
	"net/http"
	"sync/atomic"
)

var l logger.Logger
var ctx *Ctx

type Ctx struct {
	cache          Cache
	cfg            Config
	mu             sync.Mutex

	handlers       []*Handler
	currentHandler *Handler
	sem            chan struct{}

	countReqTotal uint64
	countReqFail  uint64
}

func (c *Ctx) ChooseHandler() *Handler {
	c.mu.Lock()
	h := c.currentHandler
	for {
		if h.isAvailable(true) {
			break
		}

		//l.Debug("ChooseHandler > choose next", atomic.LoadUint64(&c.countReqTotal))
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
	atomic.AddUint64(&c.countReqTotal, 1)
	if c.cfg.LimitConcurrency > 0 {
		c.sem <- struct{}{}
	}

	country, ok = c.cache.Get(ip)
	if !ok {
		h := c.ChooseHandler()
		country, ok = h.GetCountry(ip)
		if ok {
			c.cache.Set(ip, country)
		} else {
			atomic.AddUint64(&c.countReqFail, 1)
		}
	}

	if c.cfg.LimitConcurrency > 0 {
		<-c.sem
	}
	return
}

func (c *Ctx) InitializeProvHandlers() {
	var prev *Handler
	for id, pName := range c.cfg.Providers {
		p, ok := requesterMap[pName]
		if ok {
			h := newHandler(p, c.cfg.LimitRPP, id+1, c.cfg.PeriodMs)
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

func newContext(cfg Config) *Ctx{
	return &Ctx{
		cfg: cfg,
		sem: make(chan struct{}, cfg.LimitConcurrency),
		cache: newCache(cfg.CacheTTL),
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
	ctx = newContext(cfg)
	ctx.InitializeProvHandlers()

	http.HandleFunc("/", httpHandler)
	l.Info("Listening ..")
	log.Fatal(http.ListenAndServe(":8080", nil))
	return
}