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

func Float32ToUint32(f float32) uint32 {
	return math.Float32bits(f)
}

func Uint32ToFloat32(i uint32) float32 {
	return math.Float32frombits(i)
}

func Float64ToUint64(f float64) uint64 {
	return math.Float64bits(f)
}

func Uint64ToFLoat64(i uint64) float64 {
	return math.Float64frombits(i)
}
