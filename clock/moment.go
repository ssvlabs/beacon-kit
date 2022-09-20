package clock

import (
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Moment represents an instant in time of the Beacon Chain
// and provides conversions between slots, epochs and time.
type Moment struct {
	clock *clock
	slot  phase0.Slot
}

// Slot returns the slot.
func (m Moment) Slot() phase0.Slot {
	return m.slot
}

// StartSlot returns the start slot of the epoch.
func (m Moment) StartSlot() phase0.Slot {
	return phase0.Slot(m.Epoch()) * m.clock.SlotsPerEpoch
}

// EndSlot returns the end slot of the epoch.
func (m Moment) EndSlot() phase0.Slot {
	return m.StartSlot() + m.clock.SlotsPerEpoch - 1
}

// Epoch returns the epoch.
func (m Moment) Epoch() phase0.Epoch {
	return phase0.Epoch(m.slot / m.clock.SlotsPerEpoch)
}

// Time returns the time.
func (m Moment) Time() time.Time {
	return m.clock.GenesisTime.Add(time.Duration(m.slot) * m.clock.SlotDuration)
}

// Until returns the duration until m. It is a shorthand for
// time.Until(m.Time())
func (m Moment) Until() time.Duration {
	return time.Until(m.Time())
}
