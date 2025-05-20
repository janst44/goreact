package main

import (
	"net/http"

	"github.com/janst44/go-react-todo/internal/database"
	"github.com/labstack/echo/v4"
)

// GetTodos gets all todos for user
//
// @Summary Get all todos
// @Description Get all todos for user
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} database.Todo
// @Router /api/v1/todos [get]
func (app *application) handleGetTodos(c echo.Context) error {
	user := app.GetUserFromContext(c)
	todos, err := app.models.Todos.Get(user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch todos"})
	}
	return c.JSON(http.StatusOK, todos)
}

// CreateTodo creates a new todo
//
// @Summary Create a new todo
// @Description Create a new todo
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param todo body database.TodoCreate true "Todo object"
// @Success 201 {object} database.Todo
// @Router /api/v1/todos [post]
func (app *application) handleCreateTodo(c echo.Context) error {
	var input database.TodoCreate

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation failed", "details": err.Error()})
	}

	user := app.GetUserFromContext(c)
	todo, err := app.models.Todos.Insert(&input, user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create todo"})
	}

	return c.JSON(http.StatusCreated, todo)
}

// UpdateTodo updates a todo
//
// @Summary Update a todo
// @Description Update a todo
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Param todo body database.TodoPatch true "Todo object"
// @Success 200 {object} database.Todo
// @Router /api/v1/todos/{id} [patch]
func (app *application) handleUpdateTodo(c echo.Context) error {
	id := c.Param("id")
	var input database.TodoPatch

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation failed", "details": err.Error()})
	}

	user := app.GetUserFromContext(c)

	updated, err := app.models.Todos.Update(id, &input, user.Id)
	if err != nil {
		switch err.Error() {
		case "todo not found":
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Todo not found"})
		case "no updates provided":
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "No updates provided"})
		default:
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update todo"})
		}
	}

	return c.JSON(http.StatusOK, updated)
}

// DeleteTodo deletes a todo
//
// @Summary Delete a todo
// @Description Delete a todo
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 204 {object} database.Todo
// @Router /api/v1/todos/{id} [delete]
func (app *application) handleDeleteTodo(c echo.Context) error {
	id := c.Param("id")
	user := app.GetUserFromContext(c)

	if err := app.models.Todos.Delete(id, user.Id); err != nil {
		if err.Error() == "todo not found" {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Todo not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete todo"})
	}

	return c.NoContent(http.StatusNoContent)
}
