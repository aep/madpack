package madpack

import (
    "io"
    "fmt"
    "reflect"
    "encoding/binary"
)

type Encoder struct {
    W io.Writer
    Index   *Index
}

func NewEncoder(w io.Writer) *Encoder{
    return &Encoder{W: w, Index: EmptyIndex()}
}

func NewEncoderWithIndex(w io.Writer, i* Index) *Encoder{
    if i == nil {
        i = EmptyIndex()
    }
    return &Encoder{W: w, Index: i}
}


func (self *Encoder) EncodeUint(value uint64) (err error) {
    if value <= 111 {
        err = binary.Write(self.W, binary.LittleEndian,uint8(value))
        if err != nil { return err }
    } else if value <= 0xff {
        _,err = self.W.Write([]byte{0x70})
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian,uint8(value))
        if err != nil { return err }
    } else if value <= 0xffff {
        _,err = self.W.Write([]byte{0x71})
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian,uint16(value))
        if err != nil { return err }
    } else if value <= 0xffffffff{
        _,err = self.W.Write([]byte{0x72})
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian,uint32(value))
        if err != nil { return err }
    } else {
        _,err = self.W.Write([]byte{0x73})
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian,uint64(value))
        if err != nil { return err }
    }
    return nil
}

func (self *Encoder) EncodeInt(value int64) error {
    var err error
    if value > 0 {
        return self.EncodeUint(uint64(value))
    } else if value <= 127 && value >= -128 {
        _,err = self.W.Write([]byte{0x74})
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian,int8(value))
        if err != nil { return err }
    } else if value <= 32767 && value >= -32768 {
        _,err = self.W.Write([]byte{0x75})
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian,int16(value))
        if err != nil { return err }
    } else if value <= 2147483647 && value >= -2147483648 {
        _,err = self.W.Write([]byte{0x76})
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian,int32(value))
        if err != nil { return err }
    } else {
        _,err = self.W.Write([]byte{0x77})
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian,int64(value))
        if err != nil { return err }
    }
    return nil
}

func (self *Encoder) EncodeBool(value bool) (err error) {
    if value {
        _,err = self.W.Write([]byte{0x79})
    } else {
        _,err = self.W.Write([]byte{0x7a})
    }
    return
}

func (self *Encoder) EncodeBytes(value []byte, t uint8) error {
    var err error
    var l = len(value)
    if l < 11 {
        err = binary.Write(self.W, binary.LittleEndian, uint8(t | uint8(l)))
        if err != nil { return err }
    } else if l < 0xff {
        err = binary.Write(self.W, binary.LittleEndian, uint8(t | 0b00001100))
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian, uint8(l))
        if err != nil { return err }
    } else if l < 0xffff {
        err = binary.Write(self.W, binary.LittleEndian, uint8(t | 0b00001101))
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian, uint16(l))
        if err != nil { return err }
    } else if l < 0xffffffff {
        err = binary.Write(self.W, binary.LittleEndian, uint8(t | 0b00001110))
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian, uint32(l))
        if err != nil { return err }
    } else {
        err = binary.Write(self.W, binary.LittleEndian, uint8(t | 0b00001111))
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian, uint64(l))
        if err != nil { return err }
    }

    i, err := self.W.Write(value)
    if err != nil {
        return err;
    }
    if i != l {
        return fmt.Errorf("short write")
    }
    return nil

}


func (self *Encoder) EncodeEnd() (err error) {
    _,err = self.W.Write([]byte{0xff})
    return err
}

func (self *Encoder) EncodeMapStart() (err error) {
    _,err = self.W.Write([]byte{0x7b})
    return err
}

func (self *Encoder) EncodeKV(k string, v interface{}) (err error) {
    var rv = reflect.ValueOf(v)
    switch rv.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        var vv = rv.Int()
        if vv >0 && vv < 0xff {
            err = self.EncodeKey(k, 0x00);
            if err != nil { return err }
            return binary.Write(self.W, binary.LittleEndian, uint8(vv))
        } else if vv >0 && vv < 0xffff {
            err = self.EncodeKey(k, 0x20);
            if err != nil { return err }
            return binary.Write(self.W, binary.LittleEndian, uint16(vv))
        }
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        var vv = rv.Uint()
        if vv < 0xff {
            err = self.EncodeKey(k, 0x00);
            if err != nil { return err }
            return binary.Write(self.W, binary.LittleEndian, uint8(vv))
        } else if vv < 0xffff {
            err = self.EncodeKey(k, 0x20);
            if err != nil { return err }
            return binary.Write(self.W, binary.LittleEndian, uint16(vv))
        }
    case reflect.Map:
        iter := rv.MapRange()
        err = self.EncodeKey(k, 0xa0);
        if err != nil { return err }
        for iter.Next() {
            mk := iter.Key()
            mv := iter.Value()
            if mk.Kind() != reflect.String {
                return fmt.Errorf("cannot encode map with non-string key %T", v)
            }

            err = self.EncodeKV(mk.String(), mv.Interface())
            if err != nil { return err }
        }
        return self.EncodeEnd();
    }

    err = self.EncodeKey(k, 0xe0);
    if err != nil { return err }
    err = self.EncodeValue(v);
    return err
}

func (self *Encoder) EncodeKey(key string, valuetype uint8) (err error) {
    var ki = self.Index.Insert(key)
    if ki != 0 {
        if ki <= 0x1a {
            _,err = self.W.Write([]byte{uint8(ki) | valuetype})
            return err;
        } else if ki <= 0xff {
            _,err = self.W.Write([]byte{uint8(0x1b) | valuetype, uint8(ki)})
            return err;
        } else {
            _,err = self.W.Write([]byte{uint8(0x1c) | valuetype})
            if err != nil { return err }
            err = binary.Write(self.W, binary.LittleEndian,uint16(ki))
            return err
        }
    }

    var l = len(key)
    if l < 0xff {
        _,err = self.W.Write([]byte{uint8(0x1d) | valuetype})
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian,uint8(l))
    } else if l < 0xffff {
        _,err = self.W.Write([]byte{uint8(0x1e) | valuetype})
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian,uint16(l))
    } else {
        return fmt.Errorf("key too long: \"%s\"", key)
    }

    i, err := self.W.Write([]byte(key))
    if err != nil {
        return err;
    }
    if i != l {
        return fmt.Errorf("short write")
    }

    return nil
}

func isIntegral(val float64) bool {
	return val == float64(int64(val))
}

func (self *Encoder) EncodeValue(v interface{}) (err error) {
    var rv = reflect.ValueOf(v)
    switch rv.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return self.EncodeInt(rv.Int())
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return self.EncodeUint(rv.Uint())
    case reflect.Bool:
        return self.EncodeBool(rv.Bool())
    case reflect.Map:
        iter := rv.MapRange()
        err = self.EncodeMapStart();
        if err != nil { return err }
        for iter.Next() {
            mk := iter.Key()
            mv := iter.Value()
            if mk.Kind() != reflect.String {
                return fmt.Errorf("cannot encode map with non-string key %T", v)
            }

            err = self.EncodeKV(mk.String(), mv.Interface())
            if err != nil { return err }
        }
        return self.EncodeEnd();

    case reflect.Array, reflect.Slice:
        if vv, ok := v.([]byte); ok {
            return self.EncodeBytes([]byte(vv), 0x90)
        } else {
            err = self.EncodeBytes([]byte(vv), 0x7c)
            if err != nil { return err }

            for i := 0; i < rv.Len(); i++ {
                var vv = rv.Index(i)
                err = self.EncodeValue(vv.Interface());
                if err != nil { return err }
            }

            return self.EncodeEnd();
        }
    case reflect.String:
        return self.EncodeBytes([]byte(rv.String()), 0x80)
    case reflect.Float32:
        if isIntegral(float64(v.(float32))) {
            return self.EncodeValue(int64(v.(float32)))
        }
        _,err = self.W.Write([]byte{0x7d})
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian, v.(float32))
        if err != nil { return err }
        return nil
    case reflect.Float64:
        if isIntegral(v.(float64)) {
            return self.EncodeValue(int64(v.(float64)))
        }
        _,err = self.W.Write([]byte{0x7e})
        if err != nil { return err }
        err = binary.Write(self.W, binary.LittleEndian, v.(float64))
        if err != nil { return err }
        return nil
    case reflect.Invalid:
        _,err = self.W.Write([]byte{0x78})
        if err != nil { return err }
        return nil
    case reflect.Interface, reflect.Ptr:
        if rv.IsNil() {
            _,err = self.W.Write([]byte{0x78})
            if err != nil { return err }
            return nil
        }
        return fmt.Errorf("cannot encode %T", v)
    default:
        return fmt.Errorf("cannot encode %T %d", v, rv.Kind())
    }
}
