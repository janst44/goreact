package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/janst44/go-react-todo/internal/database"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

// RegisterUser registers a new user
// @Summary		Registers a new user
// @Description	Registers a new user
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		user	body	registerRequest	true	"User object"
// @Success		201		{object}	database.User
// @Router		/api/v1/auth/register [post]
func (app *application) registerUser(c echo.Context) error {
	var register registerRequest

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

// Login logs in a user
//
// @Summary		Logs in a user
// @Description	Logs in a user
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param		user	body	loginRequest	true	"User object"
// @Success		200		{object}	loginResponse
// @Router		/api/v1/auth/login [post]
func (app *application) login(c echo.Context) error {
	var auth loginRequest
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

	return c.JSON(http.StatusOK, loginResponse{Token: tokenString})
}
