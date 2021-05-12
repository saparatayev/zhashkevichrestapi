package main

import (
	"log"
	"zhashkRestApi"
	"zhashkRestApi/pkg/handler"
)

func main() {
	handlers := new(handler.Handler)

	srv := new(zhashkRestApi.Server)
	if err := srv.Run("8081", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
