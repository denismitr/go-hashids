package hashids

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_EncodeEpochAndDecodeEpoch(t *testing.T) {
	t.Parallel()

	tt := []struct {
		time   time.Time
		length int
		salt   string
	}{
		{time.Now().Add(5 * time.Second), 40, "test salt"},
		{time.Now(), 8, "test salt"},
		{time.Now().Add(-100 * time.Hour), 10, "my salt"},
		{time.Now().Add(-10000 * time.Hour), 30, "test salt"},
		{time.Now().Add(70000 * time.Hour), 30, "test salt"},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(fmt.Sprintf("Time: %s", tc.time), func(t *testing.T) {
			options := Options{
				Length:   tc.length,
				Salt:     tc.salt,
				Alphabet: DefaultAlphabet,
			}

			h, err := New(options)
			if err != nil {
				t.Fatal(err)
			}

			hash, err := h.EncodeTime(tc.time)
			if err != nil {
				t.Fatal(err)
			}

			u, err := h.Decode(hash).AsTime()
			if err != nil {
				t.Fatal(err)
			}

			assert.True(t, tc.time.Equal(u))
			assert.Equal(t, int64(0), tc.time.Sub(u).Nanoseconds())
		})
	}
}
