package eliql

import (
	"fmt"
	"testing"
	"time"

	"github.com/Learn-Along/learn-go/projects/eliql/internal"
	mockstreamlib "github.com/Learn-Along/learn-go/projects/eliql/internal/mockstream"
	restapilib "github.com/Learn-Along/learn-go/projects/eliql/internal/restapi"
	websocketlib "github.com/Learn-Along/learn-go/projects/eliql/internal/websocket"
)

// NewCollection creates new collection
func TestNewCollection(t *testing.T) {
	bufferSize := 5
	collection := NewCollection(bufferSize)
	expected := Collection{streams: map[string]Stream{}, bufferSize: bufferSize}
	if err := expectCollectionsToMatch(collection, &expected); err != nil {
		t.Fatal(err)
	}
}

// AddWebsocketStream adds a websocket stream to the collection
func TestCollection_AddWebsocketStream(t *testing.T) {
	bufferSize := 5
	cn := Collection{streams: map[string]Stream{}, bufferSize: bufferSize}
	expectedCn := Collection{streams: map[string]Stream{}, bufferSize: bufferSize}
	urls := map[string]string{
		"test-1": "ws://example.com/results/test-1",
		"test-2": "ws://example.com/results/test-2",
		"test-3": "ws://example.com/results/test-3",
		"test-4": "ws://example.com/results/test-4",
	}

	for k, v := range urls {
		err := cn.AddWebsocketStream(k, v)
		if err != nil {
			t.Fatalf("error adding websocket stream %s", err)
		}

		expectedCn.streams[k] = websocketlib.NewWebsocketStream(v)
	}

	if err := expectCollectionsToMatch(&cn, &expectedCn); err != nil {
		t.Fatal(err)
	}
}

// AddRestAPIStream adds a websocket stream to the collection
func TestCollection_AddRestAPIStream(t *testing.T) {
	bufferSize := 5
	interval := time.Minute
	cn := Collection{streams: map[string]Stream{}, bufferSize: bufferSize}
	expectedCn := Collection{streams: map[string]Stream{}, bufferSize: bufferSize}
	urls := map[string]string{
		"test-1": "https://example.com/results/test-1",
		"test-2": "https://example.com/results/test-2",
		"test-3": "https://example.com/results/test-3",
		"test-4": "https://example.com/results/test-4",
	}

	for k, v := range urls {
		err := cn.AddRestAPIStream(k, v, interval)
		if err != nil {
			t.Fatalf("error adding REST API stream %s", err)
		}

		expectedCn.streams[k] = restapilib.NewRestAPIStream(v)
	}

	if err := expectCollectionsToMatch(&cn, &expectedCn); err != nil {
		t.Fatal(err)
	}
}

// AddStream adds a stream to the collection
func TestCollection_AddStream(t *testing.T) {
	bufferSize := 5
	interval := time.Minute
	cn := Collection{streams: map[string]Stream{}, bufferSize: bufferSize}
	expectedCn := Collection{streams: map[string]Stream{}, bufferSize: bufferSize}
	testData := map[string][]internal.Record{
		"test-1": {internal.Record{"foo": 9}},
		"test-2": {internal.Record{"foo": 19}},
		"test-3": {internal.Record{"foo": 90}},
		"test-4": {internal.Record{"foo": 900}},
	}

	for k, v := range testData {
		stream := mockstreamlib.NewMockStream(v, interval)
		err := cn.AddStream(k, stream)
		if err != nil {
			t.Fatalf("error adding stream %s", err)
		}

		expectedCn.streams[k] = stream
	}

	if err := expectCollectionsToMatch(&cn, &expectedCn); err != nil {
		t.Fatal(err)
	}
}

func TestCollection_Query(t *testing.T) {
	bufferSize := 10
	interval := time.Millisecond
	testData := map[string][]internal.Record{
		"test-1": {
			internal.Record{"foo": 9, "bar": 9, "name": "John"},
			internal.Record{"foo": 78, "bar": 49, "name": "Jane"},
			internal.Record{"foo": 1, "bar": 679, "name": "Paul"},
			internal.Record{"foo": 5, "bar": 809, "name": "Richard"},
		},
		"test-2": {
			internal.Record{"beans": 19, "java": 9, "name": "Jane"},
			internal.Record{"beans": 39, "java": 900, "name": "Paul"},
			internal.Record{"beans": 29, "java": 78, "name": "Richard"},
		},
	}
	mockCn := Collection{bufferSize: bufferSize, streams: map[string]Stream{}}

	for k, v := range testData {
		mockCn.streams[k] = mockstreamlib.NewMockStream(v, interval)
	}

	t.Run("SelectFromOneDataset", func(t *testing.T) {
		// Test selecting from one dataset
	})

	t.Run("SelectAllFromOneDataset", func(t *testing.T) {
		// Test selecting all fields from one dataset
	})

	t.Run("SelectWithComputationsFromOneDataset", func(t *testing.T) {
		// Test selecting fields with some computations from one dataset
	})

	t.Run("SelectWithInnerJoin", func(t *testing.T) {
		// Test selecting with inner joins
	})

	t.Run("SelectAllWithInnerJoin", func(t *testing.T) {
		// Test selecting all fields with inner join on the 'tables'
	})

	t.Run("SelectWithComputationsWithInnerJoin", func(t *testing.T) {
		// Test selecting fields with some computations with inner join
	})

	t.Run("SelectWithUnion", func(t *testing.T) {
		// Test selecting with union
	})

	t.Run("SelectAllWithUnion", func(t *testing.T) {
		// Test selecting all fields with union
	})

	t.Run("SelectWithComputationsWithUnion", func(t *testing.T) {
		// Test selecting fields with some computations with union
	})

	t.Run("SelectWithGroupBy", func(t *testing.T) {
		// Test selecting with group by
	})

	t.Run("SelectAllWithUnion", func(t *testing.T) {
		// Test selecting all fields with union
	})

	t.Run("SelectWithComputationsWithUnion", func(t *testing.T) {
		// Test selecting fields with some computations with union
	})

	t.Run("SelectWithOrderBy", func(t *testing.T) {
		// Test selecting with order by
	})

	t.Run("SelectAllWithUnion", func(t *testing.T) {
		// Test selecting all fields with union
	})

	t.Run("SelectWithComputationsWithUnion", func(t *testing.T) {
		// Test selecting fields with some computations with union
	})

	t.Run("SelectWithFilter", func(t *testing.T) {
		// Test selecting with `where` filters
	})

	t.Run("SelectAllWithFilter", func(t *testing.T) {
		// Test selecting all fields while filtering
	})

	t.Run("SelectWithComputationsWithFilter", func(t *testing.T) {
		// Test selecting fields with some computations while filtering
	})

	t.Run("SelectWithTime", func(t *testing.T) {
		// Test selecting with time functions
	})

	t.Run("SelectAllWithTime", func(t *testing.T) {
		// Test selecting all fields with at least one of the fileds depending on time
	})

	t.Run("SelectWithComputationsWithTime", func(t *testing.T) {
		// Test selecting fields with some computations with at least one of the fileds depending on time
	})
}

func expectCollectionsToMatch(first *Collection, second *Collection) error {
	if first.bufferSize != second.bufferSize {
		return fmt.Errorf("bufferSize: %d not equal to %d", first.bufferSize, second.bufferSize)
	}

	if len(first.streams) != len(second.streams) {
		return fmt.Errorf("streams %v doesn't match %v", first.streams, second.streams)
	}

	for k := range first.streams {
		if !first.streams[k].Equals(second.streams[k]) {
			return fmt.Errorf("streams for key '%s', %v doesn't match %v", k, first.streams[k], second.streams[k])
		}
	}

	return nil
}
