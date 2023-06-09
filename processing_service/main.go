package main

import (
	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/processing_service/controller"
	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/processing_service/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	r := gin.Default()

	//r.GET("/", controller.Index)
	r.POST("/api/processing/submit", controller.Submit)
	r.GET("/api/processing/summary", controller.Summary)
	r.Run()
}
