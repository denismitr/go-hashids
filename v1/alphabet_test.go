package hashids

import (
	"testing"
)

func Test_DefaultLength(t *testing.T) {
	tt := []struct {
		name   string
		encode interface{}
		length int
	}{
		{"1 gets padded to 8 symbols", int64(1), 8},
		{"10 gets padded to 16 symbols", 10, 8},
	}

	o := New(DefaultOptions())

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := o.Encode(tc.encode)
			if err != nil {
				t.Fatal(err)
			}

			if len(hash) < tc.length {
				t.Fatalf("Expected hash to have length of %d, got %d. Hash is %s", tc.length, len(hash), hash)
			}
		})
	}
}
