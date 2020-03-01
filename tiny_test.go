package ason

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct {
	A1 string
	A2 bool
	A3 int
	A4 int8
	A5 int16
	A6 int32
	A7 int64
	A8 float64
}

type B struct {
	B1 string
	B2 bool
	B3 int
	B4 float64
}

type SS string
type u8 uint8

type C struct {
	C1 B
	C2 [3]int
	C3 []*D
	C4 map[int]int
	C5 map[string]string
	C6 map[string]D
	C7 map[string]*D
	C8 []interface{}
	C9 map[string]interface{}
}

type D struct {
	X SS
	Y u8
}

type E struct {
	S  string
	A  [3]string
	A2 [3]int
	S1 []string
	S2 []int
	M  map[int]string
	S3 D
	S4 *D
	S5 interface{}
}

func TestTinyMarshalEmpty(t *testing.T) {
	testStr := "''`''^''^''`0^0^0`[]`[]`{}`''^0`nil`nil"
	e := E{}

	result := Marshal(e)
	assert.Equal(t, testStr, result)
	// fmt.Println("result: ", result)
}

func TestTinyUnmarshalEmpty(t *testing.T) {
	testStr := "''`''^''^''`0^0^0`[]`[]`{}`''^0`nil`nil"
	e := E{
		S:  "hello",
		S5: "interface",
	}

	Unmarshal(testStr, &e)
	s, _ := json.Marshal(e)
	fmt.Println("result: ", string(s))
}

func TestTinyMarshalCommon(t *testing.T) {
	testStr := "bbb^false^2030^3.3333`1122^3344^5566`x2x2~51^xxx~19`3~34^1~12^2~23`z~zzz^x~''^y~yyy`k1~vvv111|9^k2~vvv222|7`kk1~v3|1^kk2~v4|0`1^2^0.003`kkk1~s^kkk2~2"

	c := C{
		C1: B{
			B1: "bbb",
			B2: false,
			B3: 2030,
			B4: 3.3333,
		},
		C2: [3]int{
			1122,
			3344,
			5566,
		},
		C3: []*D{
			&D{
				X: "x2x2",
				Y: 51,
			},
			&D{
				X: "xxx",
				Y: 19,
			},
		},
		C4: map[int]int{
			1: 12,
			2: 23,
			3: 34,
		},
		C5: map[string]string{
			"x": "",
			"y": "yyy",
			"z": "zzz",
		},
		C6: map[string]D{
			"k1": D{
				X: "vvv111",
				Y: 9,
			},
			"k2": D{
				X: "vvv222",
				Y: 7,
			},
		},
		C7: map[string]*D{
			"kk1": &D{
				X: "v3",
				Y: 1,
			},
			"kk2": &D{
				X: "v4",
				Y: 0,
			},
		},
		C8: []interface{}{
			float64(1), float64(2.0), float64(3e-3),
		},
		C9: map[string]interface{}{
			"kkk1": "s", "kkk2": float64(2),
		},
	}

	result := Marshal(c)
	assert.Equal(t, testStr, result)
	// fmt.Println("result: ", result)
}

func TestTinyMarshalSimple(t *testing.T) {
	testStr := "abcde`true`11`22`33`44`55`66.66"

	a := A{
		A1: "abcde",
		A2: true,
		A3: 11,
		A4: int8(22),
		A5: int16(33),
		A6: int32(44),
		A7: int64(55),
		A8: float64(66.66),
	}
	result := Marshal(a)
	assert.Equal(t, testStr, result)
}

func TestTinyUnmarshalSimple(t *testing.T) {
	testStr := "abcde`true`11`22`33`44`55`66.66"

	a := A{}
	Unmarshal(testStr, &a)
	assert.Equal(t, "abcde", a.A1)
	assert.Equal(t, true, a.A2)
	assert.Equal(t, 11, a.A3)
	assert.Equal(t, int8(22), a.A4)
	assert.Equal(t, int16(33), a.A5)
	assert.Equal(t, int32(44), a.A6)
	assert.Equal(t, int64(55), a.A7)
	assert.Equal(t, 66.66, a.A8)
}

// type Card struct {
// 	Id   int32
// 	Name string
// }

// type Person struct {
// 	Name    string
// 	Age     int
// 	IsMan   bool
// 	Count   int32
// 	Books   [3]int32
// 	Cashes  []int32
// 	Cards   map[string]*Card
// 	CurCard Card
// 	Arr2    [3]Card
// 	Slice2  []Card
// 	Map2    map[string]string
// 	Map3    map[int32]int32
// }

// func TestTinyMarshalSimple(t *testing.T) {
// testStr := "joe`30`true`999`0^0^0`[]`{}`9999^don't cry`1001~''^2001~ggg^3001~hhh`9991~zzz^9992~yyy^9993~xxx`aaa~vvvv^bbb~uuuuu`{}"

// var p Person
// p.Name = "joe"
// p.Age = 30
// p.IsMan = true
// p.Count = 999
// //p.Books = [3]int32{123, 456, 779}
// p.Books = [3]int32{}
// //p.Cashes = []int32{333, 34333, 353533, 3223332}
// p.Cashes = make([]int32, 0)
// //p.Cards = map[string]Card{"key_j":Card{1001, "jjj"}, "key_g":Card{2001, "ggg"}, "key_h":Card{3001, "hhh"}}
// p.Arr2 = [3]Card{Card{1001, ""}, Card{2001, "ggg"}, Card{3001, "hhh"}}
// p.Slice2 = []Card{Card{9991, "zzz"}, Card{9992, "yyy"}, Card{9993, "xxx"}}
// p.CurCard = Card{9999, "don't cry"}
// p.Map2 = map[string]string{"aaa": "vvvv", "bbb": "uuuuu"}

// // v := reflect.ValueOf(p)
// result := Marshal(p)
// errmsg := fmt.Sprintf("result: %s", result)
// assert.Equal(t, testStr, result, errmsg)

// fmt.Println("struct to string: ", result)
// }

// func TestASON_Unmarshal(t *testing.T) {
// 	testStr := "joe`30`true`999`123^456^779`333^34333^353533^3223332`key_j~1001|jjj^key_g~2001|ggg^key_h~3001|hhh`9999^don't cry`1001~jjj^2001~ggg^3001~hhh`9991~zzz^9992~yyy^9993~xxx`ccc~dddddd^xxx~yyyyyy`777~888"
// 	var p2 Person
// 	p2.Name = "jim"
// 	p2.Age = 999999
// 	p2.Count = 888888
// 	// v2 := reflect.ValueOf(&p2)
// 	// Unmarshal(testStr, v2.Elem())
// 	Unmarshal(testStr, &p2)

// 	assert.Equal(t, "joe", p2.Name, "")
// 	assert.Equal(t, 30, p2.Age, "")
// 	assert.Equal(t, true, p2.IsMan, "")
// 	assert.Equal(t, "dddddd", p2.Map2["ccc"], "")
// 	assert.Equal(t, int32(888), p2.Map3[777], "")

// 	p2v, _ := json.Marshal(p2)
// 	// json.Unmarshal()
// 	fmt.Println("string to struct: ", string(p2v))
// }
