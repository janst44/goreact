package main

import (
	"net/http"

	"github.com/janst44/go-react-todo/internal/database"
	"github.com/labstack/echo/v4"
)

// @Summary Get all todos
// @Description Retrieves a list of all todos for the authenticated user.
// @Tags todos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} database.Todo
// @Failure 401 {string} string "Unauthorized"
// @Router /api/v1/todos [get]
func (app *application) handleGetTodos(c echo.Context) error {
	user := app.GetUserFromContext(c)
	todos, err := app.models.Todos.Get(user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to fetch todos",
		})
	}
	return c.JSON(http.StatusOK, todos)
}

// @Summary Create a new todo
// @Description Adds a new todo for the authenticated user.
// @Tags todos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param todo body database.TodoCreate true "Todo object"
// @Success 201 {object} database.Todo
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized
// @Router /api/v1/todos [post]
func (app *application) handleCreateTodo(c echo.Context) error {
	var input database.TodoCreate

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
		})
	}

	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	user := app.GetUserFromContext(c)
	todo, err := app.models.Todos.Insert(&input, user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to create todo",
		})
	}

	return c.JSON(http.StatusCreated, todo)
}

// @Summary Update a todo
// @Description Updates the fields of a todo identified by ID.
// @Tags todos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param todo body database.TodoPatch true "Todo object"
// @Success 200 {object} database.Todo
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /api/v1/todos/{id} [patch]
func (app *application) handleUpdateTodo(c echo.Context) error {
	id := c.Param("id")
	var input database.TodoPatch

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
		})
	}

	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	user := app.GetUserFromContext(c)

	updated, err := app.models.Todos.Update(id, &input, user.Id)
	if err != nil {
		switch err.Error() {
		case "todo not found":
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "Todo not found",
			})
		case "no updates provided":
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    "VALIDATION_ERROR",
				Message: "No updates provided",
			})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to update todo",
			})
		}
	}

	return c.JSON(http.StatusOK, updated)
}

// @Summary Delete a todo
// @Description Deletes the todo with the specified ID.
// @Tags todos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 204 {string} string "No Content"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /api/v1/todos/{id} [delete]
func (app *application) handleDeleteTodo(c echo.Context) error {
	id := c.Param("id")
	user := app.GetUserFromContext(c)

	if err := app.models.Todos.Delete(id, user.Id); err != nil {
		if err.Error() == "todo not found" {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "Todo not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to delete todo",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
