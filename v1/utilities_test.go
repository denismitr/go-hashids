package hashids

import (
	"reflect"
	"strconv"
	"testing"
)

func Test_HashFunc(t *testing.T) {
	tt := []struct {
		input    int64
		alphabet []rune
		result   []rune
	}{
		{1, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'b'}},
		{1, []rune{'1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'2'}},
		{10, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'6'}},
		{10, []rune{'1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'2', '3'}},
		{124, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'5', '3'}},
		{124, []rune{'1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'2', '9', '5'}},
		{9, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'5'}},
		{98, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'3', '3'}},
		{4567, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'c', 'b', 'a', 'e'}},
		{4567, []rune{'1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'2', '1', '9', '3', '9'}},
		{998, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'1', '7', '6'}},
		{998, []rune{'1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'2', '9', '5', '7'}},
		{11, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'7'}},
		{11, []rune{'1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'2', '4'}},
	}

	for _, tc := range tt {
		t.Run(strconv.Itoa(int(tc.input)), func(t *testing.T) {
			actual := hash(tc.input, tc.alphabet)

			if !reflect.DeepEqual(tc.result, actual) {
				t.Fatalf("For input %d expected %s, got %s", tc.input, string(tc.result), string(actual))
			}
		})
	}
}

func Test_UnhashFucn(t *testing.T) {
	tt := []struct {
		result   int64
		alphabet []rune
		input    []rune
	}{
		{1, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'b'}},
		{1, []rune{'1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'2'}},
		{10, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'6'}},
		{10, []rune{'1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'2', '3'}},
		{124, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'5', '3'}},
		{124, []rune{'1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'2', '9', '5'}},
		{9, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'5'}},
		{98, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'3', '3'}},
		{4567, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'c', 'b', 'a', 'e'}},
		{4567, []rune{'1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'2', '1', '9', '3', '9'}},
		{998, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'1', '7', '6'}},
		{998, []rune{'1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'2', '9', '5', '7'}},
		{11, []rune{'a', 'b', 'c', 'd', 'e', '1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'7'}},
		{11, []rune{'1', '2', '3', '4', '5', '6', '7', '9'}, []rune{'2', '4'}},
	}

	for _, tc := range tt {
		t.Run(string(tc.input), func(t *testing.T) {
			actual, err := unhash(tc.input, tc.alphabet)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.result, actual) {
				t.Fatalf("For input %s expected %d, got %d", string(tc.input), tc.result, actual)
			}
		})
	}
}

func Test_SplitHashFunc(t *testing.T) {
	tt := []struct {
		seps     string
		input    []rune
		expected [][]rune
	}{
		{"cfhistuCFHISTU", []rune{'1', 'a', 'b'}, [][]rune{{'1', 'a', 'b'}}},
		{"cfhistuCFHISTU", []rune{'1', 'a', 'b', 'c', 'e', 'd', '1'}, [][]rune{{'1', 'a', 'b'}, {'e', 'd', '1'}}},
		{"cfhistuCFHISTU", []rune{'c', 'f', 'h', 'i', 's', 't', 'u'}, [][]rune{{}, {}, {}, {}, {}, {}, {}, {}}},
		{"cfhistuCFHISTU", []rune{'1', '2', '3', 'c'}, [][]rune{{'1', '2', '3'}, {}}},
		{"cfhistuCFHISTU", []rune{'c', '1', '2', '3', 'F'}, [][]rune{{}, {'1', '2', '3'}, {}}},
		{"cfhistuCFHISTU", []rune{'x', '1', '2', '3', 'F'}, [][]rune{{'x', '1', '2', '3'}, {}}},
		{"cfhistuCFHISTU", []rune{'y', 'b', 'v', 'x', '1', '2', '3', 'F'}, [][]rune{{'y', 'b', 'v', 'x', '1', '2', '3'}, {}}},
	}

	for _, tc := range tt {
		t.Run(string(tc.input), func(t *testing.T) {
			actual := separate(tc.input, []rune(tc.seps))

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fatalf("On input %s expected %v, got %v", string(tc.input), tc.expected, actual)
			}
		})
	}
}

func Test_ShuffleFunc(t *testing.T) {
	tt := []struct {
		in   []rune
		out  []rune
		salt []rune
	}{
		{
			in:   []rune(defaultAlphabet),
			salt: []rune("Test salt"),
			out:  []rune("sI0gJwor67dkiH4EPKvfMjRAh8uFBSZzQcG5Op3DxTNYa1Lqy92XmbtWlCneVU"),
		},
		{
			in:   []rune(defaultAlphabet),
			salt: []rune("Another salt"),
			out:  []rune("I9GNpMxeBgy5rYlzovXE31Z7nHcCVfOAtwRmaDWbuQhjkd02J84sFLKqPS6TUi"),
		},
		{
			in:   []rune("customAlphabet123"),
			salt: []rune("Test salt"),
			out:  []rune("lem3aAhtc1tuo2bsp"),
		},
		{
			in:   []rune("customAlphabet123"),
			salt: []rune("Another salt"),
			out:  []rune("mopeAatb1cth32lus"),
		},
	}

	for _, tc := range tt {
		t.Run(string(tc.in), func(t *testing.T) {
			expected := make([]rune, len(tc.out))
			copy(expected, tc.out)

			actual := shuffle(tc.in, tc.salt)

			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("On input %s expected to see %s got %s", string(tc.in), string(expected), string(actual))
			}
		})
	}
}
