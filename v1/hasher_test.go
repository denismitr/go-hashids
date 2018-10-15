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
		{[]int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5}, []int64{5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5}},
		{[]int64{1, 2, 3, 4, 5}, []int64{1, 2, 3, 4, 5}},
		{[]int64{29, 30, 26, 29, 27, 30, 30, 31}, []int64{29, 30, 26, 29, 27, 30, 30, 31}},
		{2147483647, []int64{2147483647}},
		{4294967295, []int64{4294967295}},
		{9223372036854775807, []int64{9223372036854775807}},
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
				t.Fatal(result.Err())
			}

			if result.Len() == 0 {
				t.Fatalf(
					"Expected result list to have %d items, instead there is %d",
					len(tc.out),
					result.Len())
			}

			actual := result.Int64Slice()

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
		{1, "0NV0", "this is my salt", 4},
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
		{
			[]int{1000000001, 1000000002, 1000000003, 1000000004, 1000000005},
			"jOAab4REXGWr37ZKwoYMK758XDaxjkYv9Zb4VnWJ0ak5EnoY6JR8WDXPlpOZr7dr2bdevMmX9PjBow6EKVkR5EBr8QlA3bapo12ZX7GML5JMX0VDrGxRoYWjLbk25ePnAnYRdl5KAPLbZ6Qj9rxX48VOJxZwa5oX0jekKv2bQqYALbnpG2190mlExMZqj4RwoBe9YdAElDROnP0w2xJ1pGLQBMeRYVWGp4KXDQlvOb0mEBP69kdoDm7ZrEjKM2ABJYRqV1MwnRO76bJoYBG4eQ5P0EZrlKQbV08Zvxp97o3GPOAXYOxEnAMmXDKqRPvlwL4WakbVOlLPQBqo8DZKG6nkx7e3RkL3R9m4wJbQDOWd1a0jY65GDKEqMxp54wdYLrn7mVlR9e5BMOobw4Zv6dYQaPq78VRG2paRklKJ5rqPXb3vZLVY0PLKEqnd6VxA8ZbYJB7Oa95OPWdaYvn5R479q3V8AorxlN7VoaAn8TpLraRecE0n43nfYERBkoCBJ3L80zeG1wbB0LEQ6mjMJ2KpZXkDGlmj41okRDXwe0MWrvQp3287WdwnOM9EmQjo61DxBA4epJx2LmXAG9rE103KjDlknWejak0o2QW31bP8XOAvZB6JMxElA7GVoerXZnPpB8qv2K5vd2a0A1mjwp4rW9YXMEbJZJ5e31V9Y0B87odpG26rQjeda1kj2JE4wqD6WnBRm5LMVrLadx9mv8A2KkX3WjplDq3w5epl0vQaLbWG41POxX8nknjo1w87A392dra5ZJqx6LvMZqj46kom87XKra5WbV3eYX6Qa3Lrd7JPWVAvO8Kk5DVnW496mMpDPr18RlE7d3BGOa12Ep3wekBDMqm70voWGJ4w69BQZvlKOqE8dap3Jm17ked09jKVn6WOPqmDRwYx4vxZnOp0L73YDAJQl8q1GWa4mAq0Vwe1bx2GKjBMQv493LQ2wApd3lmr16GOoeLBPRqE5kpQdLV2DBxvnml86M9",
			"my test salt",
			999,
		},
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

func Test_DecodeReturnsCorrectInput(t *testing.T) {
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

		t.Run(fmt.Sprintf("Hash %s", tc.hash), func(t *testing.T) {
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

func Test_InvalidInput(t *testing.T) {
	tt := []struct {
		input interface{}
		t     string
	}{
		{false, fmt.Sprintf("%T", false)},
		{nil, fmt.Sprintf("%T", nil)},
		{new(struct{}), fmt.Sprintf("%T", new(struct{}))},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(fmt.Sprintf("Input type %T", tc.input), func(t *testing.T) {
			h, _ := New(DefaultOptions("some salt"))

			_, err := h.Encode(tc.input)
			if err == nil {
				t.Fatal("err should not be nil")
			}

			assert.Contains(t, err.Error(), tc.t)
		})
	}
}
