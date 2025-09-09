package beacon

import (
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/beacon-kit/clock"
)

// Spec contains the network-specific Beacon Chain configuration
// and provides helper methods to access it.
type Spec struct {
	Network     Network
	GenesisTime time.Time

	GenesisSlot        phase0.Slot
	GenesisForkVersion phase0.Version
	FarFutureEpoch     phase0.Epoch

	SlotsPerEpoch  phase0.Slot
	SecondsPerSlot uint64

	MaxCommitteesPerSlot          uint64
	TargetCommitteeSize           uint64
	TargetAggregatorsPerCommittee uint64
	AttestationSubnetCount        uint64

	// AttestationPropagationSlotRange is the maximum number of slots
	// during which an attestation can be propagated, after which
	// there is no point in submitting it.
	AttestationPropagationSlotRange phase0.Slot

	SyncCommitteeSize                    uint64
	TargetAggregatorsPerSyncSubcommittee uint64
	SyncCommitteeSubnetCount             uint64
	EpochsPerSyncCommitteePeriod         phase0.Epoch

	DomainBeaconProposer              phase0.DomainType
	DomainBeaconAttester              phase0.DomainType
	DomainRandao                      phase0.DomainType
	DomainDeposit                     phase0.DomainType
	DomainVoluntaryExit               phase0.DomainType
	DomainSelectionProof              phase0.DomainType
	DomainAggregateAndProof           phase0.DomainType
	DomainSyncCommittee               phase0.DomainType
	DomainSyncCommitteeSelectionProof phase0.DomainType
	DomainContributionAndProof        phase0.DomainType
	DomainApplicationMask             phase0.DomainType
	DomainApplicationBuilder          phase0.DomainType

	AltairForkEpoch    phase0.Epoch
	BellatrixForkEpoch phase0.Epoch
}

func (s *Spec) Clock() clock.Clock {
	return clock.New(clock.Params{
		GenesisTime:   s.GenesisTime,
		SlotsPerEpoch: s.SlotsPerEpoch,
		SlotDuration:  s.SlotDuration(),
	})
}

// SlotDuration returns the time.Duration of a slot.
func (s *Spec) SlotDuration() time.Duration {
	return time.Duration(s.SecondsPerSlot) * time.Second
}

// TimeAtSlot returns the time at the start of the given slot.
func (s *Spec) TimeAtSlot(slot phase0.Slot) time.Time {
	return s.GenesisTime.Add(time.Duration(slot) * s.SlotDuration())
}

// SlotAtTime returns the slot at the given time.
func (s *Spec) SlotAtTime(t time.Time) phase0.Slot {
	return phase0.Slot(t.Sub(s.GenesisTime) / s.SlotDuration())
}

// EpochFromSlot returns the epoch at the given slot.
func (s *Spec) EpochFromSlot(slot phase0.Slot) phase0.Epoch {
	return phase0.Epoch(slot / s.SlotsPerEpoch)
}

// StartSlot returns the first slot in the given epoch.
func (s *Spec) StartSlot(epoch phase0.Epoch) phase0.Slot {
	return phase0.Slot(epoch) * s.SlotsPerEpoch
}

// EndSlot returns the last slot in the given epoch.
func (s *Spec) EndSlot(epoch phase0.Epoch) phase0.Slot {
	return s.StartSlot(epoch) + s.SlotsPerEpoch - 1
}

// CommitteesAtSlot returns the number of committees per slot.
func (s *Spec) CommitteesAtSlot(activeValidators uint64) uint64 {
	n := activeValidators / uint64(s.SlotsPerEpoch) / s.TargetCommitteeSize
	if n > s.MaxCommitteesPerSlot {
		return s.MaxCommitteesPerSlot
	}
	if n == 0 {
		return 1
	}
	return n
}

// AttestationSubnetID returns the subnet ID for an attestation.
// See https://github.com/ethereum/consensus-specs/blob/395fdd456657482b7257c8b9a9d68bea68917aaf/specs/phase0/validator.md#broadcast-attestation
func (s *Spec) AttestationSubnetID(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, committeesAtSlot uint64) SubnetID {
	slotsSinceEpochStart := slot % s.SlotsPerEpoch
	committeesSinceEpochStart := committeesAtSlot * uint64(slotsSinceEpochStart)
	return SubnetID((committeesSinceEpochStart + uint64(committeeIndex)) % s.AttestationSubnetCount)
}

// SyncCommitteeSubnetID returns the subnet ID for a sync committee.
func (s *Spec) SyncCommitteeSubnetID(index phase0.CommitteeIndex) SubnetID {
	return SubnetID(index) / SubnetID(s.SyncSubcommitteeSize())
}

// SyncSubcommitteeSize returns the size of a sync subcommittee.
func (s *Spec) SyncSubcommitteeSize() uint64 {
	return s.SyncCommitteeSize / s.SyncCommitteeSubnetCount
}
