package stream

import (
	"context"
	"errors"
	"fmt"

	"github.com/modernice/goes/aggregate"
)

var (
	// ErrClosed is returned by a Stream when trying to read from it after it
	// has been closed.
	ErrClosed = errors.New("stream closed")
)

type stream struct {
	aggregates []aggregate.Aggregate
	aggregate  aggregate.Aggregate
	pos        int
	err        error
	closed     chan struct{}
}

// New returns an in-memory Stream filled with the provided aggregates.
func New(as ...aggregate.Aggregate) aggregate.Stream {
	return &stream{
		aggregates: as,
		closed:     make(chan struct{}),
	}
}

// All iterates over the Stream s and returns its aggregates. If a call to
// cur.Next causes an error, the already fetched aggregates and that error are
// returned. All automatically calls s.Close(ctx) when done.
func All(ctx context.Context, s aggregate.Stream) (aggregates []aggregate.Aggregate, err error) {
	defer func() {
		if cerr := s.Close(ctx); cerr != nil {
			if err != nil {
				err = fmt.Errorf("[0] %w\n[1] %s", err, cerr)
				return
			}
			err = cerr
		}
	}()
	for s.Next(ctx) {
		aggregates = append(aggregates, s.Aggregate())
	}
	err = s.Err()
	return
}

func (c *stream) Next(ctx context.Context) bool {
	select {
	case <-c.closed:
		c.err = ErrClosed
		return false
	default:
		c.err = nil
	}
	if len(c.aggregates) <= c.pos {
		return false
	}
	c.aggregate = c.aggregates[c.pos]
	c.pos++
	return true
}

func (c *stream) Aggregate() aggregate.Aggregate {
	return c.aggregate
}

func (c *stream) Err() error {
	return c.err
}

func (c *stream) Close(context.Context) error {
	select {
	case <-c.closed:
	default:
		close(c.closed)
	}
	return nil
}