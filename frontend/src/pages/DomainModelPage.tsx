import { useEffect, useMemo, useState } from 'react'
import { entitiesApi, relationshipsApi } from '../api'
import { useProjects } from '../ProjectContext'
import { Modal } from '../components/Modal'
import {
  RELATIONSHIP_LABELS,
  RELATIONSHIP_TYPES,
  STEREOTYPE_LABELS,
  STEREOTYPES,
} from '../constants'
import type {
  DomainAttribute,
  DomainEntity,
  DomainRelationship,
  RelationshipType,
  Stereotype,
} from '../types'

const emptyEntity: Partial<DomainEntity> = {
  name: '',
  displayName: '',
  description: '',
  stereotype: 'entity',
  attributes: [],
}

const emptyAttr: DomainAttribute = {
  name: '',
  dataType: 'string',
  description: '',
  required: false,
  isPrimary: false,
}

export function DomainModelPage() {
  const { currentProject } = useProjects()
  const [entities, setEntities] = useState<DomainEntity[]>([])
  const [rels, setRels] = useState<DomainRelationship[]>([])

  const [entityModal, setEntityModal] = useState(false)
  const [editingEntity, setEditingEntity] = useState<DomainEntity | null>(null)
  const [entityForm, setEntityForm] = useState<Partial<DomainEntity>>(emptyEntity)

  const [relModal, setRelModal] = useState(false)
  const [editingRel, setEditingRel] = useState<DomainRelationship | null>(null)
  const [relForm, setRelForm] = useState<Partial<DomainRelationship>>({})

  const load = () => {
    if (!currentProject) return
    entitiesApi.list({ projectId: currentProject.id }).then(setEntities)
    relationshipsApi.list({ projectId: currentProject.id }).then(setRels)
  }

  useEffect(() => {
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentProject])

  // ---- Entity editing ----
  const openCreateEntity = () => {
    setEditingEntity(null)
    setEntityForm({ ...emptyEntity, attributes: [] })
    setEntityModal(true)
  }
  const openEditEntity = (e: DomainEntity) => {
    setEditingEntity(e)
    setEntityForm({ ...e, attributes: e.attributes.map((a) => ({ ...a })) })
    setEntityModal(true)
  }
  const saveEntity = async () => {
    if (!currentProject) return
    const payload = { ...entityForm, projectId: currentProject.id }
    if (editingEntity) await entitiesApi.update(editingEntity.id, payload)
    else await entitiesApi.create(payload)
    setEntityModal(false)
    load()
  }
  const removeEntity = async (e: DomainEntity) => {
    if (!confirm(`エンティティ「${e.name}」を削除しますか？関連も削除されます。`)) return
    await entitiesApi.remove(e.id)
    load()
  }

  const updateAttr = (idx: number, patch: Partial<DomainAttribute>) => {
    const attrs = [...(entityForm.attributes ?? [])]
    attrs[idx] = { ...attrs[idx], ...patch }
    setEntityForm({ ...entityForm, attributes: attrs })
  }
  const addAttr = () =>
    setEntityForm({
      ...entityForm,
      attributes: [...(entityForm.attributes ?? []), { ...emptyAttr }],
    })
  const removeAttr = (idx: number) =>
    setEntityForm({
      ...entityForm,
      attributes: (entityForm.attributes ?? []).filter((_, i) => i !== idx),
    })

  // ---- Relationship editing ----
  const openCreateRel = () => {
    setEditingRel(null)
    setRelForm({
      type: 'association',
      sourceEntityId: entities[0]?.id,
      targetEntityId: entities[1]?.id ?? entities[0]?.id,
      sourceCardinality: '1',
      targetCardinality: '*',
      label: '',
    })
    setRelModal(true)
  }
  const openEditRel = (r: DomainRelationship) => {
    setEditingRel(r)
    setRelForm({ ...r })
    setRelModal(true)
  }
  const saveRel = async () => {
    if (!currentProject) return
    const payload = { ...relForm, projectId: currentProject.id }
    if (editingRel) await relationshipsApi.update(editingRel.id, payload)
    else await relationshipsApi.create(payload)
    setRelModal(false)
    load()
  }
  const removeRel = async (r: DomainRelationship) => {
    if (!confirm('この関連を削除しますか？')) return
    await relationshipsApi.remove(r.id)
    load()
  }

  const entityName = (id: number) => entities.find((e) => e.id === id)?.name ?? `#${id}`

  // Mermaid classDiagram のテキストを生成（コピーして図として利用可能）
  const mermaid = useMemo(() => {
    const lines = ['classDiagram']
    for (const e of entities) {
      lines.push(`  class ${e.name} {`)
      for (const a of e.attributes) {
        lines.push(`    ${a.dataType || 'any'} ${a.name}${a.isPrimary ? ' 🔑' : ''}`)
      }
      lines.push('  }')
    }
    const arrow: Record<RelationshipType, string> = {
      association: '-->',
      aggregation: 'o--',
      composition: '*--',
      inheritance: '--|>',
      dependency: '..>',
    }
    for (const r of rels) {
      const s = entityName(r.sourceEntityId)
      const t = entityName(r.targetEntityId)
      const lbl = r.label ? ` : ${r.label}` : ''
      lines.push(`  ${s} ${arrow[r.type]} ${t}${lbl}`)
    }
    return lines.join('\n')
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [entities, rels])

  if (!currentProject) return null

  return (
    <div>
      <header className="page-header">
        <div>
          <h1>ドメインモデル</h1>
          <p className="muted">エンティティ・属性・関連によるドメインモデリング</p>
        </div>
        <div className="header-actions">
          <button className="btn" onClick={openCreateRel} disabled={entities.length < 1}>
            ＋ 関連を追加
          </button>
          <button className="btn btn-primary" onClick={openCreateEntity}>
            ＋ エンティティを追加
          </button>
        </div>
      </header>

      {entities.length === 0 ? (
        <div className="empty-state">エンティティがありません。</div>
      ) : (
        <div className="entity-grid">
          {entities.map((e) => (
            <div key={e.id} className="uml-class">
              <div className="uml-header">
                <div>
                  <div className="uml-stereo">«{STEREOTYPE_LABELS[e.stereotype]}»</div>
                  <div className="uml-name">
                    {e.name}
                    {e.displayName && <span className="uml-jp"> / {e.displayName}</span>}
                  </div>
                </div>
                <div className="actions">
                  <button className="icon-btn" onClick={() => openEditEntity(e)}>
                    ✏️
                  </button>
                  <button className="icon-btn" onClick={() => removeEntity(e)}>
                    🗑️
                  </button>
                </div>
              </div>
              <div className="uml-attrs">
                {e.attributes.length === 0 && <div className="muted small">（属性なし）</div>}
                {e.attributes.map((a) => (
                  <div key={a.id ?? a.name} className="uml-attr">
                    <span>
                      {a.isPrimary && <span className="pk">PK </span>}
                      {a.name}
                      {a.required && <span className="req-star">*</span>}
                    </span>
                    <span className="mono muted">{a.dataType}</span>
                  </div>
                ))}
              </div>
              {e.description && <div className="uml-desc">{e.description}</div>}
            </div>
          ))}
        </div>
      )}

      <section className="panel">
        <div className="panel-head">
          <h2>関連</h2>
        </div>
        {rels.length === 0 ? (
          <p className="muted">関連が定義されていません。</p>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>元</th>
                <th style={{ width: 130 }}>種類</th>
                <th>先</th>
                <th>ラベル</th>
                <th style={{ width: 90 }}></th>
              </tr>
            </thead>
            <tbody>
              {rels.map((r) => (
                <tr key={r.id}>
                  <td className="mono">
                    {entityName(r.sourceEntityId)}
                    {r.sourceCardinality && (
                      <span className="card-mult"> [{r.sourceCardinality}]</span>
                    )}
                  </td>
                  <td>
                    <span className="badge">{RELATIONSHIP_LABELS[r.type]}</span>
                  </td>
                  <td className="mono">
                    {entityName(r.targetEntityId)}
                    {r.targetCardinality && (
                      <span className="card-mult"> [{r.targetCardinality}]</span>
                    )}
                  </td>
                  <td>{r.label}</td>
                  <td className="actions">
                    <button className="icon-btn" onClick={() => openEditRel(r)}>
                      ✏️
                    </button>
                    <button className="icon-btn" onClick={() => removeRel(r)}>
                      🗑️
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>

      {entities.length > 0 && (
        <section className="panel">
          <div className="panel-head">
            <h2>Mermaid クラス図</h2>
            <button
              className="btn small"
              onClick={() => navigator.clipboard?.writeText(mermaid)}
            >
              コピー
            </button>
          </div>
          <pre className="code-block">{mermaid}</pre>
        </section>
      )}

      {/* Entity modal */}
      <Modal
        open={entityModal}
        title={editingEntity ? 'エンティティを編集' : 'エンティティを追加'}
        onClose={() => setEntityModal(false)}
        footer={
          <>
            <button className="btn" onClick={() => setEntityModal(false)}>
              キャンセル
            </button>
            <button
              className="btn btn-primary"
              onClick={saveEntity}
              disabled={!entityForm.name}
            >
              保存
            </button>
          </>
        }
      >
        <div className="form-grid">
          <label className="field">
            <span>名前 *（英語）</span>
            <input
              value={entityForm.name ?? ''}
              placeholder="Customer"
              onChange={(e) => setEntityForm({ ...entityForm, name: e.target.value })}
            />
          </label>
          <label className="field">
            <span>表示名（日本語）</span>
            <input
              value={entityForm.displayName ?? ''}
              placeholder="顧客"
              onChange={(e) => setEntityForm({ ...entityForm, displayName: e.target.value })}
            />
          </label>
          <label className="field">
            <span>ステレオタイプ</span>
            <select
              value={entityForm.stereotype}
              onChange={(e) =>
                setEntityForm({ ...entityForm, stereotype: e.target.value as Stereotype })
              }
            >
              {STEREOTYPES.map((s) => (
                <option key={s} value={s}>
                  {STEREOTYPE_LABELS[s]}
                </option>
              ))}
            </select>
          </label>
          <label className="field span-2">
            <span>説明</span>
            <textarea
              rows={2}
              value={entityForm.description ?? ''}
              onChange={(e) => setEntityForm({ ...entityForm, description: e.target.value })}
            />
          </label>
        </div>

        <div className="attr-editor">
          <div className="attr-editor-head">
            <span>属性</span>
            <button className="btn small" onClick={addAttr}>
              ＋ 属性
            </button>
          </div>
          {(entityForm.attributes ?? []).map((a, idx) => (
            <div key={idx} className="attr-row">
              <input
                placeholder="属性名"
                value={a.name}
                onChange={(e) => updateAttr(idx, { name: e.target.value })}
              />
              <input
                placeholder="型"
                value={a.dataType}
                onChange={(e) => updateAttr(idx, { dataType: e.target.value })}
              />
              <label className="check">
                <input
                  type="checkbox"
                  checked={a.isPrimary}
                  onChange={(e) => updateAttr(idx, { isPrimary: e.target.checked })}
                />
                PK
              </label>
              <label className="check">
                <input
                  type="checkbox"
                  checked={a.required}
                  onChange={(e) => updateAttr(idx, { required: e.target.checked })}
                />
                必須
              </label>
              <button className="icon-btn" onClick={() => removeAttr(idx)}>
                🗑️
              </button>
            </div>
          ))}
        </div>
      </Modal>

      {/* Relationship modal */}
      <Modal
        open={relModal}
        title={editingRel ? '関連を編集' : '関連を追加'}
        onClose={() => setRelModal(false)}
        footer={
          <>
            <button className="btn" onClick={() => setRelModal(false)}>
              キャンセル
            </button>
            <button
              className="btn btn-primary"
              onClick={saveRel}
              disabled={!relForm.sourceEntityId || !relForm.targetEntityId}
            >
              保存
            </button>
          </>
        }
      >
        <div className="form-grid">
          <label className="field">
            <span>元エンティティ</span>
            <select
              value={relForm.sourceEntityId ?? ''}
              onChange={(e) =>
                setRelForm({ ...relForm, sourceEntityId: Number(e.target.value) })
              }
            >
              {entities.map((e) => (
                <option key={e.id} value={e.id}>
                  {e.name}
                </option>
              ))}
            </select>
          </label>
          <label className="field">
            <span>先エンティティ</span>
            <select
              value={relForm.targetEntityId ?? ''}
              onChange={(e) =>
                setRelForm({ ...relForm, targetEntityId: Number(e.target.value) })
              }
            >
              {entities.map((e) => (
                <option key={e.id} value={e.id}>
                  {e.name}
                </option>
              ))}
            </select>
          </label>
          <label className="field">
            <span>種類</span>
            <select
              value={relForm.type}
              onChange={(e) =>
                setRelForm({ ...relForm, type: e.target.value as RelationshipType })
              }
            >
              {RELATIONSHIP_TYPES.map((t) => (
                <option key={t} value={t}>
                  {RELATIONSHIP_LABELS[t]}
                </option>
              ))}
            </select>
          </label>
          <label className="field">
            <span>ラベル</span>
            <input
              value={relForm.label ?? ''}
              placeholder="発注する 等"
              onChange={(e) => setRelForm({ ...relForm, label: e.target.value })}
            />
          </label>
          <label className="field">
            <span>元の多重度</span>
            <input
              value={relForm.sourceCardinality ?? ''}
              placeholder="1"
              onChange={(e) => setRelForm({ ...relForm, sourceCardinality: e.target.value })}
            />
          </label>
          <label className="field">
            <span>先の多重度</span>
            <input
              value={relForm.targetCardinality ?? ''}
              placeholder="*"
              onChange={(e) => setRelForm({ ...relForm, targetCardinality: e.target.value })}
            />
          </label>
        </div>
      </Modal>
    </div>
  )
}
