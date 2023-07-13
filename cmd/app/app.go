package main

import (
	"warehouse-application/pkg/logging"
	"warehouse-application/pkg/server"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("start app")
	server.StartServer()
}
