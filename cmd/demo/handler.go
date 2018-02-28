package main

import (
	"time"
	"sync"
)

type Handler struct {
	id int

	mu         sync.Mutex
	deadline   time.Time
	counterRPM int
	limitRPP   int
	periodMS   int

	next      *Handler
	requester Requester // delegate
}

func newHandler(e Requester, limitRPP int, id int, periodMS int) *Handler {
	return &Handler{
		id:        id,
		requester: e,
		limitRPP:  limitRPP,
		periodMS:  periodMS,
	}
}

func (h *Handler) isAvailable(countReq bool) (res bool) {
	h.mu.Lock()
	if time.Now().After(h.deadline) {
		//l.Debug("reset counter", h.id, h.counterRPM)
		h.deadline = time.Now().Add(time.Duration(h.periodMS) * time.Millisecond)
		h.counterRPM = h.limitRPP
		res = true
	} else {
		if countReq {
			//l.Debug("check", h.id, h.counterRPM, h.limitRPP, h.counterRPM < h.limitRPP)
		}
		if h.counterRPM > 0 {
			if countReq {
				h.counterRPM--
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
	res, ok := h.requester.GetCountry(ip)
	return res, ok
}