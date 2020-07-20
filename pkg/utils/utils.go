package utils

import (
	"os"
	"os/signal"
	"syscall"
)

type SigTermHandler struct {
	stop chan os.Signal
}

func NewSigTermHandler() *SigTermHandler {
	return &SigTermHandler{stop: make(chan os.Signal)}
}

func (s *SigTermHandler) Wait() error {
	signal.Notify(s.stop, os.Interrupt, syscall.SIGTERM)
	<-s.stop
	return nil
}
func (s *SigTermHandler) Stop() {
	close(s.stop)
}

func HasString(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}