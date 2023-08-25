//go:build mongo

package mongo_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	mongodb "go.mongodb.org/mongo-driver/mongo"

	"github.com/modernice/goes/aggregate"
	"github.com/modernice/goes/backend/mongo"
	"github.com/modernice/goes/backend/mongo/mongotest"
	"github.com/modernice/goes/backend/testing/eventstoretest"
	"github.com/modernice/goes/codec"
	"github.com/modernice/goes/event"
	etest "github.com/modernice/goes/event/test"
)

func TestEventStore(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		eventstoretest.Run(t, "mongostore", func(enc codec.Encoding) event.Store {
			return mongotest.NewEventStore(enc, mongo.URL(os.Getenv("MONGOSTORE_URL")), mongo.Database(nextEventDatabase()))
		})
	})

	t.Run("ReplicaSet", func(t *testing.T) {
		eventstoretest.Run(t, "mongostore", func(enc codec.Encoding) event.Store {
			return mongotest.NewEventStore(
				enc,
				mongo.URL(os.Getenv("MONGOREPLSTORE_URL")),
				mongo.Transactions(true),
				mongo.Database(nextEventDatabase()),
			)
		})
	})
}

var evtDBID uint64

func nextEventDatabase() string {
	id := atomic.AddUint64(&evtDBID, 1)
	return fmt.Sprintf("events_%d", id)
}

func TestEventStore_Insert_versionError(t *testing.T) {
	enc := etest.NewEncoder()
	s := mongo.NewEventStore(enc, mongo.URL(os.Getenv("MONGOSTORE_URL")), mongo.Database(nextEventDatabase()))

	if _, err := s.Connect(context.Background()); err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}

	a := aggregate.New("foo", uuid.New())

	states := s.StateCollection()
	if _, err := states.InsertOne(context.Background(), bson.M{
		"aggregateName": "foo",
		"aggregateId":   a.AggregateID(),
		"version":       5,
	}); err != nil {
		t.Fatalf("failed to insert state: %v", err)
	}

	events := []event.Event{
		event.New[any]("foo", etest.FooEventData{}, event.Aggregate(a.AggregateID(), a.AggregateName(), a.AggregateVersion()+5)),
		event.New[any]("foo", etest.FooEventData{}, event.Aggregate(a.AggregateID(), a.AggregateName(), a.AggregateVersion()+6)),
		event.New[any]("foo", etest.FooEventData{}, event.Aggregate(a.AggregateID(), a.AggregateName(), a.AggregateVersion()+7)),
	}

	err := s.Insert(context.Background(), events...)

	var versionError mongo.VersionError
	if !errors.As(err, &versionError) {
		t.Fatalf("Insert should fail a %T error; got %T", versionError, err)
	}

	if versionError.AggregateName != "foo" {
		t.Errorf("VersionError should have AggregateName %q; got %q", "foo", versionError.AggregateName)
	}

	if versionError.AggregateID != a.AggregateID() {
		t.Errorf("VersionError should have AggregateID %s; got %s", a.AggregateID(), versionError.AggregateID)
	}

	if versionError.CurrentVersion != 5 {
		t.Errorf("VersionError should have CurrentVersion %d; got %d", 5, a.AggregateVersion())
	}

	if versionError.Event != events[0] {
		t.Errorf("VersionError should have Event %v; got %v", events[0], versionError.Event)
	}
}

func TestEventStore_Insert_withTransactionHandler(t *testing.T) {
	enc := etest.NewEncoder()
	var handlerCalled int
	handler := mongo.EventHandler(func(ctx context.Context, e ...event.Event) error {
		handlerCalled++
		mongoCtx, ok := ctx.(mongodb.SessionContext)
		if !ok || mongoCtx == nil {
			t.Errorf("not mongo session context")
		}
		return nil
	})
	s := mongotest.NewEventStore(
		enc,
		mongo.URL(os.Getenv("MONGOREPLSTORE_URL")),
		mongo.Database(nextEventDatabase()),
		mongo.WithTxEventHandler(handler),
	)

	if _, err := s.Connect(context.Background()); err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}

	a := aggregate.New("foo", uuid.New())
	evt := event.New[any]("foo", etest.FooEventData{}, event.Aggregate(a.AggregateID(), a.AggregateName(), 1))
	err := s.Insert(context.Background(), evt)
	if err != nil {
		t.Errorf("error inserting %s", err)
	}

	if handlerCalled != 1 {
		t.Errorf("handler was not called once")
	}

	found, err := s.Find(context.Background(), evt.ID())
	if err != nil {
		t.Fatalf("expected store.Find not to return error; got %#v", err)
	}

	if !event.Equal(found, evt.Event()) {
		t.Errorf("found event doesn't match inserted event\ninserted: %#v\n\nfound: %#v\n\ndiff: %s", evt, found, cmp.Diff(
			evt, found, cmp.AllowUnexported(evt),
		))
	}
}

func TestEventStore_Insert_with2TransactionHandlers1Failing(t *testing.T) {
	enc := etest.NewEncoder()
	expectedErr := errors.New("some error in the handler")
	var handler1Called, handler2Called int
	handler1 := mongo.EventHandler(func(ctx context.Context, e ...event.Event) error {
		handler1Called++
		return nil
	})
	handler2 := mongo.EventHandler(func(ctx context.Context, e ...event.Event) error {
		handler2Called++
		return expectedErr
	})
	s := mongotest.NewEventStore(
		enc,
		mongo.URL(os.Getenv("MONGOREPLSTORE_URL")),
		mongo.Database(nextEventDatabase()),
		mongo.WithTxEventHandler(handler1),
		mongo.WithTxEventHandler(handler2),
	)

	if _, err := s.Connect(context.Background()); err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}

	a := aggregate.New("foo", uuid.New())
	evt := event.New[any]("foo", etest.FooEventData{}, event.Aggregate(a.AggregateID(), a.AggregateName(), 1))
	err := s.Insert(context.Background(), evt)
	if !errors.Is(err, expectedErr) {
		t.Errorf("not expected error")
	}
	if handler1Called != 1 {
		t.Errorf("handler 1 was not called once")
	}
	if handler2Called != 1 {
		t.Errorf("handler 2 was not called once")
	}
	found, err := s.Find(context.Background(), evt.ID())
	if err == nil {
		t.Errorf("expected store.Find not to return error; got %#v", err)
	}
	if found != nil {
		t.Errorf("found event, expected nil")
	}
}

func TestEventStore_WithTxEventHandler_failsWithoutTransactions(t *testing.T) {
	enc := etest.NewEncoder()
	handler := mongo.EventHandler(func(ctx context.Context, e ...event.Event) error {
		return nil
	})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()

	_ = mongotest.NewEventStore(
		enc,
		mongo.URL(os.Getenv("MONGOREPLSTORE_URL")),
		mongo.Database(nextEventDatabase()),
		mongo.WithTxEventHandler(handler),
		mongo.Transactions(false),
	)
}
