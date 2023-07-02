package main

import (
	config2 "warehouse-application/internal/config"
	"warehouse-application/internal/postgresSQL"
	"warehouse-application/internal/server"
)

var configPath = "/Users/kare/GolandProjects/warehouse-application/config/config.yaml"

func main() {
	config := config2.MustLoad(configPath)
	dataSourceName := config.PostgresConfig.GetDataSourceName()
	db := postgresSQL.Connect(dataSourceName)
	defer postgresSQL.Disconnect(db)
	server.StartWebServer(db)
}
