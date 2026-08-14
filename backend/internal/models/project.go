package models

// Project (プロジェクト) groups all requirements, terms and domain models.
type Project struct {
	Base
	Name        string `gorm:"size:200;not null" json:"name"`
	Key         string `gorm:"size:32;uniqueIndex;not null" json:"key"` // 例: SHOP
	Description string `gorm:"type:text" json:"description"`

	Requirements  []Requirement        `json:"requirements,omitempty"`
	Terms         []Term               `json:"terms,omitempty"`
	Entities      []DomainEntity       `json:"entities,omitempty"`
	Relationships []DomainRelationship `json:"relationships,omitempty"`
}
