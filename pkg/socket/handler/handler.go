package handler

import (
	"encoding/json"

	"github.com/apex/log"
)

var handlers = map[string]func(payload []byte, namespace string){}

func Handle(event string, payload interface{}, namespace string) {
	if handler, ok := handlers[event]; ok {

		jsonPayload, _ := json.Marshal(payload)

		go handler(jsonPayload, namespace)
	} else {
		log.WithField("event", event).WithField("payload", payload).Debug("Unhandled message")
	}
}
