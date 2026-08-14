import { useEffect, useState } from 'react'
import { termsApi } from '../api'
import { useProjects } from '../ProjectContext'
import { Modal } from '../components/Modal'
import type { Term } from '../types'

const emptyForm: Partial<Term> = {
  name: '',
  reading: '',
  englishName: '',
  definition: '',
  aliases: '',
  category: '',
  notes: '',
}

export function TermsPage() {
  const { currentProject } = useProjects()
  const [items, setItems] = useState<Term[]>([])
  const [keyword, setKeyword] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Term | null>(null)
  const [form, setForm] = useState<Partial<Term>>(emptyForm)

  const load = () => {
    if (!currentProject) return
    termsApi
      .list({ projectId: currentProject.id, keyword: keyword || undefined })
      .then(setItems)
  }

  useEffect(() => {
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentProject, keyword])

  const openCreate = () => {
    setEditing(null)
    setForm(emptyForm)
    setModalOpen(true)
  }

  const openEdit = (t: Term) => {
    setEditing(t)
    setForm({ ...t })
    setModalOpen(true)
  }

  const save = async () => {
    if (!currentProject) return
    const payload = { ...form, projectId: currentProject.id }
    if (editing) await termsApi.update(editing.id, payload)
    else await termsApi.create(payload)
    setModalOpen(false)
    load()
  }

  const remove = async (t: Term) => {
    if (!confirm(`用語「${t.name}」を削除しますか？`)) return
    await termsApi.remove(t.id)
    load()
  }

  if (!currentProject) return null

  return (
    <div>
      <header className="page-header">
        <div>
          <h1>用語辞書</h1>
          <p className="muted">プロジェクト内で共通利用する用語の定義</p>
        </div>
        <button className="btn btn-primary" onClick={openCreate}>
          ＋ 用語を追加
        </button>
      </header>

      <div className="toolbar">
        <input
          className="search"
          placeholder="用語・定義を検索..."
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
        />
      </div>

      {items.length === 0 ? (
        <div className="empty-state">用語がありません。</div>
      ) : (
        <div className="term-grid">
          {items.map((t) => (
            <div key={t.id} className="term-card">
              <div className="term-card-head">
                <div>
                  <h3>{t.name}</h3>
                  {(t.reading || t.englishName) && (
                    <div className="term-meta">
                      {t.reading && <span>{t.reading}</span>}
                      {t.englishName && <span className="mono">{t.englishName}</span>}
                    </div>
                  )}
                </div>
                <div className="actions">
                  <button className="icon-btn" onClick={() => openEdit(t)}>
                    ✏️
                  </button>
                  <button className="icon-btn" onClick={() => remove(t)}>
                    🗑️
                  </button>
                </div>
              </div>
              <p className="term-def">{t.definition}</p>
              {t.aliases && (
                <div className="term-aliases">
                  {t.aliases.split(',').map(
                    (a) =>
                      a.trim() && (
                        <span key={a} className="tag">
                          {a.trim()}
                        </span>
                      ),
                  )}
                </div>
              )}
              {t.category && <span className="chip small">{t.category}</span>}
            </div>
          ))}
        </div>
      )}

      <Modal
        open={modalOpen}
        title={editing ? '用語を編集' : '用語を追加'}
        onClose={() => setModalOpen(false)}
        footer={
          <>
            <button className="btn" onClick={() => setModalOpen(false)}>
              キャンセル
            </button>
            <button
              className="btn btn-primary"
              onClick={save}
              disabled={!form.name || !form.definition}
            >
              保存
            </button>
          </>
        }
      >
        <div className="form-grid">
          <label className="field">
            <span>用語 *</span>
            <input
              value={form.name ?? ''}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
            />
          </label>
          <label className="field">
            <span>読み仮名</span>
            <input
              value={form.reading ?? ''}
              onChange={(e) => setForm({ ...form, reading: e.target.value })}
            />
          </label>
          <label className="field">
            <span>英語表記</span>
            <input
              value={form.englishName ?? ''}
              onChange={(e) => setForm({ ...form, englishName: e.target.value })}
            />
          </label>
          <label className="field">
            <span>分類</span>
            <input
              value={form.category ?? ''}
              onChange={(e) => setForm({ ...form, category: e.target.value })}
            />
          </label>
          <label className="field span-2">
            <span>定義 *</span>
            <textarea
              rows={3}
              value={form.definition ?? ''}
              onChange={(e) => setForm({ ...form, definition: e.target.value })}
            />
          </label>
          <label className="field span-2">
            <span>別名・同義語（カンマ区切り）</span>
            <input
              value={form.aliases ?? ''}
              placeholder="別名1, 別名2"
              onChange={(e) => setForm({ ...form, aliases: e.target.value })}
            />
          </label>
          <label className="field span-2">
            <span>備考</span>
            <textarea
              rows={2}
              value={form.notes ?? ''}
              onChange={(e) => setForm({ ...form, notes: e.target.value })}
            />
          </label>
        </div>
      </Modal>
    </div>
  )
}
