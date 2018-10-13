package hashids

import (
	"fmt"
	"math"
)

const (
	// DefaultAlphabet - with all latin letters and all digits
	DefaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	// DefaultLength of the hash which is basically a minimal length of the hash
	// Length will grow automatically as required
	DefaultLength = 16
	// MinAlphabetLength - custome alphabet cannot be smaller than this value
	MinAlphabetLength = 16

	sepDiv      = 3.5
	guardDiv    = 12.0
	defaultSeps = "cfhistuCFHISTU"
)

// Options for the obfuscator
type Options struct {
	Alphabet string
	Length   int
	Salt     string

	alphabet []rune
	salt     []rune
	seps     []rune
	guards   []rune
}

// DefaultOptions for the obfuscator
func DefaultOptions(salt string) Options {
	return Options{
		Alphabet: DefaultAlphabet,
		Length:   DefaultLength,
		Salt:     salt,
	}
}

// initialize alphabet, seps, guards
func (o *Options) initialize() error {
	alphabet, err := o.validateAlphabet()
	if err != nil {
		return err
	}

	o.salt = []rune(o.Salt)
	o.alphabet = alphabet

	o.calculateSeps()
	o.createGuards()

	return nil
}

func (o Options) validateAlphabet() ([]rune, error) {
	var alphabetRunes []rune

	if o.Alphabet == "" {
		alphabetRunes = []rune(DefaultAlphabet)
	} else {
		alphabetRunes = []rune(o.Alphabet)
	}

	if len(alphabetRunes) < MinAlphabetLength {
		return nil, fmt.Errorf("Alphabet length must be at least %d", MinAlphabetLength)
	}

	unique := make(map[rune]bool, len(alphabetRunes))

	for _, r := range alphabetRunes {
		if _, ok := unique[r]; ok {
			return nil, fmt.Errorf("duplicate character in alphabet: %s", string([]rune{r}))
		}

		if r == ' ' {
			return nil, fmt.Errorf("alphabet may not contain empty spaces")
		}

		unique[r] = true
	}

	return alphabetRunes, nil
}

// AlphabetAsSlice to use in algotithm
func (o Options) alphabetCopy() []rune {
	cp := make([]rune, len(o.alphabet))
	copy(cp, o.alphabet)
	return cp
}

// saltCopy to use in algotithm
func (o Options) saltCopy() (out []rune) {
	out = make([]rune, len(o.salt))
	copy(out, o.salt)
	return
}

func (o *Options) calculateSeps() {
	o.seps = []rune(defaultSeps)

	// seps should contain only characters present in alphabet
	// alphabet should not contain seps
	for i := 0; i < len(o.seps); i++ {
		foundIndex := -1
		for j, a := range o.alphabet {
			if a == o.seps[i] {
				foundIndex = j
				break
			}
		}
		if foundIndex == -1 {
			o.seps = append(o.seps[:i], o.seps[i+1:]...)
			i--
		} else {
			o.alphabet = append(o.alphabet[:foundIndex], o.alphabet[foundIndex+1:]...)
		}
	}

	o.seps = shuffle(o.seps, o.saltCopy())

	if len(o.seps) == 0 || float64(len(o.alphabet))/float64(len(o.seps)) > sepDiv {
		sepsLength := int(math.Ceil(float64(len(o.alphabet)) / sepDiv))
		if sepsLength == 1 {
			sepsLength++
		}

		if sepsLength > len(o.seps) {
			diff := sepsLength - len(o.seps)
			o.seps = append(o.seps, o.alphabet[:diff]...)
			o.alphabet = o.alphabet[diff:]
		} else {
			o.seps = o.seps[:sepsLength]
		}
	}

	o.alphabet = shuffle(o.alphabet, o.saltCopy())
}

func (o *Options) createGuards() {
	count := int(math.Ceil(float64(len(o.alphabet)) / guardDiv))

	if len(o.alphabet) < 3 {
		o.guards = o.seps[:count]
		o.seps = o.seps[count:]
	} else {
		o.guards = o.alphabet[:count]
		o.alphabet = o.alphabet[count:]
	}
}
