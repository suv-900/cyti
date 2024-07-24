package api

import (
	"os"

	"github.com/op/go-logging"
)

var log *logging.Logger

func initLogger() {
	log = logging.MustGetLogger("atlas_logger")
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)

	var backend = logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	logging.SetBackend(backendFormatter)
}

func GetLogger() *logging.Logger {
	if log != nil {
		return log
	} else {
		initLogger()
		return log
	}
}
