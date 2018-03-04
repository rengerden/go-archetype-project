package foo

import (
	"database/sql/driver"
	"errors"
	"bytes"
	"encoding/json"
	"time"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
const (
	Yes YesNoEnum = true
	No            = false
)

type YesNoEnum bool

func (yne YesNoEnum) Value() (driver.Value, error) {
	return bool(yne), nil
}

func (yne *YesNoEnum) Scan(value interface{}) error {
	// if value is nil, false
	if value == nil {
		// set the value of the pointer yne to YesNoEnum(false)
		*yne = YesNoEnum(false)
		return nil
	}

	if bv, err := driver.Bool.ConvertValue(value); err == nil {
		// if this is a bool type
		if v, ok := bv.(bool); ok {
			// set the value of the pointer yne to YesNoEnum(v)
			*yne = YesNoEnum(v)
			return nil
		}
	}
	// otherwise, return an error
	return errors.New("failed to scan YesNoEnum")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type NullString struct {
	String string	//
	Valid  bool	// Valid is true if String is not NULL
}

func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

func (ns *NullString) Scan(value interface{}) error {
	//fmt.Println("Scan >", value)
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	//ns.Valid = true

	if bv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := bv.(string); ok {
			// set the value of the pointer yne to YesNoEnum(v)
			*ns = NullString{v,true}
			return nil
		}
	}
	return errors.New("failed to scan NullString")
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	//fmt.Println("MarshalJSON >", ns.String, ns.Valid)
	var buf bytes.Buffer

	if ns.Valid {
		buf.WriteByte('"')
		buf.WriteString(ns.String)
		buf.WriteByte('"')

	} else {
		buf.WriteString("null")
	}

	//fmt.Println("MarshalJSON >", buf.String())
	return buf.Bytes(), nil
}


func (ns *NullString) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)

	if err != nil {
		return err
	}
	ns.String = s
	ns.Valid = s != ""
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type JsonUnixNsToTime time.Time

func (p JsonUnixNsToTime) String() string {
	i := time.Time(p).Unix()
	return strconv.FormatInt(int64(i), 10)
}

func (ns JsonUnixNsToTime) Value() (driver.Value, error) {
	return time.Time(ns), nil
}
func (ns *JsonUnixNsToTime) UnmarshalJSON(b []byte) error {
	var i int64
	err := json.Unmarshal(b, &i)

	if err != nil {
		return err
	}

	*ns = JsonUnixNsToTime(time.Unix(i / 1000, 0))
	return nil
}