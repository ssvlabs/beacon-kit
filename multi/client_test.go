package multi

import (
	"context"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/beacon-kit"
	"github.com/bloxapp/beacon-kit/mocks"
	"github.com/bloxapp/beacon-kit/pool"
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
		c.(*mocks.Client).On("SubmitAttestations", mock.Anything, &api.SubmitAttestationsOpts{}).Maybe().Return(nil)
	}
	err := client.SubmitAttestations(context.Background(), &api.SubmitAttestationsOpts{})
	require.Equal(t, nil, err)

	for _, c := range mockClients {
		c.(*mocks.Client).AssertCalled(t, "SubmitAttestations", mock.Anything, &api.SubmitAttestationsOpts{})
	}
}

func TestBestAttestationDataSelection(t *testing.T) {
	const (
		timeout      = 1 * time.Second
		earlyTimeout = 100 * time.Millisecond
	)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create an offline client (which will never return).
	offlineClient := mocks.NewClient(t)
	offlineClient.On("Address", mock.Anything).Return("offline")
	offlineClient.On("AttestationData", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		<-args.Get(0).(context.Context).Done()
	}).Return(nil, context.Canceled)

	// Create an online client (which will return immediately).
	onlineClient := mocks.NewClient(t)
	onlineClient.On("Address", mock.Anything).Return("online")
	onlineClient.On("AttestationData", mock.Anything, mock.Anything, mock.Anything).Return(&api.Response[*phase0.AttestationData]{Data: &phase0.AttestationData{}}, nil)

	// Create a multi.Client with BestAttestationDataSelection.
	client := New(
		beacon.Mainnet,
		pool.New([]beacon.Client{offlineClient, onlineClient}, nil, pool.FirstSuccess(true), pool.SelectAll(), pool.Concurrency(2)),
		Options{},
	)
	err := client.BestAttestationDataSelection(ctx, earlyTimeout)
	require.NoError(t, err)

	// Check that earlyTimeout is respected.
	start := time.Now()
	dataResp, err := client.AttestationData(ctx, &api.AttestationDataOpts{
		Slot:           0,
		CommitteeIndex: 0,
	})
	require.NoError(t, err)
	require.NotNil(t, dataResp)
	require.NotNil(t, dataResp.Data)
	took := time.Since(start)
	require.True(t, took >= earlyTimeout, "exited too early!")
	require.True(t, took < earlyTimeout+(earlyTimeout/2), "exited too late!")
}
