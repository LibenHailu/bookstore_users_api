package users

import (
	"net/http"
	"strconv"

	"github.com/LibenHailu/bookstore_oauth_library/oauth"
	"github.com/LibenHailu/bookstore_users_api/domains/users"
	"github.com/LibenHailu/bookstore_users_api/services"
	"github.com/LibenHailu/bookstore_users_api/utils/errors"
	"github.com/gin-gonic/gin"
)

func getUserId(userIdParams string) (int64, *errors.RestErr) {
	userId, userErr := strconv.ParseInt(userIdParams, 10, 64)
	if userErr != nil {
		return 0, errors.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}

func Get(c *gin.Context) {
	// if oauth.IsPublic(c.Request){

	// }

	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}
	// if callerId:= oauth.GetCallerId(c.Request); callerId == 0{
	// 	err:=errors.RestErr{
	// 		Status: http.StatusUnauthorized,
	// 		Message: "resource not availbale",
	// 	}
	// 	c.JSON(err.Status,err)
	// }

	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
	}
	user, getErr := services.UsersService.GetUser(userId)

	if getErr != nil {

		c.JSON(getErr.Status, getErr)
		return
	}
	if oauth.GetCallerId(c.Request) == userId {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}
	// c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func Create(c *gin.Context) {

	var user users.User
	if err := c.ShouldBind(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	result, saveErr := services.UsersService.CreateUser(user)

	if saveErr != nil {

		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))

}

func Update(c *gin.Context) {

	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
	}

	var user users.User
	if err := c.ShouldBind(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId

	isPartial := c.Request.Method == http.MethodPatch
	result, err := services.UsersService.UpdateUser(isPartial, user)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))

}

func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
	}

	if err := services.UsersService.DeleteUser(userId); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted successfully"})
}

func Search(c *gin.Context) {
	status := c.Query("status")
	users, err := services.UsersService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	// users.Marshall()
	// result := make([]interface{}, len(users))
	// for index, user := range users {
	// 	result[index] = user.Marshall(c.GetHeader("X-Public") == "true")
	// }
	c.JSON(http.StatusOK, users.Marshall(c.GetHeader("X-Public") == "true"))
}

func Login(c *gin.Context) {

	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user, err := services.UsersService.LoginUser(request)

	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}
