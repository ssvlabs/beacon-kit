package multi

import (
	"encoding/binary"
	"math/rand"
	"sync"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestBlockRootSlots(t *testing.T) {
	roots := make([]phase0.Root, 10e3)
	for i := 0; i < len(roots); i++ {
		rand.Read(roots[i][:])
	}

	s := newBlockRootSlots()

	// Concurrent setters.
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				root := roots[rand.Intn(len(roots))]
				slot := phase0.Slot(binary.LittleEndian.Uint64(root[:8]))
				s.Set(root, slot)
			}
		}()
	}

	// Concurrent getters.
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				root := roots[rand.Intn(len(roots))]
				expectedSlot := phase0.Slot(binary.LittleEndian.Uint64(root[:8]))
				slot, ok := s.Get(root)
				if ok {
					require.Equal(t, expectedSlot, slot)
				}
			}
		}()
	}

	// Verify something was written.
	wg.Wait()
	require.Greater(t, s.Len(), 0)

	// Set all roots.
	for _, root := range roots {
		slot := phase0.Slot(binary.LittleEndian.Uint64(root[:8]))
		s.Set(root, slot)
	}

	// Verify length.
	require.Equal(t, s.Len(), len(roots))

	// Verify all roots.
	for _, root := range roots {
		expectedSlot := phase0.Slot(binary.LittleEndian.Uint64(root[:8]))
		slot, ok := s.Get(root)
		require.True(t, ok)
		require.Equal(t, expectedSlot, slot)
	}
}
