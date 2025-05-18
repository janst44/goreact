package main

import (
	"github.com/janst44/go-react-todo/internal/database"
	"github.com/labstack/echo/v4"
)

func (app *application) GetUserFromContext(c echo.Context) *database.User {
	contextUser := c.Get("user")
	if contextUser == nil {
		return &database.User{}
	}

	user, ok := contextUser.(*database.User)
	if !ok {
		return &database.User{}
	}

	return user
}
