package multi

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

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

	blockRootSlots   map[phase0.Root]phase0.Slot
	blockRootSlotsMu *sync.Mutex
}

func New(spec *beacon.Spec, poolClient *pool.Client, options Options) *Client {
	return &Client{
		spec:             spec,
		Client:           poolClient,
		options:          options,
		blockRootSlots:   map[phase0.Root]phase0.Slot{},
		blockRootSlotsMu: &sync.Mutex{},
	}
}

// BestAttestationDataSelection subscribes to block events to select
// the best (rather than the first) AttestationData.
func (c *Client) BestAttestationDataSelection(ctx context.Context) error {
	err := c.EventsWithClient(ctx, []string{"block"}, func(_ beacon.Client, e *api.Event) {
		if e.Data == nil {
			return
		}
		data := e.Data.(*api.BlockEvent)

		c.blockRootSlotsMu.Lock()
		defer c.blockRootSlotsMu.Unlock()
		c.blockRootSlots[data.Block] = data.Slot
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
			time.Sleep(time.Second * 30)
			func() {
				c.blockRootSlotsMu.Lock()
				defer c.blockRootSlotsMu.Lock()
				minSlot := c.spec.Clock().Now().Slot() - maxSlotAge
				for root, slot := range c.blockRootSlots {
					if slot < minSlot {
						delete(c.blockRootSlots, root)
					}
				}
			}()
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
		best        *phase0.AttestationData
		highestSlot phase0.Slot
		mu          sync.Mutex
	)
	err := c.Call(ctx, func(ctx context.Context, client beacon.Client) error {
		data, err := client.AttestationData(ctx, slot, committeeIndex)
		if err != nil {
			return err
		}
		if data == nil {
			return nil
		}
		dataSlot := c.blockRootSlots[data.BeaconBlockRoot]
		log.Printf("FoundSlotForAttestation: %#x -> %d", data.BeaconBlockRoot, dataSlot)
		if best == nil || dataSlot > highestSlot {
			mu.Lock()

			if best != nil && dataSlot > highestSlot {
				b1, err := json.Marshal(data)
				if err != nil {
					return err
				}
				b2, err := json.Marshal(best)
				if err != nil {
					return err
				}
				log.Printf("SelectedBetterAttestation: %s (((VERSUS))) %s", string(b1), string(b2))
			}

			best = data
			highestSlot = dataSlot
			mu.Unlock()
		}
		return nil
	})
	return best, err
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
