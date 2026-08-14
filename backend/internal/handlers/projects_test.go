package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ys-systemrun/require_manager/backend/internal/models"
)

func TestCreateAndGetProject(t *testing.T) {
	env := newTestEnv(t)

	var created models.Project
	rec := env.request(http.MethodPost, "/api/projects", map[string]any{
		"name":        "ECサイト",
		"key":         "SHOP",
		"description": "テスト用",
	})
	env.decode(rec, http.StatusCreated, &created)

	if created.ID == 0 {
		t.Fatal("expected non-zero ID")
	}
	if created.Key != "SHOP" {
		t.Fatalf("key = %q, want SHOP", created.Key)
	}

	var got models.Project
	rec = env.request(http.MethodGet, fmt.Sprintf("/api/projects/%d", created.ID), nil)
	env.decode(rec, http.StatusOK, &got)
	if got.Name != "ECサイト" {
		t.Fatalf("name = %q", got.Name)
	}
}

func TestListProjects(t *testing.T) {
	env := newTestEnv(t)
	env.seedProject("A", "プロジェクトA")
	env.seedProject("B", "プロジェクトB")

	var list []models.Project
	rec := env.request(http.MethodGet, "/api/projects", nil)
	env.decode(rec, http.StatusOK, &list)
	if len(list) != 2 {
		t.Fatalf("len = %d, want 2", len(list))
	}
}

func TestUpdateProject(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "旧名")

	var updated models.Project
	rec := env.request(http.MethodPut, fmt.Sprintf("/api/projects/%d", p.ID), map[string]any{
		"name":        "新名",
		"key":         "SHOP2",
		"description": "更新",
	})
	env.decode(rec, http.StatusOK, &updated)
	if updated.Name != "新名" || updated.Key != "SHOP2" {
		t.Fatalf("update not applied: %+v", updated)
	}
}

func TestDeleteProject(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "削除対象")

	rec := env.request(http.MethodDelete, fmt.Sprintf("/api/projects/%d", p.ID), nil)
	mustStatus(t, rec, http.StatusNoContent)

	rec = env.request(http.MethodGet, fmt.Sprintf("/api/projects/%d", p.ID), nil)
	mustStatus(t, rec, http.StatusNotFound)
}

func TestGetProjectNotFound(t *testing.T) {
	env := newTestEnv(t)
	rec := env.request(http.MethodGet, "/api/projects/999", nil)
	mustStatus(t, rec, http.StatusNotFound)
}

func TestProjectSummary(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "集計テスト")

	// 2 approved + 1 draft requirement, 1 term.
	for _, st := range []string{"approved", "approved", "draft"} {
		rec := env.request(http.MethodPost, "/api/requirements", map[string]any{
			"projectId": p.ID, "title": "req-" + st, "status": st, "type": "functional", "priority": "medium",
		})
		mustStatus(t, rec, http.StatusCreated)
	}
	rec := env.request(http.MethodPost, "/api/terms", map[string]any{
		"projectId": p.ID, "name": "用語", "definition": "定義",
	})
	mustStatus(t, rec, http.StatusCreated)

	var summary struct {
		Requirements    int64            `json:"requirements"`
		Terms           int64            `json:"terms"`
		StatusBreakdown map[string]int64 `json:"statusBreakdown"`
	}
	rec = env.request(http.MethodGet, fmt.Sprintf("/api/projects/%d/summary", p.ID), nil)
	env.decode(rec, http.StatusOK, &summary)

	if summary.Requirements != 3 {
		t.Fatalf("requirements = %d, want 3", summary.Requirements)
	}
	if summary.Terms != 1 {
		t.Fatalf("terms = %d, want 1", summary.Terms)
	}
	if summary.StatusBreakdown["approved"] != 2 {
		t.Fatalf("approved = %d, want 2", summary.StatusBreakdown["approved"])
	}
}

func TestInvalidProjectID(t *testing.T) {
	env := newTestEnv(t)
	rec := env.request(http.MethodGet, "/api/projects/abc", nil)
	mustStatus(t, rec, http.StatusBadRequest)
}
