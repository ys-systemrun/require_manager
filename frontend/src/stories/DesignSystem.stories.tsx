import type { ReactNode } from 'react'
import type { Meta, StoryObj } from '@storybook/react'
import {
  PRIORITIES,
  PRIORITY_LABELS,
  REQ_TYPES,
  REQ_TYPE_LABELS,
  STATUSES,
  STATUS_LABELS,
  STEREOTYPE_LABELS,
} from '../constants'

const meta: Meta = {
  title: 'Design System/Overview',
  parameters: { layout: 'padded' },
}
export default meta

type Story = StoryObj

function Section({ title, children }: { title: string; children: ReactNode }) {
  return (
    <section style={{ marginBottom: 28 }}>
      <h2 style={{ fontSize: 14, marginBottom: 10 }}>{title}</h2>
      <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', alignItems: 'center' }}>
        {children}
      </div>
    </section>
  )
}

/** アプリ全体で使われるバッジ・ボタンなどの視覚要素カタログ。 */
export const Overview: Story = {
  render: () => (
    <div style={{ maxWidth: 760 }}>
      <Section title="ステータス バッジ">
        {STATUSES.map((s) => (
          <span key={s} className={`badge badge-${s}`}>
            {STATUS_LABELS[s]}
          </span>
        ))}
      </Section>

      <Section title="要求種別 バッジ">
        {REQ_TYPES.map((t) => (
          <span key={t} className={`badge type-${t}`}>
            {REQ_TYPE_LABELS[t]}
          </span>
        ))}
      </Section>

      <Section title="優先度 バッジ">
        {PRIORITIES.map((p) => (
          <span key={p} className={`badge prio-${p}`}>
            {PRIORITY_LABELS[p]}
          </span>
        ))}
      </Section>

      <Section title="ボタン">
        <button className="btn btn-primary">プライマリ</button>
        <button className="btn">セカンダリ</button>
        <button className="btn small">スモール</button>
        <button className="btn btn-primary" disabled>
          無効
        </button>
        <button className="icon-btn">✏️</button>
        <button className="link">リンク</button>
      </Section>

      <Section title="チップ / タグ">
        <span className="chip">SHOP</span>
        <span className="chip small">在庫</span>
        <span className="tag">別名</span>
        <span className="tag active-tag">選択中</span>
      </Section>
    </div>
  ),
}

/** ダッシュボードの統計カード。 */
export const StatCards: Story = {
  render: () => (
    <div className="stat-grid" style={{ maxWidth: 760 }}>
      {[
        { label: '要求', value: 12, icon: '📋', color: 'blue' },
        { label: '用語', value: 8, icon: '📖', color: 'green' },
        { label: 'エンティティ', value: 5, icon: '🧩', color: 'purple' },
        { label: '関連', value: 4, icon: '🔗', color: 'orange' },
      ].map((c) => (
        <div key={c.label} className={`stat-card stat-${c.color}`}>
          <span className="stat-icon">{c.icon}</span>
          <div>
            <div className="stat-value">{c.value}</div>
            <div className="stat-label">{c.label}</div>
          </div>
        </div>
      ))}
    </div>
  ),
}

/** ドメインモデルの UML クラス風カード。 */
export const UmlClass: Story = {
  render: () => (
    <div className="uml-class" style={{ maxWidth: 300 }}>
      <div className="uml-header">
        <div>
          <div className="uml-stereo">«{STEREOTYPE_LABELS.aggregate_root}»</div>
          <div className="uml-name">
            Order<span className="uml-jp"> / 注文</span>
          </div>
        </div>
      </div>
      <div className="uml-attrs">
        <div className="uml-attr">
          <span>
            <span className="pk">PK </span>id
            <span className="req-star">*</span>
          </span>
          <span className="mono muted">uuid</span>
        </div>
        <div className="uml-attr">
          <span>
            orderedAt<span className="req-star">*</span>
          </span>
          <span className="mono muted">datetime</span>
        </div>
        <div className="uml-attr">
          <span>totalAmount</span>
          <span className="mono muted">money</span>
        </div>
      </div>
      <div className="uml-desc">顧客による購入の確定記録。</div>
    </div>
  ),
}

/** 用語辞書のカード。 */
export const TermCard: Story = {
  render: () => (
    <div className="term-card" style={{ maxWidth: 320 }}>
      <div className="term-card-head">
        <div>
          <h3>カート</h3>
          <div className="term-meta">
            <span>かーと</span>
            <span className="mono">Cart</span>
          </div>
        </div>
      </div>
      <p className="term-def">顧客が購入予定の商品を一時的に保持する入れ物。</p>
      <div className="term-aliases">
        <span className="tag">ショッピングカート</span>
        <span className="tag">買い物かご</span>
      </div>
      <span className="chip small">販売</span>
    </div>
  ),
}
