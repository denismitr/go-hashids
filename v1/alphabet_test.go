package hashids

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_MinLengthPadding(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		input   interface{}
		decoded []int64
		length  int
	}{
		{"1 gets padded to 8 symbols", int64(1), []int64{1}, 8},
		{"10 gets padded to 16 symbols", 10, []int64{10}, 16},
		{"22 gets padded to 22 symbols", []int{10, 23, 56}, []int64{10, 23, 56}, 22},
		{"36 gets padded to 36 symbols", []int64{10, 23, 56}, []int64{10, 23, 56}, 36},
	}

	options := DefaultOptions("test salt")

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			options.Length = tc.length
			o, err := New(options)
			if err != nil {
				t.Fatal(err)
			}

			hash, err := o.Encode(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			if len(hash) != tc.length {
				t.Fatalf("Expected hash to have length of %d, got %d. Hash is %s", tc.length, len(hash), hash)
			}

			decoded, _ := o.Decode(hash).Unwrap()

			if !reflect.DeepEqual(decoded, tc.decoded) {
				t.Fatalf("Expected decoded result %v to be equal to input %v", decoded, tc.decoded)
			}
		})
	}
}

func Test_CustomAlphabet(t *testing.T) {
	t.Parallel()

	alphs := []string{
		"abcdefghPJHGDTWQSC",
		"B)Wdkjwpouyftgnc!-hA",
		"POVSMDCKARZXIEQJLB",
		"123456789QWERTYL",
		"ROF92813jdh4kpifcC7l",
		"cCsSfFhHuUiItT01",
		"abdegjklCFHISTUc",
		"abdegjklmnopqrSF",
		"abdegjklmnopqrvwxyzABDEGJKLMNOPQRVWXYZ1234567890",
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890~!@#$%^&*()-_=+|",
		"~!@#$%^&*()-_=+{[}]",
	}

	for _, alph := range alphs {
		t.Run(fmt.Sprintf("Alphabet %s", alph), func(t *testing.T) {
			options := DefaultOptions("test salt")
			options.Length = 0
			options.Alphabet = alph

			o, err := New(options)
			if err != nil {
				t.Fatal(err)
			}

			hash, err := o.Encode(1, 2, 3)
			if err != nil {
				t.Fatal(err)
			}

			decoded, _ := o.Decode(hash).Unwrap()
			expected := []int64{1, 2, 3}

			if !reflect.DeepEqual(decoded, expected) {
				t.Fatalf("Expected decoded result %v to be equal to input 1, 2, 3", decoded)
			}
		})
	}
}
