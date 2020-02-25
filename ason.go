package ason

import (
	"reflect"
	"strconv"
	"strings"
)

var chars = []string{"`", "^", "~", "|", "[", "]", "{", "}", "::", ";;", ",,", ">>", "<<", "$$", "@@", "##", "&&"}

// struct to string
func Format(f reflect.Value, level int) (result string) {
	sep := chars[level]
	switch f.Kind() {
	case reflect.Bool: // TODO
		result += strconv.FormatBool(f.Bool())
	case reflect.String:
		result += f.String()
	case reflect.Int32, reflect.Int, reflect.Int64, reflect.Int8, reflect.Int16:
		result += strconv.FormatInt(f.Int(), 10)
	case reflect.Uint32, reflect.Uint, reflect.Uint64, reflect.Uint8, reflect.Uint16:
		result += strconv.FormatUint(f.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		result += strconv.FormatFloat(f.Float(), 'E', -1, 64)
	case reflect.Array, reflect.Slice:
		count := f.Len()
		str := ""
		for i := 0; i < count; i++ {
			if i >= count-1 {
				str += Format(f.Index(i), level+1)
			} else {
				str += Format(f.Index(i), level+1) + sep
			}
		}
		result += str
	case reflect.Struct:
		count := f.NumField()
		str := ""
		for i := 0; i < count; i++ {
			if i >= count-1 {
				str += Format(f.Field(i), level+1)
			} else {
				str += Format(f.Field(i), level+1) + sep
			}
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
				s = strconv.FormatInt(k.Int(), 10)
			case reflect.Uint32, reflect.Uint, reflect.Uint64, reflect.Uint8, reflect.Uint16:
				s = strconv.FormatUint(k.Uint(), 10)
			case reflect.Float32, reflect.Float64:
				s = strconv.FormatFloat(k.Float(), 'E', -1, 64)
			}
			if i >= count-1 {
				str += s + chars[level+1] + Format(f.MapIndex(k), level+2)
			} else {
				str += s + chars[level+1] + Format(f.MapIndex(k), level+2) + sep
			}
		}
		result += str
	}
	return
}

// string to struct
func Parse(str string, f reflect.Value, level int) {
	arr := strings.Split(str, chars[level])
	count := len(arr)

	switch f.Kind() {
	case reflect.Bool:
		val, err := strconv.ParseBool(str)
		if err == nil {
			f.SetBool(val)
		}
	case reflect.String:
		f.SetString(str)
	case reflect.Int32, reflect.Int, reflect.Int64, reflect.Int8, reflect.Int16:
		val, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			f.SetInt(val)
		}
	case reflect.Uint32, reflect.Uint, reflect.Uint64, reflect.Uint8, reflect.Uint16:
		val, err := strconv.ParseUint(str, 10, 64)
		if err == nil {
			f.SetUint(val)
		}
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(str, 64)
		if err == nil {
			f.SetFloat(val)
		}
	case reflect.Ptr:
		d := f.Elem()
		Parse(str, d, level)
	case reflect.Array:
		count := f.Len()
		for i := 0; i < count; i++ {
			Parse(arr[i], f.Index(i), level+1)
		}
	case reflect.Slice:
		slice := reflect.MakeSlice(f.Type(), count, count)
		f.Set(slice)
		for i := 0; i < count; i++ {
			Parse(arr[i], f.Index(i), level+1)
		}
	case reflect.Struct:
		count := f.NumField()
		for i := 0; i < count; i++ {
			Parse(arr[i], f.Field(i), level+1)
		}
	case reflect.Map:
		t := f.Type()
		elemType := t.Elem()

		var mapElem reflect.Value
		var subv reflect.Value

		kt := t.Key()
		var kv reflect.Value

		v := reflect.MakeMap(t)

		for i := 0; i < count; i++ {
			arr1 := strings.Split(arr[i], chars[level+1])
			ks := arr1[0]
			vs := arr1[1]
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
				n, err := strconv.ParseInt(ks, 10, 64)
				if err == nil {
					kv = reflect.ValueOf(n).Convert(kt)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				n, err := strconv.ParseUint(ks, 10, 64)
				if err == nil {
					kv = reflect.ValueOf(n).Convert(kt)
				}
			}
			Parse(vs, subv, level+2)
			if kv.IsValid() {
				v.SetMapIndex(kv, subv)
			}
		}
		f.Set(v)
	}
}
