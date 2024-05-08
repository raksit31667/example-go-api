package user

import (
	"database/sql"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/raksit31667/example-go-api/middleware"
	"go.uber.org/zap"
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
	logger := middleware.GetLogger(c)

	if err := c.Bind(&user); err != nil {
		logger.Error("failed to bind user", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(user); err != nil {
		logger.Error("failed to validate user", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	context := c.Request().Context()
	var lastInsertId int
	err := handler.db.QueryRowContext(context, createUserQuery, user.Name, user.Email).Scan(&lastInsertId)
	if err != nil {
		logger.Error("failed to insert user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	user.ID = lastInsertId
	logger.Info("user created", zap.Any("user", user))
	return c.JSON(http.StatusCreated, user)
}

func (handler *handler) GetAll(c echo.Context) error {
	context := c.Request().Context()
	rows, err := handler.db.QueryContext(context, `SELECT id, name, email FROM "user"`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user := User{}
		rows.Scan(&user.ID, &user.Name, &user.Email)
		users = append(users, user)
	}

	return c.JSON(http.StatusOK, users)
}
