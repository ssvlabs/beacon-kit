package beacon

import "errors"

// SubnetID represents a subnet in the Ethereum 2.0 networking layer.
type SubnetID uint64

var (
	// ErrBlockNotFound is returned when the requested block was not found.
	ErrBlockNotFound = errors.New("block not found")
)
