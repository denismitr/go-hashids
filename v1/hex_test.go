package hashids

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EncodeHexDecodeHex(t *testing.T) {
	tt := []struct {
		hash   string
		hex    string
		length int
		salt   string
	}{
		{"vk9bLlGDRWVzyXIkflC6IbSGfyfvqBM7m8wj04KX", "deadbeef", 40, "test salt"},
		{"GLC6SKuxIZfQh6hatOc4FKsj", "abcdef123456", 8, "test salt"},
		{"WksrIAioSYSlSQCWCMCJCESoSXHoHpHbHrHaH9H6Hn", "ABCDDD6666DDEEEEEEEEE", 10, "my salt"},
		{"DnSntKI1U4UeUNIlIVHEfAUlcpCqfqs4Izu1uAi5FJuVteUO", "507f1f77bcf86cd799439011", 40, "some salt"},
		{"4bhkf1f9fmf4fDheI4IKI7I0I3IQfyfdfefVf4FmF3F5FXFbFgFYC0SZCASOCk", "f00000fddddddeeeee4444444ababab", 30, "test salt"},
		{"RyTwU6cqfjhMipi3sAtOuZCmFkT4UGcLfrhlijiYsVtAupCgFpT0UZcvfQhOi4iys9tjuLCz", "abcdef123456abcdef123456abcdef123456", 0, ""},
		{"0yfzcOcJcecgcQc5cbcGc5cEcKcyc0cLcQc0cBcrcMcvcpcGcdcacpcJcacdcGclc8cjcKcVcacmczcNcpcdcWcecocYcecMcqc5cBcXc4", "f000000000000000000000000000000000000000000000000000f", 40, "my test"},
		{"PNUXUBU2UKUNUPUYU0U2UbU6UaUNUXUOUoUqU4UwUgUDUPU6UoUYUEULUrUoU9U2UXUGUbUeUgU2UQUKUMUnUbUjUEULUAUGUxUZU9UOUY", "fffffffffffffffffffffffffffffffffffffffffffffffffffff", 20, "salt"},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.hex, func(t *testing.T) {
			options := Options{
				Length:   tc.length,
				Salt:     tc.salt,
				Alphabet: DefaultAlphabet,
			}

			h, err := New(options)
			if err != nil {
				t.Fatal(err)
			}

			hash, err := h.EncodeHex(tc.hex)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.hash, hash)

			hex, err := h.Decode(hash).AsHex()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, strings.ToLower(tc.hex), hex)
		})
	}
}
