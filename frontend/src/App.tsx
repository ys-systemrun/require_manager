import { NavLink, Navigate, Route, Routes } from 'react-router-dom'
import { ProjectProvider, useProjects } from './ProjectContext'
import { DashboardPage } from './pages/DashboardPage'
import { RequirementsPage } from './pages/RequirementsPage'
import { TermsPage } from './pages/TermsPage'
import { DomainModelPage } from './pages/DomainModelPage'
import { ProjectsPage } from './pages/ProjectsPage'

function ProjectSwitcher() {
  const { projects, currentProject, setCurrentProjectId } = useProjects()
  if (projects.length === 0) return null
  return (
    <div className="project-switcher">
      <label>プロジェクト</label>
      <select
        value={currentProject?.id ?? ''}
        onChange={(e) => setCurrentProjectId(Number(e.target.value))}
      >
        {projects.map((p) => (
          <option key={p.id} value={p.id}>
            {p.name}
          </option>
        ))}
      </select>
    </div>
  )
}

function Layout() {
  const { loading, projects } = useProjects()

  const nav = [
    { to: '/dashboard', label: 'ダッシュボード', icon: '📊' },
    { to: '/requirements', label: '要求管理', icon: '📋' },
    { to: '/terms', label: '用語辞書', icon: '📖' },
    { to: '/domain', label: 'ドメインモデル', icon: '🧩' },
    { to: '/projects', label: 'プロジェクト設定', icon: '⚙️' },
  ]

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <span className="brand-mark">RM</span>
          <div>
            <div className="brand-title">Require Manager</div>
            <div className="brand-sub">要件定義・要求管理</div>
          </div>
        </div>
        <ProjectSwitcher />
        <nav>
          {nav.map((n) => (
            <NavLink
              key={n.to}
              to={n.to}
              className={({ isActive }) => `nav-item${isActive ? ' active' : ''}`}
            >
              <span className="nav-icon">{n.icon}</span>
              {n.label}
            </NavLink>
          ))}
        </nav>
      </aside>

      <main className="content">
        {loading ? (
          <div className="empty-state">読み込み中...</div>
        ) : projects.length === 0 ? (
          <Routes>
            <Route path="/projects" element={<ProjectsPage />} />
            <Route
              path="*"
              element={
                <div className="empty-state">
                  <p>プロジェクトがまだありません。</p>
                  <NavLink to="/projects" className="btn btn-primary">
                    プロジェクトを作成
                  </NavLink>
                </div>
              }
            />
          </Routes>
        ) : (
          <Routes>
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            <Route path="/dashboard" element={<DashboardPage />} />
            <Route path="/requirements" element={<RequirementsPage />} />
            <Route path="/terms" element={<TermsPage />} />
            <Route path="/domain" element={<DomainModelPage />} />
            <Route path="/projects" element={<ProjectsPage />} />
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
          </Routes>
        )}
      </main>
    </div>
  )
}

export default function App() {
  return (
    <ProjectProvider>
      <Layout />
    </ProjectProvider>
  )
}
