package hashids

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EndodedAndDecodedValuesAreEqual(t *testing.T) {
	t.Parallel()

	tt := []struct {
		in  interface{}
		out []int64
	}{
		{1, []int64{1}},
		{123, []int64{123}},
		{54, []int64{54}},
		{48, []int64{48}},
		{[]int{123, 88}, []int64{123, 88}},
		{100, []int64{100}},
		{555777999, []int64{555777999}},
		{10, []int64{10}},
		{[]int64{1, 2, 3, 4, 5}, []int64{1, 2, 3, 4, 5}},
		{[]int64{29, 30, 26, 29, 27, 30, 30, 31}, []int64{29, 30, 26, 29, 27, 30, 30, 31}},
	}

	o, _ := New(DefaultOptions("test salt"))

	for _, tc := range tt {
		tc := tc
		t.Run(fmt.Sprintf("Input %v output %v", tc.in, tc.out), func(t *testing.T) {

			hash, err := o.Encode(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			result := o.Decode(hash)

			if result.HasError() {
				t.Fatal(result.Error())
			}

			if result.Len() == 0 {
				t.Fatalf(
					"Expected result list to have %d items, instead there is %d",
					len(tc.out),
					result.Len())
			}

			actual := result.AsInt64Slice()

			if !reflect.DeepEqual(actual, tc.out) {
				t.Fatalf("Expected result to be %v items, instead got %v", tc.out, actual)
			}
		})
	}
}

func Test_EncodeReturnsCorrectHash(t *testing.T) {
	t.Parallel()

	tt := []struct {
		input  interface{}
		hash   string
		salt   string
		length int
	}{
		{[]int{45, 434, 1313, 99}, "7nnhzEsDkiYa", "this is my salt", 8},
		{[]int{45, 434, 1313, 99}, "nG7nnhzEsDkiYadK", "this is my salt", 16},
		{1, "B0NV05", "this is my salt", 6},
		{1, "QGQ707", "this is another salt", 6},
		{[]int{1}, "B0NV05", "this is my salt", 6},
		{2, "yLA6m0oM", "this is my salt", 8},
		{[]int64{2}, "yLA6m0oM", "this is my salt", 8},
		{1, "JEDngB0NV05ev1Ww", "this is my salt", 16},
		{1, "b9qVeQGQ707ay8Kl", "this is another salt", 16},
		{1000, "Xzjd5vJGvO", "this is my salt", 10},
		{[]int64{1000}, "Xzjd5vJGvO", "this is my salt", 10},
		{[]int64{1, 10, 1000}, "40rlHmFyQd", "this is my salt", 10},
		{[]int64{1, 10, 1000}, "303gcXFo60", "this is another salt", 10},
		{[]int64{2, 24, 234567810}, "nG2fJTDWGebV", "test salt", 12},
		{[]int64{2, 24, 234567810}, "w9XIviZljBvY", "another test salt", 12},
		{[]int64{2, 24, 234567810}, "rBwGnG2fJTDWGebVP24d", "test salt", 20},
		{[]int64{29, 30, 26, 29, 27, 30, 30, 31}, "lGDRWVzyXIkflC6IbSGfyfvqBM7m8w", "test salt", 30},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(fmt.Sprintf("%v", tc.input), func(t *testing.T) {
			options := Options{
				Length: tc.length,
				Salt:   tc.salt,
			}

			o, err := New(options)
			if err != nil {
				t.Fatal(err)
			}

			hash, err := o.Encode(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.hash, hash)
		})
	}
}

func Test_EncodeReturnsCorrectInput(t *testing.T) {
	t.Parallel()

	tt := []struct {
		output []int64
		hash   string
		salt   string
		length int
	}{
		{[]int64{45, 434, 1313, 99}, "7nnhzEsDkiYa", "this is my salt", 8},
		{[]int64{45, 434, 1313, 99}, "nG7nnhzEsDkiYadK", "this is my salt", 16},
		{[]int64{1}, "B0NV05", "this is my salt", 6},
		{[]int64{1}, "QGQ707", "this is another salt", 6},
		{[]int64{1}, "EDngB0NV05ev1W", "this is my salt", 14},
		{[]int64{2}, "yLA6m0oM", "this is my salt", 8},
		{[]int64{1}, "JEDngB0NV05ev1Ww", "this is my salt", 16},
		{[]int64{1}, "b9qVeQGQ707ay8Kl", "this is another salt", 16},
		{[]int64{1000}, "Xzjd5vJGvO", "this is my salt", 10},
		{[]int64{1, 10, 1000}, "40rlHmFyQd", "this is my salt", 10},
		{[]int64{1, 10, 1000}, "303gcXFo60", "this is another salt", 10},
		{[]int64{2, 24, 234567810}, "nG2fJTDWGebV", "test salt", 12},
		{[]int64{2, 24, 234567810}, "w9XIviZljBvY", "another test salt", 12},
		{[]int64{2, 24, 234567810}, "rBwGnG2fJTDWGebVP24d", "test salt", 20},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(fmt.Sprintf("%s", tc.hash), func(t *testing.T) {
			options := Options{
				Length: tc.length,
				Salt:   tc.salt,
			}

			o, err := New(options)
			if err != nil {
				t.Fatal(err)
			}

			result, err := o.Decode(tc.hash).Unwrap()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.output, result)
		})
	}
}

func Test_ErrorOnDecode(t *testing.T) {
	tt := []struct {
		alphabet string
		hash     string
		err      string
	}{
		{"Alphabet1234567890", "uuuiQO", "alphabet that was used for hashing was different"},
		{"Alphabet1234567890", "QQAlphabet", "alphabet that was used for hashing was different"},
		{"Alphabet1234567890", "CPQAlphabet34", "alphabet that was used for hashing was different"},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.alphabet+"_"+tc.hash, func(t *testing.T) {
			options := DefaultOptions("test salt")
			options.Alphabet = tc.alphabet
			options.Length = 8

			o, err := New(options)
			if err != nil {
				t.Fatal(err)
			}

			decoded, err := o.Decode(tc.hash).Unwrap()
			if decoded != nil {
				t.Fatalf("expected nil, got %v", decoded)
			}

			if err.Error() != tc.err {
				t.Fatalf("excpected error message to be %s, got %s", tc.err, err.Error())
			}
		})
	}
}

func Test_SaltError(t *testing.T) {
	tt := []struct {
		encodeSalt string
		decodeSalt string
		input      interface{}
	}{
		{"salt A", "salt B", 2},
		{"test salt", "wrong salt", 1345},
		{"good salt", "bad salt", []int{40, 1239, 456}},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.encodeSalt+"~>"+tc.decodeSalt, func(t *testing.T) {
			options := DefaultOptions(tc.encodeSalt)
			o, _ := New(options)

			hash, err := o.Encode(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			options = DefaultOptions(tc.decodeSalt)
			o, _ = New(options)

			result, err := o.Decode(hash).Unwrap()
			if result != nil {
				t.Fatalf("expected result to be nil, got %v", result)
			}

			if !strings.Contains(err.Error(), "mismatch between encoded and decoded values") {
				t.Fatalf("expected mismatch error, got %v", err)
			}
		})
	}
}

func Test_DefaultOptions_Length(t *testing.T) {
	options := DefaultOptions("this is my salt")

	o, _ := New(options)

	numbers := []int64{45, 434, 1313, 99}
	hash, err := o.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}

	if len(hash) != 16 {
		t.Fatalf("Expected hash length to be 16, got %d", len(hash))
	}

	result, err := o.Decode(hash).Unwrap()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, numbers, result)
}

func Test_HexEncodedAndDecodedValuesAreEqual(t *testing.T) {
	tt := []struct {
		hex      string
		expected string
		salt     string
		length   int
	}{
		{"5a74d76ac89b05000e977baa", "qmTqfesOIqHrsoCYf9UkFZixSKuBT4umuruXuMiDsVsbSrfV", "this is my salt", 30},
		{"5a74d76ac89b05000e977baa", "r6sdC0iBF5IXiZU3CQuLT1HJSofDs3fvfMfXfdHjivibS8Cw", "test salt", 18},
		{"111aff", "5JqQ5h6hYhjCyhgqjL", "test salt", 18},
		{"111affe", "Yzx1hmh6hKCyhvhNPb", "test salt", 18},
		{"1a", "XBe7QdP7Wh5PMa8Ojy", "test salt", 18},
		{"1", "O35oKBgz41PVdL9MQA", "test salt", 18},
		{"1", "gz41PV", "test salt", 6},
		{"2", "9VbdrOYnAYnxDlLEWj", "test salt", 18},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.hex, func(t *testing.T) {
			options := DefaultOptions(tc.salt)
			options.Length = tc.length

			o, _ := New(options)

			hash, err := o.Encode(tc.hex)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.expected, hash)

			hex, err := o.Decode(hash).AsHex()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.hex, hex)
		})
	}
}
