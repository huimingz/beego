package reqparse

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/astaxie/beego"
)

var reqParser *RequestParser

func init() {
	reqParser = &RequestParser{}
}

type RequestParser struct {
	HttpErrorCode    int
	DisableAutoLower bool
}

// 标签格式：
//  `parser:"c;Required;Range(1,2,3)"`
// tag formate: parser:"kname; default(); choices(a,c,f); location("json"); trim; nullable" help:"xxxsa"

func (p *RequestParser) ParseArgs(c *beego.Controller, obj interface{}) (err error) {
	if p.HttpErrorCode == 0 {
		p.HttpErrorCode = http.StatusBadRequest
	}

	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)

	if objV.Kind() == reflect.Ptr && objV.Elem().Kind() == reflect.Struct {
		objT = objT.Elem()
		objV = objV.Elem()
	} else {
		err = errors.New(objT.Name() + " must be a struct pointer")
		return
	}

	var key string
	for i := 0; i < objV.NumField(); i++ {
		// 跳过不可导出字段 (unexport struct field)
		if !objV.Field(i).CanSet() {
			continue
		}

		tag := objT.Field(i).Tag.Get("parser")
		tags := strings.Split(tag, ";")

		// 除去空白字符
		for i, v := range tags {
			tags[i] = strings.TrimSpace(v)
		}

		// 获取字段名
		if tags[0] == "" {
			if !p.DisableAutoLower {
				key = strings.ToLower(objT.Field(i).Name)
			} else {
				key = objT.Field(i).Name
			}
		} else {
			key = tags[0]
		}

		if len(tag) > 1 {
			tags = tags[1:]
		} else {
			tags = make([]string, 0)
		}

		var hasReuired = false
		for _, tag := range tags {
			if strings.Title(tag) == "Required" {
				hasReuired = true
				break
			}
		}

		// 检查key是否存在
		hasKey, err := requiredCheck(c, key, hasReuired)
		if err != nil {
			return err
		}

		if hasKey {
			// 获取key的值，并设置该值
			err = p.autoSetValue(&FromValues{ctl: c}, key, objV.Field(i))
			if err != nil {
				return err
			}
		} else {
			// 设置默认值
			err = setZeroValue(objV.Field(i))
			if err != nil {
				return &ValueError{objT.Field(i).Name, err.Error()}
			}
		}

		// 验证值是否符合指定的条件
		err = p.Valid(key, tags, objV.Field(i))
		if err != nil {
			return err
		}

		// 检查自定义验证
		err = p.validCustom(key, objV.Field(i), objV, objT)
		if err != nil {
			return &ValueError{key, err.Error()}
		}
	}
	return nil

}

func (p RequestParser) validCustom(key string, val, vs reflect.Value, t reflect.Type) error {
	params := make([]reflect.Value, 2)
	cvName := "Validate" + strings.Title(key)
	cvMethod, ok := t.MethodByName(cvName)
	if !ok {
		return nil
	}
	params[0] = vs
	params[1] = val

	// 检查参数数量
	if cvMethod.Func.Type().NumIn() != 2 {
		// 跳过参数格式不匹配的自定义验证方法
		return nil
	}
	if cvMethod.Func.Type().NumOut() != 1 {
		return nil
	}
	result := cvMethod.Func.Call(params)[0]
	if err, ok := result.Interface().(error); ok && err != nil {
		return err
	}
	return nil
}

func (p *RequestParser) Valid(k string, tags []string, v reflect.Value) error {
	// 验证条件
	validator := Validation{v, k}
	vt := reflect.TypeOf(validator)

	for _, tag := range tags {
		if tag == "Required" {
			continue
		}

		// 检查是否是函数

		vf, params, err := parseFunc(tag, vt)
		if err != nil {
			return err
		}
		params[0] = reflect.ValueOf(validator)
		result := vf.Func.Call(params)[0]
		if err, ok := result.Interface().(error); ok {
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *RequestParser) autoSetValue(geter ValueGetter, k string, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := geter.GetInt64(k)
		if err != nil {
			return ValueError{k, fmt.Sprintf("'%s' is not a valid choice", geter.GetString(k))}
		}
		v.SetInt(val)
	case reflect.String:
		v.SetString(geter.GetString(k))
	case reflect.Float32, reflect.Float64:
	}
	return nil
}
