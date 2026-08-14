package models

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
