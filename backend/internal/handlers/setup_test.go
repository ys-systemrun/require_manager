package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/ys-systemrun/require_manager/backend/internal/config"
	"github.com/ys-systemrun/require_manager/backend/internal/database"
	"github.com/ys-systemrun/require_manager/backend/internal/models"
	"github.com/ys-systemrun/require_manager/backend/internal/router"
	"gorm.io/gorm"
)

// testEnv bundles everything a handler test needs.
type testEnv struct {
	t      *testing.T
	engine *gin.Engine
	db     *gorm.DB
}

// newTestEnv spins up an isolated in-process API backed by a temporary SQLite
// database, migrated with the real schema and served through the real router.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	dsn := filepath.Join(t.TempDir(), "test.db")
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	cfg := &config.Config{CORSAllowOrigins: []string{"http://localhost:5173"}}
	return &testEnv{t: t, engine: router.New(db, cfg), db: db}
}

// request performs an HTTP request against the router and returns the recorder.
// When body is non-nil it is JSON-encoded.
func (e *testEnv) request(method, path string, body interface{}) *httptest.ResponseRecorder {
	e.t.Helper()
	var reader *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			e.t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.engine.ServeHTTP(rec, req)
	return rec
}

// requestRaw sends an arbitrary (possibly malformed) string body.
func (e *testEnv) requestRaw(method, path, body string) *httptest.ResponseRecorder {
	e.t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.engine.ServeHTTP(rec, req)
	return rec
}

// decode unmarshals a successful JSON response body into out and asserts status.
func (e *testEnv) decode(rec *httptest.ResponseRecorder, wantStatus int, out interface{}) {
	e.t.Helper()
	if rec.Code != wantStatus {
		e.t.Fatalf("status = %d, want %d; body: %s", rec.Code, wantStatus, rec.Body.String())
	}
	if out != nil {
		if err := json.Unmarshal(rec.Body.Bytes(), out); err != nil {
			e.t.Fatalf("decode response: %v; body: %s", err, rec.Body.String())
		}
	}
}

// seedProject inserts a project directly via the DB and returns it.
func (e *testEnv) seedProject(key, name string) models.Project {
	e.t.Helper()
	p := models.Project{Name: name, Key: key}
	if err := e.db.Create(&p).Error; err != nil {
		e.t.Fatalf("seed project: %v", err)
	}
	return p
}

func mustStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, want, rec.Body.String())
	}
}
