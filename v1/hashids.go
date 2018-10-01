package hashids

import (
	"fmt"
)

const (
	defaultAlphabet  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	defaultMinLength = 8

	sepDiv      = 3.5
	guardDiv    = 12.0
	defaultSeps = "cfhistuCFHISTU"
)

type Result []int64

// Options for the obfuscator
type Options struct {
	Alphabet string

	MinLength int

	Salt string

	Seps string

	Guargs string
}

// DefaultOptions for the obfuscator
func DefaultOptions() Options {
	return Options{
		Alphabet:  defaultAlphabet,
		MinLength: defaultMinLength,
		Seps:      defaultSeps,
	}
}

func (o Options) AlphabetSlice() []rune {
	return []rune(o.Alphabet)
}

type Obfuscator struct {
	options            Options
	guards             []rune
	maxLengthPerNumber int
}

func New(options Options) (Obfuscator, error) {
	guards := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
	return Obfuscator{
		options:            options,
		guards:             guards,
		maxLengthPerNumber: 5,
	}, nil
}

func (o Obfuscator) Encode(v ...interface{}) (string, error) {
	if len(v) == 0 {
		return "", fmt.Errorf("expected at least 1 value")
	}

	var slice []int64

	for _, item := range v {
		switch value := item.(type) {
		case int64:
			slice = []int64{value}
		case int:
			slice = []int64{int64(value)}
		default:
			return "", fmt.Errorf("Value must be of type int64")
		}
	}

	return o.encodeSlice(slice)
}

func (o Obfuscator) encodeSlice(slice []int64) (string, error) {
	var result []rune

	alphabetSlice := o.options.AlphabetSlice()
	numbersHash := createNumbersHashFor(slice)
	maxResultLength := o.getMaxResultLengthFor(slice)

	// buffer := make([]rune, len(o.options.Alphabet)+len(o.options.Salt)+1)

	for _, n := range slice {
		result = append(result, hash(n, alphabetSlice)...)
	}

	for len(result) < o.options.MinLength {
		result = append(result, alphabetSlice[len(result)])
	}

	return string(result), nil
}

func (o Obfuscator) getMaxResultLengthFor(slice []int64) int {
	maxLength := o.maxLengthPerNumber * len(slice)
	if maxLength < o.options.MinLength {
		return o.options.MinLength
	}
	return maxLength
}

func (o Obfuscator) Decode(hash string) Decoded {
	return Decoded{}
}

func hash(number int64, alphabet []rune) []rune {
	var hash []rune
	alphabetLength := int64(len(alphabet))

	for {
		s := alphabet[number%alphabetLength]
		hash = append(hash, s)
		number /= alphabetLength
		if number == 0 {
			break
		}
	}

	for i := len(hash)/2 - 1; i >= 0; i-- {
		j := len(hash) - 1 - i
		hash[i], hash[j] = hash[j], hash[i]
	}

	return hash
}
