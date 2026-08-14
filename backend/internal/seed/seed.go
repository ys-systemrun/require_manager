package seed

import (
	"log"

	"github.com/ys-systemrun/require_manager/backend/internal/models"
	"gorm.io/gorm"
)

// Run inserts sample data the first time the app starts against an empty DB.
func Run(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Project{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil // すでにデータあり
	}
	log.Println("seeding sample data ...")

	project := models.Project{
		Name:        "ECサイト リニューアル",
		Key:         "SHOP",
		Description: "既存ECサイトの刷新プロジェクト。要求管理・用語辞書・ドメインモデリングのサンプル。",
	}
	if err := db.Create(&project).Error; err != nil {
		return err
	}

	// 要求 (親)
	parent := models.Requirement{
		ProjectID:   project.ID,
		Code:        "REQ-001",
		Title:       "オンラインで商品を購入できること",
		Description: "顧客が24時間いつでも商品を検索し購入できる。",
		Rationale:   "実店舗の営業時間外の販売機会を獲得するため。",
		Source:      "事業責任者",
		Type:        "business",
		Priority:    "high",
		Status:      "approved",
	}
	if err := db.Create(&parent).Error; err != nil {
		return err
	}

	children := []models.Requirement{
		{
			ProjectID: project.ID, ParentID: &parent.ID,
			Code: "REQ-002", Title: "商品をキーワードで検索できる",
			Description: "商品名・カテゴリ・タグでの検索を提供する。",
			Type:        "functional", Priority: "high", Status: "approved",
		},
		{
			ProjectID: project.ID, ParentID: &parent.ID,
			Code: "REQ-003", Title: "カートに商品を追加・削除できる",
			Description: "購入前に複数商品をカートで管理できる。",
			Type:        "functional", Priority: "high", Status: "implemented",
		},
		{
			ProjectID: project.ID, ParentID: &parent.ID,
			Code: "REQ-004", Title: "ページ表示は2秒以内に完了する",
			Description: "主要ページの応答時間は95パーセンタイルで2秒以内。",
			Rationale:   "離脱率低減のため。", Type: "non_functional", Priority: "medium", Status: "proposed",
		},
		{
			ProjectID: project.ID, ParentID: &parent.ID,
			Code: "REQ-005", Title: "個人情報保護法を遵守すること",
			Description: "顧客の個人情報の取得・保管・利用は関連法令に準拠する。",
			Type:        "constraint", Priority: "high", Status: "approved",
		},
	}
	if err := db.Create(&children).Error; err != nil {
		return err
	}

	// 用語辞書
	terms := []models.Term{
		{ProjectID: project.ID, Name: "カート", Reading: "かーと", EnglishName: "Cart",
			Definition: "顧客が購入予定の商品を一時的に保持する入れ物。", Aliases: "ショッピングカート,買い物かご", Category: "販売"},
		{ProjectID: project.ID, Name: "注文", Reading: "ちゅうもん", EnglishName: "Order",
			Definition: "顧客が商品の購入を確定した記録。1件以上の注文明細を持つ。", Category: "販売"},
		{ProjectID: project.ID, Name: "SKU", Reading: "えすけーゆー", EnglishName: "Stock Keeping Unit",
			Definition: "在庫管理上の最小識別単位。色・サイズ違いなどを区別する。", Category: "在庫"},
	}
	if err := db.Create(&terms).Error; err != nil {
		return err
	}

	// ドメインモデル (エンティティ + 属性)
	customer := models.DomainEntity{
		ProjectID: project.ID, Name: "Customer", DisplayName: "顧客",
		Description: "サイトを利用して商品を購入する人。", Stereotype: "aggregate_root",
		Attributes: []models.DomainAttribute{
			{Name: "id", DataType: "uuid", Required: true, IsPrimary: true, SortOrder: 0},
			{Name: "name", DataType: "string", Description: "氏名", Required: true, SortOrder: 1},
			{Name: "email", DataType: "string", Description: "連絡先メールアドレス", Required: true, SortOrder: 2},
		},
	}
	order := models.DomainEntity{
		ProjectID: project.ID, Name: "Order", DisplayName: "注文",
		Description: "顧客による購入の確定記録。", Stereotype: "aggregate_root",
		Attributes: []models.DomainAttribute{
			{Name: "id", DataType: "uuid", Required: true, IsPrimary: true, SortOrder: 0},
			{Name: "orderedAt", DataType: "datetime", Description: "注文日時", Required: true, SortOrder: 1},
			{Name: "totalAmount", DataType: "money", Description: "合計金額", Required: true, SortOrder: 2},
		},
	}
	orderLine := models.DomainEntity{
		ProjectID: project.ID, Name: "OrderLine", DisplayName: "注文明細",
		Description: "注文に含まれる1商品の数量と価格。", Stereotype: "entity",
		Attributes: []models.DomainAttribute{
			{Name: "productId", DataType: "uuid", Required: true, SortOrder: 0},
			{Name: "quantity", DataType: "int", Description: "数量", Required: true, SortOrder: 1},
			{Name: "unitPrice", DataType: "money", Description: "単価", Required: true, SortOrder: 2},
		},
	}
	if err := db.Create(&[]*models.DomainEntity{&customer, &order, &orderLine}).Error; err != nil {
		return err
	}

	rels := []models.DomainRelationship{
		{ProjectID: project.ID, SourceEntityID: customer.ID, TargetEntityID: order.ID,
			Type: "association", Label: "発注する", SourceCardinality: "1", TargetCardinality: "*"},
		{ProjectID: project.ID, SourceEntityID: order.ID, TargetEntityID: orderLine.ID,
			Type: "composition", Label: "含む", SourceCardinality: "1", TargetCardinality: "1..*"},
	}
	if err := db.Create(&rels).Error; err != nil {
		return err
	}

	log.Println("seeding complete")
	return nil
}
