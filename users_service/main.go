package main

import (
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
	repository.Connect()
	defer repository.Close()

	// Creamos la instancia
	r := gin.Default()

	// Declaramos los endpoints
	r.GET("/users/ping", controller.Ping)
	r.POST("/users/login", controller.Login)
	r.GET("/users/validated_ping", controller.ValidatedPing)

	// Inicia la escucha
	r.Run()
}
