package beacon

import (
	"context"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Client is an interface which interacts with a Beacon node.
// NOTE: "make generate" MUST be run after changing this interface,
// in order to update any generated code which depends on it.
type Client interface {
	// Name returns the name of the client implementation.
	Name() string

	// Address returns the address of the client.
	Address() string

	Spec(ctx context.Context) (map[string]interface{}, error)
	Genesis(ctx context.Context) (*apiv1.Genesis, error)
	BeaconBlockRoot(ctx context.Context, blockID string) (*phase0.Root, error)
	SignedBeaconBlock(ctx context.Context, blockID string) (*spec.VersionedSignedBeaconBlock, error)
	BeaconBlockHeader(ctx context.Context, blockID string) (*apiv1.BeaconBlockHeader, error)
	ProposerDuties(ctx context.Context, epoch phase0.Epoch, indices []phase0.ValidatorIndex) ([]*apiv1.ProposerDuty, error)
	AttesterDuties(ctx context.Context, epoch phase0.Epoch, indices []phase0.ValidatorIndex) ([]*apiv1.AttesterDuty, error)
	SyncCommitteeDuties(ctx context.Context, epoch phase0.Epoch, indices []phase0.ValidatorIndex) ([]*apiv1.SyncCommitteeDuty, error)
	Domain(ctx context.Context, domainType phase0.DomainType, epoch phase0.Epoch) (phase0.Domain, error)
	Validators(ctx context.Context, stateID string, indices []phase0.ValidatorIndex) (map[phase0.ValidatorIndex]*apiv1.Validator, error)
	ValidatorsByPubKey(ctx context.Context, stateID string, validatorPubKeys []phase0.BLSPubKey) (map[phase0.ValidatorIndex]*apiv1.Validator, error)
	SubmitProposalPreparations(ctx context.Context, preparations []*apiv1.ProposalPreparation) error
	SubmitValidatorRegistrations(ctx context.Context, registrations []*api.VersionedSignedValidatorRegistration) error
	Proposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti [32]byte) (*api.VersionedProposal, error)
	SubmitBeaconBlock(ctx context.Context, block *spec.VersionedSignedBeaconBlock) error
	BlindedProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti [32]byte) (*api.VersionedBlindedProposal, error)
	SubmitBlindedBeaconBlock(ctx context.Context, block *api.VersionedSignedBlindedBeaconBlock) error
	SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*apiv1.BeaconCommitteeSubscription) error
	AttestationData(ctx context.Context, slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.AttestationData, error)
	SubmitAttestations(ctx context.Context, attestations []*phase0.Attestation) error
	AggregateAttestation(ctx context.Context, slot phase0.Slot, attestationDataRoot phase0.Root) (*phase0.Attestation, error)
	SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*phase0.SignedAggregateAndProof) error
	SubmitSyncCommitteeSubscriptions(ctx context.Context, subscriptions []*apiv1.SyncCommitteeSubscription) error
	SubmitSyncCommitteeMessages(ctx context.Context, messages []*altair.SyncCommitteeMessage) error
	SyncCommitteeContribution(ctx context.Context, slot phase0.Slot, subcommitteeIndex uint64, beaconBlockRoot phase0.Root) (*altair.SyncCommitteeContribution, error)
	SubmitSyncCommitteeContributions(ctx context.Context, contributionAndProofs []*altair.SignedContributionAndProof) error
	Events(ctx context.Context, topics []string, handler eth2client.EventHandlerFunc) error
}
