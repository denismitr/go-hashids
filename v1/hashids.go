package hashids

import "fmt"

const (
	defaultAlphabet  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	defaultMinLength = 8

	sepDiv      = 3.5
	guardDiv    = 12.0
	defaultSeps = "cfhistuCFHISTU"
)

type Result []int64

type Obfuscator struct {
	options            Options
	guards             []rune
	seps               []rune
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

// Encode number, numbers or slice of numbers
func (o Obfuscator) Encode(v ...interface{}) (string, error) {
	if len(v) == 0 {
		return "", fmt.Errorf("expected at least 1 value")
	}

	var slice []int64

	for _, item := range v {
		switch value := item.(type) {
		case []int64:
			slice = value
		case []int:
			for i, n := range value {
				slice[i] = int64(n)
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

func (o Obfuscator) encodeSlice(slice []int64) (string, error) {
	for _, n := range slice {
		if n < 0 {
			return "", fmt.Errorf("negative numbers like %d are not allowed", n)
		}
	}

	alphabetSlice := o.options.alphabetAsSlice()
	numbersHash := createNumbersHashFor(slice)
	maxResultLength := o.getMaxResultLengthFor(slice)
	lottery := alphabetSlice[numbersHash%int64(len(alphabetSlice))]
	result := make([]rune, 0, maxResultLength)

	buf := make([]rune, len(o.options.Alphabet)+len(o.options.Salt)+1)

	for i, n := range slice {
		buf = buf[:1]
		buf[0] = lottery
		buf = append(buf, o.options.saltAsSlice()...)
		buf = append(buf, alphabetSlice...)
		hashSlice := hash(n, alphabetSlice)
		result = append(result, hashSlice...)

		if i < len(slice)-1 {
			n %= int64(hashSlice[0]) + int64(i)
			result = append(result, o.seps[n%int64(len(o.seps))])
		}
	}

	if len(result) < o.options.MinLength {
		i := (numbersHash + int64(result[0])) % int64(len(o.guards))
		result = append([]rune{o.guards[i]}, result...)

		if len(result) < o.options.MinLength && len(result) > 2 {
			i := (numbersHash + int64(result[2])) % int64(len(o.guards))
			result = append(result, o.guards[i])
		}
	}

	middle := len(alphabetSlice) / 2
	for len(result) < o.options.MinLength {
		result = append(alphabetSlice[middle:], append(result, alphabetSlice[:middle]...)...)
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

func (o Obfuscator) Decode(hash string) Decoded {
	return Decoded{}
}
