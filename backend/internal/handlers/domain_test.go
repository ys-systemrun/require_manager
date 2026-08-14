package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ys-systemrun/require_manager/backend/internal/models"
)

func createEntity(env *testEnv, body map[string]any) models.DomainEntity {
	env.t.Helper()
	var e models.DomainEntity
	rec := env.request(http.MethodPost, "/api/entities", body)
	env.decode(rec, http.StatusCreated, &e)
	return e
}

func TestCreateEntityWithAttributes(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "P")

	e := createEntity(env, map[string]any{
		"projectId":  p.ID,
		"name":       "Customer",
		"stereotype": "aggregate_root",
		"attributes": []map[string]any{
			{"name": "id", "dataType": "uuid", "isPrimary": true, "required": true},
			{"name": "email", "dataType": "string", "required": true},
		},
	})

	var got models.DomainEntity
	rec := env.request(http.MethodGet, fmt.Sprintf("/api/entities/%d", e.ID), nil)
	env.decode(rec, http.StatusOK, &got)
	if len(got.Attributes) != 2 {
		t.Fatalf("attributes = %d, want 2", len(got.Attributes))
	}
	if !got.Attributes[0].IsPrimary || got.Attributes[0].Name != "id" {
		t.Fatalf("first attr unexpected: %+v", got.Attributes[0])
	}
}

func TestUpdateEntityReplacesAttributes(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "P")
	e := createEntity(env, map[string]any{
		"projectId": p.ID, "name": "Order",
		"attributes": []map[string]any{
			{"name": "id", "dataType": "uuid"},
			{"name": "old", "dataType": "string"},
		},
	})

	var updated models.DomainEntity
	rec := env.request(http.MethodPut, fmt.Sprintf("/api/entities/%d", e.ID), map[string]any{
		"name": "Order", "displayName": "注文",
		"attributes": []map[string]any{
			{"name": "id", "dataType": "uuid", "isPrimary": true},
			{"name": "total", "dataType": "money"},
			{"name": "orderedAt", "dataType": "datetime"},
		},
	})
	env.decode(rec, http.StatusOK, &updated)

	var got models.DomainEntity
	rec = env.request(http.MethodGet, fmt.Sprintf("/api/entities/%d", e.ID), nil)
	env.decode(rec, http.StatusOK, &got)
	if got.DisplayName != "注文" {
		t.Fatalf("displayName = %q", got.DisplayName)
	}
	if len(got.Attributes) != 3 {
		t.Fatalf("attributes = %d, want 3 (old set replaced)", len(got.Attributes))
	}
	for _, a := range got.Attributes {
		if a.Name == "old" {
			t.Fatal("stale attribute 'old' still present")
		}
	}
}

func TestRelationshipCRUD(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "P")
	src := createEntity(env, map[string]any{"projectId": p.ID, "name": "Customer"})
	dst := createEntity(env, map[string]any{"projectId": p.ID, "name": "Order"})

	var rel models.DomainRelationship
	rec := env.request(http.MethodPost, "/api/relationships", map[string]any{
		"projectId": p.ID, "sourceEntityId": src.ID, "targetEntityId": dst.ID,
		"type": "association", "label": "発注する", "sourceCardinality": "1", "targetCardinality": "*",
	})
	env.decode(rec, http.StatusCreated, &rel)
	// Preloaded entities should be present in the response.
	if rel.SourceEntity == nil || rel.SourceEntity.Name != "Customer" {
		t.Fatalf("source entity not preloaded: %+v", rel.SourceEntity)
	}

	// List
	var list []models.DomainRelationship
	rec = env.request(http.MethodGet, fmt.Sprintf("/api/relationships?projectId=%d", p.ID), nil)
	env.decode(rec, http.StatusOK, &list)
	if len(list) != 1 {
		t.Fatalf("len = %d, want 1", len(list))
	}

	// Update
	var updated models.DomainRelationship
	rec = env.request(http.MethodPut, fmt.Sprintf("/api/relationships/%d", rel.ID), map[string]any{
		"sourceEntityId": src.ID, "targetEntityId": dst.ID, "type": "composition", "label": "含む",
	})
	env.decode(rec, http.StatusOK, &updated)
	if updated.Type != "composition" {
		t.Fatalf("type = %q, want composition", updated.Type)
	}

	// Delete
	rec = env.request(http.MethodDelete, fmt.Sprintf("/api/relationships/%d", rel.ID), nil)
	mustStatus(t, rec, http.StatusNoContent)
}

func TestDeleteEntityCascadesAttributesAndRelationships(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProject("SHOP", "P")
	a := createEntity(env, map[string]any{
		"projectId": p.ID, "name": "A",
		"attributes": []map[string]any{{"name": "x", "dataType": "int"}},
	})
	b := createEntity(env, map[string]any{"projectId": p.ID, "name": "B"})

	rec := env.request(http.MethodPost, "/api/relationships", map[string]any{
		"projectId": p.ID, "sourceEntityId": a.ID, "targetEntityId": b.ID, "type": "association",
	})
	mustStatus(t, rec, http.StatusCreated)

	// Delete entity A -> its attributes and the relationship should be gone.
	rec = env.request(http.MethodDelete, fmt.Sprintf("/api/entities/%d", a.ID), nil)
	mustStatus(t, rec, http.StatusNoContent)

	var attrCount int64
	env.db.Model(&models.DomainAttribute{}).Where("entity_id = ?", a.ID).Count(&attrCount)
	if attrCount != 0 {
		t.Fatalf("attributes remaining = %d, want 0", attrCount)
	}

	var rels []models.DomainRelationship
	rec = env.request(http.MethodGet, fmt.Sprintf("/api/relationships?projectId=%d", p.ID), nil)
	env.decode(rec, http.StatusOK, &rels)
	if len(rels) != 0 {
		t.Fatalf("relationships remaining = %d, want 0", len(rels))
	}
}

func TestGetEntityNotFound(t *testing.T) {
	env := newTestEnv(t)
	rec := env.request(http.MethodGet, "/api/entities/12345", nil)
	mustStatus(t, rec, http.StatusNotFound)
}
