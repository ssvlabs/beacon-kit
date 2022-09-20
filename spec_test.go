package beacon

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestSpecTimeAtSlot(t *testing.T) {
	spec := &Spec{
		GenesisTime:    time.Unix(1606824000, 0),
		SecondsPerSlot: 12,
	}
	require.Equal(t, spec.TimeAtSlot(0), spec.GenesisTime)
	require.Equal(t, spec.TimeAtSlot(32), spec.GenesisTime.
		Add(32*time.Duration(spec.SecondsPerSlot)*time.Second))
}

// TestSpecAttestationSubnetID reproduces the output of the Python script described in
// attestationSubnetExpectedHash.
func TestSpecAttestationSubnetID(t *testing.T) {
	spec := &Spec{
		MaxCommitteesPerSlot:   64,
		AttestationSubnetCount: 64,
		SlotsPerEpoch:          32,
	}

	buf := bytes.NewBuffer(nil)
	for slot := 0; slot < 32; slot++ {
		for committeeIndex := 0; committeeIndex < 64; committeeIndex++ {
			fmt.Fprintf(
				buf,
				"%d,%d,%d;",
				slot,
				committeeIndex,
				spec.AttestationSubnetID(phase0.Slot(slot), phase0.CommitteeIndex(committeeIndex), 64),
			)
		}
	}
	hash := sha1.Sum(buf.Bytes())
	require.Equal(t, hex.EncodeToString(hash[:]), attestationSubnetExpectedHash)
}

// attestationSubnetExpectedHash is the output of the following Python snippet:
//   SLOTS_PER_EPOCH = 32
//   ATTESTATION_SUBNET_COUNT = 64
//
//   def compute_subnet_for_attestation(committees_per_slot: int, slot: int, committee_index: int) -> int:
//       slots_since_epoch_start = int(slot % SLOTS_PER_EPOCH)
//       committees_since_epoch_start = committees_per_slot * slots_since_epoch_start
//       return int((committees_since_epoch_start + committee_index) % ATTESTATION_SUBNET_COUNT)
//
//   output = ''
//   for slot in range(0, 32):
//       for committee in range(0, 64):
//           subnet = compute_subnet_for_attestation(64, slot, committee)
//           output += f'{slot},{committee},{subnet};'
//
//   import hashlib
//   hashlib.sha1(output.encode('utf-8')).hexdigest()
const attestationSubnetExpectedHash = "95f76cfe1f07c26d2d8d775cab47c47664679637"
