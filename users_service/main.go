package main

import (
	"log"

	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/users_service/controller"
	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/users_service/initializers"
	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/users_service/repository"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	// Levantamos la base de datos
	err := repository.Connect()
	if err != nil {
		log.Fatal("No se pudo conectar a la bdd")
	}
	defer repository.Close()

	// Creamos la instancia
	r := gin.Default()

	// Declaramos los endpoints
	r.GET("/users/ping", controller.Ping)
	r.POST("/users/login", controller.Login)
	r.GET("/users/validated_ping", controller.ValidatedPing)
	r.GET("/users/listall", controller.ListAll)

	// Inicia la escucha
	r.Run()
}
