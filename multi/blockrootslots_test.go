package multi

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"sync"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestBlockRootSlots(t *testing.T) {
	roots := make([]phase0.Root, 10e3)
	for i := 0; i < len(roots); i++ {
		_, err := rand.Read(roots[i][:])
		require.NoError(t, err)
	}

	s := newBlockRootSlots()

	// Concurrent setters.
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				rnd, err := rand.Int(rand.Reader, big.NewInt(int64(len(roots))))
				require.NoError(t, err)
				root := roots[rnd.Int64()]
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
				rnd, err := rand.Int(rand.Reader, big.NewInt(int64(len(roots))))
				require.NoError(t, err)
				root := roots[rnd.Int64()]
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
