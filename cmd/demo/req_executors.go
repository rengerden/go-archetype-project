package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type ReqExecutor interface {
	GetCountry(ip string) (string, bool)
}

var ReqExecutors = map[string]ReqExecutor{
	"http://geoip.nekudo.com": ReqExecutorImplA{},
	"http://freegeoip.net": ReqExecutorImplB{},
	"test1": ReqExecutorImplTest1{},
	"test2": ReqExecutorImplTest2{},
}

type ReqExecutorImplA struct {}
type ReqExecutorImplB struct {}
type ReqExecutorImplTest1 struct {}
type ReqExecutorImplTest2 struct {}

func (r ReqExecutorImplA) GetCountry(ip string) (string, bool) {
	v := struct {
		Country struct{Name string}
	}{}
	url := "http://geoip.nekudo.com/api/" + ip + "/en"
	resp, err := http.Get(url)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&v)

	return v.Country.Name, resp.StatusCode == 200
}

func (r ReqExecutorImplB) GetCountry(ip string) (string, bool) {
	v := struct {
		Country_name string
	}{}
	url := "http://freegeoip.net/json/" + ip
	resp, err := http.Get(url)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&v)

	return v.Country_name, resp.StatusCode == 200
}

func (r ReqExecutorImplTest1) GetCountry(ip string) (string, bool) {
	time.Sleep(100 * time.Millisecond)
	if ip == "1.1.1.1" {
		return "Australia", true
	} else {
		return "United States", true
	}
}

func (r ReqExecutorImplTest2) GetCountry(ip string) (string, bool) {
	time.Sleep(100 * time.Millisecond)
	if ip == "1.1.1.1" {
		return "Australia", true
	} else {
		return "United States", true
	}
}
