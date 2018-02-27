package main

import (
	"time"
	"sync"
)

type ProvHandler struct {
	mu         sync.Mutex
	deadline   time.Time
	counterRPM int

	next     *ProvHandler
	executor ReqExecutor // delegate
}

func newHandler(e ReqExecutor) *ProvHandler {
	return &ProvHandler{
		executor: e,
	}
}

func (h *ProvHandler) isAvailable(countReq bool) (res bool) {
	h.mu.Lock()
	if time.Now().After(h.deadline) {
		h.deadline = time.Now().Add(1 * time.Minute)
	} else {
		if h.counterRPM < cfg.LimitRPM {
			if countReq {
				h.counterRPM++
			}
			res = true
		}
	}
	h.mu.Unlock()
	return
}

func (h *ProvHandler) GetCountry(ip string) (string, bool) {
	if !h.isAvailable(false) {
		return "", false
	}

	return h.executor.GetCountry(ip)
}