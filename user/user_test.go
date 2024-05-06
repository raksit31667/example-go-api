package user

import (
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

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		assertResponseCode(t, response.Code, http.StatusCreated)

		wantedUser := User{ID: 1, Name: "raksit", Email: "raksit.m@ku.th"}
		gotUser := getUserFromResponse(t, response.Body)

		assertUserResponse(t, gotUser, wantedUser)
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

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

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

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

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

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

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

func assertResponseCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status code %d but want %d", got, want)
	}
}

func assertUserResponse(t testing.TB, got, want User) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v but want %v", got, want)
	}
}
