package main

import (
	"testing"
)

func Test_Main(t *testing.T) {

	R := []ReqExecutor {ReqExecutorImplA{}, ReqExecutorImplB{}}
	for _, r := range R {
		res, _ := r.GetCountry("8.8.8.8")
		if res != "United States" {
			t.Fail()
		}
	}
}

