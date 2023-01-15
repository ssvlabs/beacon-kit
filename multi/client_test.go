package multi

import (
	"context"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/beacon-kit"
	"github.com/bloxapp/beacon-kit/mocks"
	"github.com/bloxapp/beacon-kit/pool"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWith(t *testing.T) {
	// Create a client with FirstSuccess.
	with := New(
		beacon.Mainnet,
		pool.New(nil, nil, pool.FirstSuccess(true)),
		Options{},
	)

	// Fork it without FirstSuccess.
	without := with.With(pool.FirstSuccess(false))

	// Check that the scope is not shared!
	require.True(t, bool(with.Scope().FirstSuccess))
	require.False(t, bool(without.Scope().FirstSuccess))
}

func TestSubmitAttestations(t *testing.T) {
	mockClients := make([]beacon.Client, 32)
	for i := 0; i < 32; i++ {
		mockClients[i] = mocks.NewClient(t)
	}

	// Create a client with FirstSuccess.
	client := New(
		beacon.Mainnet,
		pool.New(mockClients, nil, pool.FirstSuccess(true), pool.SelectAll(), pool.Concurrency(1)),
		Options{},
	)

	// Call SubmitAttestation and expect it to override
	// the FirstSuccess and call both clients.
	for _, c := range mockClients {
		c.(*mocks.Client).On("SubmitAttestations", mock.Anything, []*phase0.Attestation{}).Maybe().Return(nil)
	}
	err := client.SubmitAttestations(context.Background(), []*phase0.Attestation{})
	require.Equal(t, nil, err)

	for _, c := range mockClients {
		c.(*mocks.Client).AssertCalled(t, "SubmitAttestations", mock.Anything, []*phase0.Attestation{})
	}
}
