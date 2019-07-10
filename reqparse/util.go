package reqparse

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

type ValueError struct {
	key string
	msg string
}

func (e ValueError) Error() string {
	return e.key + ": " + e.msg
}

func parseFunc(tag string, vft reflect.Type) (vf reflect.Method, params []reflect.Value, err error) {
	if tag == "" {
		return
	}

	compiler, err := regexp.Compile(`(\w+)\(([\w, ]*)\)`)
	if err != nil {
		return
	}

	group := compiler.FindStringSubmatch(tag)

	vf, ok := vft.MethodByName(group[1])
	if !ok {
		err = fmt.Errorf("invalid function name '%s'", group[1])
	}

	params = append(params, vf.Func)

	if group[2] == "" {
		return
	}

	args := strings.Split(group[2], ",")
	if args[0] == "" && len(args) == 1 {
		args = make([]string, 0)
	}

	// 检查参数是否是能转换为指定参数类型，否则返回错误
	for i := 1; i < vf.Func.Type().NumIn(); i++ {
		arg := args[i-1]
		arg = strings.TrimSpace(arg)

		switch vf.Func.Type().In(i).Kind() {
		case reflect.String:
			params = append(params, reflect.ValueOf(arg))
		case reflect.Int:
			val, err := strconv.Atoi(arg)
			if err != nil {
				return vf, params, err
			}
			params = append(params, reflect.ValueOf(val))
		case reflect.Slice:
			choices := make([]string, len(args[i-1:]))
			for index, arg_ := range args[i-1:] {
				arg_ = strings.TrimSpace(arg_)
				choices[index] = arg_
			}
			params = append(params, reflect.ValueOf(choices))
			break
		}
	}

	// 检查参数长度是否符合方法参数长度
	if len(params) < vf.Func.Type().NumIn() {
		err = fmt.Errorf("call func '%s' with too few input arguments", vf.Name)
		return
	} else if len(params) > vf.Func.Type().NumIn() {
		err = fmt.Errorf("call func '%s' with too many input arguments", vf.Name)
		return
	}
	return
}

func keyIsExist(ctl *beego.Controller, key string) bool {
	params := ctl.Ctx.Input.Params()
	for pKey := range params {
		if pKey == key {
			return true
		}
	}

	if ctl.Ctx.Input.Context.Request.Form == nil {
		ctl.Ctx.Input.Context.Request.ParseForm()
	}
	if _, ok := ctl.Ctx.Input.Context.Request.Form[key]; ok {
		return true
	}
	return false
}

func requiredCheck(ctl *beego.Controller, key string, required bool) (hasKey bool, err error) {
	hasKey = keyIsExist(ctl, key)
	if !hasKey {
		if required {
			err = &ValueError{key, fmt.Sprintf("missing required parameter '%s'", key)}
		}
	}
	return
}

func SetZeroValue(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(0)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(0)
	case reflect.String:
		v.SetString("")
	case reflect.Float32, reflect.Float64:
		v.SetFloat(0)
	default:
		return errors.New(fmt.Sprintf("'%s' is unsupported value type", v.Kind()))
	}
	return nil
}
