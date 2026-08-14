package models

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
