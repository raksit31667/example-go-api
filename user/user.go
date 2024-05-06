package user

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

const createUserQuery = "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id;"

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type handler struct {
	db *sql.DB
}

func New(db *sql.DB) *handler {
	return &handler{db}
}

func (handler *handler) Create(c echo.Context) error {
	user := User{}
	err := c.Bind(&user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	context := c.Request().Context()
	var lastInsertId int
	handler.db.QueryRowContext(context, createUserQuery, user.Name, user.Email).Scan(&lastInsertId)
	user.ID = lastInsertId
	return c.JSON(http.StatusCreated, user)
}
