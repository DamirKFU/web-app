package users

import "app/internal/core"

type UsersController struct {
	server *core.Server
}

func NewUsersController(server *core.Server) *UsersController {
	return &UsersController{server: server}
}
