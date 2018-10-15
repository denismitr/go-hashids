package hashids

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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

			h, err := New(options)
			if err != nil {
				t.Fatal(err)
			}

			hash, err := h.Encode(1, 2, 3)
			if err != nil {
				t.Fatal(err)
			}

			decoded, _ := h.Decode(hash).Unwrap()
			expected := []int64{1, 2, 3}

			if !reflect.DeepEqual(decoded, expected) {
				t.Fatalf("Expected decoded result %v to be equal to input 1, 2, 3", decoded)
			}
		})
	}
}

func Test_CustomAlphabetWithPrefix(t *testing.T) {
	t.Parallel()

	tt := []struct {
		value    interface{}
		hash     string
		salt     string
		alphabet string
		prefix   string
		length   int
	}{
		{[]int64{1}, "joed16", "this is my salt", LowercaseAlphabetWithDigits, "cus_", 6},
		{[]int64{156}, "2vk4e9xpeng7", "some salt", LowercaseAlphabetWithDigits, "cus_", 12},
		{[]int64{1}, "0NV0", "this is my salt", DefaultAlphabet, "user_", 4},
		{[]int64{1}, "тлпирп", "this is another salt", "абвгдежзиклмнпрсто1234", "order_", 6},
		{[]int64{1, 3, 7}, "м2н8оБоз", "this is test salt", "98АБВГДЕжзиклмнпрсто1234", "преф_", 8},
		{[]int64{1234, 33, 79}, "AB992794A6", "this is test salt", "1234567890_!&*BAZ", "baz_", 8},
		{[]int64{1234, 33, 79}, "A*_8AB992794A687", "this is test salt", "1234567890_!&*BAZ", "", 16},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(fmt.Sprintf("input %v", tc.value), func(t *testing.T) {
			options := Options{
				Length:   tc.length,
				Salt:     tc.salt,
				Alphabet: tc.alphabet,
				Prefix:   tc.prefix,
			}

			h, err := New(options)
			if err != nil {
				t.Fatal(err)
			}

			hash, err := h.Encode(tc.value)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.prefix+tc.hash, hash)

			result, err := h.Decode(hash).Unwrap()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.value, result)
		})
	}
}

func Test_CustomAlphabetWithExplicitPrefix(t *testing.T) {
	t.Parallel()

	tt := []struct {
		value    interface{}
		hash     string
		salt     string
		alphabet string
		prefix   string
		length   int
	}{
		{[]int64{1}, "joed16", "this is my salt", LowercaseAlphabetWithDigits, "cus_", 6},
		{[]int64{156}, "2vk4e9xpeng7", "some salt", LowercaseAlphabetWithDigits, "cus_", 12},
		{[]int64{1}, "0NV0", "this is my salt", DefaultAlphabet, "user_", 4},
		{[]int64{1}, "тлпирп", "this is another salt", "абвгдежзиклмнпрсто1234", "order_", 6},
		{[]int64{1, 3, 7}, "м2н8оБоз", "this is test salt", "98АБВГДЕжзиклмнпрсто1234", "преф_", 8},
		{[]int64{1234, 33, 79}, "AB992794A6", "this is test salt", "1234567890_!&*BAZ", "baz_", 8},
		{[]int64{1234, 33, 79}, "A*_8AB992794A687", "this is test salt", "1234567890_!&*BAZ", "", 16},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(fmt.Sprintf("input %v", tc.value), func(t *testing.T) {
			options := Options{
				Length:   tc.length,
				Salt:     tc.salt,
				Alphabet: tc.alphabet,
			}

			h, err := New(options)
			if err != nil {
				t.Fatal(err)
			}

			h.SetPrefix(tc.prefix)

			hash, err := h.Encode(tc.value)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.prefix+tc.hash, hash)

			result, err := h.Decode(hash).Unwrap()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.value, result)

			h.ClearPrefix()

			result, err = h.Decode(removePrefix(hash, tc.prefix)).Unwrap()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.value, result)
		})
	}
}
