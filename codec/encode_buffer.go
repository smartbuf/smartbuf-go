package codec

import (
	"github.com/smartbuf/smartbuf-go/utils"
)

type EncodeBuffer struct {
	limit  int
	offset int
	data   []byte
}

func NewEncodeBuffer(limit int) *EncodeBuffer {
	return &EncodeBuffer{limit: limit}
}

func (t *EncodeBuffer) Reset() {
	t.offset = 0
}

func (t *EncodeBuffer) WriteByte(v byte) {
	if len(t.data) == t.offset {
		t.ensureCapacity(t.offset + 1)
	}
	t.data[t.offset] = v
	t.offset++
}

func (t *EncodeBuffer) WriteVarInt(n int64) {
	t.WriteVarUint(utils.IntToUint(n))
}

func (t *EncodeBuffer) WriteVarUint(n uint64) int {
	if len(t.data) < t.offset+10 {
		t.ensureCapacity(t.offset + 10)
	}
	var oldOffset = t.offset
	for n != 0 || oldOffset == t.offset {
		if (n & 0xFFFFFFFFFFFFFF80) == 0 {
			t.data[t.offset] = byte(n)
			t.offset++
		} else {
			t.data[t.offset] = byte((n | 0x80) & 0xFF)
			t.offset++
		}
		n >>= 7
	}
	return t.offset - oldOffset
}

func (t *EncodeBuffer) WriteFloat32(f float32) {
	if len(t.data) < t.offset+4 {
		t.ensureCapacity(t.offset + 4)
	}
	var bits = utils.Float32ToUint32(f)
	for i := 0; i < 4; i++ {
		t.data[t.offset] = byte(bits & 0xFF)
		t.offset++
		bits >>= 8
	}
}

func (t *EncodeBuffer) WriteFloat64(f float64) {
	if len(t.data) < t.offset+8 {
		t.ensureCapacity(t.offset + 8)
	}
	var bits = utils.Float64ToUint64(f)
	for i := 0; i < 8; i++ {
		t.data[t.offset] = byte(bits & 0xFF)
		t.offset++
		bits >>= 8
	}
}

func (t *EncodeBuffer) WriteString(str string) {
	if minLen := t.offset + len(str)*4 + 4; len(t.data) < minLen {
		t.ensureCapacity(minLen)
	}
	bytes := utils.EncodeUTF8(str)
	t.WriteVarUint(uint64(len(bytes)))
	t.WriteByteArray(bytes)
}

func (t *EncodeBuffer) WriteBoolArray(arr []bool) {
	var l = len(arr)
	if len(t.data) < t.offset+(l+1)/8 {
		t.ensureCapacity(t.offset + (l+1)/8)
	}
	for i := 0; i < l; i += 8 {
		var b = byte(0)
		for j := 0; j < 8; j++ {
			off := i + j
			if off >= l {
				break
			}
			if arr[off] {
				b |= 1 << j
			}
		}
		t.data[t.offset] = b
		t.offset++
	}
}

func (t *EncodeBuffer) WriteByteArray(arr []byte) {
	var l = len(arr)
	if len(t.data) < t.offset+l {
		t.ensureCapacity(t.offset + l)
	}
	copy(t.data, arr)
	t.offset += l
}

func (t *EncodeBuffer) WriteShortArray(arr []int16) {
	if len(t.data) < t.offset+len(arr)*2 {
		t.ensureCapacity(t.offset + len(arr)*2)
	}
	for _, s := range arr {
		t.data[t.offset] = byte(s >> 8)
		t.data[t.offset+1] = byte(s & 0xFF)
		t.offset += 2
	}
}

func (t *EncodeBuffer) WriteInt32Array(arr []int32) {
	for _, v := range arr {
		t.WriteVarInt(int64(v))
	}
}

func (t *EncodeBuffer) WriteUint32Array(arr []uint32) {
	for _, v := range arr {
		t.WriteVarInt(int64(v))
	}
}

func (t *EncodeBuffer) WriteInt64Array(arr []int64) {
	for _, v := range arr {
		t.WriteVarInt(v)
	}
}

func (t *EncodeBuffer) WriteUint64Array(arr []uint64) {
	for _, v := range arr {
		t.WriteVarUint(v) // TODO as string?
	}
}

func (t *EncodeBuffer) WriteFloat32Array(arr []float32) {
	for _, v := range arr {
		t.WriteFloat32(v)
	}
}
func (t *EncodeBuffer) WriteFloat64Array(arr []float64) {
	for _, v := range arr {
		t.WriteFloat64(v)
	}
}

func (t *EncodeBuffer) ensureCapacity(size int) {
	var newSize = size
	if n := len(t.data) * 2; size < n {
		newSize = n
	}
	if newSize > t.limit {
		newSize = t.limit
	}
	if newSize < size {
		panic("no space")
	}
	newData := make([]byte, newSize, newSize)
	copy(newData, t.data)
	t.data = newData
}
