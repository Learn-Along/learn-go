package mockstream

import (
	"fmt"
	"time"

	"github.com/Learn-Along/learn-go/projects/eliql/internal"
)

type MockStream struct {
	values   []internal.Record
	interval time.Duration
	done     chan bool
}

func NewMockStream(values []internal.Record, interval time.Duration) *MockStream {
	return &MockStream{
		values:   values,
		interval: interval,
		done:     make(chan bool),
	}
}

func (w *MockStream) Start(recv chan internal.Record) error {
	ticker := time.NewTicker(w.interval)
	index := 0
	valueCount := len(w.values)

	for {
		if index >= valueCount {
			close(recv)
			return nil
		}

		select {
		case <-w.done:
			return nil
		case <-ticker.C:
			recv <- w.values[index]
			index++
		}
	}
}

func (w *MockStream) Close() error {
	w.done <- true
	return nil
}

func (w *MockStream) Equals(other interface{}) bool {
	otherStream, ok := other.(MockStream)
	if !ok {
		return false
	}

	for i := range w.values {
		if fmt.Sprintf("%v", w.values[i]) != fmt.Sprintf("%v", otherStream.values[i]) {
			return false
		}
	}

	return true
}
