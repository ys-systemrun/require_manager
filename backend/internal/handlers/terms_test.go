package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ys-systemrun/require_manager/backend/internal/models"
)

func createTerm(env *testEnv, body map[string]any) models.Term {
	env.t.Helper()
	var term models.Term
	rec := env.request(http.MethodPost, "/api/terms", body)
	env.decode(rec, http.StatusCreated, &term)
	return term
}

func TestCreateAndListTerms(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "P")

	createTerm(env, map[string]any{"projectId": p.ID, "name": "カート", "definition": "買い物かご", "aliases": "ショッピングカート"})
	createTerm(env, map[string]any{"projectId": p.ID, "name": "注文", "definition": "購入の記録"})

	var list []models.Term
	rec := env.request(http.MethodGet, fmt.Sprintf("/api/terms?projectId=%d", p.ID), nil)
	env.decode(rec, http.StatusOK, &list)
	if len(list) != 2 {
		t.Fatalf("len = %d, want 2", len(list))
	}
}

func TestTermKeywordSearch(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "P")
	createTerm(env, map[string]any{"projectId": p.ID, "name": "カート", "definition": "買い物かご", "aliases": "ショッピングカート"})
	createTerm(env, map[string]any{"projectId": p.ID, "name": "注文", "definition": "購入の記録"})
	createTerm(env, map[string]any{"projectId": p.ID, "name": "SKU", "englishName": "Stock Keeping Unit", "definition": "在庫の最小単位"})

	cases := []struct {
		keyword string
		want    int
	}{
		{"カート", 1},    // name match
		{"ショッピング", 1}, // alias match
		{"購入", 1},     // definition match
		{"stock", 1},  // englishName, case-insensitive
		{"存在しない語", 0},
	}
	for _, tc := range cases {
		t.Run(tc.keyword, func(t *testing.T) {
			var list []models.Term
			rec := env.request(http.MethodGet, fmt.Sprintf("/api/terms?projectId=%d&keyword=%s", p.ID, tc.keyword), nil)
			env.decode(rec, http.StatusOK, &list)
			if len(list) != tc.want {
				t.Fatalf("keyword %q: len = %d, want %d", tc.keyword, len(list), tc.want)
			}
		})
	}
}

func TestUpdateTerm(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "P")
	term := createTerm(env, map[string]any{"projectId": p.ID, "name": "旧語", "definition": "旧定義"})

	var updated models.Term
	rec := env.request(http.MethodPut, fmt.Sprintf("/api/terms/%d", term.ID), map[string]any{
		"name": "新語", "definition": "新定義", "category": "販売",
	})
	env.decode(rec, http.StatusOK, &updated)
	if updated.Name != "新語" || updated.Category != "販売" {
		t.Fatalf("update not applied: %+v", updated)
	}
}

func TestDeleteTerm(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "P")
	term := createTerm(env, map[string]any{"projectId": p.ID, "name": "削除語", "definition": "x"})

	rec := env.request(http.MethodDelete, fmt.Sprintf("/api/terms/%d", term.ID), nil)
	mustStatus(t, rec, http.StatusNoContent)

	var list []models.Term
	rec = env.request(http.MethodGet, fmt.Sprintf("/api/terms?projectId=%d", p.ID), nil)
	env.decode(rec, http.StatusOK, &list)
	if len(list) != 0 {
		t.Fatalf("len = %d, want 0", len(list))
	}
}
