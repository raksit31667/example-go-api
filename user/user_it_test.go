//go:build integration

package user

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/raksit31667/example-go-api/config"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func TestCreateUserIT(t *testing.T) {
	t.Run("create user given valid user", func(t *testing.T) {
		db := newDatabase(t)
		defer db.Close()
		t.Cleanup(func() {
			db.Exec("DELETE FROM user WHERE name = $1", "raksit")
		})

		h := NewHandler(db)
		e := echo.New()
		defer e.Close()

		e.POST("/users", h.Create)

		payload := `{"name": "raksit", "email": "raksit.m@ku.th"}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assertResponseCode(t, rec.Code, http.StatusCreated)
	})
}

func newDatabase(t testing.TB) *sql.DB {
	t.Helper()

	osGetter := &config.OsEnvGetter{}
	configProvider := config.ConfigProvider{Getter: osGetter}
	cfg := configProvider.GetConfig()
	db, err := sql.Open("postgres", cfg.Server.DBConnectionString)
	if err != nil {
		t.Fatalf("failed to open database connection %v", err)
	}
	return db
}
