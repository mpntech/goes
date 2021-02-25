package action

import (
	"context"
	"errors"

	"github.com/modernice/goes/command"
	"github.com/modernice/goes/event"
)

var (
	// ErrMissingBus is returned when trying to publish an Event without an
	// event.Bus or dispatch a Command without a command.Bus.
	ErrMissingBus = errors.New("missing bus")
)

// Context is the context for running Actions.
type Context interface {
	context.Context

	// Action returns the currently executed Action.
	Action() Action

	// Publish publishes the given Events via the underlying Event Bus. If no
	// Event Bus is available, Publish returns ErrMissingBus.
	Publish(context.Context, ...event.Event) error

	// Dispatch dispatches the given Command via the underlying Command Bus. If
	// no Command Bus is available, Dispatch returns ErrMissingBus.
	Dispatch(context.Context, command.Command) error

	// Run runs the Action with the specified name.
	Run(context.Context, string) error
}

// A ContextOption configures a Context.
type ContextOption func(*actionContext)

type actionContext struct {
	context.Context

	act        Action
	eventBus   event.Bus
	commandBus command.Bus
	run        func(context.Context, string) error
}

// WithEventBus returns a ContextOption that provides the Context with an
// event.Bus.
func WithEventBus(bus event.Bus) ContextOption {
	return func(cfg *actionContext) {
		cfg.eventBus = bus
	}
}

// WithCommandBus returns a ContextOption that provides the Context with a
// command.Bus.
func WithCommandBus(bus command.Bus) ContextOption {
	return func(cfg *actionContext) {
		cfg.commandBus = bus
	}
}

// WithRunner returns a ContextOption that specifies the runner function that is
// called when ctx.Run is called.
func WithRunner(run func(context.Context, string) error) ContextOption {
	return func(cfg *actionContext) {
		cfg.run = run
	}
}

// NewContext returns a new Context for the given Action.
func NewContext(parent context.Context, act Action, opts ...ContextOption) Context {
	ctx := &actionContext{
		Context: parent,
		act:     act,
	}
	for _, opt := range opts {
		opt(ctx)
	}
	return ctx
}

func (ctx *actionContext) Action() Action {
	return ctx.act
}

func (ctx *actionContext) Publish(c context.Context, events ...event.Event) error {
	if ctx.eventBus == nil {
		return ErrMissingBus
	}
	return ctx.eventBus.Publish(c, events...)
}

func (ctx *actionContext) Dispatch(c context.Context, cmd command.Command) error {
	if ctx.commandBus == nil {
		return ErrMissingBus
	}
	return ctx.commandBus.Dispatch(c, cmd)
}

func (ctx *actionContext) Run(c context.Context, name string) error {
	if ctx.run == nil {
		return nil
	}
	return ctx.run(c, name)
}