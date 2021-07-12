# madpack

stupid small json binary encoding.

uses mad hax (tm) to pack all valid json faster and smaller than **everything else**
(tested: msgpack, ubjson, protobuf, gzip and zstandard)


## building & usage

get the [zetz.it](https://zetz.it) compiler

    zz build --release
    ./target/release/bin/madpack < some.json > some.madpack
    ./target/release/bin/madpack --unpack < some.madpack

zetz can automatically create packages for many languages including npm, python, rust, golang.


so lets say we have this small json

```json
{
    "sha256": "beep boop yadda",
    "commitmsg": "hella",
    "stable": false,
    "contentsize": 2332
}
```

    $ wc -c < some.json
    108

compressing it with zstandard (it's really good)

    $ zstd some.json | wc -c
    99

madpack beats zstd out of the box for small messages

    $ ./target/release/bin/madpack < some.json | wc -c
    65

but preshared indices are the real magic

    $ ./target/release/bin/madpack --make-index < some.json > some.madindex
    $ ./target/release/bin/madpack --index some.madindex  < some.json | wc -c
    29




## preshared index for structured data

The way to beat zipped json is to create an index from a data sample or schema and simply not include that index in the encoded file.
This is what protobuf does, but we still beat protobuf with other mad hax;

Also unlike with protobuf, the message can be fully recovered without the index.
Desaster recovery or bug hunting is much easier with a fully intact message structure.

to create an index from any json run

    ./target/release/bin/madpack --make-index < some.json > some.madindex


you can then use the index in encoding and decoding

    ./target/release/bin/madpack --index some.madindex < some.json > some.madpack
    ./target/release/bin/madpack --index some.madindex --unpack < some.madpack

if you ever loose the index, you can still read the message and guess the key names from context

    ./target/release/bin/madpack --unpack < some.madpack
    {
        "1" : 77,
        "2" : "0.10.1-4-g2bc35c7-dirty-R3riyLGLmb",
        "3" : true,






## encoding

- there's a maximum of 65535 unique strings per file
- key strings can only be 65535 bytes long
- in maps,   every member is preceeded by a key byte
- in arrays, every member is preceeded by a value byte

### key byte

    000x xxxx 0x00 value is u8
    001x xxxx 0x20 value is u16
    010x xxxx 0x40 value is f32
    011x xxxx 0x60 value is bytes  u8
    100x xxxx 0x80 value is string u8
    101x xxxx 0xa0 value is map
    110x xxxx 0xc0 value is array
    111x xxxx 0xe0 full value byte follows


    xxx0 0000 0x00 reserved
    xxx0 0001 0x01 key number 1
    xxx1 1010 0x1a key number 26
    xxx1 1011 0x1b key number as a u8
    xxx1 1100 0x1c key number as a u16
    xxx1 1101 0x1d key is a string size u8
    xxx1 1110 0x1e key is a string size u16

    1111 1111 0xff end

### value byte

    0000 0000 0x00  literal 0
    0110 1111 0x6f  literal 111

    0111 0000 0x70  u8
    0111 0001 0x71  u16
    0111 0010 0x72  u32
    0111 0011 0x73  u64

    0111 0100 0x74  i8
    0111 0101 0x75  i16
    0111 0110 0x76  i32
    0111 0111 0x77  i64

    0111 1000 0x78  null
    0111 1001 0x79  true
    0111 1010 0x7a  false
    0111 1011 0x7b  map

    0111 1100 0x7c  array
    0111 1101 0x7d  f32
    0111 1110 0x7e  f64
    0111 1111 0x7f  ext

    dynamic size:

    1000 xxxx 0x80 string
    1001 xxxx 0x90 bytes
    1010 xxxx 0xa0 reserved
    1011 xxxx 0xb0 reserved
    1100 xxxx 0xc0 reserved
    1101 xxxx 0xd0 reserved
    1110 xxxx 0xe0 reserved
    1111 xxxx 0xf0 reserved

         0000 size 0
         1011 size 11
         1100 size see next 1 bytes
         1101 size see next 2 bytes
         1110 size see next 4 bytes
         1111 size see next 8 bytes


    1111 1111 0xff end





