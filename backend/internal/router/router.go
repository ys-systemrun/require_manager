package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ys-systemrun/require_manager/backend/internal/config"
	"github.com/ys-systemrun/require_manager/backend/internal/handlers"
	"gorm.io/gorm"
)

// New builds the Gin engine with all routes and middleware registered.
func New(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	h := handlers.New(db)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		projects := api.Group("/projects")
		{
			projects.GET("", h.ListProjects)
			projects.POST("", h.CreateProject)
			projects.GET("/:id", h.GetProject)
			projects.PUT("/:id", h.UpdateProject)
			projects.DELETE("/:id", h.DeleteProject)
			projects.GET("/:id/summary", h.ProjectSummary)
		}

		reqs := api.Group("/requirements")
		{
			reqs.GET("", h.ListRequirements)
			reqs.POST("", h.CreateRequirement)
			reqs.GET("/:id", h.GetRequirement)
			reqs.PUT("/:id", h.UpdateRequirement)
			reqs.DELETE("/:id", h.DeleteRequirement)
		}

		terms := api.Group("/terms")
		{
			terms.GET("", h.ListTerms)
			terms.POST("", h.CreateTerm)
			terms.GET("/:id", h.GetTerm)
			terms.PUT("/:id", h.UpdateTerm)
			terms.DELETE("/:id", h.DeleteTerm)
		}

		entities := api.Group("/entities")
		{
			entities.GET("", h.ListEntities)
			entities.POST("", h.CreateEntity)
			entities.GET("/:id", h.GetEntity)
			entities.PUT("/:id", h.UpdateEntity)
			entities.DELETE("/:id", h.DeleteEntity)
		}

		rels := api.Group("/relationships")
		{
			rels.GET("", h.ListRelationships)
			rels.POST("", h.CreateRelationship)
			rels.PUT("/:id", h.UpdateRelationship)
			rels.DELETE("/:id", h.DeleteRelationship)
		}
	}

	return r
}
