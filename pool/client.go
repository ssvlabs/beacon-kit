package pool

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/ssvlabs/beacon-kit"
)

// state holds properties for Client, so that a pointer to it can be shared for read/write
// between copies of Client (see Client.With and Client.SetClients)
type state struct {
	clients   []beacon.Client
	clientsMu sync.RWMutex

	desiredSubscriptions map[uuid.UUID]subscription
	// clientSubscriptions is map of Client.Address() -> subscription.uuid -> Context.Cancel()
	clientSubscriptions map[string]map[uuid.UUID]func()
	subscriptionsMu     sync.RWMutex
}

// Client implements a beacon.Client which replicates calls to
// a []beacon.Client slice with a per-call configurable
// concurrency, retries and client selection mechanism.
type Client struct {
	*state
	scope Scope
	methods
}

// New creates a new Client with the given clients and options.
func New(clients []beacon.Client, options ...interface{}) *Client {
	scope := DefaultScope()
	scope.apply(options...)

	client := &Client{
		state: &state{
			clients:              clients,
			desiredSubscriptions: map[uuid.UUID]subscription{},
			clientSubscriptions:  map[string]map[uuid.UUID]func(){},
		},
		scope: *scope,
	}
	client.methods = methods{
		defaultClient: client.defaultClient,
		callFunc:      client.Call,
	}
	return client
}

func (c *Client) defaultClient() beacon.Client {
	c.clientsMu.RLock()
	defer c.clientsMu.RUnlock()

	return c.clients[0]
}

// Size returns the number of clients in the pool.
func (c *Client) Size() int {
	c.clientsMu.RLock()
	defer c.clientsMu.RUnlock()

	return len(c.clients)
}

// Clients returns the clients in the pool.
func (c *Client) Clients() []beacon.Client {
	c.clientsMu.RLock()
	defer c.clientsMu.RUnlock()

	// Return a copy of the c.clients
	return append([]beacon.Client{}, c.clients...)
}

// SetClients swaps the clients in the pool with the given clients.
// Ongoing calls are not affected, only next calls will use
// the new clients.
func (c *Client) SetClients(clients []beacon.Client) error {
	c.clientsMu.Lock()
	// TODO: the clients slice is copied in With, so this won't work!
	c.clients = clients
	c.clientsMu.Unlock()

	return c.updateSubscriptions(context.Background())
}

// Scope returns the current Scope.
func (c *Client) Scope() Scope {
	return c.scope
}

// With returns a copy of Client with the given options applied to the Scope.
// Scope is isolated, but state (such as the underlying clients)
// is shared among copies.
func (c *Client) With(options ...interface{}) *Client {
	copy := *c
	copy.scope.apply(options...)

	// Update the methods to call the copy.
	copy.methods = methods{
		defaultClient: copy.defaultClient,
		callFunc:      copy.Call,
	}

	return &copy
}

// Call calls callFunc for each selected client in the pool
// with concurrency and retries according to the current Scope.
func (c *Client) Call(ctx context.Context, callFunc func(context.Context, beacon.Client) error) error {
	clients := c.Clients()

	selectFunc := c.scope.Select(len(clients))
	selectedClients := make([]beacon.Client, 0, len(clients))
	for clientIndex, client := range clients {
		if selectFunc(clientIndex, client) {
			selectedClients = append(selectedClients, client)
		}
	}
	if len(selectedClients) == 0 {
		return errors.New("no clients selected")
	}

	call := newCall(c.scope, selectedClients, callFunc)
	return call.Do(ctx)
}

type subscription struct {
	topics  []string
	handler EventHandlerFunc
}

type EventHandlerFunc func(beacon.Client, *v1.Event)

func (c *Client) Events(ctx context.Context, opts *api.EventsOpts) error {
	func() {
		c.subscriptionsMu.Lock()
		defer c.subscriptionsMu.Unlock()
		c.desiredSubscriptions[uuid.New()] = subscription{
			topics: opts.Topics,
			handler: func(c beacon.Client, e *v1.Event) {
				opts.Handler(e)
			},
		}
	}()
	return c.updateSubscriptions(ctx)
}

func (c *Client) EventsWithClient(ctx context.Context, topics []string, handler EventHandlerFunc) error {
	func() {
		c.subscriptionsMu.Lock()
		defer c.subscriptionsMu.Unlock()
		c.desiredSubscriptions[uuid.New()] = subscription{
			topics:  topics,
			handler: handler,
		}
	}()
	return c.updateSubscriptions(ctx)
}

func (c *Client) updateSubscriptions(ctx context.Context) error {
	clients := c.Clients()

	c.subscriptionsMu.Lock()
	defer c.subscriptionsMu.Unlock()

	for clientAddress, clientSubscriptions := range c.clientSubscriptions {
		var client beacon.Client
		for _, cc := range clients {
			if cc.Address() == clientAddress {
				client = cc
				break
			}
		}
		if client == nil {
			// Client has been removed — unsubscribe all.
			for _, cancel := range clientSubscriptions {
				log.Printf("Cancelling %s", clientAddress)
				cancel()
			}
			c.clientSubscriptions[clientAddress] = map[uuid.UUID]func(){}
			continue
		}
	}

	// Subscribe to desiredSubscriptions in each client, if not already subscribed.
	var g multierror.Group
	for _, client := range clients {
		clientSubscriptions := c.clientSubscriptions[client.Address()]
		for subscriptionUUID, sub := range c.desiredSubscriptions {
			_, clientSubcribed := clientSubscriptions[subscriptionUUID]
			if clientSubcribed {
				continue
			}
			provider, ok := client.(eth2client.EventsProvider)
			if !ok {
				continue
			}

			// Register client subscription.
			subscriptionCtx, cancel := context.WithCancel(context.Background())
			clientSubscriptions := c.clientSubscriptions[client.Address()]
			if clientSubscriptions == nil {
				clientSubscriptions = map[uuid.UUID]func(){}
			}
			clientSubscriptions[subscriptionUUID] = cancel
			c.clientSubscriptions[client.Address()] = clientSubscriptions

			// Subscribe.
			func(client beacon.Client, sub subscription, subscriptionUUID uuid.UUID) {
				g.Go(func() error {
					log.Printf("Subscribing %s to %s (UUID: %x)", strings.ToUpper(sub.topics[0]), client.Address(), subscriptionUUID[:])
					return provider.Events(subscriptionCtx, &api.EventsOpts{
						Topics: sub.topics,
						Handler: func(e *v1.Event) {
							sub.handler(client, e)
						},
					})
				})
			}(client, sub, subscriptionUUID)
		}
	}

	return g.Wait().ErrorOrNil()
}
