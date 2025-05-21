package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/janst44/go-react-todo/internal/database"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents the login payload
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret123"`
}

// LoginResponse represents the JWT response
type LoginResponse struct {
	Token string `json:"token" example:"your.jwt.token"`
}

// RegisterRequest represents the registration payload
type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret123"`
	Name     string `json:"name" example:"Jane Doe"`
}

// @Summary Registers a new user
// @Description Creates a user account with email, password, and name
// @Tags auth
// @Accept json
// @Produce json
// @Param user body main.RegisterRequest true "User registration payload"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Router /api/v1/auth/register [post]
func (app *application) registerUser(c echo.Context) error {
	var register RegisterRequest

	if err := c.Bind(&register); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Something went wrong"})
	}

	user := database.User{
		Email:    register.Email,
		Password: string(hashedPassword),
		Name:     register.Name,
	}

	if err := app.models.Users.Insert(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not create user"})
	}

	return c.JSON(http.StatusCreated, user)
}

// @Summary Logs in a user
// @Description Authenticates a user and returns a JWT token for future requests.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body main.LoginRequest true "User login payload"
// @Success 200 {string} main.LoginResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /api/v1/auth/login [post]
func (app *application) login(c echo.Context) error {
	var auth LoginRequest
	if err := c.Bind(&auth); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	existingUser, err := app.models.Users.GetByEmail(auth.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Something went wrong"})
	}
	if existingUser == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}

	if error := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password)); error != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": existingUser.Id,
		"exp":    time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenString, error := token.SignedString([]byte(app.jwtSecret))
	if error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error generating token"})
	}

	return c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}
