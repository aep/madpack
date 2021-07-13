package main

import (
    "github.com/aep/madpack/go"
    "os"
    "flag"
    "fmt"
    "encoding/json"
    "io/ioutil"
)


func main() {
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "madpack usage:\n  default\n  \tjson (stdin) to madpack (stdout)\n")
        flag.PrintDefaults()
    }

    unpack := flag.Bool("unpack", false, "madpack (stdin) to json (stdout)")
    indexFileName := flag.String("index", "", "read preshared index file for un/pack")
    makeIndex := flag.Bool("make-index", false, "create preshared index from stdin content")
    flag.Parse()

    var index *madpack.Index = madpack.EmptyIndex()

    if indexFileName != nil && *indexFileName != "" {
        dat, err := ioutil.ReadFile(*indexFileName)
        if err != nil { panic(err) }

        index, err = madpack.DecodeIndex(dat)
        if err != nil { panic(err) }
    }

    if *makeIndex {

        var dec = json.NewDecoder(os.Stdin)
        var v = make(map[string]interface{})
        var err = dec.Decode(&v)
        if err != nil { panic(err) }

        var enc = madpack.NewEncoderWithIndex(ioutil.Discard, index)
        err = enc.EncodeValue(v)
        if err != nil { panic(err) }

        _,err = os.Stdout.Write(index.Encode())
        if err != nil { panic(err) }

    } else if *unpack {
        var dec= madpack.NewDecoderWithIndex(os.Stdin, index)
        var v, err = dec.DecodeValue()
        if err != nil { panic(err) }

        var enc = json.NewEncoder(os.Stdout)
        err = enc.Encode(&v)
        if err != nil { panic(err) }

    } else {
        var dec = json.NewDecoder(os.Stdin)
        var v = make(map[string]interface{})
        var err = dec.Decode(&v)
        if err != nil { panic(err) }

        var enc = madpack.NewEncoderWithIndex(os.Stdout, index)
        err = enc.EncodeValue(v)
        if err != nil { panic(err) }

    }

}
