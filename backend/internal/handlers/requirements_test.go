package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ys-systemrun/require_manager/backend/internal/models"
)

func createRequirement(env *testEnv, body map[string]any) models.Requirement {
	env.t.Helper()
	var r models.Requirement
	rec := env.request(http.MethodPost, "/api/requirements", body)
	env.decode(rec, http.StatusCreated, &r)
	return r
}

func TestCreateRequirementDefaults(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "P")

	r := createRequirement(env, map[string]any{
		"projectId": p.ID,
		"title":     "検索できること",
	})
	if r.Type != "functional" || r.Priority != "medium" || r.Status != "draft" {
		t.Fatalf("defaults not applied: %+v", r)
	}
}

func TestListRequirementsFilters(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "P")

	createRequirement(env, map[string]any{"projectId": p.ID, "code": "REQ-001", "title": "カート機能", "type": "functional", "status": "approved", "priority": "high"})
	createRequirement(env, map[string]any{"projectId": p.ID, "code": "REQ-002", "title": "性能要件", "type": "non_functional", "status": "draft", "priority": "medium"})
	createRequirement(env, map[string]any{"projectId": p.ID, "code": "REQ-003", "title": "決済機能", "type": "functional", "status": "draft", "priority": "high"})

	cases := []struct {
		name  string
		query string
		want  int
	}{
		{"all", "", 3},
		{"status approved", "&status=approved", 1},
		{"type functional", "&type=functional", 2},
		{"keyword 機能", "&keyword=機能", 2},
		{"keyword no match", "&keyword=xyz", 0},
		{"combined", "&type=functional&status=draft", 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var list []models.Requirement
			rec := env.request(http.MethodGet, fmt.Sprintf("/api/requirements?projectId=%d%s", p.ID, tc.query), nil)
			env.decode(rec, http.StatusOK, &list)
			if len(list) != tc.want {
				t.Fatalf("len = %d, want %d", len(list), tc.want)
			}
		})
	}
}

func TestUpdateRequirement(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "P")
	r := createRequirement(env, map[string]any{"projectId": p.ID, "title": "旧", "priority": "low", "status": "draft"})

	var updated models.Requirement
	rec := env.request(http.MethodPut, fmt.Sprintf("/api/requirements/%d", r.ID), map[string]any{
		"title": "新", "priority": "high", "status": "approved", "type": "business",
	})
	env.decode(rec, http.StatusOK, &updated)
	if updated.Title != "新" || updated.Priority != "high" || updated.Status != "approved" {
		t.Fatalf("update not applied: %+v", updated)
	}
	// projectId must be preserved even though not in the payload.
	if updated.ProjectID != p.ID {
		t.Fatalf("projectId = %d, want %d", updated.ProjectID, p.ID)
	}
}

func TestRequirementHierarchyAndPreload(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "P")
	parent := createRequirement(env, map[string]any{"projectId": p.ID, "title": "親要求"})
	createRequirement(env, map[string]any{"projectId": p.ID, "title": "子要求", "parentId": parent.ID})

	var got models.Requirement
	rec := env.request(http.MethodGet, fmt.Sprintf("/api/requirements/%d", parent.ID), nil)
	env.decode(rec, http.StatusOK, &got)
	if len(got.Children) != 1 {
		t.Fatalf("children = %d, want 1", len(got.Children))
	}
	if got.Children[0].Title != "子要求" {
		t.Fatalf("child title = %q", got.Children[0].Title)
	}
}

func TestDeleteRequirementDetachesChildren(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "P")
	parent := createRequirement(env, map[string]any{"projectId": p.ID, "title": "親"})
	child := createRequirement(env, map[string]any{"projectId": p.ID, "title": "子", "parentId": parent.ID})

	rec := env.request(http.MethodDelete, fmt.Sprintf("/api/requirements/%d", parent.ID), nil)
	mustStatus(t, rec, http.StatusNoContent)

	// Child survives with its parent link cleared.
	var got models.Requirement
	rec = env.request(http.MethodGet, fmt.Sprintf("/api/requirements/%d", child.ID), nil)
	env.decode(rec, http.StatusOK, &got)
	if got.ParentID != nil {
		t.Fatalf("parentId = %v, want nil", *got.ParentID)
	}
}

func TestCreateRequirementBadJSON(t *testing.T) {
	env := newTestEnv(t)
	rec := env.requestRaw(http.MethodPost, "/api/requirements", "{not json}")
	mustStatus(t, rec, http.StatusBadRequest)
}
