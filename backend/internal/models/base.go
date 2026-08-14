package models

import (
	"time"

	"gorm.io/gorm"
)

// Base holds the common columns shared by all entities.
type Base struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// AllModels は AutoMigrate 対象の一覧。
func AllModels() []interface{} {
	return []interface{}{
		&Project{},
		&Requirement{},
		&Term{},
		&DomainEntity{},
		&DomainAttribute{},
		&DomainRelationship{},
	}
}
