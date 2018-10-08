package hashids

type ResultMapFunc func(v int64, i int) int64

type Decoded struct {
	result []int64
	err    error
}

func (d Decoded) HasError() bool {
	return d.err != nil
}

func (d Decoded) Count() int {
	return len(d.result)
}

func (d Decoded) Error() string {
	return d.err.Error()
}

// Unwrap the raw result and error
func (d Decoded) Unwrap() ([]int64, error) {
	return d.result, d.err
}

// AsInt slice
func (d Decoded) AsInt() []int {
	if d.result == nil {
		return nil
	}

	out := make([]int, 0)

	for _, v := range d.result {
		out = append(out, int(v))
	}

	return out
}

// AsInt64 slice
func (d Decoded) AsInt64() []int64 {
	return d.result
}

// Map over the results
func (d Decoded) Map(f ResultMapFunc) Decoded {
	result := make([]int64, len(d.result))

	for i, v := range d.result {
		result[i] = f(v, i)
	}

	return Decoded{result, nil}
}
