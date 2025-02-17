// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	api "github.com/attestantio/go-eth2-client/api"
	altair "github.com/attestantio/go-eth2-client/spec/altair"

	context "context"

	mock "github.com/stretchr/testify/mock"

	phase0 "github.com/attestantio/go-eth2-client/spec/phase0"

	spec "github.com/attestantio/go-eth2-client/spec"

	v1 "github.com/attestantio/go-eth2-client/api/v1"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Address provides a mock function with no fields
func (_m *Client) Address() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Address")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// AggregateAttestation provides a mock function with given fields: ctx, opts
func (_m *Client) AggregateAttestation(ctx context.Context, opts *api.AggregateAttestationOpts) (*api.Response[*spec.VersionedAttestation], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for AggregateAttestation")
	}

	var r0 *api.Response[*spec.VersionedAttestation]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.AggregateAttestationOpts) (*api.Response[*spec.VersionedAttestation], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.AggregateAttestationOpts) *api.Response[*spec.VersionedAttestation]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[*spec.VersionedAttestation])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.AggregateAttestationOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AttestationData provides a mock function with given fields: ctx, opts
func (_m *Client) AttestationData(ctx context.Context, opts *api.AttestationDataOpts) (*api.Response[*phase0.AttestationData], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for AttestationData")
	}

	var r0 *api.Response[*phase0.AttestationData]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.AttestationDataOpts) (*api.Response[*phase0.AttestationData], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.AttestationDataOpts) *api.Response[*phase0.AttestationData]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[*phase0.AttestationData])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.AttestationDataOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AttesterDuties provides a mock function with given fields: ctx, opts
func (_m *Client) AttesterDuties(ctx context.Context, opts *api.AttesterDutiesOpts) (*api.Response[[]*v1.AttesterDuty], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for AttesterDuties")
	}

	var r0 *api.Response[[]*v1.AttesterDuty]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.AttesterDutiesOpts) (*api.Response[[]*v1.AttesterDuty], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.AttesterDutiesOpts) *api.Response[[]*v1.AttesterDuty]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[[]*v1.AttesterDuty])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.AttesterDutiesOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BeaconBlockHeader provides a mock function with given fields: ctx, opts
func (_m *Client) BeaconBlockHeader(ctx context.Context, opts *api.BeaconBlockHeaderOpts) (*api.Response[*v1.BeaconBlockHeader], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for BeaconBlockHeader")
	}

	var r0 *api.Response[*v1.BeaconBlockHeader]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.BeaconBlockHeaderOpts) (*api.Response[*v1.BeaconBlockHeader], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.BeaconBlockHeaderOpts) *api.Response[*v1.BeaconBlockHeader]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[*v1.BeaconBlockHeader])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.BeaconBlockHeaderOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BeaconBlockRoot provides a mock function with given fields: ctx, opts
func (_m *Client) BeaconBlockRoot(ctx context.Context, opts *api.BeaconBlockRootOpts) (*api.Response[*phase0.Root], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for BeaconBlockRoot")
	}

	var r0 *api.Response[*phase0.Root]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.BeaconBlockRootOpts) (*api.Response[*phase0.Root], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.BeaconBlockRootOpts) *api.Response[*phase0.Root]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[*phase0.Root])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.BeaconBlockRootOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BeaconCommittees provides a mock function with given fields: ctx, opts
func (_m *Client) BeaconCommittees(ctx context.Context, opts *api.BeaconCommitteesOpts) (*api.Response[[]*v1.BeaconCommittee], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for BeaconCommittees")
	}

	var r0 *api.Response[[]*v1.BeaconCommittee]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.BeaconCommitteesOpts) (*api.Response[[]*v1.BeaconCommittee], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.BeaconCommitteesOpts) *api.Response[[]*v1.BeaconCommittee]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[[]*v1.BeaconCommittee])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.BeaconCommitteesOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Domain provides a mock function with given fields: ctx, domainType, epoch
func (_m *Client) Domain(ctx context.Context, domainType phase0.DomainType, epoch phase0.Epoch) (phase0.Domain, error) {
	ret := _m.Called(ctx, domainType, epoch)

	if len(ret) == 0 {
		panic("no return value specified for Domain")
	}

	var r0 phase0.Domain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, phase0.DomainType, phase0.Epoch) (phase0.Domain, error)); ok {
		return rf(ctx, domainType, epoch)
	}
	if rf, ok := ret.Get(0).(func(context.Context, phase0.DomainType, phase0.Epoch) phase0.Domain); ok {
		r0 = rf(ctx, domainType, epoch)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(phase0.Domain)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, phase0.DomainType, phase0.Epoch) error); ok {
		r1 = rf(ctx, domainType, epoch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Genesis provides a mock function with given fields: ctx, opts
func (_m *Client) Genesis(ctx context.Context, opts *api.GenesisOpts) (*api.Response[*v1.Genesis], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for Genesis")
	}

	var r0 *api.Response[*v1.Genesis]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.GenesisOpts) (*api.Response[*v1.Genesis], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.GenesisOpts) *api.Response[*v1.Genesis]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[*v1.Genesis])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.GenesisOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenesisDomain provides a mock function with given fields: ctx, domainType
func (_m *Client) GenesisDomain(ctx context.Context, domainType phase0.DomainType) (phase0.Domain, error) {
	ret := _m.Called(ctx, domainType)

	if len(ret) == 0 {
		panic("no return value specified for GenesisDomain")
	}

	var r0 phase0.Domain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, phase0.DomainType) (phase0.Domain, error)); ok {
		return rf(ctx, domainType)
	}
	if rf, ok := ret.Get(0).(func(context.Context, phase0.DomainType) phase0.Domain); ok {
		r0 = rf(ctx, domainType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(phase0.Domain)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, phase0.DomainType) error); ok {
		r1 = rf(ctx, domainType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsActive provides a mock function with no fields
func (_m *Client) IsActive() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsActive")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsSynced provides a mock function with no fields
func (_m *Client) IsSynced() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsSynced")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Name provides a mock function with no fields
func (_m *Client) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Proposal provides a mock function with given fields: ctx, opts
func (_m *Client) Proposal(ctx context.Context, opts *api.ProposalOpts) (*api.Response[*api.VersionedProposal], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for Proposal")
	}

	var r0 *api.Response[*api.VersionedProposal]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.ProposalOpts) (*api.Response[*api.VersionedProposal], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.ProposalOpts) *api.Response[*api.VersionedProposal]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[*api.VersionedProposal])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.ProposalOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProposerDuties provides a mock function with given fields: ctx, opts
func (_m *Client) ProposerDuties(ctx context.Context, opts *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for ProposerDuties")
	}

	var r0 *api.Response[[]*v1.ProposerDuty]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.ProposerDutiesOpts) *api.Response[[]*v1.ProposerDuty]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[[]*v1.ProposerDuty])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.ProposerDutiesOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignedBeaconBlock provides a mock function with given fields: ctx, opts
func (_m *Client) SignedBeaconBlock(ctx context.Context, opts *api.SignedBeaconBlockOpts) (*api.Response[*spec.VersionedSignedBeaconBlock], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for SignedBeaconBlock")
	}

	var r0 *api.Response[*spec.VersionedSignedBeaconBlock]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.SignedBeaconBlockOpts) (*api.Response[*spec.VersionedSignedBeaconBlock], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.SignedBeaconBlockOpts) *api.Response[*spec.VersionedSignedBeaconBlock]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[*spec.VersionedSignedBeaconBlock])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.SignedBeaconBlockOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Spec provides a mock function with given fields: ctx, opts
func (_m *Client) Spec(ctx context.Context, opts *api.SpecOpts) (*api.Response[map[string]interface{}], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for Spec")
	}

	var r0 *api.Response[map[string]interface{}]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.SpecOpts) (*api.Response[map[string]interface{}], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.SpecOpts) *api.Response[map[string]interface{}]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[map[string]interface{}])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.SpecOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubmitAggregateAttestations provides a mock function with given fields: ctx, opts
func (_m *Client) SubmitAggregateAttestations(ctx context.Context, opts *api.SubmitAggregateAttestationsOpts) error {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for SubmitAggregateAttestations")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.SubmitAggregateAttestationsOpts) error); ok {
		r0 = rf(ctx, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitAttestations provides a mock function with given fields: ctx, opts
func (_m *Client) SubmitAttestations(ctx context.Context, opts *api.SubmitAttestationsOpts) error {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for SubmitAttestations")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.SubmitAttestationsOpts) error); ok {
		r0 = rf(ctx, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitBeaconCommitteeSubscriptions provides a mock function with given fields: ctx, subscriptions
func (_m *Client) SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*v1.BeaconCommitteeSubscription) error {
	ret := _m.Called(ctx, subscriptions)

	if len(ret) == 0 {
		panic("no return value specified for SubmitBeaconCommitteeSubscriptions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*v1.BeaconCommitteeSubscription) error); ok {
		r0 = rf(ctx, subscriptions)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitBlindedProposal provides a mock function with given fields: ctx, opts
func (_m *Client) SubmitBlindedProposal(ctx context.Context, opts *api.SubmitBlindedProposalOpts) error {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for SubmitBlindedProposal")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.SubmitBlindedProposalOpts) error); ok {
		r0 = rf(ctx, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitProposal provides a mock function with given fields: ctx, opts
func (_m *Client) SubmitProposal(ctx context.Context, opts *api.SubmitProposalOpts) error {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for SubmitProposal")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.SubmitProposalOpts) error); ok {
		r0 = rf(ctx, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitProposalPreparations provides a mock function with given fields: ctx, preparations
func (_m *Client) SubmitProposalPreparations(ctx context.Context, preparations []*v1.ProposalPreparation) error {
	ret := _m.Called(ctx, preparations)

	if len(ret) == 0 {
		panic("no return value specified for SubmitProposalPreparations")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*v1.ProposalPreparation) error); ok {
		r0 = rf(ctx, preparations)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitSyncCommitteeContributions provides a mock function with given fields: ctx, contributionAndProofs
func (_m *Client) SubmitSyncCommitteeContributions(ctx context.Context, contributionAndProofs []*altair.SignedContributionAndProof) error {
	ret := _m.Called(ctx, contributionAndProofs)

	if len(ret) == 0 {
		panic("no return value specified for SubmitSyncCommitteeContributions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*altair.SignedContributionAndProof) error); ok {
		r0 = rf(ctx, contributionAndProofs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitSyncCommitteeMessages provides a mock function with given fields: ctx, messages
func (_m *Client) SubmitSyncCommitteeMessages(ctx context.Context, messages []*altair.SyncCommitteeMessage) error {
	ret := _m.Called(ctx, messages)

	if len(ret) == 0 {
		panic("no return value specified for SubmitSyncCommitteeMessages")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*altair.SyncCommitteeMessage) error); ok {
		r0 = rf(ctx, messages)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitSyncCommitteeSubscriptions provides a mock function with given fields: ctx, subscriptions
func (_m *Client) SubmitSyncCommitteeSubscriptions(ctx context.Context, subscriptions []*v1.SyncCommitteeSubscription) error {
	ret := _m.Called(ctx, subscriptions)

	if len(ret) == 0 {
		panic("no return value specified for SubmitSyncCommitteeSubscriptions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*v1.SyncCommitteeSubscription) error); ok {
		r0 = rf(ctx, subscriptions)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitValidatorRegistrations provides a mock function with given fields: ctx, registrations
func (_m *Client) SubmitValidatorRegistrations(ctx context.Context, registrations []*api.VersionedSignedValidatorRegistration) error {
	ret := _m.Called(ctx, registrations)

	if len(ret) == 0 {
		panic("no return value specified for SubmitValidatorRegistrations")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*api.VersionedSignedValidatorRegistration) error); ok {
		r0 = rf(ctx, registrations)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SyncCommitteeContribution provides a mock function with given fields: ctx, opts
func (_m *Client) SyncCommitteeContribution(ctx context.Context, opts *api.SyncCommitteeContributionOpts) (*api.Response[*altair.SyncCommitteeContribution], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for SyncCommitteeContribution")
	}

	var r0 *api.Response[*altair.SyncCommitteeContribution]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.SyncCommitteeContributionOpts) (*api.Response[*altair.SyncCommitteeContribution], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.SyncCommitteeContributionOpts) *api.Response[*altair.SyncCommitteeContribution]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[*altair.SyncCommitteeContribution])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.SyncCommitteeContributionOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SyncCommitteeDuties provides a mock function with given fields: ctx, opts
func (_m *Client) SyncCommitteeDuties(ctx context.Context, opts *api.SyncCommitteeDutiesOpts) (*api.Response[[]*v1.SyncCommitteeDuty], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for SyncCommitteeDuties")
	}

	var r0 *api.Response[[]*v1.SyncCommitteeDuty]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.SyncCommitteeDutiesOpts) (*api.Response[[]*v1.SyncCommitteeDuty], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.SyncCommitteeDutiesOpts) *api.Response[[]*v1.SyncCommitteeDuty]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[[]*v1.SyncCommitteeDuty])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.SyncCommitteeDutiesOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Validators provides a mock function with given fields: ctx, opts
func (_m *Client) Validators(ctx context.Context, opts *api.ValidatorsOpts) (*api.Response[map[phase0.ValidatorIndex]*v1.Validator], error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for Validators")
	}

	var r0 *api.Response[map[phase0.ValidatorIndex]*v1.Validator]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.ValidatorsOpts) (*api.Response[map[phase0.ValidatorIndex]*v1.Validator], error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *api.ValidatorsOpts) *api.Response[map[phase0.ValidatorIndex]*v1.Validator]); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Response[map[phase0.ValidatorIndex]*v1.Validator])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *api.ValidatorsOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
