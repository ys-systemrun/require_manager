package models

// Requirement (要求・要件) は要求管理の中核。親子で階層化できる。
type Requirement struct {
	Base
	ProjectID uint  `gorm:"index;not null" json:"projectId"`
	ParentID  *uint `gorm:"index" json:"parentId"`

	Code        string `gorm:"size:64;index" json:"code"` // 例: REQ-001
	Title       string `gorm:"size:300;not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Rationale   string `gorm:"type:text" json:"rationale"` // 根拠・理由
	Source      string `gorm:"size:300" json:"source"`     // 出典 (ステークホルダー等)

	// Type: business(業務要求) / functional(機能要件) / non_functional(非機能要件) / constraint(制約)
	Type string `gorm:"size:32;not null;default:functional" json:"type"`
	// Priority: high / medium / low
	Priority string `gorm:"size:16;not null;default:medium" json:"priority"`
	// Status: draft / proposed / approved / implemented / verified / rejected / deprecated
	Status string `gorm:"size:16;not null;default:draft" json:"status"`

	Children []Requirement `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}
