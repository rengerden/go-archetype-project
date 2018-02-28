package main

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
)

type Requester interface {
	GetCountry(ip string) (string, bool)
}

var requesterMap = map[string]Requester{
	"http://geoip.nekudo.com": RequesterImplA{},
	"http://freegeoip.net":    RequesterImplB{},
	"test1":                   RequesterImplTest1{},
	"test2":                   RequesterImplTest2{},
}

type RequesterImplA struct {}
type RequesterImplB struct {}
type RequesterImplTest1 struct {}
type RequesterImplTest2 struct {}

func (r RequesterImplA) GetCountry(ip string) (string, bool) {
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

func (r RequesterImplB) GetCountry(ip string) (string, bool) {
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

////////

func (r RequesterImplTest1) GetCountry(ip string) (string, bool) {
	url := "http://localhost:8081/s1?ip=" + ip
	resp, err := http.Get(url)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), resp.StatusCode == 200
}

func (r RequesterImplTest2) GetCountry(ip string) (string, bool) {
	url := "http://localhost:8081/s2?ip=" + ip
	resp, err := http.Get(url)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), resp.StatusCode == 200
}

//func (r RequestorImplTest3) GetCountry(ip string) (string, bool) {
//	time.Sleep(100 * time.Millisecond)
//	if ip == "1.1.1.1" {
//		return "Australia", true
//	} else {
//		return "United States", true
//	}
//}