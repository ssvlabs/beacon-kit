package pool

import (
	"context"
	"math/rand"
	"time"

	"github.com/bloxapp/beacon-kit"
)

// Scope dictates the behaviour of calls to the pool.
type Scope struct {
	Select       SelectFunc
	Retry        RetryFunc
	Timeout      Timeout
	Concurrency  Concurrency
	FirstSuccess FirstSuccess
	Trace        Trace
}

func (s *Scope) apply(options ...interface{}) {
	for _, option := range options {
		switch v := option.(type) {
		case SelectFunc:
			s.Select = v
		case RetryFunc:
			s.Retry = v
		case Timeout:
			s.Timeout = v
		case Concurrency:
			s.Concurrency = v
		case FirstSuccess:
			s.FirstSuccess = v
		case Trace:
			s.Trace = v
		}
	}
}

// DefaultScope returns the default Scope.
func DefaultScope() *Scope {
	return &Scope{
		Select:       SelectRandom(),
		Retry:        RetryEveryLimit(time.Millisecond*50, 2),
		Timeout:      Timeout(time.Second * 30),
		Concurrency:  4,
		FirstSuccess: true,
	}
}

// SelectFunc is a function that decides whether to call a client or not.
type SelectFunc func(size int) func(int, beacon.Client) bool

// SelectAll returns a SelectFunc which selects all clients.
func SelectAll() SelectFunc {
	return func(size int) func(int, beacon.Client) bool {
		return func(int, beacon.Client) bool {
			return true
		}
	}
}

// SelectRandom returns a SelectFunc which selects a random client.
func SelectRandom() SelectFunc {
	return func(size int) func(int, beacon.Client) bool {
		i := rand.Intn(size)
		return func(clientIndex int, client beacon.Client) bool {
			return clientIndex == i
		}
	}
}

// SelectRandoms returns a SelectFunc which selects the given number of random clients.
func SelectRandoms(count int) SelectFunc {
	return func(size int) func(int, beacon.Client) bool {
		if count >= size {
			return SelectAll()(size)
		}
		indices := rand.Perm(count)
		return func(clientIndex int, client beacon.Client) bool {
			for _, i := range indices {
				if clientIndex == i {
					return true
				}
			}
			return false
		}
	}
}

// SelectAdjacentRandoms returns a SelectFunc which selects the given number of adjacent random clients.
//
// This is useful if the order of the clients is significant. For example, given the following clients:
//
//	[Prysm, Lighthouse, Prysm, Lighthouse, Prysm, Lighthouse]
//
// SelectAdjacentRandoms(2) would always select both a Prysm and a Lighthouse.
func SelectAdjacentRandoms(count int) SelectFunc {
	return func(size int) func(int, beacon.Client) bool {
		if count >= size {
			return SelectAll()(size)
		}
		i := rand.Intn(size)
		return func(clientIndex int, client beacon.Client) bool {
			for n := 0; n < count; n++ {
				if (i+n)%size == clientIndex {
					return true
				}
			}
			return false
		}
	}
}

// RetryFunc determines whether to retry an individual client call or not.
type RetryFunc func(tries int, err error) (time.Duration, bool)

// RetryEvery returns a RetryFunc which retries every given duration.
func RetryEvery(every time.Duration) RetryFunc {
	return func(tries int, err error) (time.Duration, bool) {
		return every, true
	}
}

// RetryEveryLimit returns a RetryFunc which retries every given duration
// up to a given limit.
func RetryEveryLimit(every time.Duration, limit int) RetryFunc {
	return func(tries int, err error) (time.Duration, bool) {
		if tries >= limit {
			return 0, false
		}
		return every, true
	}
}

// Timeout is the timeout for each individual call a client.
type Timeout time.Duration

// Concurrency is the limit of concurrent calls each call
// to the pool can make.
type Concurrency int

// FirstSuccess quits after the first successful call.
type FirstSuccess bool

// Trace is a function that receives traces after each call.
type Trace func(context.Context, CallTrace)
