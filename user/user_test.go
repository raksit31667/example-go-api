package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
)

func TestCreateUser(t *testing.T) {
	t.Run("create user given valid user", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "raksit", "email": "raksit.m@ku.th"}`))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		row := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(createUserQuery).WithArgs("raksit", "raksit.m@ku.th").WillReturnRows(row)
		handler := NewHandler(db)
		err := handler.Create(c)

		assertNoError(t, err)
		assertResponseCode(t, response.Code, http.StatusCreated)

		want := User{ID: 1, Name: "raksit", Email: "raksit.m@ku.th"}
		got := getUserFromResponse(t, response.Body)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("create user given invalid user", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "raksit"}`))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		handler := NewHandler(nil)
		err := handler.Create(c)

		assertNoError(t, err)
		assertResponseCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("create user given error during user binding", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(``))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		handler := NewHandler(nil)
		err := handler.Create(c)

		assertNoError(t, err)
		assertResponseCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("create user given error during query", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "raksit", "email": "raksit.m@ku.th"}`))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		mock.ExpectQuery(createUserQuery).WithArgs("raksit", "raksit.m@ku.th").WillReturnError(errors.New("query error"))
		handler := NewHandler(db)
		err := handler.Create(c)

		assertNoError(t, err)
		assertResponseCode(t, response.Code, http.StatusInternalServerError)
	})
}

func TestGetAllUsers(t *testing.T) {
	t.Run("get all users given users exist in the database", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "name", "email"})
		rows.AddRow(1, "raksit", "raksit.m@ku.th")
		rows.AddRow(2, "earth", "rak-sit@hotmail.com")
		mock.ExpectQuery(getAllUsersQuery).WillReturnRows(rows)

		handler := NewHandler(db)
		err := handler.GetAll(c)

		assertNoError(t, err)
		assertResponseCode(t, response.Code, http.StatusOK)

		want := []User{
			{ID: 1, Name: "raksit", Email: "raksit.m@ku.th"},
			{ID: 2, Name: "earth", Email: "rak-sit@hotmail.com"},
		}
		got := getUsersFromResponse(t, response.Body)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("get all users given error during query", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		mock.ExpectQuery(getAllUsersQuery).WillReturnError(errors.New("query error"))
		handler := NewHandler(db)
		err := handler.GetAll(c)

		assertNoError(t, err)
		assertResponseCode(t, response.Code, http.StatusInternalServerError)
	})
}

func TestGetUserById(t *testing.T) {
	t.Run("get user by id given a user exists in the database", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("2")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		row := sqlmock.NewRows([]string{"id", "name", "email"})
		row.AddRow(2, "earth", "rak-sit@hotmail.com")
		mock.ExpectQuery(getByIdQuery).WithArgs("2").WillReturnRows(row)

		handler := NewHandler(db)
		err := handler.GetById(c)

		assertNoError(t, err)
		assertResponseCode(t, response.Code, http.StatusOK)

		want := User{
			ID: 2, Name: "earth", Email: "rak-sit@hotmail.com",
		}
		got := getUserFromResponse(t, response.Body)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("get user by id given user does not exist", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		mock.ExpectQuery(getByIdQuery).WithArgs("1").WillReturnError(sql.ErrNoRows)
		handler := NewHandler(db)
		err := handler.GetById(c)

		assertNoError(t, err)
		assertResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("get user by id given error during query", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		mock.ExpectQuery(getByIdQuery).WithArgs("1").WillReturnError(errors.New("query error"))
		handler := NewHandler(db)
		err := handler.GetById(c)

		assertNoError(t, err)
		assertResponseCode(t, response.Code, http.StatusInternalServerError)
	})
}

func getUserFromResponse(t testing.TB, body io.Reader) (user User) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(&user); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return
}

func getUsersFromResponse(t testing.TB, body io.Reader) (users []User) {
	t.Helper()

	if err := json.NewDecoder(body).Decode(&users); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
}

func assertResponseCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status code %d but want %d", got, want)
	}
}
