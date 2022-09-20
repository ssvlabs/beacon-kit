package multi

import (
	"context"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/beacon-kit"
	"github.com/bloxapp/beacon-kit/pool"
)

type CallTrace struct {
	pool.CallTrace
	SubnetIDs []beacon.SubnetID
}

type Options struct{}

// Client implements a protocol-aware beacon.Client on top of pool.Client
// with ideal behaviour for the different calls.
type Client struct {
	*pool.Client
	spec    *beacon.Spec
	options Options
}

func New(spec *beacon.Spec, poolClient *pool.Client, options Options) *Client {
	return &Client{
		spec:    spec,
		Client:  poolClient,
		options: options,
	}
}

func (c *Client) With(options ...interface{}) *Client {
	copy := *c
	copy.Client = c.Client.With(options...)
	return &copy
}

func (c *Client) SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*api.BeaconCommitteeSubscription) error {
	return c.submitter().SubmitBeaconCommitteeSubscriptions(ctx, subscriptions)
}

func (c *Client) SubmitAttestations(ctx context.Context, attestations []*phase0.Attestation) error {
	return c.submitter().SubmitAttestations(ctx, attestations)
}

func (c *Client) SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*phase0.SignedAggregateAndProof) error {
	return c.submitter().SubmitAggregateAttestations(ctx, aggregateAndProofs)
}

func (c *Client) SubmitSyncCommitteeSubscriptions(ctx context.Context, subscriptions []*api.SyncCommitteeSubscription) error {
	return c.submitter().SubmitSyncCommitteeSubscriptions(ctx, subscriptions)
}

func (c *Client) submitter() *pool.Client {
	return c.Client.With(
		pool.FirstSuccess(false),
	)
}
