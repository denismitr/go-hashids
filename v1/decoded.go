package hashids

import (
	"fmt"
	"time"
)

// ResultMapFunc to go through all numbers in result
type ResultMapFunc func(v int64, i int) int64

// DecodedResult of the hash input
type DecodedResult struct {
	numbers []int64
	err     error
}

// NewDecodedResult result
func NewDecodedResult(numbers []int64, err error) *DecodedResult {
	d := new(DecodedResult)
	d.err = err

	if numbers != nil {
		d.numbers = make([]int64, len(numbers))
		copy(d.numbers, numbers)
	}

	return d
}

// HasError - whether result contains error
func (d DecodedResult) HasError() bool {
	return d.err != nil
}

// Len of result
func (d DecodedResult) Len() int {
	return len(d.numbers)
}

// Err - get the error of the result
func (d DecodedResult) Err() error {
	return d.err
}

// Unwrap the raw result and error
func (d DecodedResult) Unwrap() ([]int64, error) {
	return d.numbers, d.err
}

// IntSlice from result
func (d DecodedResult) IntSlice() ([]int, error) {
	if d.err != nil {
		return nil, d.err
	}

	out := make([]int, 0)

	for _, v := range d.numbers {
		out = append(out, int(v))
	}

	return out, nil
}

// FirstInt of the result
func (d DecodedResult) FirstInt() (int, error) {
	if d.err != nil {
		return 0, d.err
	}

	if len(d.numbers) > 0 {
		return int(d.numbers[0]), nil
	}

	return 0, fmt.Errorf("empty result")
}

// FirstInt64 of the result
func (d DecodedResult) FirstInt64() (int64, error) {
	if d.err != nil {
		return 0, d.err
	}

	if len(d.numbers) > 0 {
		return d.numbers[0], nil
	}

	return 0, fmt.Errorf("empty result")
}

// Int64Slice slice
func (d DecodedResult) Int64Slice() ([]int64, error) {
	return d.Unwrap()
}

// AsTime transform result into time object and return it
func (d DecodedResult) AsTime() (time.Time, error) {
	if d.err != nil {
		return time.Unix(0, 0), d.err
	}

	if len(d.numbers) != 1 {
		return time.Unix(0, 0), fmt.Errorf("valid timestamp must be contained in a int64 slice as single value, got %v", d.numbers)
	}

	t := time.Unix(0, d.numbers[0])

	return t, nil
}

// AsHex returns result converted to hexidecimal format
func (d DecodedResult) AsHex() (string, error) {
	if d.err != nil {
		return "", d.err
	}
	return numsToHex(d.numbers)
}

// Map over the results
func (d DecodedResult) Map(f ResultMapFunc) DecodedResult {
	result := make([]int64, len(d.numbers))

	for i, v := range d.numbers {
		result[i] = f(v, i)
	}

	return DecodedResult{result, d.err}
}
