package multi

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/beacon-kit"
	"github.com/bloxapp/beacon-kit/logging"
	"github.com/bloxapp/beacon-kit/pool"
	"go.uber.org/zap"
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

	blockRootSlots *blockRootSlots
}

func New(spec *beacon.Spec, poolClient *pool.Client, options Options) *Client {
	return &Client{
		spec:           spec,
		Client:         poolClient,
		options:        options,
		blockRootSlots: newBlockRootSlots(),
	}
}

// BestAttestationDataSelection subscribes to block events to select
// the best (rather than the first) AttestationData.
func (c *Client) BestAttestationDataSelection(ctx context.Context) error {
	err := c.EventsWithClient(ctx, []string{"block"}, func(client beacon.Client, e *api.Event) {
		// log.Printf("GotBlockEventData: %#v", e.Data)
		if e.Data == nil {
			return
		}
		data := e.Data.(*api.BlockEvent)

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

func (c *Client) AttestationData(ctx context.Context, slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.AttestationData, error) {
	var (
		bestData       *phase0.AttestationData
		bestDataSlot   phase0.Slot
		bestDataClient string
		mu             sync.Mutex
	)
	err := c.With(pool.FirstSuccess(false)).
		Call(ctx, func(ctx context.Context, client beacon.Client) error {
			data, err := client.AttestationData(ctx, slot, committeeIndex)
			if err != nil {
				return err
			}
			if data == nil {
				return nil
			}

			dataSlot, _ := c.blockRootSlots.Get(data.BeaconBlockRoot)
			// log.Printf("FoundSlotForAttestation: %#x -> %d", data.BeaconBlockRoot, dataSlot)

			logging.FromContext(ctx).Debug("Scoring AttestationData",
				zap.String("client", client.Address()),
				zap.String("block_root", fmt.Sprintf("%#x", data.BeaconBlockRoot)),
				zap.Uint64("derived_slot", uint64(dataSlot)))

			func() {
				mu.Lock()
				defer mu.Unlock()
				if bestData == nil || dataSlot > bestDataSlot {
					if bestData != nil && dataSlot > bestDataSlot {
						// log.Printf("SelectedBetterAttestation: %d>%d — %#x (((VERSUS))) %#x", dataSlot, bestDataSlot, data.BeaconBlockRoot, bestData.BeaconBlockRoot)

						logging.FromContext(ctx).Debug("Better AttestationData detected",
							zap.String("client", bestDataClient),
							zap.String("block_root", fmt.Sprintf("%#x", bestData.BeaconBlockRoot)),
							zap.Uint64("slot", uint64(bestDataSlot)),
							zap.String("better_client", client.Address()),
							zap.String("better_block_root", fmt.Sprintf("%#x", data.BeaconBlockRoot)),
							zap.Uint64("better_slot", uint64(dataSlot)))
					}

					bestData = data
					bestDataSlot = dataSlot
					bestDataClient = client.Address()
				}
			}()
			return nil
		})
	return bestData, err
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
