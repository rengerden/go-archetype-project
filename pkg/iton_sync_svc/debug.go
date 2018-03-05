package iton_sync_svc

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

func printRows(rows []DataObj) {
	for _, r := range rows {
		log.Debug(">", r)
	}
}

func printRow(r DataObj) {
	log.Debug(">", r)
}

func printAll(rows []interface{}) {
	for _, r := range rows {
		q := r.(DataObj)
		log.Debug(">", r, q)
	}
}

func printParams(Params gin.Params) {
	for _, v := range Params {
		log.Info("k,v>", v.Key, v.Value)
	}
}

func QueryDB_simple(q string) error {
	rows, err := db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	log.Debug(columns)

	var values = make([]interface{}, len(columns))
	for i := range values {
		var ii interface{}
		values[i] = &ii
	}

	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			log.Info("row err >", err.Error())
			return err
		}
		for i, colName := range columns {
			var rawValue = *(values[i].(*interface{}))
			var rawType = reflect.TypeOf(rawValue)
			log.Info(colName, rawType, rawValue)
		}
	}
	return nil
}
