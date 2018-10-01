package hashids

func createNumbersHashFor(slice []int64) int64 {
	nh := int64(0)
	for i, n := range slice {
		nh += (n % int64(i+100))
	}
	return nh
}
