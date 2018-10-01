package hashids

func createNumbersHashFor(slice []int64) int64 {
	nh := int64(0)
	for i, n := range slice {
		nh += (n % int64(i+100))
	}
	return nh
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
