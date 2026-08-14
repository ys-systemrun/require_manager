package models

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
