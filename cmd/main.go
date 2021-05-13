package main

import (
	"log"
	"zhashkRestApi"
	"zhashkRestApi/pkg/handler"
	"zhashkRestApi/pkg/repository"
	"zhashkRestApi/pkg/service"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(zhashkRestApi.Server)
	if err := srv.Run("8081", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
