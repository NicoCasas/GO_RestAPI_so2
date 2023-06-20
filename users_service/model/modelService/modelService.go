package modelService

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/users_service/model"
	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/users_service/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwd_path        = "/etc/passwd"
	group_path         = "/etc/group"
	ssh_group_env_name = "SSH_CLIENTS_GROUP_NAME"
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

/*
Crea un usuario en el sistema operativo con username = user.Username y password = user.Password
en el grupo indicado por la variable de entorno SSH_CLIENTS_GROUP_NAME
*/
func CreateOSUser(user model.User) (u model.OSUser, err error) {

	if osUserExists(user.Username) {
		err = model.ErrOSUserAlreadyExists
		return
	}

	var group_id string = getGroupId(os.Getenv(ssh_group_env_name))

	cmdCreateUser := exec.Command("sudo", "useradd", "-g", group_id, "-s", "/bin/bash", user.Username)
	cmdPassWd := exec.Command("sudo", "passwd", user.Username)

	var pass []string = []string{user.Password, user.Password}
	cmdPassWd.Stdin = strings.NewReader(strings.Join(pass, "\n"))

	err = cmdCreateUser.Run()
	if err != nil {
		fmt.Println(err.Error(), ": No se pudo crear el usuario")
		return
	}

	err = cmdPassWd.Run()
	if err != nil {
		fmt.Println(err.Error(), ": No se pudo crear la contraseña")
		return
	}

	u = getOSUserByUsername(user.Username)
	return
}

func getOSUserByUsername(username string) (u model.OSUser) {
	OSUsers := GetOSUsers()
	for _, user := range OSUsers {
		if username == user.Username {
			u = user
			break
		}
	}
	return
}

func osUserExists(username string) bool {
	OSUsers := GetOSUsers()
	for _, user := range OSUsers {
		if username == user.Username {
			return true
		}
	}
	return false
}

func getGroupId(groupName string) string {
	fp, err := os.Open(group_path)

	if err != nil {
		fmt.Println("Error abriendo archivo")
		return ""
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		slice := strings.Split(scanner.Text(), ":")
		if slice[0] != groupName {
			continue
		}
		return slice[2]
	}
	return ""
}

func CreateSSHClientGroupIfNotExists() error {
	if groupExists(os.Getenv(ssh_group_env_name)) {
		return nil
	}
	cmdCreateGroup := exec.Command("sudo", "groupadd", os.Getenv(ssh_group_env_name))
	err := cmdCreateGroup.Run()
	if err != nil {
		fmt.Println(err.Error(), ": No se pudo crear el grupo")
		return err
	}
	return nil
}

func groupExists(groupName string) bool {
	return getGroupId(groupName) != ""
}
