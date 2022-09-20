package pool

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelectRandoms(t *testing.T) {
	tests := []struct {
		randoms  int
		size     int
		expected int
	}{
		{0, 0, 0},
		{2, 3, 2},
		{3, 3, 3},
		{4, 3, 3},
	}
	for i, test := range tests {
		for methodIndex, method := range []func(int) SelectFunc{SelectRandoms, SelectAdjacentRandoms} {
			selected := method(test.randoms)(test.size)
			actual := 0
			for j := 0; j < test.size; j++ {
				if selected(j, nil) {
					actual++
				}
			}
			require.Equal(t, test.expected, actual, "test %d: %d(%d)", i, methodIndex, test.randoms)
		}
	}
}
