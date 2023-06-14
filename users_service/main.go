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
	initializers.CheckSSHGroup()
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
	//r.GET("/api/users/ping", controller.Ping)
	r.POST("/api/users/login", controller.Login)
	//r.GET("/api/users/validated_ping", controller.ValidatedPing)
	r.GET("/api/users/listall", controller.ListAll)
	r.POST("/api/users/createuser", controller.Create_user)

	// Inicia la escucha
	r.Run()
}
