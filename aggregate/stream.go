package aggregate

import (
	"context"

	"github.com/modernice/goes/internal/xerror"
)

// Drain drains the given History channel and returns its Histories.
//
// Drain accepts optional error channels which will cause Drain to fail on any
// error. When Drain encounters an error from any of the error channels, the
// already drained Histories and that error are returned. Similarly, when ctx is
// canceled, the drained Histories and ctx.Err() are returned.
//
// Drain returns when the provided History channel is closed or it encounters an
// error from an error channel and does not wait for the error channels to be
// closed.
func Drain(ctx context.Context, str <-chan History, errs ...<-chan error) ([]History, error) {
	out := make([]History, 0, len(str))
	err := Walk(ctx, func(h History) { out = append(out, h) }, str, errs...)
	return out, err
}

// Walk retrieves from the given History channel until it is closed, ctx is
// closed or any of the provided error channels receives an error. For every
// History h that is received from the History channel, walkFn(h) is called.
// Should ctx be canceled before the History channel is closed, ctx.Err() is
// returned. Should an error be received from one of the optional error
// channels, that error is returned. Otherwise Walk returns nil.
//
// Example:
//
//	var repo aggregate.Repository
//	str, errs, err := repo.Query(context.TODO(), query.New())
//	// handle err
//	err := stream.Walk(context.TODO(), func(h aggregate.History) {
//		log.Println(fmt.Sprintf("Received History: %v", h))
//	}, str, errs)
//	// handle err
func Walk(
	ctx context.Context,
	walkFn func(History),
	str <-chan History,
	errs ...<-chan error,
) error {
	errChan, stop := xerror.FanIn(errs...)
	defer stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err, ok := <-errChan:
			if ok {
				return err
			}
			errChan = nil
		case his, ok := <-str:
			if !ok {
				return nil
			}
			walkFn(his)
		}
	}
}