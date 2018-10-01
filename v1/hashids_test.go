package hashids

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndodeAndDecodeValuesAreEqual(t *testing.T) {
	tt := []struct {
		name    string
		encode  interface{}
		decoded Decoded
	}{
		{"1 -> 1", 1, Decoded{[]int64{1}, nil}},
		{"123 -> 123", 123, Decoded{[]int64{1}, nil}},
	}

	o := New(DefaultOptions())

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := o.Encode(tc.encode)
			if err != nil {
				t.Fatal(err)
			}

			decoded := o.Decode(hash)

			if decoded.Count() == 0 {
				t.Fatalf(
					"Expected result list to have %d items, instead there is %d",
					tc.decoded.Count(),
					decoded.Count())
			}

			decoded.Map(func(v int64, i int) int64 {
				log.Fatal(v, i)
				assert.Equal(t, v, tc.decoded.result[i])
				return v
			})
		})
	}
}
