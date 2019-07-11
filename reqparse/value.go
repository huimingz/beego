package reqparse

import "github.com/astaxie/beego"

type ValueGetter interface {
	GetString(key string, def ...string) string
	GetStrings(key string, def ...[]string) []string
	GetInt64(key string, def ...int64) (val int64, err error)
	GetBool(key string, def ...bool) (val bool, err error)
}

type FromValues struct {
	ctl *beego.Controller
}

func (v *FromValues) GetString(key string, def ...string) string {
	return v.ctl.GetString(key)
}

func (v *FromValues) GetStrings(k string, def ...[]string) []string {
	return v.ctl.GetStrings(k)
}

func (v *FromValues) GetInt64(key string, def ...int64) (val int64, err error) {
	return v.ctl.GetInt64(key)
}

func (v *FromValues) GetBool(key string, def ...bool) (val bool, err error) {
	return v.ctl.GetBool(key, def...)
}
