package codec

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

}

func (t *EncodeBuffer) WriteVarInt(n int64) {
	//t.WriteVarUint(n)
}

func (t *EncodeBuffer) WriteVarUint(n uint64) {
	if len(t.data) < t.offset+10 {
		//this.ensureCapacity(this.offset + 10);
	}
	//var oldOffset = t.offset

	//do {
	//	if ((n & 0xFFFFFFFFFFFFFF80) == 0) {
	//	this.data[this.offset++] = n;
	//} else {
	//	this.data[this.offset++] = ((n | 0x80) & 0xFF);
	//}
	//	n >>>= 7;
	//} while (n != 0);
	//return this.offset - oldOffset;
}

func (t *EncodeBuffer) ensureCapacity(size int) {
	//var  newSize = math.min(math.max(len(t.data) * 2, size), t.limit)
	//if (newSize < size) {
	//	throw new Error("no space");
	//}
	//let newData = new Uint8Array(newSize);
	//newData.set(t.data);
	//t.data = newData;
}

/**
public writeVarInt(n: number): void {
        this.writeVarUint(NumberUtils.intToUint(n));
    }

    public writeVarUint(n: number): number {
        if (this.data.length < this.offset + 10) {
            this.ensureCapacity(this.offset + 10);
        }
        let oldOffset = this.offset;
        do {
            if ((n & 0xFFFFFFFFFFFFFF80) == 0) {
                this.data[this.offset++] = n;
            } else {
                this.data[this.offset++] = ((n | 0x80) & 0xFF);
            }
            n >>>= 7;
        } while (n != 0);
        return this.offset - oldOffset;
    }

    public writeFloat(f: number): void {
        if (this.data.length < this.offset + 4) {
            this.ensureCapacity(this.offset + 4);
        }
        let bits = NumberUtils.floatToBits(f);
        for (let i = 0; i < 4; i++) {
            this.data[this.offset++] = (bits & 0xFF);
            bits >>>= 8;
        }
    }

    public writeDouble(d: number): void {
        if (this.data.length < this.offset + 8) {
            this.ensureCapacity(this.offset + 8);
        }
        let bits = NumberUtils.doubleToBits(d);
        for (let i = 0; i < 8; i++) {
            this.data[this.offset++] = (bits & 0xFF);
            bits >>>= 8;
        }
    }

    public writeString(str: string): void {
        let bytes = UTF8Utils.encodeUTF8(str);
        this.writeVarUint(bytes.length);
        this.writeByteArray(bytes);
    }

    public writeBooleanArray(arr: boolean[]): void {
        let len = arr.length;
        if (this.data.length < this.offset + (len + 1) / 8) {
            this.ensureCapacity(this.offset + (len + 1) / 8);
        }
        let off = 0;
        for (let i = 0; i < len; i += 8) {
            let b = 0;
            for (let j = 0; j < 8; j++) {
                if ((off = i + j) >= len) {
                    break;
                }
                if (arr[off]) {
                    b |= 1 << j;
                }
            }
            this.data[this.offset++] = b;
        }
    }

    public writeByteArray(arr: Uint8Array): void {
        let len = arr.length;
        if (this.data.length < this.offset + len) {
            this.ensureCapacity(this.offset + len);
        }
        this.data.set(arr, this.offset);
        this.offset += len;
    }

    public writeShortArray(arr: Uint16Array): void {
        if (this.data.length < this.offset + arr.length * 2) {
            this.ensureCapacity(this.offset + arr.length * 2);
        }
        for (let i = 0; i < arr.length; i++) {
            let s = arr[i];
            this.data[this.offset++] = (s >> 8);
            this.data[this.offset++] = (s & 0xFF);
        }
    }

    public writeIntArray(arr: Uint32Array): void {
        for (let i = 0; i < arr.length; i++) {
            this.writeVarInt(arr[i]);
        }
    }

    public writeFloatArray(arr: Float32Array): void {
        for (let i = 0; i < arr.length; i++) {
            this.writeFloat(arr[i]);
        }
    }

    public writeDoubleArray(arr: Float64Array): void {
        for (let i = 0; i < arr.length; i++) {
            this.writeDouble(arr[i]);
        }
    }

    private ensureCapacity(size: number): void {
        let newSize = Math.min(Math.max(this.data.length * 2, size), this.limit);
        if (newSize < size) {
            throw new Error("no space");
        }
        let newData = new Uint8Array(newSize);
        newData.set(this.data);
        this.data = newData;
    }

*/
