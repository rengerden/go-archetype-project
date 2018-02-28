package main

import (
	"time"
	"sync"
)

type Handler struct {
	mu         sync.Mutex
	deadline   time.Time
	counterRPM int

	next      *Handler
	requester Requester // delegate

	limitRPM  int
	//sem       chan struct{}
}

func newHandler(e Requester, limitRPM int) *Handler {
	const concurrencyLevel = 2

	return &Handler{
		requester: e,
		limitRPM: limitRPM,
		//sem: make(chan struct{}, concurrencyLevel),
	}
}

func (h *Handler) isAvailable(countReq bool) (res bool) {
	h.mu.Lock()
	if time.Now().After(h.deadline) {
		h.deadline = time.Now().Add(1 * time.Minute)
	} else {
		if h.counterRPM < h.limitRPM {
			if countReq {
				h.counterRPM++
			}
			res = true
		}
	}
	h.mu.Unlock()
	return
}

func (h *Handler) GetCountry(ip string) (string, bool) {
	if !h.isAvailable(false) {
		return "", false
	}

	//h.sem <- struct{}{}
	res, ok := h.requester.GetCountry(ip)
	//<- h.sem

	return res, ok
}