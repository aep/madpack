package madpack

import (
    "io"
    "fmt"
    "reflect"
    "encoding/binary"
    "errors"
)

type Decoder struct {
    R       io.Reader
    Index   *Index
}

func NewDecoder(r io.Reader) *Decoder {
    return &Decoder{R: r, Index: EmptyIndex()}
}

func NewDecoderWithIndex(r io.Reader, i* Index) *Decoder{
    if i == nil {
        i = EmptyIndex()
    }
    return &Decoder{R: r, Index: i}
}

type End struct {}

func setUint(v reflect.Value, vv uint64) {
    if v.Kind() == reflect.Interface {
        var vn = reflect.ValueOf(vv)
        v.Set(vn)
    } else {
        v.SetUint(vv)
    }
}
func setInt(v reflect.Value, vv int64) {
    if v.Kind() == reflect.Interface {
        var vn = reflect.ValueOf(vv)
        v.Set(vn)
    } else {
        v.SetInt(vv)
    }
}

func setFloat(v reflect.Value, vv float64) {
    if v.Kind() == reflect.Interface {
        var vn = reflect.ValueOf(vv)
        v.Set(vn)
    } else {
        v.SetFloat(vv)
    }
}
func setBool(v reflect.Value, vv bool) {
    if v.Kind() == reflect.Interface {
        var vn = reflect.ValueOf(vv)
        v.Set(vn)
    } else {
        v.SetBool(vv)
    }
}
func setString(v reflect.Value, vv string) {
    if v.Kind() == reflect.Interface {
        var vn = reflect.ValueOf(vv)
        v.Set(vn)
    } else {
        v.SetString(vv)
    }
}
func setBytes(v reflect.Value, vv []byte) {
    if v.Kind() == reflect.Interface {
        var vn = reflect.ValueOf(vv)
        v.Set(vn)
    } else {
        v.SetBytes(vv)
    }
}

func (self *Decoder) DecodeArray() (r []interface{}, eee error) {
    defer func() {
        if r := recover(); r != nil {
            eee  = fmt.Errorf("%w", r)
        }
    }()

    r = make([]interface{}, 0)

    for ;; {
        n, err := self.DecodeValue()
        if _, end := n.(End); end {
            break;
        }
        if err != nil { return r, err }
        r = append(r, n);
    }
    return r, nil
}

func (self *Decoder) DecodeMap() (r map[string]interface{}, eee error) {
    defer func() {
        if r := recover(); r != nil {
            eee  = fmt.Errorf("%w", r)
        }
    }()

    r = make(map[string]interface{})


    for ;; {

        var eb [1]byte
        _, err := io.ReadFull(self.R, eb[:])
        if err != nil {
            if errors.Is(err, io.EOF) {
                err = nil
            }
            return r , err
        }

        var key_str = "";
        var value interface{}

        var idxt = eb[0] & 0x1f;
        switch idxt {
            case 0x00 :{
                key_str = "invalid key 0"
            }
            case 0x1b : {
                var vv uint8;
                err = binary.Read(self.R, binary.LittleEndian, &vv)
                if err != nil { return nil, err }
                key_str = self.Index.Lookup(uint(vv))
            }
            case 0x1c : {
                var vv uint16;
                err = binary.Read(self.R, binary.LittleEndian, &vv)
                if err != nil { return nil, err }
                key_str = self.Index.Lookup(uint(vv))
            }
            case 0x1d: {
                var vv uint8;
                err = binary.Read(self.R, binary.LittleEndian, &vv)
                if err != nil { return nil, err }
                var bs = make([]byte, vv)
                _, err := io.ReadFull(self.R, bs[:])
                if err != nil { return nil, err }
                key_str = string(bs)
            }
            case 0x1e: {
                var vv uint16;
                err = binary.Read(self.R, binary.LittleEndian, &vv)
                if err != nil { return nil, err }
                var bs = make([]byte, vv)
                _, err := io.ReadFull(self.R, bs[:])
                if err != nil { return nil, err }
                key_str = string(bs)
            }
            case 0x1f: {
                return r, nil
            }
            default:
                key_str = self.Index.Lookup(uint(eb[0] & 0x1f))

        }

        var vt = eb[0] & 0xe0;
        switch vt {
            case 0x00: {
                var vv uint8;
                err = binary.Read(self.R, binary.LittleEndian, &vv)
                if err != nil { return nil, err }
                value = uint64(vv);
            }
            case 0x20: {
                var vv uint16;
                err = binary.Read(self.R, binary.LittleEndian, &vv)
                if err != nil { return nil, err }
                value = uint64(vv);
            }
            case 0x40: {
                var vv float32;
                err = binary.Read(self.R, binary.LittleEndian, &vv)
                if err != nil { return nil, err }
                value = float64(vv);
            }
            case 0x60: {
                var vv uint8;
                err = binary.Read(self.R, binary.LittleEndian, &vv)
                if err != nil { return nil, err }
                var bs = make([]byte, vv)
                _, err := io.ReadFull(self.R, bs[:])
                if err != nil { return nil, err }
                value = bs
            }
            case 0x80: {
                var vv uint8;
                err = binary.Read(self.R, binary.LittleEndian, &vv)
                if err != nil { return nil, err }
                var bs = make([]byte, vv)
                _, err := io.ReadFull(self.R, bs[:])
                if err != nil { return nil, err }
                value = string(bs)
            }
            case 0xa0: {
                value, err = self.DecodeMap();
                if err != nil { return nil, err }
            }
            case 0xc0: {
                value, err = self.DecodeArray();
                if err != nil { return nil, err }
            }
            case 0xe0: {
                value, err = self.DecodeValue();
                if err != nil { return nil, err }
            }
        }

        r[key_str] = value
    }
}

func (self *Decoder) DecodeValue() (r interface{}, eee error) {
    defer func() {
        if r := recover(); r != nil {
            eee  = fmt.Errorf("%w", r)
        }
    }()

    var eb [1]byte
    _, err := io.ReadFull(self.R, eb[:])
    if err != nil { return nil, err }


    switch eb[0] {
        case 0x70 : {
            var vv uint8;
            err = binary.Read(self.R, binary.LittleEndian, &vv)
            if err != nil { return nil, err }
            return uint64(vv), nil
        }
        case 0x74 : {
            var vv int8;
            err = binary.Read(self.R, binary.LittleEndian, &vv)
            if err != nil { return nil, err }
            return int64(vv), nil
        }
        case 0x71 : {
            var vv uint16;
            err = binary.Read(self.R, binary.LittleEndian, &vv)
            if err != nil { return nil, err }
            return uint64(vv), nil
        }
        case 0x75 : {
            var vv int16;
            err = binary.Read(self.R, binary.LittleEndian, &vv)
            if err != nil { return nil, err }
            return int64(vv), nil
        }
        case 0x72 : {
            var vv uint32;
            err = binary.Read(self.R, binary.LittleEndian, &vv)
            if err != nil { return nil, err }
            return uint64(vv), nil
        }
        case 0x76 : {
            var vv int32;
            err = binary.Read(self.R, binary.LittleEndian, &vv)
            if err != nil { return nil, err }
            return int64(vv), nil
        }
        case 0x7d : {
            var vv float32;
            err = binary.Read(self.R, binary.LittleEndian, &vv)
            if err != nil { return nil, err }
            return float64(vv), nil
        }
        case 0x73 : {
            var vv uint64;
            err = binary.Read(self.R, binary.LittleEndian, &vv)
            if err != nil { return nil, err }
            return uint64(vv), nil
        }
        case 0x77 : {
            var vv int64;
            err = binary.Read(self.R, binary.LittleEndian, &vv)
            if err != nil { return nil, err }
            return int64(vv), nil
        }
        case 0x7e : {
            var vv float64;
            err = binary.Read(self.R, binary.LittleEndian, &vv)
            if err != nil { return nil, err }
            return float64(vv), nil
        }
        case 0x78 : {
            return nil, nil
        }
        case 0x79 : {
            return true, nil
        }
        case 0x7a : {
            return false , nil
        }
        case 0x7b : {
            return self.DecodeMap()
        }
        case 0x7c : {
            return self.DecodeArray()
        }
        case 0xff : {
            return End{}, nil
        }
        case 0x80: {
            return "", nil
        }
        case 0x90: {
            return []byte{}, nil
        }
        case 0x8c, 0x9c: {
            var l uint8;
            err = binary.Read(self.R, binary.LittleEndian, &l)
            if err != nil { return nil, err }

            var vv = make([]byte, l)
            _, err := io.ReadFull(self.R, vv[:])
            if err != nil { return nil, err }

            switch eb[0] {
                case 0x8c: { return string(vv), nil }
                case 0x9c: { return vv, nil }
            }
        }
        case 0x8d, 0x9d : {
            var l uint16;
            err = binary.Read(self.R, binary.LittleEndian, &l)
            if err != nil { return nil, err }

            var vv = make([]byte, l)
            _, err := io.ReadFull(self.R, vv[:])
            if err != nil { return nil, err }

            switch eb[0] {
                case 0x8d: { return string(vv), nil }
                case 0x9d: { return vv, nil }
            }
        }
        case 0x8e, 0x9e : {
            var l uint32;
            err = binary.Read(self.R, binary.LittleEndian, &l)
            if err != nil { return nil, err }

            var vv = make([]byte, l)
            _, err := io.ReadFull(self.R, vv[:])
            if err != nil { return nil, err }

            switch eb[0] {
                case 0x8e: { return string(vv), nil }
                case 0x9e: { return vv, nil }
            }
        }
        case 0x8f, 0x9f : {
            var l uint64;
            err = binary.Read(self.R, binary.LittleEndian, &l)
            if err != nil { return nil, err }

            var vv = make([]byte, l)
            _, err := io.ReadFull(self.R, vv[:])
            if err != nil { return nil, err }

            switch eb[0] {
                case 0x8f: { return string(vv), nil }
                case 0x9f: { return vv, nil }
            }
        }
        default: {
            if eb[0] <= 0x6f {
                return uint64(eb[0]), nil
            } else {
                var l = (eb[0] & 0x0f);
                var vv = make([]byte, l)
                _, err := io.ReadFull(self.R, vv[:])
                if err != nil { return nil, err }

                switch eb[0] & 0xf0 {
                    case 0x80: { return string(vv), nil }
                    case 0x90: { return vv, nil }
                }
            }
        }
    }
    return nil, nil
}

