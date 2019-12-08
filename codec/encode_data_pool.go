package codec

import "github.com/go-eden/common/emath"

const (
	HasFloat         = 1
	HasDouble        = 1 << 1
	HasVarint        = 1 << 2
	HasString        = 1 << 3
	HasSymbolAdded   = 1 << 4
	HasSymbolExpired = 1 << 5
)

type EncodeDataPool struct {
	symbolLimit int
	flags       int

	floats      []float32
	floatIndex  map[float32]int
	doubles     []float64
	doubleIndex map[float64]int
	varints     []int64
	varintIndex map[int64]int
	strings     []string
	stringIndex map[string]int

	symbolId      IDAllocator
	symbols       []*symbol
	symbolAdded   []*symbol
	symbolExpired []int
	symbolIndex   map[string]symbol
}

func newEncodeDataPool(limit int) *EncodeDataPool {
	return &EncodeDataPool{
		symbolLimit: limit,
	}
}

func (t *EncodeDataPool) registerFloat(f float32) int {
	if f == 0 {
		return 1
	}
	// TODO
	return 0
}

func (t *EncodeDataPool) registerDouble(f float64) int {
	if f == 0 {
		return 1
	}
	// TODO
	return 0
}

func (t *EncodeDataPool) registerVarint(v int64) int {
	if v == 0 {
		return 1
	}
	// TODO
	return 0
}

func (t *EncodeDataPool) registerString(s string) int {
	if s == "" {
		return 1
	}
	// TODO
	return 0
}

func (t *EncodeDataPool) registerSymbol(s string) int {
	if s == "" {
		panic("invalid symbol")
	}
	// TODO
	return 1
}

func (t *EncodeDataPool) isNeedOutput() bool {
	var flags int
	if len(t.floats) > 0 {
		flags |= HasFloat
	}
	if len(t.doubles) > 0 {
		flags |= HasDouble
	}
	if len(t.varints) > 0 {
		flags |= HasVarint
	}
	if len(t.strings) > 0 {
		flags |= HasString
	}
	if len(t.symbolAdded) > 0 {
		flags |= HasSymbolAdded
	}
	if len(t.symbolExpired) > 0 {
		flags |= HasSymbolExpired
	}

	t.flags = flags
	return flags != 0
}

func (t *EncodeDataPool) isNeedSequence() bool {
	return t.flags&(HasSymbolAdded|HasSymbolExpired) != 0
}

func (t *EncodeDataPool) write(buf *EncodeBuffer) {
	var flags = t.flags
	var size = 0
	if (flags & HasFloat) != 0 {
		size = len(t.floats)
		flags ^= HasFloat
		if flags == 0 {
			buf.WriteVarUint(uint64(size<<4) | FlagDataFloat)
		} else {
			buf.WriteVarUint(uint64(size<<4) | FlagDataFloat | 1)
		}
		for i := 0; i < size; i++ {
			buf.WriteFloat32(t.floats[i])
		}
	}
	if (flags & HasDouble) != 0 {
		size = len(t.doubles)
		flags ^= HasDouble
		if flags == 0 {
			buf.WriteVarUint(uint64(size<<4) | FlagDataDouble)
		} else {
			buf.WriteVarUint(uint64(size<<4) | FlagDataDouble | 1)
		}
		for i := 0; i < size; i++ {
			buf.WriteFloat64(t.doubles[i])
		}
	}
	if (flags & HasVarint) != 0 {
		size = len(t.varints)
		flags ^= HasVarint
		if flags == 0 {
			buf.WriteVarUint(uint64(size<<4) | FlagDataVarint)
		} else {
			buf.WriteVarUint(uint64(size<<4) | FlagDataVarint | 1)
		}
		for i := 0; i < size; i++ {
			buf.WriteVarInt(t.varints[i])
		}
	}
	if (flags & HasString) != 0 {
		size = len(t.strings)
		flags ^= HasString
		if flags == 0 {
			buf.WriteVarUint(uint64(size<<4) | FlagDataString)
		} else {
			buf.WriteVarUint(uint64(size<<4) | FlagDataString | 1)
		}
		for i := 0; i < size; i++ {
			buf.WriteString(t.strings[i])
		}
	}
	if (flags & HasSymbolExpired) != 0 {
		size = len(t.symbolExpired)
		flags ^= HasSymbolExpired
		if flags == 0 {
			buf.WriteVarUint(uint64(size<<4) | FlagDataSymbolExpired)
		} else {
			buf.WriteVarUint(uint64(size<<4) | FlagDataSymbolExpired | 1)
		}
		for i := 0; i < size; i++ {
			buf.WriteVarUint(uint64(t.symbolExpired[i]))
		}
	}
	if (flags & HasSymbolAdded) != 0 {
		size = len(t.symbolAdded)
		buf.WriteVarUint(uint64(size<<4) | FlagDataSymbolAdded)
		for i := 0; i < size; i++ {
			buf.WriteString(t.symbolAdded[i].value)
		}
	}
}

func (t *EncodeDataPool) reset() {
	t.floats = t.floats[:0]
	t.floatIndex = make(map[float32]int)
	t.doubles = t.doubles[:0]
	t.doubleIndex = make(map[float64]int)
	t.varints = t.varints[:0]
	t.varintIndex = make(map[int64]int)
	t.strings = t.strings[:0]
	t.stringIndex = make(map[string]int)
	t.symbolAdded = t.symbolAdded[:0]
	t.symbolExpired = t.symbolExpired[:0]

	// check and expire symbols if thay are too many
	var expireNum = len(t.symbolIndex) - t.symbolLimit
	if expireNum <= 0 {
		return
	}
	var symbols = make([]*symbol, 4)
	// TODO symbols.sort((a, b) => a.lastTime - b.lastTime);
	l := emath.MinInt(expireNum, len(symbols))
	for i := 0; i < l; i++ {
		expiredSymbol := symbols[i]
		delete(t.symbolIndex, expiredSymbol.value)
		t.symbolId.Release(expiredSymbol.index)
		t.symbols[expiredSymbol.index] = nil
		t.symbolExpired = append(t.symbolExpired, expiredSymbol.index)
	}
}

//////////////////////////////////////////////////////////
type symbol struct {
	value string
	index int

	lastTime uint32
}

func newSymbol(value string, index int) *symbol {
	return &symbol{value: value, index: index}
}
