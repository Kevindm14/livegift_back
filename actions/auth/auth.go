package authuser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"livegift_back/libraries/jwt"
	"livegift_back/libraries/response"
	"livegift_back/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
)

func AuthLogin(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("error trying to connnect to database")
	}

	var user models.User
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		response.HTTPError(c.Response(), c.Request(), http.StatusBadRequest, err.Error())

		return errors.Unwrap(err)
	}

	if err := tx.Where("email = ?", user.Email).First(&user); err != nil {
		response.HTTPError(c.Response(), c.Request(), http.StatusBadRequest, "Usuario no encontrado")

		return errors.Unwrap(err)
	}

	if err := user.PasswordMatch(user.Password); !err {
		return response.HTTPError(c.Response(), c.Request(), http.StatusBadRequest, "Password not match")
	}

	claim := jwt.Claim{ID: user.ID}
	token, err := claim.GetToken(os.Getenv("SIGNING_STRING"))
	if err != nil {
		return response.HTTPError(c.Response(), c.Request(), http.StatusInternalServerError, err.Error())
	}

	mapJson := response.Map{
		"token":   token,
		"user":    user,
		"message": fmt.Sprintf("Bienvenido nuevamente %v", user.FirstName),
	}

	return response.JSON(c.Response(), c.Request(), http.StatusCreated, mapJson)
}

func AuthRegister(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("error trying to connnect to database")
	}

	var user models.User
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		response.HTTPError(c.Response(), c.Request(), http.StatusBadRequest, err.Error())

		return errors.Unwrap(err)
	}

	if err := tx.Create(&user); err != nil {
		response.HTTPError(c.Response(), c.Request(), http.StatusBadRequest, err.Error())

		return errors.Unwrap(err)
	}

	claim := jwt.Claim{ID: user.ID}
	token, err := claim.GetToken(os.Getenv("SIGNING_STRING"))
	if err != nil {
		return response.HTTPError(c.Response(), c.Request(), http.StatusInternalServerError, err.Error())
	}

	mapJson := response.Map{
		"token":   token,
		"user":    user,
		"message": fmt.Sprintf("El usuario %v ha sido creado correctamente", user.FirstName),
	}

	return response.JSON(c.Response(), c.Request(), http.StatusCreated, mapJson)
}
