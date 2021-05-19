package app

import (
	"github.com/LibenHailu/bookstore_users_api/controllers/ping"
	"github.com/LibenHailu/bookstore_users_api/controllers/users"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)
	router.GET("/users/:user_id", users.GetUser)
	router.GET("/users/search", users.SearchUser)
	router.POST("/users", users.CreateUser)
}
