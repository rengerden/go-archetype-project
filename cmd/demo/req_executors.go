package main

import (
	"encoding/json"
	"net/http"
)

type ReqExecutor interface {
	GetCountry(ip string) (string, bool)
}

var ReqExecutors = map[string]ReqExecutor{
	"nekudo": ReqExecutorImplA{},
	"freegeoip": ReqExecutorImplB{},
}

type ReqExecutorImplA struct {}
type ReqExecutorImplB struct {}

func (r ReqExecutorImplA) GetCountry(ip string) (string, bool) {
	v := struct {
		City bool
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
	return v.Country.Name, true
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
	return v.Country_name, true
}