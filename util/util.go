package util

func IntervalSplits(begin, end, split int64) [][2]int64 {
	result := [][2]int64{}

	c := (end - begin) / split
	for i := begin; i < end; i++ {
		b := i + c - 1
		if b > end {
			b = end
		}
		result = append(result, [2]int64{i, b})
		i += c - 1
	}

	return result
}
