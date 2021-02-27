package snapshot_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/modernice/goes/aggregate"
	"github.com/modernice/goes/aggregate/snapshot"
)

type mockAggregate struct {
	aggregate.Aggregate
	mockState
}

type mockSnapshot mockAggregate

type mockState struct {
	A bool
	B int
	C string
}

func TestMarshal_default(t *testing.T) {
	a := &mockAggregate{
		Aggregate: aggregate.New("foo", uuid.New()),
		mockState: mockState{
			A: true,
			B: -10,
			C: "foo",
		},
	}

	if _, err := snapshot.Marshal(a); err != nil {
		t.Fatalf("Marshal shouldn't fail; failed with %q", err)
	}
}

func TestMarshal_marshaler(t *testing.T) {
	a := &mockSnapshot{
		Aggregate: aggregate.New("foo", uuid.New()),
		mockState: mockState{
			A: true,
			B: -10,
			C: "foo",
		},
	}

	if _, err := snapshot.Marshal(a); err != nil {
		t.Fatalf("Marshal shouldn't fail; failed with %q", err)
	}
}

func TestUnmarshal(t *testing.T) {
	a := &mockAggregate{
		Aggregate: aggregate.New("foo", uuid.New()),
		mockState: mockState{
			A: true,
			B: -10,
			C: "foo",
		},
	}

	b, err := snapshot.Marshal(a)
	if err != nil {
		t.Fatalf("Marshal shouldn't fail; failed with %q", err)
	}

	snap := &mockAggregate{Aggregate: aggregate.New("foo", uuid.New())}

	if err = snapshot.Unmarshal(b, snap); err != nil {
		t.Fatalf("Unmarshal shouldn't fail; failed with %q", err)
	}

	var want mockState
	if snap.mockState != want {
		t.Errorf("unmarshaled state should be zero value. want=%v got=%v", want, snap.mockState)
	}
}

func TestUnmarshal_unmarshaler(t *testing.T) {
	a := &mockSnapshot{
		Aggregate: aggregate.New("foo", uuid.New()),
		mockState: mockState{
			A: true,
			B: -10,
			C: "foo",
		},
	}

	b, err := snapshot.Marshal(a)
	if err != nil {
		t.Fatalf("Marshal shouldn't fail; failed with %q", err)
	}

	snap := &mockSnapshot{Aggregate: aggregate.New("foo", uuid.New())}

	if err = snapshot.Unmarshal(b, snap); err != nil {
		t.Fatalf("Unmarshal shouldn't fail; failed with %q", err)
	}

	if snap.mockState != a.mockState {
		t.Errorf("unmarshaled state differs from original. want=%v got=%v", a.mockState, snap.mockState)
	}
}

func (a *mockSnapshot) MarshalSnapshot() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(a.mockState); err != nil {
		return nil, fmt.Errorf("gob: %w", err)
	}
	return buf.Bytes(), nil
}

func (a *mockSnapshot) UnmarshalSnapshot(p []byte) error {
	if err := gob.NewDecoder(bytes.NewReader(p)).Decode(&a.mockState); err != nil {
		return fmt.Errorf("gob: %w", err)
	}
	return nil
}
