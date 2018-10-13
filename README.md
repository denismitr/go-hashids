### GO Hashids

#### Version 1

#### Author
[Denis Mitrofanov](https://thecollection.ru)

## COMING SOON!

### Usage

```go get https://github.com/denismitr/go-hashids/v1```

```go
options := hashids.Options{
    Length:   16,
    Salt:     "some salt",
    Alphabet: hashids.DefaultAlphabet,
}

h, err := hashids.New(options)
if err != nil {
    t.Fatal(err)
}

hash, err := h.Encode(1)
if err != nil {
    t.Fatal(err)
}
```