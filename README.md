# Tiny 

Tiny is another small object notation.

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
	fmt.Println("result: ", result)

    output:
    bbb^false^2030^3.3333`1122^3344^5566`x2x2~51^xxx~19`3~34^1~12^2~23`z~zzz^x~''^y~yyy`k1~vvv111|9^k2~vvv222|7`kk1~v3|1^kk2~v4|0`1^2^0.003`kkk1~s^kkk2~2