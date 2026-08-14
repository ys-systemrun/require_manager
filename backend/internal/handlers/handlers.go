package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler carries shared dependencies for all HTTP handlers.
type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

// parseID extracts a uint path parameter, writing a 400 response on failure.
func parseID(c *gin.Context, name string) (uint, bool) {
	raw := c.Param(name)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return 0, false
	}
	return uint(id), true
}

// projectFilter reads an optional ?projectId= query parameter.
func projectFilter(c *gin.Context) (uint, bool) {
	raw := c.Query("projectId")
	if raw == "" {
		return 0, false
	}
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(id), true
}

// keywordSearch adds a portable, case-insensitive substring match across the
// given columns. Using LOWER(col) LIKE LOWER(?) keeps the query working on both
// PostgreSQL (production) and SQLite (tests) without dialect-specific ILIKE.
func keywordSearch(db *gorm.DB, keyword string, columns ...string) *gorm.DB {
	if keyword == "" || len(columns) == 0 {
		return db
	}
	like := "%" + strings.ToLower(keyword) + "%"
	clauses := make([]string, len(columns))
	args := make([]interface{}, len(columns))
	for i, col := range columns {
		clauses[i] = "LOWER(" + col + ") LIKE ?"
		args[i] = like
	}
	return db.Where(strings.Join(clauses, " OR "), args...)
}

func badRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func serverError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
