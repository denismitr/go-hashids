package hashids

import (
	"fmt"
	"log"
	"math"
)

// Obfuscator is responsible for the encoding and decoding
type Obfuscator struct {
	options            Options
	guards             []rune
	seps               []rune
	maxLengthPerNumber int
}

// New obfuscator
func New(options Options) (*Obfuscator, error) {
	err := options.Initialize()
	if err != nil {
		return nil, err
	}

	o := &Obfuscator{
		options: options,
		guards:  options.guards,
		seps:    options.seps,
	}

	// Calculate the maximum possible string length by hashing the maximum possible id
	encoded, err := o.Encode(math.MaxInt64)
	if err != nil {
		return nil, fmt.Errorf("unable to encode maximum int64 to find max encoded value length: %s", err)
	}

	o.maxLengthPerNumber = len(encoded)

	return o, nil
}

// Encode number, numbers or slice of numbers
func (o Obfuscator) Encode(v ...interface{}) (string, error) {
	if len(v) == 0 {
		return "", fmt.Errorf("expected at least 1 value")
	}

	slice := make([]int64, 0)

	for _, item := range v {
		switch value := item.(type) {
		case []int64:
			slice = value
		case []int:
			for _, n := range value {
				slice = append(slice, int64(n))
			}
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

// Decode string hash
func (o Obfuscator) Decode(in string) Decoded {
	hashRunes := separate([]rune(in), o.guards)
	i := 0
	if len(hashRunes) > 1 && len(hashRunes) < 4 {
		i = 1
	}

	result := make([]int64, 0, 10)
	breakdown := hashRunes[i]

	if len(breakdown) > 0 {
		lottery := breakdown[0]
		breakdown = breakdown[1:]
		hashRunes = separate(breakdown, o.seps)
		alphabet := o.options.alphabet
		buf := make([]rune, len(alphabet)+len(o.options.salt))
		for _, rs := range hashRunes {
			buf = buf[:1]
			buf[0] = lottery
			buf = append(buf, o.options.salt...)
			buf = append(buf, alphabet...)
			alphabet = shuffle(alphabet, buf[:len(alphabet)])
			number, err := unhash(rs, alphabet)
			if err != nil {
				return Decoded{nil, err}
			}
			result = append(result, number)
		}
	}

	check, _ := o.Encode(result)
	if check != in {
		return Decoded{
			result: nil,
			err:    fmt.Errorf("mismatch between encode and decode: %s -> %s, obtained result %v", check, in, result),
		}
	}

	return Decoded{result, nil}
}

func (o Obfuscator) encodeSlice(slice []int64) (string, error) {
	if len(slice) == 0 {
		return "", fmt.Errorf("cannot encode an empty slice")
	}

	for _, n := range slice {
		if n < 0 {
			return "", fmt.Errorf("negative numbers like %d are not allowed", n)
		}
	}

	alphabet := o.options.alphabetCopy()
	numbersHash := createNumbersHashFor(slice)
	maxResultLength := o.getMaxResultLengthFor(slice)
	lottery := alphabet[numbersHash%int64(len(alphabet))]
	result := make([]rune, 0, maxResultLength)
	buf := make([]rune, len(alphabet)+len(o.options.salt)+1)

	log.Printf("\nAlph: %s\n numsHash: %d\n maxResultLength: %d\n lottery: %d", string(alphabet), numbersHash, maxResultLength, lottery)

	for i, n := range slice {
		buf = buf[:1]
		buf[0] = lottery
		buf = append(buf, o.options.saltAsSlice()...)
		buf = append(buf, alphabet...)
		alphabet = shuffle(alphabet, buf[:len(alphabet)])

		hashSlice := hash(n, alphabet)
		result = append(result, hashSlice...)

		if i < len(slice) {
			n %= int64(hashSlice[0]) + int64(i)
			result = append(result, o.seps[n%int64(len(o.seps))])
		}
	}

	if len(result) < o.options.MinLength {
		i := (numbersHash + int64(result[0])) % int64(len(o.guards))
		result = append([]rune{o.guards[i]}, result...)

		if len(result) < o.options.MinLength {
			i := (numbersHash + int64(result[2])) % int64(len(o.guards))
			result = append(result, o.guards[i])
		}
	}

	middle := len(alphabet) / 2
	for len(result) < o.options.MinLength {
		alphabet = shuffle(alphabet, alphabet)
		result = append(alphabet[middle:], append(result, alphabet[:middle]...)...)
		excess := len(result) - o.options.MinLength
		if excess > 0 {
			result = result[excess/2 : excess/2+o.options.MinLength]
		}
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
