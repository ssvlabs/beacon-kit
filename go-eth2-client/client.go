package eth2client

import (
	"context"
	"errors"
	"strings"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/beacon-kit"
)

// ErrCallNotSupported is returned when the implementation does not support the requested call.
var ErrCallNotSupported = errors.New("call not supported")

// ErrEmptyResponse is returned when the client did not receive a response.
var ErrEmptyResponse = errors.New("empty response")

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

func (c *Client) Spec(ctx context.Context, opts *api.SpecOpts) (*api.Response[map[string]interface{}], error) {
	resp, err := c.service.(eth2client.SpecProvider).Spec(ctx, opts)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}

func (c *Client) Genesis(ctx context.Context, opts *api.GenesisOpts) (*api.Response[*apiv1.Genesis], error) {
	resp, err := c.service.(eth2client.GenesisProvider).Genesis(ctx, opts)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}

func (c *Client) BeaconBlockRoot(ctx context.Context, opts *api.BeaconBlockRootOpts) (*api.Response[*phase0.Root], error) {
	provider, ok := c.service.(eth2client.BeaconBlockRootProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	resp, err := provider.BeaconBlockRoot(ctx, opts)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}

func (c *Client) SignedBeaconBlock(ctx context.Context, opts *api.SignedBeaconBlockOpts) (*api.Response[*spec.VersionedSignedBeaconBlock], error) {
	provider, ok := c.service.(eth2client.SignedBeaconBlockProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	resp, err := provider.SignedBeaconBlock(ctx, opts)
	if err != nil {
		var apiErr *api.Error
		if errors.As(err, &apiErr) {
			if apiErr.StatusCode == 404 {
				return nil, beacon.ErrBlockNotFound
			}
		}

		// Hack to gracefully handle missing blocks from Prysm.
		notFound := false
		errString := err.Error()
		for _, s := range []string{
			"Could not get block from block ID: rpc error: code = NotFound",
			"rpc error: code = NotFound desc = Could not find requested block: signed beacon block can't be nil", // v2.1.0
		} {
			if strings.Contains(errString, s) {
				notFound = true
				break
			}
		}
		if !notFound {
			return nil, err
		}
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	if resp.Data == nil {
		return nil, beacon.ErrBlockNotFound
	}
	return resp, nil
}

func (c *Client) BeaconBlockHeader(ctx context.Context, opts *api.BeaconBlockHeaderOpts) (*api.Response[*apiv1.BeaconBlockHeader], error) {
	provider, ok := c.service.(eth2client.BeaconBlockHeadersProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	resp, err := provider.BeaconBlockHeader(ctx, opts)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}

func (c *Client) ProposerDuties(ctx context.Context, opts *api.ProposerDutiesOpts) (*api.Response[[]*apiv1.ProposerDuty], error) {
	provider, ok := c.service.(eth2client.ProposerDutiesProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	resp, err := provider.ProposerDuties(ctx, opts)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}

func (c *Client) AttesterDuties(ctx context.Context, opts *api.AttesterDutiesOpts) (*api.Response[[]*apiv1.AttesterDuty], error) {
	provider, ok := c.service.(eth2client.AttesterDutiesProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	resp, err := provider.AttesterDuties(ctx, opts)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}

func (c *Client) SyncCommitteeDuties(ctx context.Context, opts *api.SyncCommitteeDutiesOpts) (*api.Response[[]*apiv1.SyncCommitteeDuty], error) {
	provider, ok := c.service.(eth2client.SyncCommitteeDutiesProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	resp, err := provider.SyncCommitteeDuties(ctx, opts)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}

func (c *Client) Domain(ctx context.Context, domainType phase0.DomainType, epoch phase0.Epoch) (phase0.Domain, error) {
	provider, ok := c.service.(eth2client.DomainProvider)
	if !ok {
		return phase0.Domain{}, ErrCallNotSupported
	}
	return provider.Domain(ctx, domainType, epoch)
}

func (c *Client) GenesisDomain(ctx context.Context, domainType phase0.DomainType) (phase0.Domain, error) {
	provider, ok := c.service.(eth2client.DomainProvider)
	if !ok {
		return phase0.Domain{}, ErrCallNotSupported
	}
	return provider.GenesisDomain(ctx, domainType)
}

func (c *Client) Validators(ctx context.Context, opts *api.ValidatorsOpts) (*api.Response[map[phase0.ValidatorIndex]*apiv1.Validator], error) {
	provider, ok := c.service.(eth2client.ValidatorsProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	resp, err := provider.Validators(ctx, opts)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	return resp, nil
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

func (c *Client) Proposal(ctx context.Context, opts *api.ProposalOpts) (*api.Response[*api.VersionedProposal], error) {
	provider, ok := c.service.(eth2client.ProposalProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	resp, err := provider.Proposal(ctx, opts)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}

func (c *Client) SubmitProposal(ctx context.Context, proposal *api.VersionedSignedProposal) error {
	provider, ok := c.service.(eth2client.ProposalSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitProposal(ctx, proposal)
}

func (c *Client) BlindedProposal(ctx context.Context, opts *api.BlindedProposalOpts) (*api.Response[*api.VersionedBlindedProposal], error) {
	provider, ok := c.service.(eth2client.BlindedProposalProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	resp, err := provider.BlindedProposal(ctx, opts)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}

// SubmitBlindedBeaconBlock provides a mock function with given fields: ctx, block
func (c *Client) SubmitBlindedProposal(ctx context.Context, block *api.VersionedSignedBlindedProposal) error {
	provider, ok := c.service.(eth2client.BlindedProposalSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitBlindedProposal(ctx, block)
}

func (c *Client) SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*apiv1.BeaconCommitteeSubscription) error {
	provider, ok := c.service.(eth2client.BeaconCommitteeSubscriptionsSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitBeaconCommitteeSubscriptions(ctx, subscriptions)
}

func (c *Client) AttestationData(ctx context.Context, opts *api.AttestationDataOpts) (*api.Response[*phase0.AttestationData], error) {
	provider, ok := c.service.(eth2client.AttestationDataProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	resp, err := provider.AttestationData(ctx, opts)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	return resp, nil
}

func (c *Client) SubmitAttestations(ctx context.Context, attestations []*phase0.Attestation) error {
	provider, ok := c.service.(eth2client.AttestationsSubmitter)
	if !ok {
		return ErrCallNotSupported
	}
	return provider.SubmitAttestations(ctx, attestations)
}

func (c *Client) AggregateAttestation(ctx context.Context, opts *api.AggregateAttestationOpts) (*api.Response[*phase0.Attestation], error) {
	provider, ok := c.service.(eth2client.AggregateAttestationProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	resp, err := provider.AggregateAttestation(ctx, opts)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	return resp, nil
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

func (c *Client) SyncCommitteeContribution(ctx context.Context, opts *api.SyncCommitteeContributionOpts) (*api.Response[*altair.SyncCommitteeContribution], error) {
	provider, ok := c.service.(eth2client.SyncCommitteeContributionProvider)
	if !ok {
		return nil, ErrCallNotSupported
	}
	resp, err := provider.SyncCommitteeContribution(ctx, opts)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrEmptyResponse
	}
	return resp, nil
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
