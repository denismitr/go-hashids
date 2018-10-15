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
h, err := hashids.New(hashids.DefaultOptions("my salt"))
if err != nil {
    log.Fatal(err)
}

hash, _ := h.Encode(1, 2, 3)

numbers, _ := h.Decode(hash).Unwrap()
// numbers == []int64{1, 2, 3}
```