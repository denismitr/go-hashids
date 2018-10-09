package hashids

import (
	"errors"
	"fmt"
)

func createNumbersHashFor(slice []int64) int64 {
	nh := int64(0)
	for i, n := range slice {
		nh += (n % int64(i+100))
	}
	return nh
}

func hash(in int64, alphabet []rune) []rune {
	out := make([]rune, 0)
	alphabetLength := int64(len(alphabet))

	for {
		s := alphabet[in%alphabetLength]
		out = append(out, s)
		in /= alphabetLength
		if in == 0 {
			break
		}
	}

	for i := len(out)/2 - 1; i >= 0; i-- {
		j := len(out) - 1 - i
		out[i], out[j] = out[j], out[i]
	}

	return out
}

func unhash(in, alphabet []rune) (out int64, err error) {
	for _, r := range in {
		pos := -1
		for i, s := range alphabet {
			if r == s {
				pos = i
				break
			}
		}

		if pos == -1 {
			err = fmt.Errorf("alphabet that was used for hashing was different")
			return
		}

		out = out*int64(len(alphabet)) + int64(pos)
	}

	return
}

func separate(in, seps []rune) (out [][]rune) {
	indicies := make([]int, 0)
	for i, r := range in {
		for _, s := range seps {
			if r == s {
				indicies = append(indicies, i)
			}
		}
	}

	out = make([][]rune, 0, len(indicies)+1)
	left := in[:]
	for _, idx := range indicies {
		idx -= len(in) - len(left)
		out = append(out, left[:idx])
		left = left[idx+1:]
	}

	out = append(out, left)
	return
}

func shuffle(in, salt []rune) (out []rune) {
	if len(salt) == 0 {
		out = in
		return
	}

	out = make([]rune, len(in))
	copy(out, in)
	p, v := 0, 0

	for i := len(in) - 1; i > 0; i-- {
		p += int(salt[v])
		j := (int(salt[v]) + v + p) % i
		out[i], out[j] = out[j], out[i]
		v = (v + 1) % len(salt)
	}

	return
}

func shuffleInPlace(alphabet []rune, salt []rune) {
	if len(salt) == 0 {
		return
	}

	for i, v, p := len(alphabet)-1, 0, 0; i > 0; i-- {
		p += int(salt[v])
		j := (int(salt[v]) + v + p) % i
		alphabet[i], alphabet[j] = alphabet[j], alphabet[i]
		v = (v + 1) % len(salt)
	}
	return
}

func encode(numbers []int64, alphabet, salt, seps, guards []rune, maxLengthPerNumber, minLength int) (string, error) {
	if len(numbers) == 0 {
		return "", errors.New("encoding empty array of numbers makes no sense")
	}
	for _, n := range numbers {
		if n < 0 {
			return "", errors.New("negative number not supported")
		}
	}

	numbersHash := int64(0)
	for i, n := range numbers {
		numbersHash += (n % int64(i+100))
	}

	maxRuneLength := maxLengthPerNumber * len(numbers)
	if maxRuneLength < minLength {
		maxRuneLength = minLength
	}

	result := make([]rune, 0, maxRuneLength)
	lottery := alphabet[numbersHash%int64(len(alphabet))]
	result = append(result, lottery)
	hashBuf := make([]rune, maxRuneLength)
	buffer := make([]rune, len(alphabet)+len(salt)+1)

	for i, n := range numbers {
		buffer = buffer[:1]
		buffer[0] = lottery
		buffer = append(buffer, salt...)
		buffer = append(buffer, alphabet...)
		shuffleInPlace(alphabet, buffer[:len(alphabet)])
		hashBuf = hashInPlace(n, alphabet, hashBuf)
		result = append(result, hashBuf...)

		if i+1 < len(numbers) {
			n %= int64(hashBuf[0]) + int64(i)
			result = append(result, seps[n%int64(len(seps))])
		}
	}

	if len(result) < minLength {
		guardIndex := (numbersHash + int64(result[0])) % int64(len(guards))
		result = append([]rune{guards[guardIndex]}, result...)

		if len(result) < minLength {
			guardIndex = (numbersHash + int64(result[2])) % int64(len(guards))
			result = append(result, guards[guardIndex])
		}
	}

	halfLength := len(alphabet) / 2
	for len(result) < minLength {
		shuffleInPlace(alphabet, duplicateRuneSlice(alphabet))
		result = append(alphabet[halfLength:], append(result, alphabet[:halfLength]...)...)
		excess := len(result) - minLength
		if excess > 0 {
			result = result[excess/2 : excess/2+minLength]
		}
	}

	return string(result), nil
}

func duplicateRuneSlice(data []rune) []rune {
	result := make([]rune, len(data))
	copy(result, data)
	return result
}

func hashInPlace(input int64, alphabet []rune, result []rune) []rune {
	result = result[:0]
	for {
		r := alphabet[input%int64(len(alphabet))]
		result = append(result, r)
		input /= int64(len(alphabet))
		if input == 0 {
			break
		}
	}
	for i := len(result)/2 - 1; i >= 0; i-- {
		opp := len(result) - 1 - i
		result[i], result[opp] = result[opp], result[i]
	}
	return result
}
