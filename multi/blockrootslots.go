package multi

import (
	"sync"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

type blockRootSlots struct {
	data map[phase0.Root]phase0.Slot
	mu   sync.RWMutex
}

func newBlockRootSlots() *blockRootSlots {
	return &blockRootSlots{
		data: map[phase0.Root]phase0.Slot{},
	}
}

func (r *blockRootSlots) Set(root phase0.Root, slot phase0.Slot) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[root] = slot
}

func (r *blockRootSlots) Get(root phase0.Root) (slot phase0.Slot, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	slot, ok = r.data[root]
	return
}

func (r *blockRootSlots) Purge(minSlot phase0.Slot) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for root, slot := range r.data {
		if slot < minSlot {
			delete(r.data, root)
			n++
		}
	}
	return n
}

func (r *blockRootSlots) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.data)
}
