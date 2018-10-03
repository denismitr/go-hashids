package hashids

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndodeAndDecodeValuesAreEqual(t *testing.T) {
	tt := []struct {
		name   string
		encode interface{}
		result Decoded
	}{
		{"1 -> 1", 1, Decoded{[]int64{1}, nil}},
		{"123 -> 123", 123, Decoded{[]int64{123}, nil}},
	}

	o, _ := New(DefaultOptions())

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := o.Encode(tc.encode)
			if err != nil {
				t.Fatal(err)
			}

			result := o.Decode(hash)

			if result.HasError() {
				t.Fatal(result.Error())
			}

			if result.Count() == 0 {
				t.Fatalf(
					"Expected result list to have %d items, instead there is %d",
					tc.result.Count(),
					result.Count())
			}

			result.Map(func(v int64, i int) int64 {
				log.Fatal(v, i)
				assert.Equal(t, v, tc.result.result[i])
				return v
			})
		})
	}
}
