package hashids

import (
	"fmt"
)

func createNumbersHashFor(slice []int64) int64 {
	nh := int64(0)
	for i, n := range slice {
		nh += (n % int64(i+100))
	}
	return nh
}

func hash(in int64, alphabet []rune) (out []rune) {
	out = make([]rune, 0)
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
			err = fmt.Errorf("alphabet used when hashing was different")
			return
		}

		out = out*int64(len(alphabet)) + int64(pos)
	}

	return
}

func splitHash(in, seps []rune) (out [][]rune) {
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

	out = in
	p, v := 0, 0

	for i := len(in) - 1; i > 0; i-- {
		p += int(salt[v])
		j := (int(salt[v]) + v + p) % i
		out[i], out[j] = out[j], out[i]
		v = (v + 1) % len(salt)
	}

	return
}
