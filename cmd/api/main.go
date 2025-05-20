package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/janst44/go-react-todo/internal/database"
	"github.com/janst44/go-react-todo/internal/database/env"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// @title Go Web API
// @version 1.0
// @description A simple Go web API for managing todos.
// @schemes http
// @produce application/json
// @consumes application/json
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format "Bearer <token>".

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: Error loading .env file")
	}

	dbURL := env.GetEnv("SUPABASE_DB_URL", "")
	if dbURL == "" {
		fmt.Println("Error: Database URL not found in environment variables")
		return
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		return
	}
	fmt.Println("Successfully connected to database")

	models := database.NewModels(db)

	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnv("JWT_SECRET", "default_secret"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
