package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TodoModel struct {
	DB *sql.DB
}

type Todo struct {
	Id          string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title       string    `json:"title" example:"Buy groceries"`
	Description *string   `json:"description,omitempty" example:"Milk, eggs, and bread"`
	Completed   bool      `json:"completed" example:"false"`
	CreatedAt   time.Time `json:"createdAt" example:"2025-05-20T14:28:23Z"`
	UserId      string    `json:"userId" example:"user-abc-123"`
}

type TodoCreate struct {
	Title       string  `json:"title" validate:"required,min=3" example:"Call the doctor"`
	Description *string `json:"description,omitempty" example:"Schedule annual check-up"`
}

type TodoPatch struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=3" example:"Call the dentist"`
	Description *string `json:"description,omitempty" example:"Ask about whitening treatment"`
	Completed   *bool   `json:"completed,omitempty" example:"true"`
}

func (m *TodoModel) Get(userId string) ([]Todo, error) {
	rows, err := m.DB.Query(
		`SELECT id, title, description, is_completed, created_at, user_id 
		 FROM todos 
		 WHERE user_id = $1`, userId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UserId)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		todos = append(todos, t)
	}

	return todos, nil
}

func (m *TodoModel) Insert(input *TodoCreate, userId string) (*Todo, error) {
	id := uuid.New().String()
	now := time.Now()
	completed := false

	_, err := m.DB.Exec(
		`INSERT INTO todos (id, title, description, is_completed, created_at, user_id)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		id, input.Title, input.Description, completed, now, userId,
	)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return &Todo{
		Id:          id,
		Title:       input.Title,
		Description: input.Description,
		Completed:   completed,
		CreatedAt:   now,
		UserId:      userId,
	}, nil
}

func (m *TodoModel) Update(id string, patch *TodoPatch, userId string) (*Todo, error) {
	if patch == nil {
		return nil, fmt.Errorf("nil patch")
	}

	setClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	if patch.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *patch.Title)
		argIndex++
	}
	if patch.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *patch.Description)
		argIndex++
	}
	if patch.Completed != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_completed = $%d", argIndex))
		args = append(args, *patch.Completed)
		argIndex++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	setClauses = append(setClauses, fmt.Sprintf("user_id = $%d", argIndex))
	args = append(args, userId)
	argIndex++

	query := fmt.Sprintf(`UPDATE todos SET %s WHERE id = $%d AND user_id = $%d 
	RETURNING id, title, description, is_completed, created_at, user_id`,
		joinWithComma(setClauses), argIndex, argIndex+1,
	)
	args = append(args, id, userId)

	var todo Todo
	err := m.DB.QueryRow(query, args...).Scan(
		&todo.Id, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UserId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, fmt.Errorf("update failed: %w", err)
	}
	return &todo, nil
}

func joinWithComma(parts []string) string {
	return fmt.Sprintf("%s", stringJoin(parts, ", "))
}

func stringJoin(parts []string, sep string) string {
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += sep
		}
		result += part
	}
	return result
}

func (m *TodoModel) Delete(id string, userId string) error {
	res, err := m.DB.Exec("DELETE FROM todos WHERE id = $1 AND user_id = $2", id, userId)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rowsAffected failed: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("todo not found")
	}
	return nil
}
