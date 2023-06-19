package controller

import (
	"fmt"
	"net/http"
	"os"
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

func Create_user(ctx *gin.Context) {
	var code int
	var info string
	var info_key string

	stringToken := getTokenFromRequest(ctx)

	if err := validateToken(stringToken); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "No autorizado",
		})
		return
	}

	var user model.User
	if err := ctx.Bind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Credenciales incompletas",
		})
		return
	}

	err := modelService.CreateOSUser(user)
	switch err {
	case nil:
		code = http.StatusOK
		info = "Usuario creado exitosamente"
		info_key = "message"
	case model.ErrOSUserAlreadyExists:
		code = http.StatusBadRequest
		info_key = "error"
		info = model.ErrOSUserAlreadyExists.Error()
	default:
		code = http.StatusBadRequest
		info_key = "error"
		info = "No se pudo completar la solicitud"
	}

	ctx.JSON(code, gin.H{info_key: info})

}

/**
*
 */
func getTokenFromRequest(ctx *gin.Context) string {
	return ctx.GetHeader("Authentification")

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
	if !modelService.UserExists(claims["iss"].(string)) {
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
		"iss": queryUser.Username,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
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
