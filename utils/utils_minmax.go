package utils

func MinInt(nums ...int) int {
	var ret = nums[0]
	for _, val := range nums {
		if ret > val {
			ret = val
		}
	}
	return ret
}

func MaxInt(nums ...int) int {
	var ret = nums[0]
	for _, val := range nums {
		if ret < val {
			ret = val
		}
	}
	return ret
}
