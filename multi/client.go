package multi

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"go.uber.org/zap"

	"github.com/bloxapp/beacon-kit"
	"github.com/bloxapp/beacon-kit/logging"
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

	blockRootSlots                  *blockRootSlots
	bestAttestationSelectionTimeout time.Duration
}

func New(spec *beacon.Spec, poolClient *pool.Client, options Options) *Client {
	return &Client{
		spec:           spec,
		Client:         poolClient,
		options:        options,
		blockRootSlots: newBlockRootSlots(),
	}
}

// BestAttestationDataSelection subscribes to block events to
// select the best (rather than the first) AttestationData.
//
// earlyTimeout is fired once the first AttestationData is received,
// cancelling the context for any ongoing calls, and preventing
// us from waiting for slow/stuck clients.
func (c *Client) BestAttestationDataSelection(ctx context.Context, earlyTimeout time.Duration) error {
	c.bestAttestationSelectionTimeout = earlyTimeout

	err := c.EventsWithClient(ctx, []string{"block"}, func(client beacon.Client, e *apiv1.Event) {
		// log.Printf("GotBlockEventData: %#v", e.Data)
		if e.Data == nil {
			return
		}
		data := e.Data.(*apiv1.BlockEvent)

		// log.Printf("GotBlockEvent root %#x for slot %d from %s", data.Block, data.Slot, client.Address())
		c.blockRootSlots.Set(data.Block, data.Slot)
	})
	if err != nil {
		return err
	}

	// Periodically remove old entries from blockRootSlots.
	go func() {
		// Remove block roots for slots that are more than 75 epochs old. (8 hours)
		maxSlotAge := c.spec.SlotsPerEpoch * 75

		// Every 30 seconds.
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(30 * time.Second):
				minSlot := c.spec.Clock().Now().Slot() - maxSlotAge
				deleted := c.blockRootSlots.Purge(minSlot)
				log.Printf("Purging blockRootSlots: %d slots deleted", deleted)
			}
		}
	}()

	return nil
}

func (c *Client) With(options ...interface{}) *Client {
	copy := *c
	copy.Client = c.Client.With(options...)
	return &copy
}

func (c *Client) AttestationData(ctx context.Context, opts *api.AttestationDataOpts) (*api.Response[*phase0.AttestationData], error) {
	var (
		bestData       *phase0.AttestationData
		bestDataSlot   phase0.Slot
		bestDataClient string
		mu             sync.Mutex
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	parentCtx := ctx

	err := c.With(pool.FirstSuccess(false)).
		Call(ctx, func(ctx context.Context, client beacon.Client) error {
			resp, err := client.AttestationData(ctx, opts)
			if err != nil {
				return err
			}
			if resp == nil {
				return nil
			}

			dataSlot, _ := c.blockRootSlots.Get(resp.Data.BeaconBlockRoot)

			logging.FromContext(ctx).Debug("Scoring AttestationData",
				zap.String("client", client.Address()),
				zap.String("block_root", fmt.Sprintf("%#x", resp.Data.BeaconBlockRoot)),
				zap.Uint64("derived_slot", uint64(dataSlot)))

			func() {
				mu.Lock()
				defer mu.Unlock()

				// Now that we've got the first AttestationData, we're less eager to
				// wait for the other clients to respond.
				//
				// Start a timer to cancel other calls early.
				if bestData == nil {
					go func() {
						select {
						case <-parentCtx.Done():
						case <-time.After(c.bestAttestationSelectionTimeout):
							cancel()
						}
					}()
				}

				if bestData == nil || dataSlot > bestDataSlot {
					// Log that we found a better AttestationData.
					if bestData != nil && dataSlot > bestDataSlot {
						logging.FromContext(ctx).Debug("Better AttestationData detected",
							zap.String("client", bestDataClient),
							zap.String("block_root", fmt.Sprintf("%#x", bestData.BeaconBlockRoot)),
							zap.Uint64("slot", uint64(bestDataSlot)),
							zap.String("better_client", client.Address()),
							zap.String("better_block_root", fmt.Sprintf("%#x", resp.Data.BeaconBlockRoot)),
							zap.Uint64("better_slot", uint64(dataSlot)))
					}

					bestData = resp.Data
					bestDataSlot = dataSlot
					bestDataClient = client.Address()

					// TODO: research if this is a good idea.
					// if bestDataSlot == slot {
					// 	// Cancel all other calls, we found the best possible AttestationData.
					// 	cancel()
					// }
				}
			}()
			return nil
		})

	// If at least one of the calls succeeded, return the best AttestationData, ignoring any errors.
	if bestData != nil {
		return &api.Response[*phase0.AttestationData]{Data: bestData}, nil
	}

	return nil, err
}

func (c *Client) SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*apiv1.BeaconCommitteeSubscription) error {
	return c.submitter().SubmitBeaconCommitteeSubscriptions(ctx, subscriptions)
}

func (c *Client) SubmitAttestations(ctx context.Context, attestations []*phase0.Attestation) error {
	return c.submitter().SubmitAttestations(ctx, attestations)
}

func (c *Client) SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*phase0.SignedAggregateAndProof) error {
	return c.submitter().SubmitAggregateAttestations(ctx, aggregateAndProofs)
}

func (c *Client) SubmitSyncCommitteeSubscriptions(ctx context.Context, subscriptions []*apiv1.SyncCommitteeSubscription) error {
	return c.submitter().SubmitSyncCommitteeSubscriptions(ctx, subscriptions)
}

func (c *Client) submitter() *pool.Client {
	return c.Client.With(
		pool.FirstSuccess(false),
	)
}
