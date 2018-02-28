# go-archetype-project

To build the project place it in ~/go/src/dev.rubitek.com/go-archetype-project/
To run
* make get-deps
* make demo
* make run-demo

```Params:
{
  // cache expire TTL, [minutes]
  "cacheTTL": 5,
  
  // only 2 providers available, and 2 stubs for tests: test1, test2
  "providers": ["freegeoip.net", "geoip.nekudo.com"],
  
  // limit requests per period
  "limitRPP": 60
  
  // period of time, ms
  "periodMs": 250,
  
  // limit global concurrency
  "limitConcurrency": 10
}
```

Changelog:
* Limit global concurrency using semaphore emulation technique. Lower concurrency yields lowers cache misses, but higher latency
* Simulation of geoip-provider for test

Proposals / further considerations:
* ограничивать конкуррентность на бек-ендах, т.к. могут быть ограничения на одновременное число подключений у провайдера. 
* распределять запросы по провайдерам в соотв. с установленными оганичениями
* отложенный резолвинг запросова, если нет доступного резолвера.
* аггрегировать одинаковые одновременные запросы, если такое встречается на практике. 
 