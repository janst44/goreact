package main

import (
	"net/http"

	_ "github.com/janst44/go-react-todo/docs"
	"github.com/janst44/go-react-todo/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (app *application) routes() http.Handler {
	e := echo.New()
	e.Validator = utils.NewValidator()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	v1 := e.Group("/api/v1")
	{
		v1.POST("/auth/register", app.registerUser)
		v1.POST("/auth/login", app.login)
	}

	authGroup := v1.Group("")
	authGroup.Use(app.AuthMiddleware())
	{
		authGroup.GET("/todos", app.handleGetTodos)
		authGroup.POST("/todos", app.handleCreateTodo)
		authGroup.PATCH("/todos/:id", app.handleUpdateTodo)
		authGroup.DELETE("/todos/:id", app.handleDeleteTodo)
	}
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}
