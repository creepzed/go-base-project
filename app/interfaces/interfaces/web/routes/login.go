package routes

import (
	"bitbucket.org/walmartdigital/hermes/app/interfaces/web/middleware/authentication"
	"bitbucket.org/walmartdigital/hermes/app/interfaces/web/models"
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/config"
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type loginHandler struct{}

func NewLoginHandler(e *echo.Echo) *loginHandler {
	loginHandler := &loginHandler{}

	e.POST("/login", loginHandler.Login)
	return loginHandler
}

func (l *loginHandler) Login(c echo.Context) error {
	if !config.GetBool("feature.flags.login") {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Description: "login disabled"})
	}

	username := c.FormValue("username")
	password := c.FormValue("password")
	client := c.FormValue("client")

	jwtKey := authentication.GetJwtKey()

	usernameValid := config.GetString("jwt.username")
	passwordValid := config.GetString("jwt.password")
	duration := config.GetDuration("jwt.duration")

	if len(username) == 0 || len(password) == 0 || len(client) == 0 {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Description: "error username, password and client fields are required"})
	}

	if len(usernameValid) == 0 || len(passwordValid) == 0 || duration == 0 {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Description: "undefined admin credentials error"})
	}

	if username != usernameValid || password != passwordValid {
		return echo.ErrUnauthorized
	}

	claims := &authentication.JwtCustomClaims{
		Client: client,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenGenerated, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		log.WithError(err).Error("error trying to generate token")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Description: "error trying to generate token"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenGenerated,
	})
}
