package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ys-systemrun/require_manager/backend/internal/models"
	"gorm.io/gorm"
)

// ---------- Entities ----------

func (h *Handler) ListEntities(c *gin.Context) {
	q := h.DB.Preload("Attributes", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order asc, id asc")
	}).Order("name asc")
	if pid, ok := projectFilter(c); ok {
		q = q.Where("project_id = ?", pid)
	}
	var entities []models.DomainEntity
	if err := q.Find(&entities).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, entities)
}

func (h *Handler) GetEntity(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var entity models.DomainEntity
	if err := h.DB.Preload("Attributes", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order asc, id asc")
	}).First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(c)
			return
		}
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *Handler) CreateEntity(c *gin.Context) {
	var entity models.DomainEntity
	if err := c.ShouldBindJSON(&entity); err != nil {
		badRequest(c, err)
		return
	}
	if err := h.DB.Create(&entity).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entity)
}

func (h *Handler) UpdateEntity(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var entity models.DomainEntity
	if err := h.DB.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(c)
			return
		}
		serverError(c, err)
		return
	}
	var input models.DomainEntity
	if err := c.ShouldBindJSON(&input); err != nil {
		badRequest(c, err)
		return
	}
	entity.Name = input.Name
	entity.DisplayName = input.DisplayName
	entity.Description = input.Description
	entity.Stereotype = input.Stereotype

	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&entity).Error; err != nil {
			return err
		}
		// Replace attributes wholesale when provided.
		if input.Attributes != nil {
			if err := tx.Where("entity_id = ?", entity.ID).Delete(&models.DomainAttribute{}).Error; err != nil {
				return err
			}
			for i := range input.Attributes {
				input.Attributes[i].ID = 0
				input.Attributes[i].EntityID = entity.ID
				input.Attributes[i].SortOrder = i
				if err := tx.Create(&input.Attributes[i]).Error; err != nil {
					return err
				}
			}
			entity.Attributes = input.Attributes
		}
		return nil
	}); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *Handler) DeleteEntity(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("entity_id = ?", id).Delete(&models.DomainAttribute{}).Error; err != nil {
			return err
		}
		if err := tx.Where("source_entity_id = ? OR target_entity_id = ?", id, id).
			Delete(&models.DomainRelationship{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.DomainEntity{}, id).Error
	}); err != nil {
		serverError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ---------- Relationships ----------

func (h *Handler) ListRelationships(c *gin.Context) {
	q := h.DB.Preload("SourceEntity").Preload("TargetEntity").Order("created_at asc")
	if pid, ok := projectFilter(c); ok {
		q = q.Where("project_id = ?", pid)
	}
	var rels []models.DomainRelationship
	if err := q.Find(&rels).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, rels)
}

func (h *Handler) CreateRelationship(c *gin.Context) {
	var rel models.DomainRelationship
	if err := c.ShouldBindJSON(&rel); err != nil {
		badRequest(c, err)
		return
	}
	if err := h.DB.Create(&rel).Error; err != nil {
		serverError(c, err)
		return
	}
	h.DB.Preload("SourceEntity").Preload("TargetEntity").First(&rel, rel.ID)
	c.JSON(http.StatusCreated, rel)
}

func (h *Handler) UpdateRelationship(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var rel models.DomainRelationship
	if err := h.DB.First(&rel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(c)
			return
		}
		serverError(c, err)
		return
	}
	var input models.DomainRelationship
	if err := c.ShouldBindJSON(&input); err != nil {
		badRequest(c, err)
		return
	}
	rel.SourceEntityID = input.SourceEntityID
	rel.TargetEntityID = input.TargetEntityID
	rel.Type = input.Type
	rel.Label = input.Label
	rel.SourceCardinality = input.SourceCardinality
	rel.TargetCardinality = input.TargetCardinality
	if err := h.DB.Save(&rel).Error; err != nil {
		serverError(c, err)
		return
	}
	h.DB.Preload("SourceEntity").Preload("TargetEntity").First(&rel, rel.ID)
	c.JSON(http.StatusOK, rel)
}

func (h *Handler) DeleteRelationship(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.DB.Delete(&models.DomainRelationship{}, id).Error; err != nil {
		serverError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
