package pool

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ssvlabs/beacon-kit"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func TestPoolUpdate(t *testing.T) {
	pool := NewPool(zap.NewNop(), func(ctx context.Context, address string) (beacon.Client, error) {
		time.Sleep(time.Millisecond * 100)
		return &Client{}, nil
	})

	tests := []struct {
		addresses   []string
		expectedLen int
	}{
		{
			addresses:   []string{"1", "2"},
			expectedLen: 0,
		},
		{
			addresses:   []string{"1", "2", "3"},
			expectedLen: 2,
		},
		{
			addresses:   []string{"2", "4"},
			expectedLen: 1,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var g errgroup.Group
			g.Go(func() error {
				return pool.Update(context.Background(), test.addresses)
			})
			time.Sleep(time.Millisecond * 15)
			require.Len(t, pool.Clients(), test.expectedLen, "clients updated before connection")
			err := g.Wait()
			require.NoError(t, err)
			require.Len(t, pool.Clients(), len(test.addresses), "clients not updated after connection")
		})
	}
}
