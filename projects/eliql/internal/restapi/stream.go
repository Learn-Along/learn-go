package restapi

import (
	"github.com/Learn-Along/learn-go/projects/eliql/internal"
)

/*
The stream of map[string]interface{} that comes from the REST API URL
*/
type RestAPIStream struct {
	Url string
}

func NewRestAPIStream(url string) *RestAPIStream {
	return &RestAPIStream{Url: url}
}

func (w *RestAPIStream) Start(recv chan internal.Record) error {
	return nil
}

func (w *RestAPIStream) Close() error {
	return nil
}

func (w *RestAPIStream) Equals(other interface{}) bool {
	otherStream, ok := other.(RestAPIStream)
	if !ok {
		return false
	}

	return otherStream.Url == w.Url
}
