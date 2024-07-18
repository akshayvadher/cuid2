# cuid2

Next generation guids.
Secure, collision-resistant ids optimized for horizontal scaling and
performance.

Original author and reference: https://github.com/paralleldrive/cuid2

# package

```shell
go get github.com/akshayvadher/cuid2
```

## usage

```go
package main

import (
	"fmt"
	"github.com/akshayvadher/cuid2"
)

func main() {
	id := cuid2.CreateId()
	fmt.Println(id)
	// to create default length id
	// us1hfvvf2uyzmh031bav6skw

	idWithLength := cuid2.CreateIdOf(10)
	// zev57ezp7c

	createId := cuid2.Init(customRandomFunction, customCounterFunction, length, customFingerprintString)
	createId()
	// this generates id with custom parameters
}
````

###

# cli

```shell
go install github.com/akshayvadher/cuid2/cmd/cuid2@latest
```

## usage

```shell
cuid2
# wldu51x7wn6baulkeq49qfm7

cuid2 -n 5
# generates 5 ids
# s7dseq8y5ti85c02eptzia1p
# wz4rk8nj39dpyd01gddsp9rz
# kwezj4wa69d6ta7jxg3b6lnz
# pzuju2hk01xpev6ixnnnsqba
# enxpfer2u7c00xa24li0jghc

cuid2 -len 10
# generates id with length 10
# h8crfzyp6q

cuid2 -n 3 -len 11
# generates 3 ids with length 11
# ijlk68norem
# redx1s0adbb
# mk5zg8dxgi1

cuid2 validate 123
# not a valid CUID2 "123"

cuid2 validate qf9183tmebd
# Valid CUID2 "qf9183tmebd"
```

# understanding
```mermaid
flowchart TD
    A[Init] -->PARTS[Parts]
    A -->I[Static]
    I --> IC(Counter)
    IC --> ICI(Random number\nfrom 1 to 476782367\nusing math.random\nas initial for counter\nExample: 284777857)
    ICI --> ICV(Get next int\nfrom initial counter\nExample: 284777858)
    ICV --> ICVS(String base 36 of counter\nExample: 4pjs02)
    I --> FI(Fingerprint)
    FI --> GLOBAL(Using ENV variable's keys \n+ process id \n+ entropy\nExample: PATH_JAVA_HOME_...32732joj7hxvuerwy2m8ogjjvya1sxnhryizp)
    GLOBAL --> GLOBAL_HASH(That creates hash\nExample: p051cbjeehq8k4srz81f1k27qy5ug4xspl9qurcloe9v0ibzxygy39imhqzm6zuq3gxl9758k525jm9m1zy5zdb7b5mkrjbh3h)
    GLOBAL_HASH --> GLOBAL_HASH_TRUNCATE(Trucate to the max len\nExample: p051cbjeehq8k4srz81f1k27qy5ug4xs)
    PARTS --> RANDOM_CHAR(Random char from a to z)
    PARTS --> HASH_PARTS(Hash parts)
    RANDOM_CHAR --> RANDOM_CHAR_E(Using math.random\nExample: u)
    HASH_PARTS --> TIME(Time millis\nExmaple: 1721287835377)
    TIME --> TIME_STR(To string represenataion\nwith base 36 = 0-9 and a-z\nExample: lyqyc1yp)
    HASH_PARTS --> SALT(Salt generation)
    SALT --> ENTROPY(Creating entropy\nUsing chars of base36\nUsing math.random\nExample: gxrohhjf8t0c6j5zjoj05o9k)
    TIME_STR --> COMBINE_HASH_INPUT(Combined Hash input\ntimeString + salt + count + fingerprint\nExample: lyqyc1ypgxrohhjf8t0c6j5zjoj05o9k4pjs02p051cbjeehq8k4srz81f1k27qy5ug4xs)
    ENTROPY --> COMBINE_HASH_INPUT
    ICVS --> COMBINE_HASH_INPUT
    GLOBAL_HASH_TRUNCATE --> COMBINE_HASH_INPUT
    COMBINE_HASH_INPUT --> SHA3(SHA3_512)
    SHA3 --> BIG_INT(Big int of SHA3_512\nExample: 42ybcc5cg8oormscwq75cvckeprc888472wq54acxqcdj6sno2tgnsuz4dslxpib2fftsy0gax4va9vaw5oem2vosqrw82glqgh)
    BIG_INT --> BIG_INT_1(Dropping first char due to bias\nExample: 2ybcc5cg8oormscwq75cvckeprc888472wq54acxqcdj6sno2tgnsuz4dslxpib2fftsy0gax4va9vaw5oem2vosqrw82glqgh)
    BIG_INT_1 --> HASH_TAIL(Tail\nDropping first char and trimming till leng\nExample: ybcc5cg8oormscwq75cvcke)
    RANDOM_CHAR_E --> OUTPUT(Output CUID2\nRandomLetter + hash of timeString + salt + count + fingerprint\nExample: uybcc5cg8oormscwq75cvcke)
    HASH_TAIL --> OUTPUT
```
