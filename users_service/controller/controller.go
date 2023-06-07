package controller

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/users_service/model"
	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/users_service/model/modelService"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	secret_env_name = "SECRET"
)

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func ValidatedPing(ctx *gin.Context) {
	stringToken := getTokenFromRequest(ctx)

	if err := validateToken(stringToken); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "No autorizado",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})

}

/**
*
 */
func getTokenFromRequest(ctx *gin.Context) string {
	authHeader := strings.Split(ctx.GetHeader("Authorization"), " ")
	if len(authHeader) < 2 || authHeader[0] != "Bearer" {
		return ""
	}
	return authHeader[1]

}

func ListAll(ctx *gin.Context) {
	stringToken := getTokenFromRequest(ctx)

	if err := validateToken(stringToken); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "No autorizado",
		})
		return
	}

	OS_userlist := modelService.GetOSUsers()

	ctx.JSON(http.StatusOK, gin.H{
		"data": OS_userlist,
	})
}

func validateToken(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv(secret_env_name)), nil
	})

	if err != nil {
		return ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ErrInvalidToken

	}

	// Comprobamos la expiracion del token
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return ErrInvalidToken
	}
	//Comprobamos que el usuario exista
	if !modelService.UserExists(claims["user"].(string)) {
		return ErrInvalidToken
	}

	return nil
}

func Login(ctx *gin.Context) {
	// Obtenemos credenciales
	var queryUser model.User

	if err := ctx.Bind(&queryUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error leyendo body",
		})
		return
	}

	// Validamos las credenciales
	err := modelService.ValidateUser(queryUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Usuario o contraseña incorrectos",
		})
		return
	}

	// Creamos el token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": queryUser.Username,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv(secret_env_name)))

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No se pudo crear token",
		})
		return
	}
	// Generamos la resuesta

	ctx.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
