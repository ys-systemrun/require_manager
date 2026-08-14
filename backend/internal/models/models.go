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

// Project (プロジェクト) groups all requirements, terms and domain models.
type Project struct {
	Base
	Name        string `gorm:"size:200;not null" json:"name"`
	Key         string `gorm:"size:32;uniqueIndex;not null" json:"key"` // 例: SHOP
	Description  string `gorm:"type:text" json:"description"`

	Requirements  []Requirement        `json:"requirements,omitempty"`
	Terms         []Term               `json:"terms,omitempty"`
	Entities      []DomainEntity       `json:"entities,omitempty"`
	Relationships []DomainRelationship `json:"relationships,omitempty"`
}

// Requirement (要求・要件) は要求管理の中核。親子で階層化できる。
type Requirement struct {
	Base
	ProjectID uint `gorm:"index;not null" json:"projectId"`
	ParentID  *uint `gorm:"index" json:"parentId"`

	Code        string `gorm:"size:64;index" json:"code"` // 例: REQ-001
	Title       string `gorm:"size:300;not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Rationale   string `gorm:"type:text" json:"rationale"`   // 根拠・理由
	Source      string `gorm:"size:300" json:"source"`       // 出典 (ステークホルダー等)

	// Type: business(業務要求) / functional(機能要件) / non_functional(非機能要件) / constraint(制約)
	Type string `gorm:"size:32;not null;default:functional" json:"type"`
	// Priority: high / medium / low
	Priority string `gorm:"size:16;not null;default:medium" json:"priority"`
	// Status: draft / proposed / approved / implemented / verified / rejected / deprecated
	Status string `gorm:"size:16;not null;default:draft" json:"status"`

	Children []Requirement `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// Term (用語) は用語辞書の1エントリ。
type Term struct {
	Base
	ProjectID   uint   `gorm:"index;not null" json:"projectId"`
	Name        string `gorm:"size:200;not null" json:"name"`        // 用語
	Reading     string `gorm:"size:200" json:"reading"`              // 読み仮名
	EnglishName string `gorm:"size:200" json:"englishName"`          // 英語表記
	Definition  string `gorm:"type:text;not null" json:"definition"` // 定義
	Aliases     string `gorm:"type:text" json:"aliases"`             // 別名・同義語 (カンマ区切り)
	Category    string `gorm:"size:100" json:"category"`             // 分類
	Notes       string `gorm:"type:text" json:"notes"`               // 備考
}

// DomainEntity (エンティティ) はドメインモデルの構成要素。
type DomainEntity struct {
	Base
	ProjectID   uint   `gorm:"index;not null" json:"projectId"`
	Name        string `gorm:"size:200;not null" json:"name"`
	DisplayName string `gorm:"size:200" json:"displayName"` // 日本語名
	Description string `gorm:"type:text" json:"description"`
	// Stereotype: entity / value_object / aggregate_root / service / event
	Stereotype string `gorm:"size:32;not null;default:entity" json:"stereotype"`

	Attributes []DomainAttribute `gorm:"foreignKey:EntityID" json:"attributes,omitempty"`
}

// DomainAttribute はエンティティの属性。
type DomainAttribute struct {
	Base
	EntityID    uint   `gorm:"index;not null" json:"entityId"`
	Name        string `gorm:"size:200;not null" json:"name"`
	DataType    string `gorm:"size:100" json:"dataType"` // string / int / date ...
	Description string `gorm:"type:text" json:"description"`
	Required    bool   `gorm:"default:false" json:"required"`
	IsPrimary   bool   `gorm:"default:false" json:"isPrimary"`
	SortOrder   int    `gorm:"default:0" json:"sortOrder"`
}

// DomainRelationship はエンティティ間の関連。
type DomainRelationship struct {
	Base
	ProjectID      uint   `gorm:"index;not null" json:"projectId"`
	SourceEntityID uint   `gorm:"index;not null" json:"sourceEntityId"`
	TargetEntityID uint   `gorm:"index;not null" json:"targetEntityId"`
	// Type: association / aggregation / composition / inheritance / dependency
	Type              string `gorm:"size:32;not null;default:association" json:"type"`
	Label             string `gorm:"size:200" json:"label"`
	SourceCardinality string `gorm:"size:16" json:"sourceCardinality"` // 例: 1, 0..1, *
	TargetCardinality string `gorm:"size:16" json:"targetCardinality"`

	SourceEntity *DomainEntity `gorm:"foreignKey:SourceEntityID" json:"sourceEntity,omitempty"`
	TargetEntity *DomainEntity `gorm:"foreignKey:TargetEntityID" json:"targetEntity,omitempty"`
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
