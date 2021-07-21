package handler

import "github.com/apex/log"

func handleBuild(payload interface{}) {
	log.WithField("payload", payload).Info("Handling build")
}

func init() {
	handlers["build"] = handleBuild
}
