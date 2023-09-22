package main

import (
	"whatsApp/core"
	"whatsApp/core/memory"
	"whatsApp/initialize"
)

func main() {
	memory.Init()
	core.New().Viper = core.Viper("config.yaml")
	core.New().Db = initialize.GormMysql()
	gin := initialize.Routers()
	gin.Run(":23000")
}
