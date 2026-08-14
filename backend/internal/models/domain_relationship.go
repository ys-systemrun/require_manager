package models

// DomainRelationship はエンティティ間の関連。
type DomainRelationship struct {
	Base
	ProjectID      uint `gorm:"index;not null" json:"projectId"`
	SourceEntityID uint `gorm:"index;not null" json:"sourceEntityId"`
	TargetEntityID uint `gorm:"index;not null" json:"targetEntityId"`
	// Type: association / aggregation / composition / inheritance / dependency
	Type              string `gorm:"size:32;not null;default:association" json:"type"`
	Label             string `gorm:"size:200" json:"label"`
	SourceCardinality string `gorm:"size:16" json:"sourceCardinality"` // 例: 1, 0..1, *
	TargetCardinality string `gorm:"size:16" json:"targetCardinality"`

	SourceEntity *DomainEntity `gorm:"foreignKey:SourceEntityID" json:"sourceEntity,omitempty"`
	TargetEntity *DomainEntity `gorm:"foreignKey:TargetEntityID" json:"targetEntity,omitempty"`
}
