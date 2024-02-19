package main

import (
	"fmt"
	"rinha-backend/api/controller"
	"rinha-backend/api/repository"
)

func main() {
	fmt.Println("Rinha de Backend!")
	database := repository.NewDatabase()
	server := controller.NewHttpServer(database)
	server.Start()
}
