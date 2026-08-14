import { useEffect, useState } from 'react'
import { projectsApi } from '../api'
import { useProjects } from '../ProjectContext'
import { STATUS_LABELS } from '../constants'
import type { ProjectSummary, RequirementStatus } from '../types'

export function DashboardPage() {
  const { currentProject } = useProjects()
  const [summary, setSummary] = useState<ProjectSummary | null>(null)

  useEffect(() => {
    if (!currentProject) return
    projectsApi.summary(currentProject.id).then(setSummary)
  }, [currentProject])

  if (!currentProject) return null

  const cards = [
    { label: '要求', value: summary?.requirements ?? 0, icon: '📋', color: 'blue' },
    { label: '用語', value: summary?.terms ?? 0, icon: '📖', color: 'green' },
    { label: 'エンティティ', value: summary?.entities ?? 0, icon: '🧩', color: 'purple' },
    { label: '関連', value: summary?.relationships ?? 0, icon: '🔗', color: 'orange' },
  ]

  const breakdown = summary?.statusBreakdown ?? {}
  const totalReq = summary?.requirements ?? 0

  return (
    <div>
      <header className="page-header">
        <div>
          <h1>{currentProject.name}</h1>
          <p className="muted">{currentProject.description || 'プロジェクト概要'}</p>
        </div>
        <span className="chip">{currentProject.key}</span>
      </header>

      <div className="stat-grid">
        {cards.map((c) => (
          <div key={c.label} className={`stat-card stat-${c.color}`}>
            <span className="stat-icon">{c.icon}</span>
            <div>
              <div className="stat-value">{c.value}</div>
              <div className="stat-label">{c.label}</div>
            </div>
          </div>
        ))}
      </div>

      <section className="panel">
        <h2>要求ステータスの内訳</h2>
        {totalReq === 0 ? (
          <p className="muted">要求がまだ登録されていません。</p>
        ) : (
          <div className="status-bars">
            {(Object.keys(STATUS_LABELS) as RequirementStatus[])
              .filter((s) => (breakdown[s] ?? 0) > 0)
              .map((s) => {
                const count = breakdown[s] ?? 0
                const pct = Math.round((count / totalReq) * 100)
                return (
                  <div key={s} className="status-bar-row">
                    <span className={`badge badge-${s}`}>{STATUS_LABELS[s]}</span>
                    <div className="status-bar-track">
                      <div className={`status-bar-fill fill-${s}`} style={{ width: `${pct}%` }} />
                    </div>
                    <span className="status-bar-count">
                      {count} ({pct}%)
                    </span>
                  </div>
                )
              })}
          </div>
        )}
      </section>
    </div>
  )
}
