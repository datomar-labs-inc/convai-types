package ctypes

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

const ctLayout = "2006-01-02T15:04:05.000000"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(ctLayout, s)
	return
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(ctLayout))), nil
}

var nilTime = (time.Time{}).UnixNano()

func (ct *CustomTime) IsSet() bool {
	return ct.UnixNano() != nilTime
}

func (g CustomTime) Value() (driver.Value, error) {
	return g.Time, nil
}

func (g *CustomTime) Scan(src interface{}) error {
	*g = CustomTime{src.(time.Time)}
	return nil
}

var (
	_ driver.Valuer = &CustomTime{}
	_ sql.Scanner   = &CustomTime{}
)
