import { useEffect, useMemo, useState } from 'react'
import { requirementsApi } from '../api'
import { useProjects } from '../ProjectContext'
import { Modal } from '../components/Modal'
import {
  PRIORITIES,
  PRIORITY_LABELS,
  REQ_TYPES,
  REQ_TYPE_LABELS,
  STATUSES,
  STATUS_LABELS,
} from '../constants'
import type { Priority, Requirement, RequirementStatus, RequirementType } from '../types'

const emptyForm: Partial<Requirement> = {
  code: '',
  title: '',
  description: '',
  rationale: '',
  source: '',
  type: 'functional',
  priority: 'medium',
  status: 'draft',
  parentId: null,
}

export function RequirementsPage() {
  const { currentProject } = useProjects()
  const [items, setItems] = useState<Requirement[]>([])
  const [keyword, setKeyword] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [typeFilter, setTypeFilter] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Requirement | null>(null)
  const [form, setForm] = useState<Partial<Requirement>>(emptyForm)

  const load = () => {
    if (!currentProject) return
    requirementsApi
      .list({
        projectId: currentProject.id,
        keyword: keyword || undefined,
        status: statusFilter || undefined,
        type: typeFilter || undefined,
      })
      .then(setItems)
  }

  useEffect(() => {
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentProject, keyword, statusFilter, typeFilter])

  const parentOptions = useMemo(
    () => items.filter((r) => r.id !== editing?.id),
    [items, editing],
  )

  const openCreate = () => {
    setEditing(null)
    setForm(emptyForm)
    setModalOpen(true)
  }

  const openEdit = (r: Requirement) => {
    setEditing(r)
    setForm({ ...r })
    setModalOpen(true)
  }

  const save = async () => {
    if (!currentProject) return
    const payload = { ...form, projectId: currentProject.id }
    if (editing) {
      await requirementsApi.update(editing.id, payload)
    } else {
      await requirementsApi.create(payload)
    }
    setModalOpen(false)
    load()
  }

  const remove = async (r: Requirement) => {
    if (!confirm(`「${r.title}」を削除しますか？`)) return
    await requirementsApi.remove(r.id)
    load()
  }

  // 親→子の順で並べ替え（簡易ツリー表示）
  const ordered = useMemo(() => {
    const byParent = new Map<number | null, Requirement[]>()
    for (const r of items) {
      const key = r.parentId ?? null
      if (!byParent.has(key)) byParent.set(key, [])
      byParent.get(key)!.push(r)
    }
    const result: { req: Requirement; depth: number }[] = []
    const walk = (parentId: number | null, depth: number) => {
      for (const r of byParent.get(parentId) ?? []) {
        result.push({ req: r, depth })
        walk(r.id, depth + 1)
      }
    }
    walk(null, 0)
    // 親が絞り込みで消えている子も表示（孤立分）
    const shown = new Set(result.map((x) => x.req.id))
    for (const r of items) if (!shown.has(r.id)) result.push({ req: r, depth: 0 })
    return result
  }, [items])

  if (!currentProject) return null

  return (
    <div>
      <header className="page-header">
        <div>
          <h1>要求管理</h1>
          <p className="muted">要件定義・要求の登録とトレーサビリティ管理</p>
        </div>
        <button className="btn btn-primary" onClick={openCreate}>
          ＋ 要求を追加
        </button>
      </header>

      <div className="toolbar">
        <input
          className="search"
          placeholder="キーワード検索..."
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
        />
        <select value={typeFilter} onChange={(e) => setTypeFilter(e.target.value)}>
          <option value="">種別: すべて</option>
          {REQ_TYPES.map((t) => (
            <option key={t} value={t}>
              {REQ_TYPE_LABELS[t]}
            </option>
          ))}
        </select>
        <select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)}>
          <option value="">状態: すべて</option>
          {STATUSES.map((s) => (
            <option key={s} value={s}>
              {STATUS_LABELS[s]}
            </option>
          ))}
        </select>
      </div>

      <div className="panel no-pad">
        <table className="data-table">
          <thead>
            <tr>
              <th style={{ width: 110 }}>ID</th>
              <th>タイトル</th>
              <th style={{ width: 110 }}>種別</th>
              <th style={{ width: 70 }}>優先度</th>
              <th style={{ width: 90 }}>状態</th>
              <th style={{ width: 110 }}></th>
            </tr>
          </thead>
          <tbody>
            {ordered.length === 0 && (
              <tr>
                <td colSpan={6} className="muted center">
                  要求がありません。
                </td>
              </tr>
            )}
            {ordered.map(({ req, depth }) => (
              <tr key={req.id}>
                <td className="mono">{req.code || `#${req.id}`}</td>
                <td>
                  <div style={{ paddingLeft: depth * 20 }}>
                    {depth > 0 && <span className="tree-branch">└ </span>}
                    <button className="link" onClick={() => openEdit(req)}>
                      {req.title}
                    </button>
                    {req.description && (
                      <div className="cell-sub">{req.description}</div>
                    )}
                  </div>
                </td>
                <td>
                  <span className={`badge type-${req.type}`}>
                    {REQ_TYPE_LABELS[req.type]}
                  </span>
                </td>
                <td>
                  <span className={`badge prio-${req.priority}`}>
                    {PRIORITY_LABELS[req.priority]}
                  </span>
                </td>
                <td>
                  <span className={`badge badge-${req.status}`}>
                    {STATUS_LABELS[req.status]}
                  </span>
                </td>
                <td className="actions">
                  <button className="icon-btn" onClick={() => openEdit(req)}>
                    ✏️
                  </button>
                  <button className="icon-btn" onClick={() => remove(req)}>
                    🗑️
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Modal
        open={modalOpen}
        title={editing ? '要求を編集' : '要求を追加'}
        onClose={() => setModalOpen(false)}
        footer={
          <>
            <button className="btn" onClick={() => setModalOpen(false)}>
              キャンセル
            </button>
            <button className="btn btn-primary" onClick={save} disabled={!form.title}>
              保存
            </button>
          </>
        }
      >
        <div className="form-grid">
          <label className="field">
            <span>ID / コード</span>
            <input
              value={form.code ?? ''}
              placeholder="REQ-001"
              onChange={(e) => setForm({ ...form, code: e.target.value })}
            />
          </label>
          <label className="field">
            <span>親要求</span>
            <select
              value={form.parentId ?? ''}
              onChange={(e) =>
                setForm({ ...form, parentId: e.target.value ? Number(e.target.value) : null })
              }
            >
              <option value="">（なし）</option>
              {parentOptions.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.code ? `${p.code} ` : ''}
                  {p.title}
                </option>
              ))}
            </select>
          </label>
          <label className="field span-2">
            <span>タイトル *</span>
            <input
              value={form.title ?? ''}
              onChange={(e) => setForm({ ...form, title: e.target.value })}
            />
          </label>
          <label className="field span-2">
            <span>説明</span>
            <textarea
              rows={3}
              value={form.description ?? ''}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
            />
          </label>
          <label className="field">
            <span>種別</span>
            <select
              value={form.type}
              onChange={(e) => setForm({ ...form, type: e.target.value as RequirementType })}
            >
              {REQ_TYPES.map((t) => (
                <option key={t} value={t}>
                  {REQ_TYPE_LABELS[t]}
                </option>
              ))}
            </select>
          </label>
          <label className="field">
            <span>優先度</span>
            <select
              value={form.priority}
              onChange={(e) => setForm({ ...form, priority: e.target.value as Priority })}
            >
              {PRIORITIES.map((p) => (
                <option key={p} value={p}>
                  {PRIORITY_LABELS[p]}
                </option>
              ))}
            </select>
          </label>
          <label className="field">
            <span>状態</span>
            <select
              value={form.status}
              onChange={(e) =>
                setForm({ ...form, status: e.target.value as RequirementStatus })
              }
            >
              {STATUSES.map((s) => (
                <option key={s} value={s}>
                  {STATUS_LABELS[s]}
                </option>
              ))}
            </select>
          </label>
          <label className="field">
            <span>出典</span>
            <input
              value={form.source ?? ''}
              placeholder="ステークホルダー等"
              onChange={(e) => setForm({ ...form, source: e.target.value })}
            />
          </label>
          <label className="field span-2">
            <span>根拠・理由</span>
            <textarea
              rows={2}
              value={form.rationale ?? ''}
              onChange={(e) => setForm({ ...form, rationale: e.target.value })}
            />
          </label>
        </div>
      </Modal>
    </div>
  )
}
