package users

import (
	"fmt"
	"strings"

	"github.com/LibenHailu/bookstore_users_api/datasources/mysql/users_db"
	"github.com/LibenHailu/bookstore_users_api/logger"
	"github.com/LibenHailu/bookstore_users_api/utils/errors"
	"github.com/LibenHailu/bookstore_users_api/utils/mysql_utils"
	// "gorm.io/gorm/logger"
)

const (
	// indexUniqueEmail = "email_UNIQUE"
	// errorNoRows      = "no rows in result set"
	queryInsertUser             = "INSERT INTO users(first_name, last_name,email,date_created,password,status) VALUES (?,?,?,?,?,?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, date_created,status FROM users WHERE id=?;"
	queryUpdateUser             = "UPDATE users SET first_name = ?, last_name = ?, email =? WHERE id=?;"
	queryDeleteUser             = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus       = "SELECT id ,first_name,last_name ,email,date_created,status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id ,first_name,last_name ,email,date_created,status FROM users WHERE email=? AND password=?;"
)

// var (
// 	usersDB = make(map[int64]*User)
// )

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get  user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		// sqlErr, ok := getErr.(*mysql.MySQLError)
		// if !ok {
		// 	return errors.NewInternalServerError(fmt.Sprintf("error when trying to get user: %s", getErr.Error()))
		// }
		// if strings.Contains(getErr.Error(), errorNoRows) {
		// 	return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
		// }
		// fmt.Println(err)
		// return errors.NewInternalServerError(fmt.Sprintf("error when trying to get user %d : %s", user.Id, getErr.Error()))
		logger.Error("error when trying to get user by id", err)
		return mysql_utils.ParseError(getErr)
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)

		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Password, user.Status)

	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)

		return errors.NewInternalServerError("database error")
		// sqlErr, ok := saveErr.(*mysql.MySQLError)
		// if !ok {
		// 	return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", saveErr.Error()))
		// }
		// switch sqlErr.Number {
		// case 1062:
		// 	return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", user.Email))

		// }
		// if strings.Contains(err.Error(), indexUniqueEmail) {
		// 	return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", user.Email))
		// }
		// return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", saveErr.Error()))
	}
	userId, err := insertResult.LastInsertId()

	if err != nil {
		// return errors.NewInternalServerError("error when trying to get last inserted id")
		logger.Error("error when trying to get the last insert id after creating a new user", err)

		return errors.NewInternalServerError("database error")
	}

	user.Id = userId
	return nil
}

func (user *User) UpdateUser() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)

	if err != nil {
		logger.Error("error when trying update user", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user statement", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) SearchUser(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying to find users statment", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by status", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when trying to scan users to User struct", err)
			return nil, errors.NewInternalServerError("database error")
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil

}

func (user *User) FindByEmailAndPassword() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get  user by email and password statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return errors.NewInternalServerError("invaild user credentials")
		}
		logger.Error("error when trying to get user by email and password", err)
		return mysql_utils.ParseError(getErr)
	}
	return nil
}
