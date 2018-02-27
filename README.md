# go-archetype-project

Golang archetype project with following features:
 * Makefile
 * statical code analyzers & checkers, 
 * local GOPATH and workplace
 ** dependecies got & stored locally and separately from sources
 * use of go dep to automatically find dependencies

 * stringer generator
 * logger helper with levels of logging, string formatting 

Makefile rules
* make get-deps
* make demo
* make lint



To build the project place it in ~/go/src/dev.rubitek.com/go-archetype-project/
To run
* make get-deps
* make demo
* make run-demo


Params:
{
  // cache expire TTL, [minutes]
  "cacheTTL": 5,

  // only 2 providers available, and 2 stubs for tests: test1, test2
  "providers": ["http://freegeoip.net", "http://geoip.nekudo.com"],

  // limit requests per minute per provider
  "limitRPM": 60
}


TODO:
* limit concurrency of Handler  using semaphore emulation technique
* create ReqExecutor for simulation of geoip provider