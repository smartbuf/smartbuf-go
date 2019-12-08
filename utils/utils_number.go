package utils

import "math"

func IntToUint(l int64) uint64 {
	if l >= 0 {
		return uint64(l << 1)
	} else {
		return uint64((-l)<<1) | 1
	}
}

func UintToInt(l uint64) int64 {
	if l&1 == 0 {
		return int64(l >> 1)
	} else if l == 1 {
		return math.MinInt64
	} else {
		return -int64(l >> 1)
	}
}
