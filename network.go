package beacon

import (
	"math"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Network is an Ethereum 2.0 network name.
type Network string

func (n Network) String() string {
	return string(n)
}

var Networks = map[Network]*Spec{
	Mainnet.Network: Mainnet,
	Holesky.Network: Holesky,
	Sepolia.Network: Sepolia,
	Hoodi.Network:   Hoodi,
}

var (
	Mainnet = &Spec{
		Network:                              "mainnet",
		GenesisTime:                          time.Unix(1606824023, 0),
		GenesisSlot:                          0,
		GenesisForkVersion:                   phase0.Version{0x0, 0x0, 0x0, 0x0},
		FarFutureEpoch:                       phase0.Epoch(math.MaxUint64),
		SlotsPerEpoch:                        32,
		SecondsPerSlot:                       12,
		MaxCommitteesPerSlot:                 64,
		TargetCommitteeSize:                  128,
		TargetAggregatorsPerCommittee:        16,
		AttestationSubnetCount:               64,
		AttestationPropagationSlotRange:      32,
		SyncCommitteeSize:                    512,
		SyncCommitteeSubnetCount:             4,
		TargetAggregatorsPerSyncSubcommittee: 16,
		EpochsPerSyncCommitteePeriod:         256,
		AltairForkEpoch:                      74240,
		BellatrixForkEpoch:                   144896,

		DomainBeaconProposer:              [4]byte{0, 0, 0, 0},
		DomainBeaconAttester:              [4]byte{1, 0, 0, 0},
		DomainRandao:                      [4]byte{2, 0, 0, 0},
		DomainDeposit:                     [4]byte{3, 0, 0, 0},
		DomainVoluntaryExit:               [4]byte{4, 0, 0, 0},
		DomainSelectionProof:              [4]byte{5, 0, 0, 0},
		DomainAggregateAndProof:           [4]byte{6, 0, 0, 0},
		DomainSyncCommittee:               [4]byte{7, 0, 0, 0},
		DomainSyncCommitteeSelectionProof: [4]byte{8, 0, 0, 0},
		DomainContributionAndProof:        [4]byte{9, 0, 0, 0},
		DomainApplicationMask:             [4]byte{0, 0, 0, 1},
		DomainApplicationBuilder:          [4]byte{0, 0, 0, 1},
	}

	Holesky = &Spec{
		Network:                              "holesky",
		GenesisTime:                          time.Unix(1695902400, 0),
		GenesisSlot:                          0,
		GenesisForkVersion:                   phase0.Version{0x01, 0x01, 0x70, 0x00},
		FarFutureEpoch:                       phase0.Epoch(math.MaxUint64),
		SlotsPerEpoch:                        32,
		SecondsPerSlot:                       12,
		MaxCommitteesPerSlot:                 64,
		TargetCommitteeSize:                  128,
		TargetAggregatorsPerCommittee:        16,
		AttestationSubnetCount:               64,
		AttestationPropagationSlotRange:      32,
		SyncCommitteeSize:                    512,
		SyncCommitteeSubnetCount:             4,
		TargetAggregatorsPerSyncSubcommittee: 16,
		EpochsPerSyncCommitteePeriod:         256,
		AltairForkEpoch:                      0,
		BellatrixForkEpoch:                   0,

		DomainBeaconProposer:              [4]byte{0, 0, 0, 0},
		DomainBeaconAttester:              [4]byte{1, 0, 0, 0},
		DomainRandao:                      [4]byte{2, 0, 0, 0},
		DomainDeposit:                     [4]byte{3, 0, 0, 0},
		DomainVoluntaryExit:               [4]byte{4, 0, 0, 0},
		DomainSelectionProof:              [4]byte{5, 0, 0, 0},
		DomainAggregateAndProof:           [4]byte{6, 0, 0, 0},
		DomainSyncCommittee:               [4]byte{7, 0, 0, 0},
		DomainSyncCommitteeSelectionProof: [4]byte{8, 0, 0, 0},
		DomainContributionAndProof:        [4]byte{9, 0, 0, 0},
		DomainApplicationMask:             [4]byte{0, 0, 0, 1},
		DomainApplicationBuilder:          [4]byte{0, 0, 0, 1},
	}

	Sepolia = &Spec{
		Network:                              "sepolia",
		GenesisTime:                          time.Unix(1655733600, 0),
		GenesisSlot:                          0,
		GenesisForkVersion:                   phase0.Version{0x90, 0x0, 0x0, 0x69},
		FarFutureEpoch:                       phase0.Epoch(math.MaxUint64),
		SlotsPerEpoch:                        32,
		SecondsPerSlot:                       12,
		MaxCommitteesPerSlot:                 64,
		TargetCommitteeSize:                  128,
		TargetAggregatorsPerCommittee:        16,
		AttestationSubnetCount:               64,
		AttestationPropagationSlotRange:      32,
		SyncCommitteeSize:                    512,
		SyncCommitteeSubnetCount:             4,
		TargetAggregatorsPerSyncSubcommittee: 16,
		EpochsPerSyncCommitteePeriod:         256,
		AltairForkEpoch:                      50,
		BellatrixForkEpoch:                   100,

		DomainBeaconProposer:              [4]byte{0, 0, 0, 0},
		DomainBeaconAttester:              [4]byte{1, 0, 0, 0},
		DomainRandao:                      [4]byte{2, 0, 0, 0},
		DomainDeposit:                     [4]byte{3, 0, 0, 0},
		DomainVoluntaryExit:               [4]byte{4, 0, 0, 0},
		DomainSelectionProof:              [4]byte{5, 0, 0, 0},
		DomainAggregateAndProof:           [4]byte{6, 0, 0, 0},
		DomainSyncCommittee:               [4]byte{7, 0, 0, 0},
		DomainSyncCommitteeSelectionProof: [4]byte{8, 0, 0, 0},
		DomainContributionAndProof:        [4]byte{9, 0, 0, 0},
		DomainApplicationMask:             [4]byte{0, 0, 0, 1},
		DomainApplicationBuilder:          [4]byte{0, 0, 0, 1},
	}

	Hoodi = &Spec{
		Network:                              "hoodi",
		GenesisTime:                          time.Unix(1742212800+600, 0),
		GenesisSlot:                          0,
		GenesisForkVersion:                   phase0.Version{0x10, 0x00, 0x09, 0x10},
		FarFutureEpoch:                       phase0.Epoch(math.MaxUint64),
		SlotsPerEpoch:                        32,
		SecondsPerSlot:                       12,
		MaxCommitteesPerSlot:                 64,
		TargetCommitteeSize:                  128,
		TargetAggregatorsPerCommittee:        16,
		AttestationSubnetCount:               64,
		AttestationPropagationSlotRange:      32,
		SyncCommitteeSize:                    512,
		SyncCommitteeSubnetCount:             4,
		TargetAggregatorsPerSyncSubcommittee: 16,
		EpochsPerSyncCommitteePeriod:         256,
		AltairForkEpoch:                      0,
		BellatrixForkEpoch:                   0,

		DomainBeaconProposer:              [4]byte{0, 0, 0, 0},
		DomainBeaconAttester:              [4]byte{1, 0, 0, 0},
		DomainRandao:                      [4]byte{2, 0, 0, 0},
		DomainDeposit:                     [4]byte{3, 0, 0, 0},
		DomainVoluntaryExit:               [4]byte{4, 0, 0, 0},
		DomainSelectionProof:              [4]byte{5, 0, 0, 0},
		DomainAggregateAndProof:           [4]byte{6, 0, 0, 0},
		DomainSyncCommittee:               [4]byte{7, 0, 0, 0},
		DomainSyncCommitteeSelectionProof: [4]byte{8, 0, 0, 0},
		DomainContributionAndProof:        [4]byte{9, 0, 0, 0},
		DomainApplicationMask:             [4]byte{0, 0, 0, 1},
		DomainApplicationBuilder:          [4]byte{0, 0, 0, 1},
	}
)
