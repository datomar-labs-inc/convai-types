package ctypes

import (
	"database/sql"
	"database/sql/driver"
	"strings"

	"github.com/blang/semver"
)

type Semver struct {
	semver.Version
}

func (ct *Semver) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return
	}

	v, err := semver.Parse(s)
	if err != nil {
		return err
	}

	ct.Version = v
	return
}

func (ct *Semver) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ct.String() + "\""), nil
}

func (g Semver) Value() (driver.Value, error) {
	return g.String(), nil
}

func (g *Semver) Scan(src interface{}) error {
	s := src.(string)
	v, err := semver.Parse(s)
	if err != nil {
		return err
	}

	g.Build = v.Build
	g.Major = v.Major
	g.Minor = v.Minor
	g.Patch = v.Patch
	g.Pre = v.Pre

	return nil
}

var (
	_ driver.Valuer = &DBBlueprint{}
	_ sql.Scanner   = &DBBlueprint{}
)
