### GO Hashids

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
```

#### Hexidecimal strings
Another supported format is hexidecimal strings
```go
// you can use a special method for it
// that is designed aspecially 
hash, _ := h.EncodeHex("ABCDDD6666DDEEEEEEEEE")

hex, err := h.Decode(hash).AsHex()
if err != nil {
    log.Fatal(err)
}

// hex = "abcddd6666ddeeeeeeeee" ATTENTION!!! all lower case
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

h.SetPrefix("cus_")

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

h.ClearPrefix() // you can use this method to clear the prefix
```

