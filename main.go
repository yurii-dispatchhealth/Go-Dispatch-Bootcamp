package main

import (
	"log"
	"net/http"

	"Go-Dispatch-Bootcamp/controller"
	"Go-Dispatch-Bootcamp/router"
	"Go-Dispatch-Bootcamp/service"
	"Go-Dispatch-Bootcamp/usecase"
)

func main() {
	readerService := service.New()
	readerUsecase := usecase.New(readerService)
	readerController := controller.New(readerUsecase)
	httpRouter := router.Setup(readerController)

	err := http.ListenAndServe("localhost:8080", httpRouter)

	if err != nil {
		log.Printf("Server start error: %v", err)
	}

	log.Println("Server is started on port 8080")
}
