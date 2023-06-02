package modelService

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/users_service/model"
	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/users_service/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwd_path = "/etc/passwd"
)

func ValidateUser(u model.User) error {
	queryUser := repository.GetUser(u)

	if queryUser.Username == "" {
		return model.ErrUserNotFound
	}

	err := bcrypt.CompareHashAndPassword([]byte(queryUser.Password), []byte(u.Password))
	if err != nil {
		return model.ErrInvalidPass
	}

	return nil
}

func UserExists(username string) bool {
	var u model.User = model.User{Username: username}

	queryUser := repository.GetUser(u)
	return (queryUser.Username != "" && username == queryUser.Username)
}

func GetOSUsers() []model.OSUser {
	var users []model.OSUser

	fp, err := os.Open(passwd_path)

	if err != nil {
		fmt.Println("Error abriendo archivo")
		return nil
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		var u model.OSUser

		slice := strings.Split(scanner.Text(), ":")
		u.Username = slice[0]
		u.UserID, _ = strconv.Atoi(slice[2])

		users = append(users, u)
	}

	return users

}