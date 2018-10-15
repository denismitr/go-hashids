### GO Hashids
This library can generate hashed/obfuscated ids from numbers. Usually this kind of functionality is required to create shorter slugs that don't reveal the DB incremental ids. The algorithm is reversable but you can use *salt* to make it more secure. However this algorithm is not suitable for cryptographical purpuses.

#### Version 1

#### Author
[Denis Mitrofanov](https://thecollection.ru)

### Usage

```go get https://github.com/denismitr/go-hashids/v1```

##### Public API

Import
```go
import (
	hashids "github.com/denismitr/go-hashids/v1"
)
```

Available constants:

```go
// DefaultAlphabet - with all latin letters and all digits
DefaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
// LowercaseAlphabetWithDigits all latin lowercase characters and digits
LowercaseAlphabetWithDigits = "abcdefghijklmnopqrstuvwxyz1234567890"
// DefaultLength of the hash which is basically a minimal length of the hash
// Length will grow automatically as required
DefaultLength = 16
// MinAlphabetLength - custome alphabet cannot be smaller than this value
MinAlphabetLength = 16
```

Custom options
```go
options := hashids.Options{
    Length:   16,
    Salt:     "some salt",
    Alphabet: hashids.DefaultAlphabet,
}

h, err := hashids.New(options)
if err != nil {
    log.Fatal(err)
}

hash, err := h.Encode(1)
if err != nil {
    log.Fatal(err)
}

numbers, err := h.Decode(hash).Unwrap()
// numbers == []int64{1}
```

Default options
```go
h, _ := hashids.New(hashids.DefaultOptions("my salt"))

hash, _ := h.Encode(1, 2, 3)

numbers, _ := h.Decode(hash).Unwrap()
// numbers == []int64{1, 2, 3}
```

Available input formats
* []int64
* []int
* int64
* int
* hexidecimal string
* time.Time

```go
options := hashids.Options{
    Length:   8,
    Salt:     "my salt",
    Alphabet: hashids.DefaultAlphabet,
}


h, err := hashids.New(options)
if err != nil {
    log.Fatal(err)
}

// single value
hash, _ := h.Encode(1) 
// or
hash, _ := h.Encode(int64(1)) 
// or multiple values
hash, _ := h.Encode(1, 2, 3)
// or
hash, _ := h.Encode([]int{1, 2, 3})
// or
hash, _ := h.Encode([]int64{1, 2, 3})
// hexidecimal string
hash, _ := h.Encode("ab1f")
// time.Time
hash, _ := h.Encode(time.Now())
```

#### Hexidecimal strings
Another supported format is hexidecimal strings
```go
// you can use a special method for it
// that is designed aspecially 
hash, _ := h.EncodeHex("ABECDF53")

hex, err := h.Decode(hash).AsHex()
if err != nil {
    log.Fatal(err)
}

// hex = "abecdf53" ATTENTION!!! all lower case
```

#### Optional prefixing - making a Stripe style slug
```go
options := hashids.Options{
    Length:   12,
    Salt:     "some salt",
    Alphabet: hashids.LowercaseAlphabetWithDigits,
    Prefix:   "cus_",
}

h, err := New(options)
if err != nil {
    t.Fatal(err)
}

hash, err := h.Encode(156)
if err != nil {
    log.Fatal(err)
}

// hash == cus_2vk4e9xpeng7

numbers, err := h.Decode(hash).Unwrap()
if err != nil {
    log.Fatal(err)
}

// numbers == []int64{156}
// prefix will be stripped automatically during decode
// as long as it was specified in the options or via a setter before decode
```

You may not always want to specify prefix when creating a new hasher (even though it is recommended). You have an option to set prefix and clear it with dedicated methods.

```go

options := hashids.Options{
    Length:   12,
    Salt:     "some salt",
    Alphabet: hashids.LowercaseAlphabetWithDigits,
}

h, err := New(options)
if err != nil {
    t.Fatal(err)
}

hash, err := h.SetPrefix("cus_").Encode(156)
if err != nil {
    log.Fatal(err)
}

// hash == cus_2vk4e9xpeng7

numbers, err := h.Decode(hash).Unwrap()
if err != nil {
    log.Fatal(err)
}

// numbers == []int64{156}
// prefix will be stripped automatically during decode
// as long as it was specified in the options or via a setter before decode

h.ClearPrefix() // you can use this method to clear the prefix
```

#### Working with timestamps
ATTENTION!!! Use this feature with caution. If you wany to create hashid from a timestamp, there is always a chance that in a concurrent application two timestamps generated in two different processes, goroutines or simply web requests may actually turn out to be totally identical up to a nanosecond.

```go
t := time.Now() // just for example

h, err := New(hashids.DefaultOptions("salt"))
if err != nil {
    t.Fatal(err)
}

hash, err := h.EncodeTime(t)
if err != nil {
    t.Fatal(err)
}

u, err := h.Decode(hash).AsTime()
if err != nil {
    t.Fatal(err)
}

t.Equal(u) // true
t.Sub(u).Nanoseconds() // 0 delta in nanoseconds 
```

