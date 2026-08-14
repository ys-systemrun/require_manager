import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from 'react'
import { projectsApi } from './api'
import type { Project } from './types'

interface ProjectContextValue {
  projects: Project[]
  currentProject: Project | null
  setCurrentProjectId: (id: number) => void
  reloadProjects: () => Promise<void>
  loading: boolean
}

const ProjectContext = createContext<ProjectContextValue | undefined>(undefined)

const STORAGE_KEY = 'reqmgr.currentProjectId'

export function ProjectProvider({ children }: { children: ReactNode }) {
  const [projects, setProjects] = useState<Project[]>([])
  const [currentId, setCurrentId] = useState<number | null>(() => {
    const saved = localStorage.getItem(STORAGE_KEY)
    return saved ? Number(saved) : null
  })
  const [loading, setLoading] = useState(true)

  const reloadProjects = async () => {
    const list = await projectsApi.list()
    setProjects(list)
    setCurrentId((prev) => {
      if (prev && list.some((p) => p.id === prev)) return prev
      return list.length > 0 ? list[0].id : null
    })
  }

  useEffect(() => {
    reloadProjects().finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    if (currentId) localStorage.setItem(STORAGE_KEY, String(currentId))
  }, [currentId])

  const currentProject = projects.find((p) => p.id === currentId) ?? null

  return (
    <ProjectContext.Provider
      value={{
        projects,
        currentProject,
        setCurrentProjectId: setCurrentId,
        reloadProjects,
        loading,
      }}
    >
      {children}
    </ProjectContext.Provider>
  )
}

export function useProjects() {
  const ctx = useContext(ProjectContext)
  if (!ctx) throw new Error('useProjects must be used within ProjectProvider')
  return ctx
}
