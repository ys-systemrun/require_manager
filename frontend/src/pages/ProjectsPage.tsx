import { useState } from 'react'
import { projectsApi } from '../api'
import { useProjects } from '../ProjectContext'
import { Modal } from '../components/Modal'
import type { Project } from '../types'

const emptyForm: Partial<Project> = { name: '', key: '', description: '' }

export function ProjectsPage() {
  const { projects, reloadProjects, setCurrentProjectId, currentProject } = useProjects()
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Project | null>(null)
  const [form, setForm] = useState<Partial<Project>>(emptyForm)

  const openCreate = () => {
    setEditing(null)
    setForm(emptyForm)
    setModalOpen(true)
  }
  const openEdit = (p: Project) => {
    setEditing(p)
    setForm({ ...p })
    setModalOpen(true)
  }
  const save = async () => {
    if (editing) {
      await projectsApi.update(editing.id, form)
    } else {
      const created = await projectsApi.create(form)
      setCurrentProjectId(created.id)
    }
    setModalOpen(false)
    await reloadProjects()
  }
  const remove = async (p: Project) => {
    if (!confirm(`プロジェクト「${p.name}」を削除しますか？関連データもすべて削除されます。`))
      return
    await projectsApi.remove(p.id)
    await reloadProjects()
  }

  return (
    <div>
      <header className="page-header">
        <div>
          <h1>プロジェクト設定</h1>
          <p className="muted">プロジェクトの作成・編集・切り替え</p>
        </div>
        <button className="btn btn-primary" onClick={openCreate}>
          ＋ プロジェクトを作成
        </button>
      </header>

      <div className="panel no-pad">
        <table className="data-table">
          <thead>
            <tr>
              <th style={{ width: 100 }}>キー</th>
              <th>名前</th>
              <th>説明</th>
              <th style={{ width: 140 }}></th>
            </tr>
          </thead>
          <tbody>
            {projects.length === 0 && (
              <tr>
                <td colSpan={4} className="muted center">
                  プロジェクトがありません。
                </td>
              </tr>
            )}
            {projects.map((p) => (
              <tr key={p.id} className={currentProject?.id === p.id ? 'row-active' : ''}>
                <td>
                  <span className="chip small">{p.key}</span>
                </td>
                <td>
                  <button className="link" onClick={() => setCurrentProjectId(p.id)}>
                    {p.name}
                  </button>
                  {currentProject?.id === p.id && <span className="tag active-tag">選択中</span>}
                </td>
                <td className="muted">{p.description}</td>
                <td className="actions">
                  <button className="icon-btn" onClick={() => openEdit(p)}>
                    ✏️
                  </button>
                  <button className="icon-btn" onClick={() => remove(p)}>
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
        title={editing ? 'プロジェクトを編集' : 'プロジェクトを作成'}
        onClose={() => setModalOpen(false)}
        footer={
          <>
            <button className="btn" onClick={() => setModalOpen(false)}>
              キャンセル
            </button>
            <button
              className="btn btn-primary"
              onClick={save}
              disabled={!form.name || !form.key}
            >
              保存
            </button>
          </>
        }
      >
        <div className="form-grid">
          <label className="field">
            <span>キー *</span>
            <input
              value={form.key ?? ''}
              placeholder="SHOP"
              onChange={(e) => setForm({ ...form, key: e.target.value.toUpperCase() })}
            />
          </label>
          <label className="field">
            <span>名前 *</span>
            <input
              value={form.name ?? ''}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
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
        </div>
      </Modal>
    </div>
  )
}
