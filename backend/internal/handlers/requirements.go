package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ys-systemrun/require_manager/backend/internal/models"
	"gorm.io/gorm"
)

func (h *Handler) ListRequirements(c *gin.Context) {
	q := h.DB.Order("code asc, created_at asc")
	if pid, ok := projectFilter(c); ok {
		q = q.Where("project_id = ?", pid)
	}
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if typ := c.Query("type"); typ != "" {
		q = q.Where("type = ?", typ)
	}
	if kw := c.Query("keyword"); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("title ILIKE ? OR description ILIKE ? OR code ILIKE ?", like, like, like)
	}
	var reqs []models.Requirement
	if err := q.Find(&reqs).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, reqs)
}

func (h *Handler) GetRequirement(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req models.Requirement
	if err := h.DB.Preload("Children").First(&req, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(c)
			return
		}
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *Handler) CreateRequirement(c *gin.Context) {
	var req models.Requirement
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	if err := h.DB.Create(&req).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *Handler) UpdateRequirement(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req models.Requirement
	if err := h.DB.First(&req, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(c)
			return
		}
		serverError(c, err)
		return
	}
	var input models.Requirement
	if err := c.ShouldBindJSON(&input); err != nil {
		badRequest(c, err)
		return
	}
	input.ID = req.ID
	input.ProjectID = req.ProjectID
	input.CreatedAt = req.CreatedAt
	if err := h.DB.Model(&req).Select(
		"ParentID", "Code", "Title", "Description", "Rationale",
		"Source", "Type", "Priority", "Status",
	).Updates(&input).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, input)
}

func (h *Handler) DeleteRequirement(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	// Detach children so they are not orphaned to a missing parent.
	h.DB.Model(&models.Requirement{}).Where("parent_id = ?", id).Update("parent_id", nil)
	if err := h.DB.Delete(&models.Requirement{}, id).Error; err != nil {
		serverError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
