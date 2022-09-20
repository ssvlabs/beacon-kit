package clock

import (
	"context"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Clock provides a clock for the Beacon Chain.
type Clock interface {
	// Now returns the current Moment.
	Now() Moment

	// AtSlot returns the Moment at the given slot.
	AtSlot(phase0.Slot) Moment

	// AtEpoch returns the Moment at the given epoch.
	AtEpoch(phase0.Epoch) Moment

	// AtTime returns the Moment at the given time.
	AtTime(time.Time) Moment

	// EverySlot returns a channel that emits the current Moment at each slot.
	EverySlot(context.Context) <-chan Moment

	// EveryEpoch returns a channel that emits the current Moment at each epoch.
	EveryEpoch(context.Context) <-chan Moment
}

// Params contains the required parameters for Clock.
type Params struct {
	GenesisTime   time.Time
	SlotsPerEpoch phase0.Slot
	SlotDuration  time.Duration
}

type clock struct {
	*Params
}

// New returns a new Clock.
func New(params Params) Clock {
	return &clock{
		Params: &params,
	}
}

// Now returns the current Moment. It is a shorthand for
// AtTime(time.Now())
func (c clock) Now() Moment {
	return c.AtTime(time.Now())
}

// AtSlot returns the Moment at the given slot.
func (c clock) AtSlot(slot phase0.Slot) Moment {
	return Moment{
		clock: &c,
		slot:  slot,
	}
}

// AtEpoch returns the Moment at the start of the given epoch.
func (c clock) AtEpoch(epoch phase0.Epoch) Moment {
	return Moment{
		clock: &c,
		slot:  phase0.Slot(epoch) * c.SlotsPerEpoch,
	}
}

// AtTime returns the Moment at the given time.
func (c clock) AtTime(t time.Time) Moment {
	return Moment{
		clock: &c,
		slot:  phase0.Slot(t.Sub(c.GenesisTime) / c.SlotDuration),
	}
}

// EverySlot returns a channel that emits the current Moment at each slot.
func (c clock) EverySlot(ctx context.Context) <-chan Moment {
	return c.every(ctx, func() time.Time {
		return c.AtSlot(c.Now().Slot() + 1).Time()
	})
}

// EveryEpoch returns a channel that emits the current Moment at each epoch.
func (c clock) EveryEpoch(ctx context.Context) <-chan Moment {
	return c.every(ctx, func() time.Time {
		return c.AtEpoch(c.Now().Epoch() + 1).Time()
	})
}

// every emits to the returned channel at the times returned by the given
// function, until the context is cancelled.
func (c *clock) every(ctx context.Context, next func() time.Time) <-chan Moment {
	ch := make(chan Moment)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Until(next())):
				ch <- c.Now()
			}
		}
	}()
	return ch
}
