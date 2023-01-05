package pool

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bloxapp/beacon-kit"
	"github.com/bloxapp/beacon-kit/logging"
	"go.uber.org/zap"
)

type CallLog struct {
	Client      beacon.Client
	ClientIndex int
	Attempt     int // 0 for first attempt, 1 for second, etc.
	Start       time.Time
	End         time.Time
	Err         error
}

func (l *CallLog) String() string {
	icon := "✓"
	errorLabel := ""
	if l.Err != nil {
		icon = "⨉"
		errorLabel = fmt.Sprintf(" -> %s", l.Err)
	}
	return fmt.Sprintf("%s %s (#%d attempt) (took %s)%s",
		icon, l.Client.Address(), l.Attempt, l.End.Sub(l.Start), errorLabel)
}

type CallTrace []CallLog

func (t CallTrace) String() string {
	if len(t) == 0 {
		return "CallTrace{}"
	}
	var b strings.Builder
	b.WriteString("CallTrace:\n")
	for i, call := range t {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString("\t" + call.String())
	}
	return b.String()
}

func (t CallTrace) Errors() CallTrace {
	if len(t) == 0 {
		return nil
	}
	errors := make(CallTrace, 0, len(t))
	for _, call := range t {
		if call.Err != nil {
			errors = append(errors, call)
		}
	}
	return errors
}

type Error struct {
	Trace CallTrace
}

func (e *Error) Error() string {
	return e.Trace.Errors().String()
}

type callFunc func(context.Context, beacon.Client) error

type call struct {
	scope    Scope
	clients  []beacon.Client
	callFunc callFunc
}

func newCall(scope Scope, clients []beacon.Client, callFunc callFunc) *call {
	return &call{
		scope:    scope,
		clients:  clients,
		callFunc: callFunc,
	}
}

func (c *call) Do(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Populate `jobs` with the indices of the clients.
	jobs := make(chan int, len(c.clients)*2)
	defer close(jobs)
	for clientIndex := range c.clients {
		jobs <- clientIndex
	}

	// Spawn workers to call the clients from `jobs`.
	var (
		wg      sync.WaitGroup
		calls   = make(chan *CallLog, len(c.clients)*2)
		callers = int(c.scope.Concurrency)
	)
	if callers > len(c.clients) {
		callers = len(c.clients)
	}
	defer close(calls)
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.caller(ctx, jobs, calls)
		}()
	}

	// Receive from `calls` until jobs are done or until
	// the first success (if scope.FirstSuccess is true)
	err := c.receiveCalls(ctx, jobs, calls)

	// Wait for any remaining workers.
	cancel()
	wg.Wait()

	return err
}

func (c *call) caller(ctx context.Context, jobs <-chan int, errors chan<- *CallLog) {
	for {
		select {
		case clientIndex := <-jobs:
			start := time.Now()
			err := c.callWithTimeout(ctx, c.clients[clientIndex])
			errors <- &CallLog{
				ClientIndex: clientIndex,
				Client:      c.clients[clientIndex],
				Start:       start,
				End:         time.Now(),
				Err:         err,
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *call) callWithTimeout(ctx context.Context, client beacon.Client) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.scope.Timeout))
	defer cancel()

	return c.callFunc(ctx, client)
}

func (c *call) receiveCalls(ctx context.Context, jobs chan<- int, logs <-chan *CallLog) error {
	var (
		trace            CallTrace
		clientTries      = make([]int, len(c.clients))
		exhaustedClients int
	)

	// TODO: find a way to do this nicely.
	// defer func() {
	// 	log.Printf("Call summary: %d clients, %d responses, %d errors, %v tries",
	// 		len(c.clients), exhaustedClients, len(trace.Errors()), clientTries)
	// }()

	for {
		select {
		case <-ctx.Done():
			return nil
		case log, ok := <-logs:
			if !ok {
				return nil
			}
			log.Attempt = clientTries[log.ClientIndex]
			trace = append(trace, *log)

			// TODO: don't log like that :D
			logging.FromContext(ctx).Debug(
				fmt.Sprintf("ClientCall/%s", methodFromContext(ctx)),
				zap.Int("client_index", log.ClientIndex),
				zap.Any("log", log),
			)

			if log.Err != nil {
				if c.shouldRetryError(log.Err) {
					delay, retry := c.scope.Retry(clientTries[log.ClientIndex], log.Err)
					if retry {
						clientTries[log.ClientIndex]++
						time.Sleep(delay)
						jobs <- log.ClientIndex
						continue
					}
				}
			} else if c.scope.FirstSuccess {
				// Quit on first success.
				// TODO: we ignore previous errors here.
				return nil
			}

			// If we got here, we either succeeded or we're not retrying.
			exhaustedClients++

			if exhaustedClients == len(c.clients) {
				if len(trace.Errors()) < len(c.clients) {
					// TODO: we ignore errors here.
					return nil
				}
				return &Error{trace}
			}
		}
	}
}

func (c *call) shouldRetryError(err error) bool {
	return err != beacon.ErrBlockNotFound
}
