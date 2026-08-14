import type { Meta, StoryObj } from '@storybook/react'
import { useState } from 'react'
import { Modal } from './Modal'

const meta: Meta<typeof Modal> = {
  title: 'Components/Modal',
  component: Modal,
  parameters: { layout: 'fullscreen' },
  tags: ['autodocs'],
}
export default meta

type Story = StoryObj<typeof Modal>

/** 開いた状態のモーダル。フォームとフッターのボタンを含む。 */
export const Open: Story = {
  render: () => (
    <Modal
      open
      title="要求を追加"
      onClose={() => {}}
      footer={
        <>
          <button className="btn">キャンセル</button>
          <button className="btn btn-primary">保存</button>
        </>
      }
    >
      <div className="form-grid">
        <label className="field">
          <span>ID / コード</span>
          <input defaultValue="REQ-001" />
        </label>
        <label className="field">
          <span>優先度</span>
          <select defaultValue="high">
            <option value="high">高</option>
            <option value="medium">中</option>
            <option value="low">低</option>
          </select>
        </label>
        <label className="field span-2">
          <span>タイトル *</span>
          <input defaultValue="オンラインで商品を購入できること" />
        </label>
        <label className="field span-2">
          <span>説明</span>
          <textarea rows={3} defaultValue="顧客が24時間いつでも商品を購入できる。" />
        </label>
      </div>
    </Modal>
  ),
}

/** ボタンで開閉できるインタラクティブな例。 */
export const Interactive: Story = {
  render: () => {
    const [open, setOpen] = useState(false)
    return (
      <div style={{ padding: 40 }}>
        <button className="btn btn-primary" onClick={() => setOpen(true)}>
          モーダルを開く
        </button>
        <Modal
          open={open}
          title="確認"
          onClose={() => setOpen(false)}
          footer={
            <>
              <button className="btn" onClick={() => setOpen(false)}>
                閉じる
              </button>
              <button className="btn btn-primary" onClick={() => setOpen(false)}>
                OK
              </button>
            </>
          }
        >
          <p>オーバーレイのクリック、×ボタン、フッターのボタンで閉じられます。</p>
        </Modal>
      </div>
    )
  },
}
