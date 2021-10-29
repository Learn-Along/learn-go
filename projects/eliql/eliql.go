package eliql

import (
	"time"

	"github.com/Learn-Along/learn-go/projects/eliql/internal"
)

/*
Stream The stream of map[string]interface{} that comes via its channel
*/
type Stream interface {
	Start(recv chan internal.Record) error
	Close() error
	Equals(other interface{}) bool
}

/*
Collection The collection of all initial streams from which the actual queries can be done to create new streams
*/
type Collection struct {
	streams    map[string]Stream
	bufferSize int
}

// NewCollection Creates a new collection of the given buffer size
func NewCollection(bufferSize int) *Collection {
	return nil
}

// AddWebsocketStream Add a stream from a Websocket URL
func (c *Collection) AddWebsocketStream(name string, url string) error {
	return nil
}

// AddRestAPIStream Add a stream from a REST API URL
func (c *Collection) AddRestAPIStream(name string, url string, interval time.Duration) error {
	return nil
}

// AddStream Add any generic stream that implements Stream
func (c *Collection) AddStream(name string, stream Stream) error {
	return nil
}

// Query Executes the sql string and returns a new stream that returns the appropriate records
func (c *Collection) Query(sql string) (*Stream, error) {
	return nil, nil
}
