package user

import (
	"database/sql"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

const createUserQuery = `INSERT INTO "user" (name, email) VALUES ($1, $2) RETURNING id;`

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (c *CustomValidator) Validate(i interface{}) error {
	if err := c.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

type handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *handler {
	return &handler{db}
}

func (handler *handler) Create(c echo.Context) error {
	user := User{}
	c.Echo().Validator = &CustomValidator{validator: validator.New()}

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	context := c.Request().Context()
	var lastInsertId int
	err := handler.db.QueryRowContext(context, createUserQuery, user.Name, user.Email).Scan(&lastInsertId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	user.ID = lastInsertId
	return c.JSON(http.StatusCreated, user)
}
