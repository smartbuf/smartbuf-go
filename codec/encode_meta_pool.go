package codec

import (
	"github.com/go-eden/common/emath"
	"github.com/go-eden/etime"
)

const (
	HasNameTmp        = 1
	HasNameAdded      = 1 << 1
	HasNameExpired    = 1 << 2
	HasStructTmp      = 1 << 3
	HasStructAdded    = 1 << 4
	HasStructExpired  = 1 << 5
	HasStructReferred = 1 << 6
)

// EncodeMetaPool represents an area holds struct for sharing, which support temporary and context using.
type EncodeMetaPool struct {
	cxtStructLimit int
	flags          int

	tmpNames     []string
	tmpNameIndex map[string]int

	tmpStructs     []*structModel
	tmpStructIndex *structTree

	cxtIdAlloc     *IDAllocator
	cxtNames       []*nameModel
	cxtNameAdded   []*nameModel
	cxtNameExpired []int
	cxtNameIndex   map[string]*nameModel

	cxtStructCount    int
	cxtStructIdAlloc  *IDAllocator
	cxtStructs        []*structModel
	cxtStructAdded    []*structModel
	cxtStructExpired  []int
	cxtStructReferred []*structModel
	cxtStructIndex    *structTree
}

func newEncodeMetaPool(limit int) *EncodeMetaPool {
	return &EncodeMetaPool{
		cxtStructLimit: limit,
		tmpNames:       make([]string, 4),
		tmpNameIndex:   make(map[string]int),
	}
}

// Register an struct for temporary using, should return an unique ID
func (t *EncodeMetaPool) registerTmpStruct(names []string) int {
	if names == nil {
		panic("names is nil")
	}
	if len(names) == 0 {
		return 0
	}
	node := t.tmpStructIndex.getNode(names)
	model := node.body
	if model == nil {
		nameIds := make([]int, len(names))
		for i, name := range names {
			nameId, ok := t.tmpNameIndex[name]
			if !ok {
				t.tmpNames = append(t.tmpNames, name)
				nameId = len(t.tmpNames)
			}
			nameIds[i] = nameId
		}
		model = newStructModel(names, nameIds)
		model.index = len(t.tmpStructs)
		model.id = (model.index + 1) << 1
		t.tmpStructs = append(t.tmpStructs, model)

		node.body = model
	}
	return model.id
}

// Register the specified struct into pool by its field-names.
func (t *EncodeMetaPool) registerCxtStruct(names []string) int {
	if names == nil {
		panic("names is nil")
	}
	if len(names) == 0 {
		return 0
	}
	node := t.cxtStructIndex.getNode(names)
	model := node.body
	if model == nil {
		nameIds := make([]int, len(names))
		for _, str := range names {
			name := t.cxtNameIndex[str]
			if name == nil {
				index := t.cxtIdAlloc.Require()
				name = newNameModel(str, index)
				t.cxtNames = append(t.cxtNames, name)
				t.cxtNameAdded = append(t.cxtNameAdded, name) // record for outter using
				t.cxtNameIndex[str] = name
			}
			nameIds = append(nameIds, name.index)
			name.refCount++ // update ref counter
		}
		model = newStructModel(names, nameIds)
		model.index = t.cxtStructIdAlloc.Require()
		model.id = ((model.index + 1) << 1) | 1 // identify context struct by suffixed 1
		t.cxtStructs[model.index] = model
		t.cxtStructAdded = append(t.cxtStructAdded, model)
		t.cxtStructCount++

		node.body = model
	}
	model.lastTime = int(etime.NowSecond())
	if !model.referred {
		model.referred = true
		t.cxtStructReferred = append(t.cxtStructReferred, model)
	}
	return model.id
}

func (t *EncodeMetaPool) isNeedOutput() bool {
	var flags = 0
	if len(t.tmpNames) > 0 {
		flags |= HasNameTmp
	}
	if len(t.cxtNameAdded) > 0 {
		flags |= HasNameAdded
	}
	if len(t.cxtNameExpired) > 0 {
		flags |= HasNameExpired
	}
	if len(t.tmpStructs) > 0 {
		flags |= HasStructTmp
	}
	if len(t.cxtStructAdded) > 0 {
		flags |= HasStructAdded
	}
	if len(t.cxtStructExpired) > 0 {
		flags |= HasStructExpired
	}
	if len(t.cxtStructReferred) > 0 {
		flags |= HasStructReferred
	}
	t.flags = flags
	return flags > 0
}

func (t *EncodeMetaPool) isNeedSequence() bool {
	return t.flags&(HasNameAdded|HasNameExpired|HasStructAdded|HasStructExpired) > 0
}

func (t *EncodeMetaPool) write(buf *EncodeBuffer) {
	var flags = t.flags
	var size = 0
	if (flags & HasNameTmp) > 0 {
		flags ^= HasNameTmp
		size = len(t.tmpNames)
		buf.WriteVarUint(uint64(size<<4) | FlagMetaNameTmp | 1) // must have structTmp
		for i := 0; i < size; i++ {
			buf.WriteString(t.tmpNames[i])
		}
	}
	if (flags & HasNameExpired) > 0 {
		flags ^= HasNameExpired
		size = len(t.cxtNameExpired)
		buf.WriteVarUint(uint64(size<<4) | FlagMetaNameExpired | 1) // must have structExpired
		for i := 0; i < size; i++ {
			buf.WriteVarUint(uint64(t.cxtNameExpired[i]))
		}
	}
	if (flags & HasNameAdded) > 0 {
		flags ^= HasNameAdded
		size = len(t.cxtNameAdded)
		buf.WriteVarUint(uint64(size<<4) | FlagMetaNameAdded | 1) // must have structAdded
		for i := 0; i < size; i++ {
			buf.WriteString(t.cxtNameAdded[i].name)
		}
	}
	if (flags & HasStructTmp) > 0 {
		flags ^= HasStructTmp
		size = len(t.tmpStructs)
		if flags == 0 {
			buf.WriteVarUint(uint64(size<<4) | FlagMetaStructTmp)
		}
		for i := 0; i < size; i++ {
			item := t.tmpStructs[i]
			buf.WriteVarUint(uint64(len(item.nameIds)))
			for _, nameId := range item.nameIds {
				buf.WriteVarUint(uint64(nameId))
			}
		}
	}
	if (flags & HasStructExpired) > 0 {
		size = len(t.cxtStructAdded)
		flags ^= HasStructExpired
		if flags == 0 {
			buf.WriteVarUint(uint64(size<<4) | FlagMetaStructExpired)
		}
		for i := 0; i < size; i++ {
			buf.WriteVarUint(uint64(t.cxtStructExpired[i]))
		}
	}
	if (flags & HasStructAdded) > 0 {
		size = len(t.cxtStructAdded)
		flags ^= HasStructAdded
		buf.WriteVarUint(uint64(size<<4) | FlagMetaStructAdded | 1) // must has HAS_STRUCT_REFERRED suffixed
		for i := 0; i < size; i++ {
			nameIds := t.cxtStructAdded[i].nameIds
			buf.WriteVarUint(uint64(len(nameIds)))
			for _, nameId := range nameIds {
				buf.WriteVarUint(uint64(nameId))
			}
		}
	}
	if (flags & HasStructReferred) > 0 {
		size = len(t.cxtStructReferred)
		buf.WriteVarUint(uint64(size<<4) | FlagMetaStructReferred)
		for i := 0; i < size; i++ {
			item := t.cxtStructReferred[i]
			buf.WriteVarUint(uint64(item.index))
			buf.WriteVarUint(uint64(len(item.names)))
		}
	}
}

// reset this pool
func (t *EncodeMetaPool) reset() {
	t.tmpNames = t.tmpNames[:0]
	t.tmpNameIndex = make(map[string]int)
	t.cxtNameAdded = t.cxtNameAdded[:0]
	t.cxtNameExpired = t.cxtNameExpired[:0]

	t.tmpStructs = t.tmpStructs[:0]
	t.tmpStructIndex = newStructTree()
	t.cxtStructAdded = t.cxtStructAdded[:0]
	t.cxtStructExpired = t.cxtStructExpired[:0]
	t.cxtStructReferred = t.cxtStructReferred[:0]

	for _, s := range t.cxtStructs {
		s.referred = false
	}

	expireCounnt := t.cxtStructCount - t.cxtStructLimit
	if expireCounnt <= 0 {
		return
	}
	structs := make([]*structModel, len(t.cxtStructs))
	copy(structs, t.cxtStructs)

	// pick the earlier structs to release
	// TODO asc sort structs
	count := emath.MinInt(expireCounnt, len(structs))
	for i := 0; i < count; i++ {
		expiredStruct := structs[i]
		if node := t.cxtStructIndex.getNode(expiredStruct.names); node != nil {
			node.body = nil
		}
		t.cxtStructIdAlloc.Release(expiredStruct.index)
		t.cxtStructExpired = append(t.cxtStructExpired, expiredStruct.index)
		t.cxtStructs[expiredStruct.index] = nil
		// synchronize cxtNames
		for _, nameId := range expiredStruct.nameIds {
			meta := t.cxtNames[nameId]
			meta.refCount--
			if meta.refCount == 0 {
				t.cxtNames[meta.index] = nil
				t.cxtIdAlloc.Release(meta.index)
				t.cxtNameExpired = append(t.cxtNameExpired, meta.index)
				delete(t.cxtNameIndex, meta.name)
			}
		}
	}
}

///////////////////////////////////////////////////////
// field-nameModel's metadata
type nameModel struct {
	name  string
	index int

	refCount int
}

func newNameModel(str string, index int) *nameModel {
	return &nameModel{name: str, index: index}
}

///////////////////////////////////////////////////////
// Struct model for inner usage
type structModel struct {
	id       int
	index    int
	lastTime int
	names    []string
	nameIds  []int
	referred bool
}

func newStructModel(names []string, nameIds []int) *structModel {
	return &structModel{names: names, nameIds: nameIds}
}

///////////////////////////////////////////////////////
type structTree struct {
	root *structNode
}

func newStructTree() *structTree {
	return &structTree{root: newStructNode("")}
}

func (t *structTree) getNode(names []string) *structNode {
	var node = t.root
	for _, name := range names {
		node = node.findNode(name)
	}
	return node
}

///////////////////////////////////////////////////////
type structNode struct {
	name     string
	subNodes map[string]*structNode

	body *structModel
}

func newStructNode(name string) *structNode {
	return &structNode{name: name, subNodes: make(map[string]*structNode)}
}

func (t *structNode) findNode(name string) *structNode {
	var node = t.subNodes[name]
	if node == nil {
		node = newStructNode(name)
		t.subNodes[name] = node
	}
	return node
}
