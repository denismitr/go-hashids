package hashids

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tt = []struct {
	name    string
	numbers []int64
	err     error
}{
	{"normal result", []int64{1, 2, 3}, nil},
	{"normal long result", []int64{1, 205, 3, 5, 3434543, 314423, 1234}, nil},
	{"error result", nil, errors.New("an error")},
	{"single result", []int64{10}, nil},
	{"max single result", []int64{math.MaxInt64}, nil},
}

func Test_ItCanUnwrapResult(t *testing.T) {
	t.Parallel()

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := NewDecodedResult(tc.numbers, tc.err)

			numbers, err := result.Unwrap()

			assert.Equal(t, numbers, tc.numbers)
			assert.Equal(t, err, tc.err)
		})
	}
}

func Test_ItCanReportAnError(t *testing.T) {
	t.Parallel()

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := NewDecodedResult(tc.numbers, tc.err)

			hasError := result.HasError()

			if tc.err != nil {
				assert.True(t, hasError)
			} else {
				assert.False(t, hasError)
			}
		})
	}
}

func Test_ItCanReturnErrorMsg(t *testing.T) {
	t.Parallel()

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := NewDecodedResult(tc.numbers, tc.err)

			err := result.Err()

			if tc.err != nil {
				assert.Equal(t, tc.err, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_ItCanGetLengthOfResult(t *testing.T) {
	t.Parallel()

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := NewDecodedResult(tc.numbers, tc.err)

			length := result.Len()

			if tc.err == nil {
				assert.Equal(t, len(tc.numbers), length)
			} else {
				assert.Equal(t, 0, length)
			}
		})
	}
}

func Test_ItCanGetResultAsInt64Slice(t *testing.T) {
	t.Parallel()

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := NewDecodedResult(tc.numbers, tc.err)

			numbers := result.Int64Slice()
			assert.IsType(t, *new([]int64), numbers)

			if tc.err != nil {
				assert.Nil(t, numbers)
			}
		})
	}
}

func Test_ItCanGetResultAsIntSlice(t *testing.T) {
	t.Parallel()

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := NewDecodedResult(tc.numbers, tc.err)

			numbers := result.IntSlice()
			assert.IsType(t, *new([]int), numbers)

			if tc.err != nil {
				assert.Nil(t, numbers)
			}
		})
	}
}

func Test_ItCanGetResultAsFirstInt64(t *testing.T) {
	t.Parallel()

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := NewDecodedResult(tc.numbers, tc.err)

			number := result.FirstInt64()
			assert.IsType(t, *new(int64), number)

			if tc.err != nil {
				assert.Equal(t, int64(0), number)
			} else {
				assert.Equal(t, tc.numbers[0], number)
			}
		})
	}
}

func Test_ItCanGetResultAsFirstInt(t *testing.T) {
	t.Parallel()

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := NewDecodedResult(tc.numbers, tc.err)

			number := result.FirstInt()
			assert.IsType(t, *new(int), number)

			if tc.err != nil {
				assert.Equal(t, 0, number)
			} else {
				assert.Equal(t, int(tc.numbers[0]), number)
			}
		})
	}
}

func Test_ItCanMapThroughTheResult(t *testing.T) {
	t.Parallel()

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := NewDecodedResult(tc.numbers, tc.err)

			newResult := result.Map(func(v int64, i int) int64 {
				if tc.err == nil {
					assert.Equal(t, v, tc.numbers[i])
				}

				return v * 2
			})

			numbers, err := newResult.Unwrap()
			if tc.err != nil {
				assert.Equal(t, tc.err, err)
			}

			for i, n := range tc.numbers {
				assert.True(t, numbers[i] == n*2)
			}
		})
	}
}
