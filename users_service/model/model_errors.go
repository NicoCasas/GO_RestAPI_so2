package model

const (
	ErrUserNotFound      = ModelErr("El usuario no se encuentra en la base de datos")
	ErrInvalidPass       = ModelErr("La contraseña es invalida")
	ErrUserAlreadyExists = ModelErr("UNIQUE constraint failed: USERS.USERNAME")
)

type ModelErr string

func (e ModelErr) Error() string {
	return string(e)
}
