// Package env implements tags that allow data to be filled in from the environment.
package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Decode takes a type with `env` tags and will fill fields in with values
// pulled from the environment.
func Decode(v interface{}) error {
	var (
		value reflect.Value
		vtype reflect.Type
	)

	value = reflect.ValueOf(v)
	vtype = value.Type()

	if vtype.Kind() == reflect.Ptr {
		vtype = vtype.Elem()

		if value.IsNil() {
			value.Set(reflect.New(vtype))
		}

		value = value.Elem()
	}

	for i := 0; i < value.NumField(); i++ {
		fval := value.Field(i)
		ftype := vtype.Field(i)

		// get the tagged OS env var name.
		tag, ok := ftype.Tag.Lookup("env")
		if !ok || !fval.CanSet() { // skip
			continue
		}

		if fval.Kind() == reflect.Ptr { // *struct
			if !fval.CanSet() {
				continue
			}

			if fval.IsNil() {
				// struct, then fval becomes *struct
				val := reflect.New(fval.Type().Elem())
				fval.Set(val)
			}

			if err := Decode(fval.Interface()); err != nil {
				return err
			}

			continue
		}

		if fval.Kind() == reflect.Struct { // struct
			if err := Decode(fval.Addr().Interface()); err != nil {
				return err
			}
		}

		str := os.Getenv(strings.ToUpper(tag))
		if len(str) == 0 {
			continue
		}

		if err := setField(fval, str); err != nil {
			return err
		}
	}

	return nil
}

// right now, we only allow maps that are convertible to map<string,string>
var allowedMapType = reflect.TypeOf(map[string]string{})

func setField(field reflect.Value, value string) error {
	switch field.Type().Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int32, reflect.Int64:
		val, _ := strconv.ParseInt(value, 0, field.Type().Bits())

		field.SetInt(val)
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		val, _ := strconv.ParseUint(value, 0, field.Type().Bits())

		field.SetUint(val)
	case reflect.Slice:
		// just assuming it's CSV.
		csv := strings.Split(value, ",")

		slc := reflect.MakeSlice(field.Type(), len(csv), len(csv))
		for i, v := range csv {
			if err := setField(slc.Index(i), v); err != nil {
				return err
			}
		}

		field.Set(slc)
	case reflect.Map:
		if !field.Type().ConvertibleTo(allowedMapType) {
			return fmt.Errorf("env: map type not convertible to map[string]string, is %s", field.Type())
		}

		m := reflect.MakeMap(field.Type())

		// csv with = pairs, eg: a=b,x=y,z=a
		csv := strings.Split(value, ",")
		for i := range csv {
			kv := strings.Split(csv[i], "=")
			if len(kv) != 2 {
				continue
			}

			m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(kv[1]))
		}

		field.Set(m)
	case reflect.Bool:
		b, _ := strconv.ParseBool(value)
		field.SetBool(b)
	default:
		panic(fmt.Sprintf("env: unsupported type(%s)", field.Type().Kind()))
	}

	return nil
}
