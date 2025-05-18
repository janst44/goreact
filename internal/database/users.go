package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Add debug logging
	fmt.Printf("Attempting to insert user with data:\nEmail: %s\nName: %s\nPassword length: %d\n",
		user.Email, user.Name, len(user.Password))

	query := `
		INSERT INTO users (email, password, name) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at`

	err := m.DB.QueryRowContext(ctx, query,
		user.Email,
		user.Password,
		user.Name,
	).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		fmt.Printf("Error inserting user: %v\n", err)
		return err
	}

	fmt.Printf("Successfully inserted user with ID: %s\n", user.Id)
	return nil
}

func (m *UserModel) getUser(query string, args ...interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Get(id string) (*User, error) {
	query := `
		SELECT id, email, name, password, created_at, updated_at 
		FROM users 
		WHERE id = $1`
	return m.getUser(query, id)
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, name, password, created_at, updated_at 
		FROM users 
		WHERE email = $1`
	return m.getUser(query, email)
}
