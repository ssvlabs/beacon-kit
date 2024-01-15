package beacon

import (
	eth2client "github.com/attestantio/go-eth2-client"
)

// Client is an interface which interacts with a Beacon node.
// NOTE: "make generate" MUST be run after changing this interface,
// in order to update any generated code which depends on it.
type Client interface {
	// Name returns the name of the client implementation.
	Name() string

	// Address returns the address of the client.
	Address() string

	eth2client.SpecProvider
	eth2client.GenesisProvider
	eth2client.BeaconBlockRootProvider
	eth2client.SignedBeaconBlockProvider
	eth2client.BeaconBlockHeadersProvider
	eth2client.DomainProvider

	eth2client.ValidatorsProvider

	eth2client.ProposerDutiesProvider
	eth2client.AttesterDutiesProvider
	eth2client.SyncCommitteeDutiesProvider

	eth2client.ProposalPreparationsSubmitter
	eth2client.ValidatorRegistrationsSubmitter
	eth2client.ProposalProvider
	eth2client.ProposalSubmitter
	eth2client.BlindedProposalProvider
	eth2client.BlindedProposalSubmitter

	eth2client.BeaconCommitteeSubscriptionsSubmitter
	eth2client.AttestationDataProvider
	eth2client.AttestationsSubmitter
	eth2client.AggregateAttestationProvider
	eth2client.AggregateAttestationsSubmitter

	eth2client.SyncCommitteeSubscriptionsSubmitter
	eth2client.SyncCommitteeMessagesSubmitter
	eth2client.SyncCommitteeContributionProvider
	eth2client.SyncCommitteeContributionsSubmitter
}
