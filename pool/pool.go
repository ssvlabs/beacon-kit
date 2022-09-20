package pool

import (
	"context"
	"sync"

	"github.com/bloxapp/beacon-kit"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
)

// Pool is a mapping of Client addresses to instances.
type Pool struct {
	logger    *zap.Logger
	addresses []string
	clients   map[string]beacon.Client
	connectFn func(ctx context.Context, address string) (beacon.Client, error)
	mu        sync.RWMutex
}

func NewPool(logger *zap.Logger, connectFn func(ctx context.Context, address string) (beacon.Client, error)) *Pool {
	return &Pool{
		logger:    logger,
		addresses: []string{},
		clients:   map[string]beacon.Client{},
		connectFn: connectFn,
	}
}

func (p *Pool) Update(ctx context.Context, addresses []string) error {
	newAddresses := make([]string, 0, len(addresses))
	p.mu.Lock()
	p.addresses = addresses
	for _, address := range addresses {
		if _, ok := p.clients[address]; !ok {
			newAddresses = append(newAddresses, address)
		}
	}
	p.mu.Unlock()

	if len(newAddresses) == 0 {
		return nil
	}

	var g multierror.Group
	for _, address := range newAddresses {
		addressCopy := address
		g.Go(func() error {
			p.logger.Debug("Connecting to new client", zap.String("address", addressCopy))

			client, err := p.connectFn(ctx, addressCopy)
			if err != nil {
				p.logger.Error("Failed to connect to client",
					zap.String("address", addressCopy), zap.Error(err))
				return err
			}
			p.mu.Lock()
			defer p.mu.Unlock()
			p.clients[addressCopy] = client
			return nil
		})
	}
	return g.Wait().ErrorOrNil()
}

func (p *Pool) Clients() []beacon.Client {
	p.mu.RLock()
	defer p.mu.RUnlock()

	clients := make([]beacon.Client, 0, len(p.addresses))
	for _, address := range p.addresses {
		if client, ok := p.clients[address]; ok {
			clients = append(clients, client)
		}
	}
	return clients
}
