package model

const (
	ErrUserNotFound        = ModelErr("El usuario no se encuentra en la base de datos")
	ErrInvalidPass         = ModelErr("La contraseña es invalida")
	ErrUserAlreadyExists   = ModelErr("UNIQUE constraint failed: USERS.USERNAME")
	ErrOSUserAlreadyExists = ModelErr("El nombre de usuario indicado ya es un usuario valido en el sistema operativo")
)

type ModelErr string

func (e ModelErr) Error() string {
	return string(e)
}
