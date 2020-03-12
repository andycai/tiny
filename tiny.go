package tiny

import (
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

var chars = []string{"^", "`", "|", "~", ",", ";", ":", "'", "\"", "\\", "/", "[", "]", "+", "=", "@", "#", "$", "&", "?", "<", ">", "-", "_", "(", ")", "!", "*"}
var ecChars = []string{"%5e", "%60", "%7c", "%7e", "%2c", "%3b", "%3a", "%27", "%22", "%5c", "%2f", "%5b", "%5d", "%2b", "%3d", "%40", "%23", "%24", "%26", "%3f", "%3c", "%3e", "%9a", "%9b", "%9c", "%9d", "%9e", "%9f"}

const (
	Zero = "{}"
)

func toBlank(s string) string {
	if s == Zero {
		return ""
	}
	return s
}

func toZero(s string) string {
	if s == "" {
		return Zero
	}
	return s
}

func isZero(s string) bool {
	if s == Zero {
		return true
	}
	return false
}

func encodeChar(s string) string {
	for i, v := range chars {
		s = strings.ReplaceAll(s, v, ecChars[i])
	}
	return s
}

func decodeChar(s string) string {
	for i, v := range ecChars {
		s = strings.ReplaceAll(s, v, chars[i])
	}
	return s
}

// Marshal struct to string
func Marshal(v interface{}) string {
	return marshal(reflect.ValueOf(v), 0)
}

func marshal(v reflect.Value, level int) (result string) {
	sep := chars[level]
	switch v.Kind() {
	case reflect.Bool:
		result += cast.ToString(v.Bool())
	case reflect.String:
		str := encodeChar(v.String())
		result += toZero(str)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result += cast.ToString(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result += cast.ToString(v.Uint())
	case reflect.Float32, reflect.Float64:
		result += cast.ToString(v.Float())
	case reflect.Ptr:
		if v.IsNil() {
			result += Zero
		} else {
			result += marshal(v.Elem(), level)
		}
	case reflect.Array, reflect.Slice:
		count := v.Len()
		str := ""
		for i := 0; i < count; i++ {
			if i >= count-1 {
				str += marshal(v.Index(i), level+1)
			} else {
				str += marshal(v.Index(i), level+1) + sep
			}
		}
		result += toZero(str)
	case reflect.Struct:
		count := v.NumField()
		str := ""
		for i := 0; i < count; i++ {
			if i >= count-1 {
				str += marshal(v.Field(i), level+1)
			} else {
				str += marshal(v.Field(i), level+1) + sep
			}
		}
		result += toZero(str)
	case reflect.Map:
		keys := v.MapKeys()
		count := len(keys)
		str := ""
		for i, k := range keys {
			s := ""
			switch k.Type().Kind() {
			case reflect.String:
				// 如果是 ""，转换为 {}
				s = toZero(k.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				s = cast.ToString(k.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				s = cast.ToString(k.Uint())
			case reflect.Float32, reflect.Float64:
				s = cast.ToString(k.Float())
			}
			if i >= count-1 {
				str += s + chars[level+1] + marshal(v.MapIndex(k), level+2)
			} else {
				str += s + chars[level+1] + marshal(v.MapIndex(k), level+2) + sep
			}
		}
		result += toZero(str)
	case reflect.Interface:
		result += marshal(v.Elem(), level)
	default:
		// fmt.Println("v.kind(): ", v.Kind())
		result += Zero
	}
	return
}

// Unmarshal string to struct
func Unmarshal(str string, v interface{}) {
	rv := reflect.ValueOf(v)
	unmarshal(str, rv.Elem(), 0)
}

func unmarshal(str string, v reflect.Value, level int) {
	fields := strings.Split(str, chars[level])
	count := len(fields)

	switch v.Kind() {
	case reflect.Bool:
		if isZero(str) {
			str = "false"
		}
		v.SetBool(cast.ToBool(str))
	case reflect.String:
		str = decodeChar(str)
		v.SetString(toBlank(str))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if isZero(str) {
			str = "0"
		}
		v.SetInt(cast.ToInt64(str))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if isZero(str) {
			str = "0"
		}
		v.SetUint(cast.ToUint64(str))
	case reflect.Float32, reflect.Float64:
		if isZero(str) {
			str = "0"
		}
		v.SetFloat(cast.ToFloat64(str))
	case reflect.Ptr:
		if isZero(str) {
			v.Set(reflect.New(v.Type()).Elem())
		} else {
			unmarshal(str, v.Elem(), level)
		}
	case reflect.Interface:
		if isZero(str) {
			v.Set(reflect.New(v.Type()).Elem())
		} else {
			unmarshal(str, v.Elem(), level)
		}
	case reflect.Array:
		count := v.Len()
		for i := 0; i < count; i++ {
			s := Zero
			if !isZero(str) {
				s = fields[i]
			}
			unmarshal(s, v.Index(i), level+1)
		}
	case reflect.Slice:
		slice := reflect.MakeSlice(v.Type(), count, count)
		v.Set(slice)

		if !isZero(str) {
			for i := 0; i < count; i++ {
				unmarshal(fields[i], v.Index(i), level+1)
			}
		}
	case reflect.Struct:
		n := v.NumField()
		for i := 0; i < n; i++ {
			// 支持 struct 后增加字段
			if i >= count {
				unmarshal(Zero, v.Field(i), level+1)
			} else {
				unmarshal(fields[i], v.Field(i), level+1)
			}
		}
	case reflect.Map:
		t := v.Type()
		elemType := t.Elem()

		var mapElem reflect.Value
		var subv reflect.Value

		kt := t.Key()
		var kv reflect.Value

		tv := reflect.MakeMap(t)

		if !isZero(str) {
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
					// 原来 key 是 {} 就不做转换了
					kv = reflect.ValueOf(ks).Convert(kt)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					kv = reflect.ValueOf(cast.ToInt64(ks)).Convert(kt)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					kv = reflect.ValueOf(cast.ToUint64(ks)).Convert(kt)
				}
				unmarshal(vs, subv, level+2)
				if kv.IsValid() {
					tv.SetMapIndex(kv, subv)
				}
			}
		}
		v.Set(tv)
	}
}
