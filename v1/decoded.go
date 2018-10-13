package hashids

type ResultMapFunc func(v int64, i int) int64

// DecodedNumbers hash result
type DecodedNumbers struct {
	numbers []int64
	err     error
}

// NewDecodedNumbers result
func NewDecodedNumbers(numbers []int64, err error) *DecodedNumbers {
	d := new(DecodedNumbers)
	d.err = err

	if numbers != nil {
		d.numbers = make([]int64, len(numbers))
		copy(d.numbers, numbers)
	}

	return d
}

// HasError - whether result contains error
func (d DecodedNumbers) HasError() bool {
	return d.err != nil
}

// Len of result
func (d DecodedNumbers) Len() int {
	return len(d.numbers)
}

func (d DecodedNumbers) Error() string {
	return d.err.Error()
}

// Unwrap the raw result and error
func (d DecodedNumbers) Unwrap() ([]int64, error) {
	return d.numbers, d.err
}

// IntSlice from result
func (d DecodedNumbers) IntSlice() []int {
	if d.numbers == nil {
		return nil
	}

	out := make([]int, 0)

	for _, v := range d.numbers {
		out = append(out, int(v))
	}

	return out
}

func (d DecodedNumbers) FirstInt() (first int) {
	if len(d.numbers) > 0 {
		first = int(d.numbers[0])
	}

	return
}

func (d DecodedNumbers) FirstInt64() (first int64) {
	if len(d.numbers) > 0 {
		first = d.numbers[0]
	}

	return
}

// Int64Slice slice
func (d DecodedNumbers) Int64Slice() []int64 {
	return d.numbers
}

// AsHex returns result converted to hexidecimal format
func (d DecodedNumbers) AsHex() (string, error) {
	if d.err != nil {
		return "", d.err
	}
	return numsToHex(d.numbers)
}

// Map over the results
func (d DecodedNumbers) Map(f ResultMapFunc) DecodedNumbers {
	result := make([]int64, len(d.numbers))

	for i, v := range d.numbers {
		result[i] = f(v, i)
	}

	return DecodedNumbers{result, nil}
}
