package osgi

import (
	"github.com/wttech/aemc/pkg/common/stringsx"
)

type EventList struct {
	Status string  `json:"status"`
	List   []Event `json:"data"`
}

// Event represents single OSGi event
type Event struct {
	ID       string `json:"id"`
	Topic    string `json:"topic"`
	Received int64  `json:"received"`
	Category string `json:"category"`
	Info     string `json:"info"`
}

func (el *EventList) StatusUnknown() bool {
	return len(el.List) == 0
}

func (e Event) Service() string {
	return stringsx.Between(e.Info, ", objectClass=", ", bundle=")
}

func (e Event) Details() string {
	service := e.Service()
	if len(service) > 0 {
		return service
	}
	if len(e.Info) > 0 {
		return e.Info
	}
	return e.Topic
}
