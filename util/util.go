package util

func IntervalSplitsByNumber(begin, end, split int64) [][2]int64 {
	result := [][2]int64{}

	for i := begin; i <= end; i += split {
		b := i + split - 1
		if b > end {
			b = end
		}
		result = append(result, [2]int64{i, b})
	}

	return result
}

func IntervalSplitBySize(begin, end, split int64) [][2]int64 {
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
