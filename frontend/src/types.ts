export interface Project {
  id: number
  name: string
  key: string
  description: string
  createdAt: string
  updatedAt: string
}

export type RequirementType = 'business' | 'functional' | 'non_functional' | 'constraint'
export type Priority = 'high' | 'medium' | 'low'
export type RequirementStatus =
  | 'draft'
  | 'proposed'
  | 'approved'
  | 'implemented'
  | 'verified'
  | 'rejected'
  | 'deprecated'

export interface Requirement {
  id: number
  projectId: number
  parentId: number | null
  code: string
  title: string
  description: string
  rationale: string
  source: string
  type: RequirementType
  priority: Priority
  status: RequirementStatus
  createdAt: string
  updatedAt: string
  children?: Requirement[]
}

export interface Term {
  id: number
  projectId: number
  name: string
  reading: string
  englishName: string
  definition: string
  aliases: string
  category: string
  notes: string
  createdAt: string
  updatedAt: string
}

export type Stereotype = 'entity' | 'value_object' | 'aggregate_root' | 'service' | 'event'

export interface DomainAttribute {
  id?: number
  entityId?: number
  name: string
  dataType: string
  description: string
  required: boolean
  isPrimary: boolean
  sortOrder?: number
}

export interface DomainEntity {
  id: number
  projectId: number
  name: string
  displayName: string
  description: string
  stereotype: Stereotype
  attributes: DomainAttribute[]
  createdAt: string
  updatedAt: string
}

export type RelationshipType =
  | 'association'
  | 'aggregation'
  | 'composition'
  | 'inheritance'
  | 'dependency'

export interface DomainRelationship {
  id: number
  projectId: number
  sourceEntityId: number
  targetEntityId: number
  type: RelationshipType
  label: string
  sourceCardinality: string
  targetCardinality: string
  sourceEntity?: DomainEntity
  targetEntity?: DomainEntity
}

export interface ProjectSummary {
  requirements: number
  terms: number
  entities: number
  relationships: number
  statusBreakdown: Record<string, number>
}
