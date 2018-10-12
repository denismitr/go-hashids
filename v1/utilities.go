package hashids

import (
	"encoding/hex"
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

func hexToNums(hex string) ([]int64, error) {
	nums := make([]int64, len(hex))

	for i := 0; i < len(hex); i++ {
		b := hex[i]
		switch {
		case (b >= '0') && (b <= '9'):
			b -= '0'
		case (b >= 'a') && (b <= 'f'):
			b -= 'a' - 'A'
			fallthrough
		case (b >= 'A') && (b <= 'F'):
			b -= ('A' - 0xA)
		default:
			return nil, fmt.Errorf("invalid hex digit")
		}
		// Each int is in range [16, 31]
		nums[i] = 0x10 + int64(b)
	}

	return nums, nil
}

func numsToHex(nums []int64) (string, error) {
	const hex = "0123456789abcdef"

	b := make([]byte, len(nums))

	for i, n := range nums {
		if n < 0x10 || n > 0x1f {
			return "", fmt.Errorf("invalid number")
		}
		b[i] = hex[n-0x10]
	}

	return string(b), nil
}

func isHex(s string) bool {
	_, err := hex.DecodeString(s)
	if err != nil {
		_, err := hex.DecodeString(s + "1")
		if err != nil {
			return false
		}
	}

	return true
}
