package reqparse

import (
	"errors"
	"fmt"
	"reflect"
)

// var (
// 	unFunc = map[string]bool{}
// )

type Validation struct {
	val reflect.Value
	key string
}

func (v Validation) Required() error {
	msg := fmt.Sprintf("missing required parameter '%s'", v.key)
	switch v.val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.val.Int() == 0 {
			return errors.New(msg)
		}
	case reflect.Float32, reflect.Float64:
		if v.val.Float() == 0 {
			return errors.New(msg)
		}
	case reflect.String:
		if v.val.String() == "" {
			return errors.New(msg)
		}
		v.val.Uint()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.val.Uint() == 0 {
			return errors.New(msg)
		}
	default:
		return &ValueError{v.key, "Unsupported value type"}
	}
	return nil
}

func (v Validation) Range(min, max int) error {
	msg := fmt.Sprintf("'%v' out of range[%d, %d]", v.val.Interface(), min, max)
	switch v.val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.val.Int() < int64(min) || v.val.Int() > int64(max) {
			return errors.New(msg)
		}
	case reflect.Float32, reflect.Float64:
		if v.val.Float() < float64(min) || v.val.Float() > float64(max) {
			return errors.New(msg)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.val.Uint() < uint64(min) || v.val.Uint() > uint64(max) {
			return errors.New(msg)
		}
	default:
		return errors.New("unsupported value type")
	}
	return nil
}

func (v Validation) Choices(choices []string) error {
	hasExists := false
	for _, choice := range choices {
		if fmt.Sprintf("%v", v.val.Interface()) == choice {
			hasExists = true
		}
	}

	if !hasExists {
		return fmt.Errorf("'%v' is not a valid choice", v.val.Interface())
	}
	return nil
}
