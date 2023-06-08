package repository

import (
	"database/sql"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/users_service/model"
	_ "github.com/mattn/go-sqlite3"
)

const (
	db_path_envname         = "DB_PATH"
	admin_username_env_name = "ADMIN_USERNAME"
	admin_password_env_name = "ADMIN_PASSWORD"
)

var database_p *sql.DB

func Connect() error {
	localDB_p, err := sql.Open("sqlite3", os.Getenv(db_path_envname))

	if err != nil {
		return err
	}

	database_p = localDB_p

	initializeUsersTableIfNotExists()

	return nil
}

func initializeUsersTableIfNotExists() {
	createUsersTableIfNotExists()
	createAdminUser()
}

func createUsersTableIfNotExists() {
	statement, err := database_p.Prepare("CREATE TABLE IF NOT EXISTS USERS(id INTEGER PRIMARY KEY AUTOINCREMENT, USERNAME TEXT UNIQUE, PASSWORD TEXT);")
	if err != nil {
		fmt.Println(err)
	}

	statement.Exec()
	statement.Close()
}

func createAdminUser() {
	var admin model.User = model.User{
		Username: os.Getenv(admin_username_env_name),
		Password: os.Getenv(admin_password_env_name),
	}

	err := SaveUser(admin)
	if err != nil && err.Error() != "UNIQUE constraint failed: USERS.USERNAME" {
		fmt.Println(err)
	}
}

func SaveUser(user model.User) error {
	safePass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stm, err := database_p.Prepare("INSERT INTO USERS (USERNAME, PASSWORD) VALUES (?,?);")
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(user.Username, safePass)

	return err
}

func GetUser(u model.User) model.User {
	var result model.User
	var id int

	stm, err := database_p.Prepare("SELECT * FROM USERS WHERE USERNAME=?;")
	if err != nil {
		fmt.Println(err)
		return result
	}
	defer stm.Close()

	err = stm.QueryRow(u.Username).Scan(&id, &result.Username, &result.Password)
	if err != nil {
		fmt.Println(err)
	}
	return result

}

func Close() {
	err := database_p.Close()
	if err != nil {
		fmt.Println(err)
	}
}
