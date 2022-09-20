package pool

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/beacon-kit"
	"github.com/bloxapp/beacon-kit/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	_ "net/http/pprof"
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	os.Exit(m.Run())
}

func TestPoolWithoutFirstSuccess(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	knobs := TestPoolKnobs{
		NumClients: 20,
		MinSleep:   0 * time.Millisecond,
		MaxSleep:   50 * time.Millisecond,
		ErrorRate:  0,
	}
	const retryLimit = 3
	options := []interface{}{
		SelectAll(),
		Concurrency(2),
		FirstSuccess(false),
		RetryEveryLimit(time.Millisecond, retryLimit),
	}
	pool := CreateTestPool(t, knobs, options...)
	require.Equal(t, pool.NumClients, pool.Size(), "pool size should equal the number of clients")

	// Test that a call is received by all clients.
	poolResp, err := pool.BeaconBlockHeader(ctx, "32")
	require.NoError(t, err)

	for _, client := range pool.Mocks {
		client.AssertCalled(t, "BeaconBlockHeader", mock.Anything, "32")
	}

	// Test that the pool returns the same data as the clients.
	clientResp, err := pool.Mocks[0].BeaconBlockHeader(ctx, "32")
	require.NoError(t, err)
	require.Equal(t, clientResp, poolResp, "pool should return the same data as the client")

	// Test that a call returns an error if all clients fail.
	knobs.ErrorRate = 1
	pool = CreateTestPool(t, knobs, options...)
	_, err = pool.BeaconBlockHeader(ctx, "32")
	errStruct, ok := err.(*Error)
	require.True(t, ok, "error should be a CallTrace")
	require.Len(
		t,
		errStruct.Trace,
		(retryLimit+1)*pool.Size(),
		"number of calls in error should be (RetryLimit + 1) * NumberOfClients",
	)
}

// TestPoolFaultTolerance tests that the pool can handle 50% of the clients failing.
func TestPoolFaultTolerance(t *testing.T) {
	ctx := context.Background()

	// Create a pool with some mock clients.
	knobs := TestPoolKnobs{
		NumClients: 20,
		MinSleep:   0 * time.Millisecond,
		MaxSleep:   50 * time.Millisecond,
		ErrorRate:  0.5,
	}
	pool := CreateTestPool(t, knobs, SelectAll(), Concurrency(20), FirstSuccess(true))
	require.Equal(t, pool.NumClients, pool.Size(), "pool size should be equal to number of clients")

	// Test that the pool returns the same data as the clients.
	poolBlock, err := pool.BeaconBlockHeader(ctx, "32")
	require.NoError(t, err)

	clientBlock, err := CreateTestClient(0, 0, 0).BeaconBlockHeader(ctx, "32")
	require.NoError(t, err)
	require.Equal(t, clientBlock, poolBlock, "pool should return the same data as the client")
}

type TestPoolKnobs struct {
	NumClients int
	MinSleep   time.Duration
	MaxSleep   time.Duration
	ErrorRate  float64
}

type TestPool struct {
	TestPoolKnobs

	*Client
	Mocks []*mocks.Client
}

func CreateTestPool(t *testing.T, knobs TestPoolKnobs, options ...interface{}) *TestPool {
	pool := &TestPool{
		TestPoolKnobs: knobs,
		Mocks:         make([]*mocks.Client, knobs.NumClients),
	}
	clients := make([]beacon.Client, knobs.NumClients)
	for i := 0; i < knobs.NumClients; i++ {
		errorRate := 0.0
		if knobs.ErrorRate*float64(knobs.NumClients) > float64(i) {
			errorRate = 1.0
		}
		client := CreateTestClient(knobs.MinSleep, knobs.MaxSleep, errorRate)
		pool.Mocks[i] = client
		clients[i] = client
	}

	pool.Client = New(clients, options...)
	require.Equal(
		t,
		knobs.NumClients,
		pool.Size(),
		"pool size should be equal to number of clients",
	)
	return pool
}

func CreateTestClient(minSleep, maxSleep time.Duration, errorRate float64) *mocks.Client {
	client := &mocks.Client{}
	client.On("Address").Maybe().Return("http://mock")
	client.On("BeaconBlockHeader", mock.Anything, mock.AnythingOfType("string")).
		Maybe().
		Return(func(ctx context.Context, blockRoot string) *api.BeaconBlockHeader {
			if maxSleep > 0 {
				time.Sleep(minSleep + time.Duration(rand.Intn(int(maxSleep-minSleep+1))))
			}
			return &api.BeaconBlockHeader{Root: phase0.Root{}}
		}, func(ctx context.Context, blockRoot string) error {
			if rand.Float64() <= errorRate {
				return fmt.Errorf("error")
			}
			return nil
		})
	return client
}
