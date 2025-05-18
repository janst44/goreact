package database

import (
	"database/sql"
)

// Models holds all models for the application
type Models struct {
	Todos TodoModel
	Users UserModel
}

// NewModels initializes all models with a database connection
func NewModels(db *sql.DB) Models {
	return Models{
		Todos: TodoModel{DB: db},
		Users: UserModel{DB: db},
	}
}
