package hashids

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

// AlphabetAsSlice to use in algotithm
func (o Options) alphabetAsSlice() []rune {
	return []rune(o.Alphabet)
}

// SaltAsSlice to use in algotithm
func (o Options) saltAsSlice() []rune {
	return []rune(o.Salt)
}
