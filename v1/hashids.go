package hashids

import (
	"fmt"
	"math"
)

// Hasher is responsible for the encoding and decoding
type Hasher struct {
	options            Options
	maxLengthPerNumber int

	numbers []int64
	hash    []rune
	buf     []rune
}

// New obfuscator
func New(options Options) (*Hasher, error) {
	err := options.Initialize()
	if err != nil {
		return nil, err
	}

	h := &Hasher{
		options: options,
	}

	// Calculate the maximum possible string length by hashing the maximum possible id
	encoded, err := h.Encode(math.MaxInt64)
	if err != nil {
		return nil, fmt.Errorf("unable to encode maximum int64 to find max encoded value length: %s", err)
	}

	h.maxLengthPerNumber = len(encoded)

	return h, nil
}

// Encode number, numbers or slice of numbers
func (h *Hasher) Encode(v ...interface{}) (string, error) {
	h.reset()

	if len(v) == 0 {
		return "", fmt.Errorf("expected at least 1 value")
	}

	for _, item := range v {
		switch value := item.(type) {
		case []int64:
			h.numbers = value
		case []int:
			for _, n := range value {
				h.numbers = append(h.numbers, int64(n))
			}
		case int64:
			h.numbers = append(h.numbers, int64(value))
		case int:
			h.numbers = append(h.numbers, int64(value))
		case string:
			if isHex(value) {
				nums, err := hexToNums(value)
				if err != nil {
					return "", err
				}
				h.numbers = nums
			} else {
				return "", fmt.Errorf("unkown format of string")
			}
		default:
			return "", fmt.Errorf("Value must be of type int, int64 or slice of ints, or hex")
		}
	}

	return h.encodeNumbers()
}

// Decode string hash
func (h *Hasher) Decode(input string) *DecodedNumbers {
	h.reset()

	hashRunes := separate([]rune(input), h.options.guards)
	i := 0
	if len(hashRunes) == 2 || len(hashRunes) == 3 {
		i = 1
	}

	breakdown := hashRunes[i]

	if len(breakdown) == 0 {
		breakdown = hashRunes[0]
	}

	if len(breakdown) > 0 {
		lottery := breakdown[0]
		breakdown = breakdown[1:]
		hashRunes = separate(breakdown, h.options.seps)
		alphabet := h.options.alphabetCopy()
		for _, rs := range hashRunes {
			h.buf = h.buf[:1]
			h.buf[0] = lottery
			h.buf = append(h.buf, h.options.salt...)
			h.buf = append(h.buf, alphabet...)
			alphabet = shuffle(alphabet, h.buf[:len(alphabet)])
			number, err := unhash(rs, alphabet)
			if err != nil {
				return NewDecodedNumbers(nil, err)
			}
			h.numbers = append(h.numbers, number)
		}
	}

	check, _ := h.Encode(h.numbers)
	if check != input {
		return NewDecodedNumbers(nil,
			fmt.Errorf("mismatch between encoded and decoded values: %s -> %s, obtained result %v", check, input, h.numbers))
	}

	return NewDecodedNumbers(h.numbers, nil)
}

func (h *Hasher) encodeNumbers() (string, error) {
	if len(h.numbers) == 0 {
		return "", fmt.Errorf("cannot encode an empty slice of numbers")
	}

	for _, n := range h.numbers {
		if n < 0 {
			return "", fmt.Errorf("negative numbers like %d are not allowed", n)
		}
	}

	alphabet := h.options.alphabetCopy()
	numbersHash := createNumbersHashFor(h.numbers)
	lottery := alphabet[numbersHash%int64(len(alphabet))]

	h.hash = append(h.hash, lottery)

	for i, n := range h.numbers {
		h.buf = h.buf[:1]
		h.buf[0] = lottery
		h.buf = append(h.buf, h.options.saltCopy()...)
		h.buf = append(h.buf, alphabet...)
		alphabet = shuffle(alphabet, h.buf[:len(alphabet)])

		hashSlice := hash(n, alphabet)
		h.hash = append(h.hash, hashSlice...)

		if i < len(h.numbers)-1 {
			n %= int64(hashSlice[0]) + int64(i)
			h.hash = append(h.hash, h.options.seps[n%int64(len(h.options.seps))])
		}
	}

	if len(h.hash) < h.options.MinLength {
		i := (numbersHash + int64(h.hash[0])) % int64(len(h.options.guards))
		h.hash = append([]rune{h.options.guards[i]}, h.hash...)

		if len(h.hash) < h.options.MinLength {
			i := (numbersHash + int64(h.hash[2])) % int64(len(h.options.guards))
			h.hash = append(h.hash, h.options.guards[i])
		}
	}

	middle := len(alphabet) / 2
	for len(h.hash) < h.options.MinLength {
		alphabet = shuffle(alphabet, alphabet)
		h.hash = append(alphabet[middle:], append(h.hash, alphabet[:middle]...)...)
		excess := len(h.hash) - h.options.MinLength
		if excess > 0 {
			h.hash = h.hash[excess/2 : excess/2+h.options.MinLength]
		}
	}

	return string(h.hash), nil
}

func (h Hasher) getMaxResultLengthFor(slice []int64) int {
	maxLength := h.maxLengthPerNumber * len(slice)
	if maxLength < h.options.MinLength {
		return h.options.MinLength
	}
	return maxLength
}

// reset the hash
func (h *Hasher) reset() {
	h.numbers = make([]int64, 0)
	h.hash = make([]rune, 0, h.options.MinLength)
	h.buf = make([]rune, 0, len(h.options.alphabet)+len(h.options.salt)+1)
}
