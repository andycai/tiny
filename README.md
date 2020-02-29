# Ason

Ason is another small object notation.

    type Card struct {
        Id   int32
        Name string
    }

    type Person struct {
        Name    string
        Age     int
        IsMan   bool
        Count   int32
        Books   [3]int32
        Cashes  []int32
        Cards   map[string]*Card
        CurCard Card
        Arr2    [3]Card
        Slice2  []Card
        Map2    map[string]string
        Map3    map[int32]int32
    }

    func main() {
        var p Person
        p.Name = "joe"
        p.Age = 30
        p.IsMan = true
        p.Count = 999
        //p.Books = [3]int32{123, 456, 779}
        p.Books = [3]int32{}
        //p.Cashes = []int32{333, 34333, 353533, 3223332}
        p.Cashes = make([]int32, 0)
        //p.Cards = map[string]Card{"key_j":Card{1001, "jjj"}, "key_g":Card{2001, "ggg"}, "key_h":Card{3001, "hhh"}}
        p.Arr2 = [3]Card{Card{1001, ""}, Card{2001, "ggg"}, Card{3001, "hhh"}}
        p.Slice2 = []Card{Card{9991, "zzz"}, Card{9992, "yyy"}, Card{9993, "xxx"}}
        p.CurCard = Card{9999, "don't cry"}
        p.Map2 = map[string]string{"aaa": "vvvv", "bbb": "uuuuu"}

        v := reflect.ValueOf(p)
        result := Marshal(v)
        fmt.Println("struct to string: ", result)

        var p2 Person
        p2.Name = "jim"
        p2.Age = 999999
        p2.Count = 888888
        v2 := reflect.ValueOf(&p2)
        Unmarshal("joe`30`true`999`123^456^779`333^34333^353533^3223332`key_j~1001|jjj^key_g~2001|ggg^key_h~3001|hhh`9999^don't cry`1001~jjj^2001~ggg^3001~hhh`9991~zzz^9992~yyy^9993~xxx`ccc~dddddd^xxx~yyyyyy`777~888", v2.Elem(), 0)
        p2v, _ := json.Marshal(p2)
        fmt.Println("string to struct: ", string(p2v))
    }