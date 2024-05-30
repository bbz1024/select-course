package main

import (
	"fmt"
	"select-course/demo1/src/storage/database"
	"select-course/demo1/src/utils/logger"
)

func main() {
	err := database.Client.Ping()
	fmt.Println(err)
	logger.Logger.Info("ok")
}
