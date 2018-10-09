package hashids

import (
	"fmt"
	"math"
	"strings"
)

const (
	defaultAlphabet   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	defaultMinLength  = 16
	absoluteMinLength = 6

	sepDiv      = 3.5
	guardDiv    = 12.0
	defaultSeps = "cfhistuCFHISTU"
)

// Options for the obfuscator
type Options struct {
	Alphabet  string
	MinLength int
	Salt      string

	alphabet []rune
	salt     []rune
	seps     []rune
	guards   []rune
}

// DefaultOptions for the obfuscator
func DefaultOptions(salt string) Options {
	return Options{
		Alphabet:  defaultAlphabet,
		MinLength: defaultMinLength,
		Salt:      salt,
	}
}

// Initialize alphabet, seps, guards
func (o *Options) Initialize() error {
	if o.MinLength < absoluteMinLength {
		return fmt.Errorf("Min length must be not less than %d", absoluteMinLength)
	}

	if o.Alphabet == "" {
		o.Alphabet = defaultAlphabet
	}

	o.salt = []rune(o.Salt)

	err := o.initializeAlphabet()
	if err != nil {
		return err
	}

	o.createSeps()
	o.createGuards()

	return nil
}

func (o *Options) initializeAlphabet() error {

	if len(o.Alphabet) < o.MinLength {
		return fmt.Errorf("alphabet must bt at least %d characters long", o.MinLength)
	}

	if strings.Contains(o.Alphabet, " ") {
		return fmt.Errorf("alphabet may not contain empty spaces")
	}

	o.alphabet = []rune(o.Alphabet)

	unique := make(map[rune]bool, len(o.alphabet))
	for _, r := range o.alphabet {
		if _, ok := unique[r]; ok {
			return fmt.Errorf("duplicate character in alphabet: %s", string([]rune{r}))
		}
		unique[r] = true
	}

	return nil
}

func (o Options) alphabetCopy() []rune {
	cp := make([]rune, len(o.alphabet))
	copy(cp, o.alphabet)
	return cp
}

// AlphabetAsSlice to use in algotithm
func (o Options) alphabetAsSlice() []rune {
	return o.alphabet
}

// SaltAsSlice to use in algotithm
func (o Options) saltCopy() (out []rune) {
	out = make([]rune, len(o.salt))
	copy(out, o.salt)
	return
}

func (o *Options) createSeps() {
	o.seps = []rune(defaultSeps)

	// seps should contain only characters present in alphabet; alphabet should not contains seps
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
