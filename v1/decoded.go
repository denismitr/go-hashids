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

// Map over the results
func (d Decoded) Map(f ResultMapFunc) Decoded {
	result := make([]int64, len(d.result))

	for i, v := range d.result {
		result[i] = f(v, i)
	}

	return Decoded{result, nil}
}
