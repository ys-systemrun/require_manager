import type {
  Priority,
  RelationshipType,
  RequirementStatus,
  RequirementType,
  Stereotype,
} from './types'

export const REQ_TYPE_LABELS: Record<RequirementType, string> = {
  business: '業務要求',
  functional: '機能要件',
  non_functional: '非機能要件',
  constraint: '制約',
}

export const PRIORITY_LABELS: Record<Priority, string> = {
  high: '高',
  medium: '中',
  low: '低',
}

export const STATUS_LABELS: Record<RequirementStatus, string> = {
  draft: 'ドラフト',
  proposed: '提案',
  approved: '承認済',
  implemented: '実装済',
  verified: '検証済',
  rejected: '却下',
  deprecated: '廃止',
}

export const STEREOTYPE_LABELS: Record<Stereotype, string> = {
  entity: 'エンティティ',
  value_object: '値オブジェクト',
  aggregate_root: '集約ルート',
  service: 'サービス',
  event: 'イベント',
}

export const RELATIONSHIP_LABELS: Record<RelationshipType, string> = {
  association: '関連',
  aggregation: '集約',
  composition: 'コンポジション',
  inheritance: '継承',
  dependency: '依存',
}

export const REQ_TYPES = Object.keys(REQ_TYPE_LABELS) as RequirementType[]
export const PRIORITIES = Object.keys(PRIORITY_LABELS) as Priority[]
export const STATUSES = Object.keys(STATUS_LABELS) as RequirementStatus[]
export const STEREOTYPES = Object.keys(STEREOTYPE_LABELS) as Stereotype[]
export const RELATIONSHIP_TYPES = Object.keys(RELATIONSHIP_LABELS) as RelationshipType[]
