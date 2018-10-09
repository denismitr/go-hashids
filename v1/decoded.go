package hashids

type ResultMapFunc func(v int64, i int) int64

// Decoded hash result
type Decoded struct {
	result []int64
	err    error
}

// HasError - whether result contains error
func (d Decoded) HasError() bool {
	return d.err != nil
}

// Len of result
func (d Decoded) Len() int {
	return len(d.result)
}

func (d Decoded) Error() string {
	return d.err.Error()
}

// Unwrap the raw result and error
func (d Decoded) Unwrap() ([]int64, error) {
	return d.result, d.err
}

// AsIntSlice slice
func (d Decoded) AsIntSlice() []int {
	if d.result == nil {
		return nil
	}

	out := make([]int, 0)

	for _, v := range d.result {
		out = append(out, int(v))
	}

	return out
}

// AsInt64Slice slice
func (d Decoded) AsInt64Slice() []int64 {
	return d.result
}

// AsHex returns result converted to hexidecimal format
func (d Decoded) AsHex() (string, error) {
	return numsToHex(d.result)
}

// Map over the results
func (d Decoded) Map(f ResultMapFunc) Decoded {
	result := make([]int64, len(d.result))

	for i, v := range d.result {
		result[i] = f(v, i)
	}

	return Decoded{result, nil}
}
