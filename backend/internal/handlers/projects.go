package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ys-systemrun/require_manager/backend/internal/models"
	"gorm.io/gorm"
)

func (h *Handler) ListProjects(c *gin.Context) {
	var projects []models.Project
	if err := h.DB.Order("created_at asc").Find(&projects).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, projects)
}

func (h *Handler) GetProject(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var project models.Project
	if err := h.DB.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(c)
			return
		}
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *Handler) CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		badRequest(c, err)
		return
	}
	if err := h.DB.Create(&project).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, project)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var project models.Project
	if err := h.DB.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(c)
			return
		}
		serverError(c, err)
		return
	}
	var input models.Project
	if err := c.ShouldBindJSON(&input); err != nil {
		badRequest(c, err)
		return
	}
	project.Name = input.Name
	project.Key = input.Key
	project.Description = input.Description
	if err := h.DB.Save(&project).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.Project{}, id).Error; err != nil {
		serverError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ProjectSummary returns counts used by the dashboard.
func (h *Handler) ProjectSummary(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var reqCount, termCount, entityCount, relCount int64
	h.DB.Model(&models.Requirement{}).Where("project_id = ?", id).Count(&reqCount)
	h.DB.Model(&models.Term{}).Where("project_id = ?", id).Count(&termCount)
	h.DB.Model(&models.DomainEntity{}).Where("project_id = ?", id).Count(&entityCount)
	h.DB.Model(&models.DomainRelationship{}).Where("project_id = ?", id).Count(&relCount)

	// Requirement status breakdown.
	type statusRow struct {
		Status string
		Count  int64
	}
	var rows []statusRow
	h.DB.Model(&models.Requirement{}).
		Select("status, count(*) as count").
		Where("project_id = ?", id).
		Group("status").
		Scan(&rows)
	statusBreakdown := map[string]int64{}
	for _, r := range rows {
		statusBreakdown[r.Status] = r.Count
	}

	c.JSON(http.StatusOK, gin.H{
		"requirements":    reqCount,
		"terms":           termCount,
		"entities":        entityCount,
		"relationships":   relCount,
		"statusBreakdown": statusBreakdown,
	})
}
