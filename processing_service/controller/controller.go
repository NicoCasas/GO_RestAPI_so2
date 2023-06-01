package controller

import (
	"net/http"

	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/processing_service/model"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome! "})
}

func Submit(c *gin.Context) {
	model.Increment_value()
	c.JSON(http.StatusOK, gin.H{"message": "The counter has been incremented"})
}

func Summary(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"counter": model.Get_value()})

}
