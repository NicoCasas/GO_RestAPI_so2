package controller

const (
	ErrInvalidToken = ControllerErr("Token invalido")
)

type ControllerErr string

func (e ControllerErr) Error() string {
	return string(e)
}
