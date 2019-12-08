package transport

const (
	Ver        = 0b0001_0000
	VerStream  = 0b0000_1000
	VerHasData = 0b0000_0100
	VerHasMeta = 0b0000_0010
	VerHasSeq  = 0b0000_0001
)

const (
	FlagMetaNameTmp        = 1 << 1
	FlagMetaNameAdded      = 2 << 1
	FlagMetaNameExpired    = 3 << 1
	FlagMetaStructTmp      = 4 << 1
	FlagMetaStructAdded    = 5 << 1
	FlagMetaStructExpired  = 6 << 1
	FlagMetaStructReferred = 7 << 1
)

const (
	FlagDataFloat         = 1 << 1
	FlagDataDouble        = 2 << 1
	FlagDataVarint        = 3 << 1
	FlagDataString        = 4 << 1
	FlagDataSymbolAdded   = 5 << 1
	FlagDataSymbolExpired = 6 << 1
)

const (
	ConstNull      = 0x00
	ConstFalse     = 0x01
	ConstTrue      = 0x02
	ConstZeroArray = 0x03
)

const (
	TypeConst  = -1
	TypeVarint = 0
	TypeFloat  = 1
	TypeDouble = 2
	TypeString = 3
	TypeSymbol = 4
	TypeObject = 5
	TypeArray  = 6
	TypeNArray = 7

	TypeNArrayBool   = 1<<3 | TypeNArray
	TypeNArrayByte   = 2<<3 | TypeNArray
	TypeNArrayShort  = 3<<3 | TypeNArray
	TypeNArrayInt    = 4<<3 | TypeNArray
	TypeNArrayLong   = 5<<3 | TypeNArray
	TypeNArrayFloat  = 6<<3 | TypeNArray
	TypeNArrayDouble = 7<<3 | TypeNArray
)

const (
	TypeSliceNull    = 0x00
	TypeSliceBool    = 0x01
	TypeSliceFloat   = 0x02
	TypeSliceDouble  = 0x03
	TypeSliceByte    = 0x04
	TypeSliceShort   = 0x05
	TypeSliceInt     = 0x06
	TypeSliceLong    = 0x07
	TypeSliceString  = 0x08
	TypeSliceSymbol  = 0x09
	TypeSliceObject  = 0x0A
	TypeSliceUnknown = 0x0B
)
