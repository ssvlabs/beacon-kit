package eth2client

import (
	"context"
	"errors"
	"time"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// ErrCallNotSupported is returned when the implementation does not support the requested call.
var ErrCallNotSupported = errors.New("call not supported")

// Client implements beacon.Client using the go-eth2-client package.
type Client struct {
	service eth2client.Service
}

// New creates a new Client.
func New(service eth2client.Service) *Client {
	return &Client{
		service: service,
	}
}

// Service returns the underlying go-eth2-client.Service
func (c *Client) Service() eth2client.Service {
	return c.service
}

// Name returns the name of the client implementation.
func (c *Client) Name() string {
	return c.service.Name()
}

// Address returns the address of the client.
func (c *Client) Address() string {
	return c.service.Address()
}

func (c *Client) Spec(ctx context.Context) (map[string]interface{}, error) {
	return c.service.(eth2client.SpecProvider).Spec(ctx)
}

func (c *Client) GenesisTime(ctx context.Context) (time.Time, error) {
	return c.service.(eth2client.GenesisTimeProvider).GenesisTime(ctx)
}

func (c *Client) BeaconBlockRoot(ctx context.Context, blockID string) (*phase0.Root, error) {
	provider, ok := c.service.(eth2client.BeaconBlockRootProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	return provider.BeaconBlockRoot(ctx, blockID)
}

func (c *Client) SignedBeaconBlock(ctx context.Context, blockID string) (*spec.VersionedSignedBeaconBlock, error) {
	provider, ok := c.service.(eth2client.SignedBeaconBlockProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	return provider.SignedBeaconBlock(ctx, blockID)
}

func (c *Client) BeaconBlockHeader(ctx context.Context, blockID string) (*apiv1.BeaconBlockHeader, error) {
	provider, ok := c.service.(eth2client.BeaconBlockHeadersProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	return provider.BeaconBlockHeader(ctx, blockID)
}

func (c *Client) ProposerDuties(ctx context.Context, epoch phase0.Epoch, indices []phase0.ValidatorIndex) ([]*apiv1.ProposerDuty, error) {
	provider, ok := c.service.(eth2client.ProposerDutiesProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	return provider.ProposerDuties(ctx, epoch, indices)
}

func (c *Client) AttesterDuties(ctx context.Context, epoch phase0.Epoch, indices []phase0.ValidatorIndex) ([]*apiv1.AttesterDuty, error) {
	provider, ok := c.service.(eth2client.AttesterDutiesProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	return provider.AttesterDuties(ctx, epoch, indices)
}

func (c *Client) SyncCommitteeDuties(ctx context.Context, epoch phase0.Epoch, indices []phase0.ValidatorIndex) ([]*apiv1.SyncCommitteeDuty, error) {
	provider, ok := c.service.(eth2client.SyncCommitteeDutiesProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	return provider.SyncCommitteeDuties(ctx, epoch, indices)
}

func (c *Client) Domain(ctx context.Context, domainType phase0.DomainType, epoch phase0.Epoch) (phase0.Domain, error) {
	provider, ok := c.service.(eth2client.DomainProvider)
	if !ok {
		return phase0.Domain{}, ErrCallNotSupported
	}
	return provider.Domain(ctx, domainType, epoch)
}

func (c *Client) Validators(ctx context.Context, stateID string, indices []phase0.ValidatorIndex) (map[phase0.ValidatorIndex]*apiv1.Validator, error) {
	provider, ok := c.service.(eth2client.ValidatorsProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	return provider.Validators(ctx, stateID, indices)
}

func (c *Client) ValidatorsByPubKey(ctx context.Context, stateID string, validatorPubKeys []phase0.BLSPubKey) (map[phase0.ValidatorIndex]*apiv1.Validator, error) {
	provider, ok := c.service.(eth2client.ValidatorsProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	return provider.ValidatorsByPubKey(ctx, stateID, validatorPubKeys)
}

func (c *Client) SubmitProposalPreparations(ctx context.Context, preparations []*apiv1.ProposalPreparation) error {
	provider, ok := c.service.(eth2client.ProposalPreparationsSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitProposalPreparations(ctx, preparations)
}

func (c *Client) SubmitValidatorRegistrations(ctx context.Context, registrations []*api.VersionedSignedValidatorRegistration) error {
	provider, ok := c.service.(eth2client.ValidatorRegistrationsSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitValidatorRegistrations(ctx, registrations)
}

func (c *Client) BeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*spec.VersionedBeaconBlock, error) {
	provider, ok := c.service.(eth2client.BeaconBlockProposalProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	return provider.BeaconBlockProposal(ctx, slot, randaoReveal, graffiti)
}

func (c *Client) SubmitBeaconBlock(ctx context.Context, block *spec.VersionedSignedBeaconBlock) error {
	provider, ok := c.service.(eth2client.BeaconBlockSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitBeaconBlock(ctx, block)
}

func (c *Client) BlindedBeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*api.VersionedBlindedBeaconBlock, error) {
	provider, ok := c.service.(eth2client.BlindedBeaconBlockProposalProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	return provider.BlindedBeaconBlockProposal(ctx, slot, randaoReveal, graffiti)
}

// SubmitBlindedBeaconBlock provides a mock function with given fields: ctx, block
func (c *Client) SubmitBlindedBeaconBlock(ctx context.Context, block *api.VersionedSignedBlindedBeaconBlock) error {
	provider, ok := c.service.(eth2client.BlindedBeaconBlockSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitBlindedBeaconBlock(ctx, block)
}

func (c *Client) SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*apiv1.BeaconCommitteeSubscription) error {
	provider, ok := c.service.(eth2client.BeaconCommitteeSubscriptionsSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitBeaconCommitteeSubscriptions(ctx, subscriptions)
}

func (c *Client) AttestationData(ctx context.Context, slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.AttestationData, error) {
	provider, ok := c.service.(eth2client.AttestationDataProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	return provider.AttestationData(ctx, slot, committeeIndex)
}

func (c *Client) SubmitAttestations(ctx context.Context, attestations []*phase0.Attestation) error {
	provider, ok := c.service.(eth2client.AttestationsSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitAttestations(ctx, attestations)
}

func (c *Client) AggregateAttestation(ctx context.Context, slot phase0.Slot, attestationDataRoot phase0.Root) (*phase0.Attestation, error) {
	provider, ok := c.service.(eth2client.AggregateAttestationProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	return provider.AggregateAttestation(ctx, slot, attestationDataRoot)
}

func (c *Client) SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*phase0.SignedAggregateAndProof) error {
	provider, ok := c.service.(eth2client.AggregateAttestationsSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitAggregateAttestations(ctx, aggregateAndProofs)
}

func (c *Client) SubmitSyncCommitteeSubscriptions(ctx context.Context, subscriptions []*apiv1.SyncCommitteeSubscription) error {
	provider, ok := c.service.(eth2client.SyncCommitteeSubscriptionsSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitSyncCommitteeSubscriptions(ctx, subscriptions)
}

func (c *Client) SubmitSyncCommitteeMessages(ctx context.Context, messages []*altair.SyncCommitteeMessage) error {
	provider, ok := c.service.(eth2client.SyncCommitteeMessagesSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitSyncCommitteeMessages(ctx, messages)
}

func (c *Client) SyncCommitteeContribution(ctx context.Context, slot phase0.Slot, subcommitteeIndex uint64, beaconBlockRoot phase0.Root) (*altair.SyncCommitteeContribution, error) {
	provider, ok := c.service.(eth2client.SyncCommitteeContributionProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	return provider.SyncCommitteeContribution(ctx, slot, subcommitteeIndex, beaconBlockRoot)
}

func (c *Client) SubmitSyncCommitteeContributions(ctx context.Context, contributionAndProofs []*altair.SignedContributionAndProof) error {
	provider, ok := c.service.(eth2client.SyncCommitteeContributionsSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitSyncCommitteeContributions(ctx, contributionAndProofs)
}

func (c *Client) Events(ctx context.Context, topics []string, handler eth2client.EventHandlerFunc) error {
	provider, ok := c.service.(eth2client.EventsProvider)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.Events(ctx, topics, handler)
}
