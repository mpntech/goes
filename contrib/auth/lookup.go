package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/modernice/goes/event"
	"github.com/modernice/goes/projection/lookup"
	"github.com/modernice/goes/projection/schedule"
)

const (
	// LookupActor looks up the aggregate id of an actor from a given actor id.
	LookupActor = "actor"

	// LookupRole looks up the aggregate id of a role from a given role name.
	LookupRole = "role"
)

// Lookup provides lookups from actor ids to aggregate ids of those actors.
type Lookup struct {
	*lookup.Lookup
}

var lookupEvents = [...]string{ActorIdentified, RoleIdentified}

// NewLookup returns a new lookup for aggregate ids of actors.
func NewLookup(store event.Store, bus event.Bus, opts ...schedule.ContinuousOption) *Lookup {
	return &Lookup{Lookup: lookup.New(store, bus, lookupEvents[:], opts...)}
}

// Actor returns the aggregate id of the actor with the given formatted actor id.
func (l *Lookup) Actor(ctx context.Context, id string) (uuid.UUID, bool) {
	return l.Reverse(ctx, ActorAggregate, LookupActor, id)
}

// Role returns the aggregate id of the role with the given name.
func (l *Lookup) Role(ctx context.Context, name string) (uuid.UUID, bool) {
	return l.Reverse(ctx, RoleAggregate, LookupRole, name)
}

// ProvideLookup implements lookup.Event.
func (data ActorIdentifiedData) ProvideLookup(p lookup.Provider) {
	p.Provide(LookupActor, string(data))
}

// ProvideLookup implements lookup.Event.
func (data RoleIdentifiedData) ProvideLookup(p lookup.Provider) {
	p.Provide(LookupRole, string(data))
}
