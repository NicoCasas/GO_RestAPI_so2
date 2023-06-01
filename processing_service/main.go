package main

import (
	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/processing_service/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", controller.Index)
	r.POST("/processing/submit", controller.Submit)
	r.GET("/processing/summary", controller.Summary)
	r.Run()
}
