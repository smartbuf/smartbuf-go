package utils

// Fast-Sort arithmetic which support descend sort for int[]
func DescFastSort(arr []int) {
	if len(arr) <= 1 {
		return
	}
	var low = 0
	var high = len(arr) - 1
	var baseVal = arr[low]
	for low < high {
		for ; high > low; high-- {
			if arr[high] > baseVal {
				arr[low] = arr[high]
				low++
				break
			}
		}
		for ; high > low; low++ {
			if arr[low] < baseVal {
				arr[high] = arr[low]
				high--
				break
			}
		}
	}

	arr[low] = baseVal
	if low > 1 {
		DescFastSort(arr[:low])
	}
	if (len(arr) - low) > 1 {
		DescFastSort(arr[low+1:])
	}
}
