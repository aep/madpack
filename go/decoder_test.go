package madpack;

import (
    "testing"
    "bytes"
    "fmt"
)

func TestMap(t *testing.T) {
    var dec = NewDecoder(bytes.NewReader([]byte{
        0x7b, 0x1d, 0x03, 0x6c, 0x6f, 0x6c, 0x01, 0xfd, 0x07, 0x62, 0x69, 0x72, 0x64, 0x69, 0x65, 0x73,
        0x78, 0xfd, 0x03, 0x79, 0x65, 0x73, 0x7a, 0xff, 0xff,
    }))

    out,err := dec.DecodeValue();
    if err != nil {panic(err)}

    fmt.Println(out)
}

func TestArray(t *testing.T) {
    var dec = NewDecoder(bytes.NewReader([]byte{0x7c, 0x7c, 0x79, 0x78, 1, 0xff, 2,3, 0xff}))

    out,err := dec.DecodeValue();
    if err != nil {panic(err)}

    fmt.Println(out)
}

func TestFloat(t *testing.T) {
    var dec  = NewDecoder(bytes.NewReader([]byte{0x7d, 0xfa, 0x3e, 0xf6, 0xc2}))

    out,err := dec.DecodeValue();
    if err != nil {panic(err)}

    var expect float64 = float64(float32(-123.123))
    if out != expect {
        panic(fmt.Errorf("%f != %f", out, expect))
    }
}

func TestBool(t *testing.T) {
    var dec  = NewDecoder(bytes.NewReader([]byte{0x79}))

    out,err := dec.DecodeValue();
    if err != nil {panic(err)}

    if out != true {
        panic(fmt.Errorf("%t != %t", out, true))
    }

    dec = NewDecoder(bytes.NewReader([]byte{0x7a}))

    out,err = dec.DecodeValue();
    if err != nil {panic(err)}

    if out != false{
        panic(fmt.Errorf("%t != %t", out, false))
    }
}

func test_uint(inp []byte, expect uint64) {
    var dec  = NewDecoder(bytes.NewReader(inp))

    out, err := dec.DecodeValue();
    if err != nil {panic(err)}

    if out != expect {
        panic(fmt.Errorf("%d != %d", out, expect))
    }
}

func test_sint(inp []byte, expect int64) {
    var dec  = NewDecoder(bytes.NewReader(inp))

    out,err := dec.DecodeValue();
    if err != nil {panic(err)}

    if out != expect {
        panic(fmt.Errorf("%d != %d", out, expect))
    }
}
func TestUint0Max(t *testing.T) {
    test_uint(
        []byte{0x6f},
        0x6f,
    )
}

func TestUint0Null(t *testing.T) {
    test_uint(
        []byte{0x00},
        0x00,
    )
}

func TestUint1(t *testing.T) {
    test_uint(
        []byte{0x70,0xff},
        0xff,
    )
}

func TestUint2(t *testing.T) {
    test_uint(
        []byte{0x71,0xff,0x01},
        0x01ff,
    )
}

func TestUint4(t *testing.T) {
    test_uint(
        []byte{0x72,0xff,0x01,0x01,0x01},
        0x010101ff,
    )
}

func TestUint8(t *testing.T) {
    test_uint(
        []byte{0x73,0xff,0x01,0x01,0x01,0x01,0x01,0x01,0x01},
        0x01010101010101ff,
    )
}

func TestSint1(t *testing.T) {
    test_sint(
        []byte{0x74,0xfe},
        -2,
    )
}

func TestSint2(t *testing.T) {
    test_sint(
        []byte{0x75,0xff,0xf1},
        -0x0e01,
    )
}

func TestSint4(t *testing.T) {
    test_sint(
        []byte{0x76,0xff,0x01,0x01,0x01},
        0x010101ff,
    )
}

func TestSint8(t *testing.T) {
    test_sint(
        []byte{0x77,0xff,0x01,0x01,0x01,0x01,0x01,0x01,0x01},
        0x01010101010101ff,
    )
}

func test_string(inp []byte, expect string) {
    var dec  = NewDecoder(bytes.NewReader(inp))

    out, err := dec.DecodeValue();
    if err != nil {panic(err)}

    if out != expect {
        panic(fmt.Errorf("%s != %s", out, expect))
    }
}

func TestString1Size1(t *testing.T) {
    test_string(
        []byte{0x8c, 1, 'h' },
        "h",
    )
}

func TestString1Size0(t *testing.T) {
    test_string(
        []byte{0x8c, 0 },
        "",
    )
}

func TestString0Size1(t *testing.T) {
    test_string(
        []byte{0x81, 'h' },
        "h",
    )
}

func TestString0Size0(t *testing.T) {
    test_string(
        []byte{0x80 },
        "",
    )
}

func TestString0Size11(t *testing.T) {
    test_string(
        []byte{0x8b, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'},
        "hello world",
    )
}

func TestString1Size11(t *testing.T) {
    test_string(
        []byte{0x8d, 11, 0, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'},
        "hello world",
    )
}

func TestString4Size11(t *testing.T) {
    test_string(
        []byte{0x8e, 11, 0, 0, 0, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'},
        "hello world",
    )
}

func TestString8Size11(t *testing.T) {
    test_string(
        []byte{0x8f, 11,0,0,0,0,0,0,0, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'},
        "hello world",
    )
}

func TestString2Size0(t *testing.T) {
    test_string(
        []byte{0x8d, 0,0},
        "",
    )
}

func TestString4Size0(t *testing.T) {
    test_string(
        []byte{0x8e, 0,0,0,0},
        "",
    )
}

func TestString8Size0(t *testing.T) {
    test_string(
        []byte{0x8f, 0,0,0,0,0,0,0,0},
        "",
    )
}

//test map_kvstr {
//    stdin = {0x7b, 0x9d, 0x05, 'h', 'o', 'l', 'l', 'a', 0x05, 'w', 'u', 'r', 's', 't'}
//    stdout = "{\"holla\":\"wurst\"}"
//}
