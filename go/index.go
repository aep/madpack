package madpack


import (
    "fmt"
    "bytes"
    "errors"
    "io"
)

type Index struct  {
    Items           []string
}

func EmptyIndex() *Index{
    return &Index{}
}

func (self *Index) Lookup(index uint) string {
    index -=1;
    if index >= uint(len(self.Items)) {
        self.Items = append(self.Items, fmt.Sprintf("%d", index + 1))
    }
    if index >= uint(len(self.Items)) {
        return fmt.Sprintf("%d", index + 1)
    }
    return self.Items[index]
}

func (self *Index) Insert(key string) uint {
    for i,k := range self.Items {
        if k == key {
            return uint(i + 1) ;
        }
    }
    if len(self.Items) < 0xffff {
        self.Items = append(self.Items, key)
    }
    return 0
}

func (self *Index) Encode() []byte {
    var b bytes.Buffer

    var enc = NewEncoder(&b)
    for _, item  := range self.Items {
        var err = enc.EncodeValue(item)
        if err != nil { panic(err) }
    }

    return b.Bytes()
}

func DecodeIndex(b []byte) (*Index, error) {
    var self = &Index{};

    var bb = bytes.NewBuffer(b)
    var dec = NewDecoder(bb)
    for ;; {
        var v, err = dec.DecodeValue()
        if errors.Is(err, io.EOF) {
            break
        }
        if err != nil { return nil,err }
        if vv, ok := v.(string); ok {
            self.Items = append(self.Items, vv)
        }
    }
    return self, nil
}
