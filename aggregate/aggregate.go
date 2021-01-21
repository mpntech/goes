package aggregate

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/modernice/goes/aggregate/consistency"
	"github.com/modernice/goes/event"
)

// Aggregate is an event-sourced Aggregate.
type Aggregate interface {
	// AggregateID return the UUID of the aggregate.
	AggregateID() uuid.UUID

	// AggregateName returns the name of the aggregate.
	AggregateName() string

	// AggregateVersion returns the version of the aggregate.
	AggregateVersion() int

	// AggregateChanges returns the uncommited Events of the Aggregate.
	AggregateChanges() []event.Event

	// TrackChange adds the Events to the changes of the Aggregate.
	TrackChange(...event.Event) error

	// FlushChanges increases the version of the Aggregate by the number of
	// changes and empties the changes.
	FlushChanges()

	// ApplyEvent applies the Event on the Aggregate.
	ApplyEvent(event.Event)
}

// Option is an aggregate option.
type Option func(*base)

type base struct {
	id      uuid.UUID
	name    string
	version int
	changes []event.Event
}

// New returns a new base aggregate.
func New(name string, id uuid.UUID, opts ...Option) Aggregate {
	b := &base{
		id:   id,
		name: name,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// Version returns an Option that sets the version of an Aggregate.
func Version(v int) Option {
	return func(b *base) {
		b.version = v
	}
}

// Sort sorts aggregates and returns the sorted aggregates.
func Sort(as []Aggregate, s Sorting, dir SortDirection) []Aggregate {
	sorted := make([]Aggregate, len(as))
	copy(sorted, as)

	sort.Slice(sorted, func(i, j int) bool {
		switch s {
		case SortName:
			return dir.Bool(
				sorted[i].AggregateName() <= sorted[j].AggregateName(),
			)
		case SortID:
			return dir.Bool(
				sorted[i].AggregateID().String() <=
					sorted[j].AggregateID().String(),
			)
		case SortVersion:
			return dir.Bool(
				sorted[i].AggregateVersion() <= sorted[j].AggregateVersion(),
			)
		default:
			return true
		}
	})

	return sorted
}

func (b *base) AggregateID() uuid.UUID {
	return b.id
}

func (b *base) AggregateName() string {
	return b.name
}

func (b *base) AggregateVersion() int {
	return b.version
}

func (b *base) AggregateChanges() []event.Event {
	return b.changes
}

func (b *base) TrackChange(events ...event.Event) error {
	if err := consistency.Validate(b, events...); err != nil {
		return fmt.Errorf("validate consistency: %w", err)
	}
	b.changes = append(b.changes, events...)
	return nil
}

func (b *base) FlushChanges() {
	if len(b.changes) == 0 {
		return
	}
	// b.TrackChange guarantees a correct event order, so we can safely assume
	// the last element has the highest version.
	b.version = b.changes[len(b.changes)-1].AggregateVersion()
	b.changes = nil
}

// ApplyEvent does nothing. Structs that embed base should implement ApplyEvent.
func (*base) ApplyEvent(event.Event) {}