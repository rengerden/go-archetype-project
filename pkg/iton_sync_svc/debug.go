package iton_sync_svc

import "github.com/gin-gonic/gin"

func printRows(rows []DataObj) {
	for _,r := range rows {
		log.Debug(">", r)
	}
}

func printRow(r DataObj) {
	log.Debug(">", r)
}

func printAll(rows []interface{}) {
	for _,r := range rows {
		q := r.(DataObj)
		log.Debug(">", r, q)
	}
}

func printParams(Params gin.Params) {
	for _,v := range Params {
		log.Info("k,v>", v.Key, v.Value)
	}
}