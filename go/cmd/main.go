package main

import (
    "github.com/aep/madpack/go"
    "os"
)


func main() {

    var enc = madpack.NewEncoder(os.Stdout)
    var err = enc.EncodeValue(map[string]interface{}{
        "_": int64(-922337),
        "fub": []string{"yolooooo"},
        "sub": map[string]interface{}{
            "hi": "lol",
        },
    })
    if err != nil { panic(err) }


    /*
    var inp  = []byte{0x8f, 11,0,0,0,0,0,0,0, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}
    var dec  = madpack.NewDecoder(bytes.NewReader(inp))

    val, err := dec.DecodeValue();
    if err != nil {panic(err)}

    fmt.Println(val);
    */
}
