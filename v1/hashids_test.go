package hashids

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/speps/go-hashids"
	"github.com/stretchr/testify/assert"
)

func Test_EndodeAndDecodeValuesAreEqual(t *testing.T) {
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
	}

	o, _ := New(DefaultOptions("test salt"))

	for _, tc := range tt {
		t.Run(fmt.Sprintf("Input %v output %v", tc.in, tc.out), func(t *testing.T) {

			hash, err := o.Encode(tc.in)
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
					len(tc.out),
					result.Count())
			}

			actual := result.AsInt64()

			if !reflect.DeepEqual(actual, tc.out) {
				t.Fatalf("Expected result to be %v items, instead got %v", tc.out, actual)
			}
		})
	}
}

func Test_DecodeWithKnownHash(t *testing.T) {
	options := DefaultOptions("this is my salt")
	options.MinLength = 0

	obfuscator, _ := New(options)

	hash := "7nnhzEsDkiYa"
	result, err := obfuscator.Decode(hash).Unwrap()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v -> %v", hash, result)

	expected := []int64{45, 434, 1313, 99}
	assert.Equal(t, expected, result)
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

func Test_ComperativeEncode(t *testing.T) {
	hd := hashids.NewData()
	hd.Salt = "this is my salt"
	hd.MinLength = 30
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{45, 434, 1313, 99})

	options := DefaultOptions("this is my salt")
	options.MinLength = 30

	o, _ := New(options)

	numbers := []int{45, 434, 1313, 99}
	hash, _ := o.Encode(numbers)

	if hash != e {
		t.Fatalf("\nActual: %s, expected %s", hash, e)
	}
}

func Test_ComperativeDecode(t *testing.T) {
	hash := "woQ2vqjnG7nnhzEsDkiYadKa3O71br"

	hd := hashids.NewData()
	hd.Salt = "this is my salt"
	hd.MinLength = 30
	h, _ := hashids.NewWithData(hd)
	expected := h.Decode(hash)

	options := DefaultOptions("this is my salt")
	options.MinLength = 30

	o, _ := New(options)
	actual := o.Decode(hash).AsInt()

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\nOn input %s expected: %#v, got %#v", hash, expected, actual)
	}
}
