import axios from 'axios'
import type {
  DomainEntity,
  DomainRelationship,
  Project,
  ProjectSummary,
  Requirement,
  Term,
} from './types'

const baseURL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api'

export const http = axios.create({ baseURL })

// ---------- Projects ----------
export const projectsApi = {
  list: () => http.get<Project[]>('/projects').then((r) => r.data),
  get: (id: number) => http.get<Project>(`/projects/${id}`).then((r) => r.data),
  summary: (id: number) =>
    http.get<ProjectSummary>(`/projects/${id}/summary`).then((r) => r.data),
  create: (data: Partial<Project>) =>
    http.post<Project>('/projects', data).then((r) => r.data),
  update: (id: number, data: Partial<Project>) =>
    http.put<Project>(`/projects/${id}`, data).then((r) => r.data),
  remove: (id: number) => http.delete(`/projects/${id}`).then(() => undefined),
}

// ---------- Requirements ----------
export const requirementsApi = {
  list: (params: { projectId?: number; status?: string; type?: string; keyword?: string }) =>
    http.get<Requirement[]>('/requirements', { params }).then((r) => r.data),
  get: (id: number) => http.get<Requirement>(`/requirements/${id}`).then((r) => r.data),
  create: (data: Partial<Requirement>) =>
    http.post<Requirement>('/requirements', data).then((r) => r.data),
  update: (id: number, data: Partial<Requirement>) =>
    http.put<Requirement>(`/requirements/${id}`, data).then((r) => r.data),
  remove: (id: number) => http.delete(`/requirements/${id}`).then(() => undefined),
}

// ---------- Terms ----------
export const termsApi = {
  list: (params: { projectId?: number; keyword?: string }) =>
    http.get<Term[]>('/terms', { params }).then((r) => r.data),
  create: (data: Partial<Term>) => http.post<Term>('/terms', data).then((r) => r.data),
  update: (id: number, data: Partial<Term>) =>
    http.put<Term>(`/terms/${id}`, data).then((r) => r.data),
  remove: (id: number) => http.delete(`/terms/${id}`).then(() => undefined),
}

// ---------- Domain entities ----------
export const entitiesApi = {
  list: (params: { projectId?: number }) =>
    http.get<DomainEntity[]>('/entities', { params }).then((r) => r.data),
  get: (id: number) => http.get<DomainEntity>(`/entities/${id}`).then((r) => r.data),
  create: (data: Partial<DomainEntity>) =>
    http.post<DomainEntity>('/entities', data).then((r) => r.data),
  update: (id: number, data: Partial<DomainEntity>) =>
    http.put<DomainEntity>(`/entities/${id}`, data).then((r) => r.data),
  remove: (id: number) => http.delete(`/entities/${id}`).then(() => undefined),
}

// ---------- Domain relationships ----------
export const relationshipsApi = {
  list: (params: { projectId?: number }) =>
    http.get<DomainRelationship[]>('/relationships', { params }).then((r) => r.data),
  create: (data: Partial<DomainRelationship>) =>
    http.post<DomainRelationship>('/relationships', data).then((r) => r.data),
  update: (id: number, data: Partial<DomainRelationship>) =>
    http.put<DomainRelationship>(`/relationships/${id}`, data).then((r) => r.data),
  remove: (id: number) => http.delete(`/relationships/${id}`).then(() => undefined),
}
