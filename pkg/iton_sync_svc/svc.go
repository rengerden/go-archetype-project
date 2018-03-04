package iton_sync_svc

import (
	"time"
	"os"
	"os/signal"
	"reflect"
	"database/sql"
	"net/http"
	"bufio"
	"github.com/gin-gonic/gin"
	_ "github.com/denisenkom/go-mssqldb"
	"gopkg.in/ini.v1"
	"runtime/debug"
	"strconv"
	"dev.rubetek.com/go-archetype-project/pkg/iton_sync_svc/foo"
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	"runtime"
)

var (
	db  *sql.DB
	log logger.Logger
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type DataObj interface {
	CopyVal() DataObj
}

type TimeSheet_record struct {
	Date      	foo.JsonUnixNsToTime		`json:"date"`       // of 1st enter (w/o time) 	//  compound key  !null
	AccountId 	int				`json:"accountId"`              //  compound key, !null
	EnterTime 	foo.JsonUnixNsToTime		`json:"enterTime"`	//
	ExitTime  	foo.JsonUnixNsToTime		`json:"exitTime"`	//
	TotalTime 	int				`json:"totalTime"`	            // sum of (TExit-TLastEnter)
}

func (p TimeSheet_record) CopyVal() DataObj {
	return p
}

type Person struct {
	Id         int			`json:"id"`
	FirstName  foo.NullString	`json:"firstName"`
	LastName   foo.NullString	`json:"lastName"`
	MiddleName foo.NullString	`json:"middleName"`
	CardNumber *int			`json:"cardNumber"`
	Fired	   bool			`json:"fired"`
}

func (p Person) CopyVal() DataObj {
	return p
}

type LTimeSheet []TimeSheet_record

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func handlerInterrupt() {
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		s := <- sigChan
		log.Info("Signal", s)
		log.Info("Cleanup service")

		closeDB()
		os.Exit(1)
	}()
}

func closeDB () {
	err := db.Close()
	if err != nil {
		log.Info("closeDB err", err)
	}
}

func PingDB() {
	for {
		db.Ping()
		//QueryDB_simple(`select 1 number`)
		time.Sleep(10 * time.Second)
	}
}

func structToSliceOfFieldAddress(obj DataObj) []interface{} {
	fieldArr := reflect.ValueOf(obj).Elem()
	fieldAddrArr := make([]interface{}, fieldArr.NumField())

	for index := 0; index < fieldArr.NumField(); index++ {
		f := fieldArr.Field(index).Addr().Interface()
		fieldAddrArr[index] = f
	}
	return fieldAddrArr
}

func QueryDB_(q string, do DataObj) ([]DataObj, error) {
	sqlFieldArrPtr := structToSliceOfFieldAddress(do)

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aRows []DataObj
	for rows.Next() {
		err := rows.Scan(sqlFieldArrPtr...)
		if err != nil {
			return nil, err
		}
		copyP := do.CopyVal()
		aRows = append(aRows, copyP)
	}

	err = rows.Err()
	return aRows, err
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

func mkdir_ine(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir|0777)
	}
}


func Init() {
	mkdir_ine("logs")
	log,_ = logger.NewLogger(logger.LogLevelDebug, "svc", "file")

	defer func () {
		if r := recover(); r != nil {
			log.Error("Recover", r)
			log.Error(string(debug.Stack()))
			bufio.NewReader(os.Stdin).ReadByte()
		}
	}()

	cfg, _ := ini.InsensitiveLoad("config.ini")
	adr := cfg.Section("MSSQL").Key("ADDRESS").String()

	r := gin.Default()
	r.GET("/person", getPersonList)
	r.POST("/timesheet", postTimeSheets)
	r.GET("/person/:id/img", getPersonPhoto)

	handlerInterrupt()
	db, _ = sql.Open("mssql", adr)
	maxConn := runtime.NumCPU() * 4
	if maxConn > 50 {
		maxConn = 50
	}
	db.SetMaxOpenConns(maxConn)

	//QueryDB_simple("select @@version v")
	//rows, err := QueryDB_(qPersonList, &Person{})
	//log.Info(err)
	//printRows(rows)
	r.Run(":3000")
}

func getPersonList(c *gin.Context) {
	rows, err := QueryDB_(qPersonList, &Person{})
	if err != nil {
		log.Info("getPersonList >", err)
		c.AbortWithStatus(500)
		return
	}
	//printParams(c.Params)
	var result gin.H = gin.H {
		"rows": rows,
	}
	c.JSON(http.StatusOK, result)
}

func getPersonPhoto(c *gin.Context) {
	var img []byte
	id,_ := strconv.Atoi(c.Param("id"))
	err := db.QueryRow("select Файл from ФайлыНабор where Сотрудник=:1", id).Scan(&img)
	if err != nil {
		log.Info("getPersonPhoto >", err)
		c.AbortWithStatus(404)
		return
	}
	c.Status(200)
	c.Writer.Write(img)
}

func postTimeSheets(c *gin.Context) {
	var l []TimeSheet_record
	err := c.BindJSON(&l)

	if err == nil {
		tx, _ := db.Begin()
		stmt, err := tx.Prepare(qsetTimeSh)
		if err != nil {
			log.Info("postTimeSheets >", err)
			tx.Rollback()
			return
		}

		for _,r := range l {
			_, err := stmt.Exec(r.EnterTime, r.ExitTime, r.TotalTime, r.Date, r.AccountId)
			
			if err != nil {
				log.Info("postTimeSheets >", err)
				tx.Rollback()
				return
			}
		}
		tx.Commit()
	}
}

