package madpack


import (
    "fmt"
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
    return self.Items[index]
}

func (self *Index) Insert(key string) uint {
    for i,k := range self.Items {
        if k == key {
            return uint(i);
        }
    }
    if len(self.Items) < 0xffff {
        self.Items = append(self.Items, key)
    }
    return 0
}
