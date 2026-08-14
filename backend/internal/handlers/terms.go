package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ys-systemrun/require_manager/backend/internal/models"
	"gorm.io/gorm"
)

func (h *Handler) ListTerms(c *gin.Context) {
	q := h.DB.Order("name asc")
	if pid, ok := projectFilter(c); ok {
		q = q.Where("project_id = ?", pid)
	}
	if kw := c.Query("keyword"); kw != "" {
		q = keywordSearch(q, kw, "name", "reading", "english_name", "definition", "aliases")
	}
	var terms []models.Term
	if err := q.Find(&terms).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, terms)
}

func (h *Handler) GetTerm(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var term models.Term
	if err := h.DB.First(&term, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(c)
			return
		}
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, term)
}

func (h *Handler) CreateTerm(c *gin.Context) {
	var term models.Term
	if err := c.ShouldBindJSON(&term); err != nil {
		badRequest(c, err)
		return
	}
	if err := h.DB.Create(&term).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, term)
}

func (h *Handler) UpdateTerm(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var term models.Term
	if err := h.DB.First(&term, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(c)
			return
		}
		serverError(c, err)
		return
	}
	var input models.Term
	if err := c.ShouldBindJSON(&input); err != nil {
		badRequest(c, err)
		return
	}
	term.Name = input.Name
	term.Reading = input.Reading
	term.EnglishName = input.EnglishName
	term.Definition = input.Definition
	term.Aliases = input.Aliases
	term.Category = input.Category
	term.Notes = input.Notes
	if err := h.DB.Save(&term).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, term)
}

func (h *Handler) DeleteTerm(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.Term{}, id).Error; err != nil {
		serverError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
