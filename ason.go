package ason

import (
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

var chars = []string{"`", "^", "~", "|", ">>", "<<", "[[", "]]", "{{", "}}", "::", ";;", ",,", "$$", "@@", "##", "&&"}

// struct to string
func Marshal(f reflect.Value, level int) (result string) {
	sep := chars[level]
	switch f.Kind() {
	case reflect.Bool:
		result += cast.ToString(f.Bool())
	case reflect.String:
		str := f.String()
		if str == "" {
			str = "''"
		}
		result += str
	case reflect.Int32, reflect.Int, reflect.Int64, reflect.Int8, reflect.Int16:
		result += cast.ToString(f.Int())
	case reflect.Uint32, reflect.Uint, reflect.Uint64, reflect.Uint8, reflect.Uint16:
		result += cast.ToString(f.Uint())
	case reflect.Float32, reflect.Float64:
		result += cast.ToString(f.Float())
	case reflect.Array, reflect.Slice:
		count := f.Len()
		str := ""
		for i := 0; i < count; i++ {
			if i >= count-1 {
				str += Marshal(f.Index(i), level+1)
			} else {
				str += Marshal(f.Index(i), level+1) + sep
			}
		}
		if str == "" {
			str = "[]"
		}
		result += str
	case reflect.Struct:
		count := f.NumField()
		str := ""
		for i := 0; i < count; i++ {
			if i >= count-1 {
				str += Marshal(f.Field(i), level+1)
			} else {
				str += Marshal(f.Field(i), level+1) + sep
			}
		}
		if str == "" {
			str = "{}"
		}
		result += str
	case reflect.Map:
		keys := f.MapKeys()
		count := len(keys)
		str := ""
		for i, k := range keys {
			s := ""
			switch k.Type().Kind() {
			case reflect.String:
				s = k.String()
			case reflect.Int32, reflect.Int, reflect.Int64, reflect.Int8, reflect.Int16:
				s = cast.ToString(k.Int())
			case reflect.Uint32, reflect.Uint, reflect.Uint64, reflect.Uint8, reflect.Uint16:
				s = cast.ToString(k.Uint())
			case reflect.Float32, reflect.Float64:
				s = cast.ToString(k.Float())
			}
			if i >= count-1 {
				str += s + chars[level+1] + Marshal(f.MapIndex(k), level+2)
			} else {
				str += s + chars[level+1] + Marshal(f.MapIndex(k), level+2) + sep
			}
		}
		if str == "" {
			str = "{}"
		}
		result += str
	}
	return
}

// string to struct
func Unmarshal(str string, f reflect.Value, level int) {
	fields := strings.Split(str, chars[level])
	count := len(fields)

	switch f.Kind() {
	case reflect.Bool:
		f.SetBool(cast.ToBool(str))
	case reflect.String:
		if str == "''" {
			str = ""
		}
		f.SetString(str)
	case reflect.Int32, reflect.Int, reflect.Int64, reflect.Int8, reflect.Int16:
		f.SetInt(cast.ToInt64(str))
	case reflect.Uint32, reflect.Uint, reflect.Uint64, reflect.Uint8, reflect.Uint16:
		f.SetUint(cast.ToUint64(str))
	case reflect.Float32, reflect.Float64:
		f.SetFloat(cast.ToFloat64(str))
	case reflect.Ptr:
		Unmarshal(str, f.Elem(), level)
	case reflect.Array:
		count := f.Len()
		if str == "[]" {
			str = ""
		}
		for i := 0; i < count; i++ {
			Unmarshal(fields[i], f.Index(i), level+1)
		}
	case reflect.Slice:
		if str == "[]" {
			str = ""
		}
		slice := reflect.MakeSlice(f.Type(), count, count)
		f.Set(slice)
		for i := 0; i < count; i++ {
			Unmarshal(fields[i], f.Index(i), level+1)
		}
	case reflect.Struct:
		if str == "{}" {
			str = ""
		}
		count := f.NumField()
		for i := 0; i < count; i++ {
			Unmarshal(fields[i], f.Field(i), level+1)
		}
	case reflect.Map:
		if str == "{}" {
			str = ""
		}
		t := f.Type()
		elemType := t.Elem()

		var mapElem reflect.Value
		var subv reflect.Value

		kt := t.Key()
		var kv reflect.Value

		v := reflect.MakeMap(t)

		for i := 0; i < count; i++ {
			mapFields := strings.Split(fields[i], chars[level+1])
			ks := mapFields[0]
			vs := mapFields[1]
			if !mapElem.IsValid() {
				mapElem = reflect.New(elemType).Elem()
			} else {
				mapElem.Set(reflect.Zero(elemType))
			}
			subv = mapElem

			if subv.Kind() == reflect.Ptr {
				if subv.IsNil() {
					if subv.CanSet() {
						subv.Set(reflect.New(subv.Type().Elem()))
					}
				}
				//subv = subv.Elem()
			}

			switch t.Key().Kind() {
			case reflect.String:
				kv = reflect.ValueOf(ks).Convert(kt)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				kv = reflect.ValueOf(cast.ToInt64(ks)).Convert(kt)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				kv = reflect.ValueOf(cast.ToUint64(ks)).Convert(kt)
			}
			Unmarshal(vs, subv, level+2)
			if kv.IsValid() {
				v.SetMapIndex(kv, subv)
			}
		}
		f.Set(v)
	}
}
