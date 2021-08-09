package schedule

import (
	"context"
	"fmt"
	"sync"

	"github.com/modernice/goes/event"
	"github.com/modernice/goes/event/query"
	"github.com/modernice/goes/projection"
)

type schedule struct {
	store      event.Store
	eventNames []string

	triggersMux sync.RWMutex
	triggers    []chan projection.Trigger
}

func newSchedule(store event.Store, eventNames []string) *schedule {
	return &schedule{
		store:      store,
		eventNames: eventNames,
	}
}

// Trigger manually triggers the schedule. When triggering a schedule, a
// projection Job is created and passed to subscribers of the schedule. Trigger
// does not wait for the created Job to be applied. The only error ever returned
// by Trigger is ctx.Err(), if ctx is canceled before the trigger was accepted
// by every susbcriber.
//
// Queried Events
//
// By default, when a Job is created by a trigger, the event query for the Job
// queries the configured events from the beginning of time until now, sorted by
// time. This query can be overriden using the projection.Query TriggerOption:
//
//	err := schedule.Trigger(context.TODO(), projection.Query(query.New(...)))
//
// Filter Events
//
// Events can be further filtered using additional event queries. Fetched Events
// are tested against the provided Queries to determine whether they should be
// included in the created Job:
//
//	err := schedule.Trigger(context.TODO(), projection.Filter(query.New(...), query.New(...)))
//
// Difference between filters and the base query of a Job is that a Job may have
// multiple filters but only one query. The query is always used to actually
// fetch the Events from the event store while filters are applied afterwards
// (in-memory). Events must test against every provided filter to be included in
// the projection Job.
//
// Projection guards
//
// A Projection may provide a projection guard, which is just an event query.
// When a Projection provides a guard (a `ProjectionFilter() []event.Query`
// method), that guard is automatically added as a filter when a Job queries
// Events for that Projection:
//
//	type guardedProjection struct {
//		projection.Guard
//	}
//
//	schedule := schedule.Continuously(bus, store, []string{"foo", "bar", "baz"})
//
//	schedule.Subscribe(context.TODO(), func(job projection.Job) error {
//		proj := &guardedProjection{
//			Guard: projection.Guard(query.New(query.Name("foo", "bar"))),
//		}
//
//		// job.Apply queries "foo", "bar" & "baz" Events, then filters them
//		// using the projection.Guard so that only "foo" & "bar" are applied.
//		return job.Apply(job, proj)
//	})
//
//	schedule.Trigger(context.TODO())
func (schedule *schedule) Trigger(ctx context.Context, opts ...projection.TriggerOption) error {
	schedule.triggersMux.RLock()
	triggers := make([]chan projection.Trigger, len(schedule.triggers))
	copy(triggers, schedule.triggers)
	schedule.triggersMux.RUnlock()

	for _, triggers := range triggers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case triggers <- schedule.newTrigger(opts...):
		}
	}

	return nil
}

func (schedule *schedule) newTriggers() <-chan projection.Trigger {
	triggers := make(chan projection.Trigger)

	schedule.triggersMux.Lock()
	schedule.triggers = append(schedule.triggers, triggers)
	schedule.triggersMux.Unlock()

	return triggers
}

func (schedule *schedule) newTrigger(opts ...projection.TriggerOption) projection.Trigger {
	t := projection.NewTrigger(opts...)
	if t.Query == nil {
		t.Query = query.New(query.Name(schedule.eventNames...), query.SortBy(event.SortTime, event.SortAsc))
	}
	return t
}

func (schedule *schedule) removeTriggers(triggers <-chan projection.Trigger) {
	schedule.triggersMux.Lock()
	defer schedule.triggersMux.Unlock()
	for i, striggers := range schedule.triggers {
		if striggers == triggers {
			schedule.triggers = append(schedule.triggers[:i], schedule.triggers[i+1:]...)
			return
		}
	}
}

func (schedule *schedule) handleTriggers(
	ctx context.Context,
	triggers <-chan projection.Trigger,
	jobs chan<- projection.Job,
	out chan<- error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case trigger := <-triggers:
			opts := []projection.JobOption{projection.WithFilter(trigger.Filter...)}
			if trigger.Reset {
				opts = append(opts, projection.WithReset())
			}
			job := projection.NewJob(ctx, schedule.store, trigger.Query, opts...)
			select {
			case <-ctx.Done():
				return
			case jobs <- job:
			}
		}
	}
}

func (schedule *schedule) applyJobs(
	ctx context.Context,
	apply func(projection.Job) error,
	jobs <-chan projection.Job,
	out chan<- error,
	done chan struct{},
) {
	defer close(done)
	defer close(out)
	for job := range jobs {
		if err := apply(job); err != nil {
			select {
			case <-ctx.Done():
				return
			case out <- fmt.Errorf("apply Job: %w", err):
			}
		}
	}
}