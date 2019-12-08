package codec

import "github.com/smartbuf/smartbuf-go/utils"

type IDAllocator struct {
	nextId   int
	reuseIds []int
}

// Acquire an unique and incremental ID
func (t *IDAllocator) Require() (id int) {
	if l := len(t.reuseIds); l == 0 {
		id = t.nextId
		t.nextId++
	} else {
		id = t.reuseIds[l-1]
		t.reuseIds = t.reuseIds[:l-1]
	}
	return
}

// Release the specified id, It will be used in high priority.
func (t *IDAllocator) Release(id int) {
	if id >= t.nextId {
		log.Error("id[%v] is invalid to release", id)
		return
	}
	t.reuseIds = append(t.reuseIds, id)
	utils.DescFastSort(t.reuseIds)
}
