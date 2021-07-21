package handler

import "github.com/apex/log"

var handlers = map[string]func(payload interface{}){}

func Handle(event string, payload interface{}) {
	if handler, ok := handlers[event]; ok {
		go handler(payload)
	} else {
		log.WithField("event", event).WithField("payload", payload).Debug("Unhandled message")
	}
}
