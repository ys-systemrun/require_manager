import type { Meta, StoryObj } from '@storybook/react'
import { Badge } from './Badge'
import { STATUS_LABELS } from '../constants'
import type { RequirementStatus } from '../types'

const meta: Meta<typeof Badge> = {
  title: 'Components/Badge',
  component: Badge,
  tags: ['autodocs'],
  args: {
    label: 'ラベル',
  },
  argTypes: {
    variant: {
      control: 'select',
      options: [
        undefined,
        'draft',
        'proposed',
        'approved',
        'implemented',
        'verified',
        'rejected',
        'deprecated',
      ],
    },
  },
}
export default meta

type Story = StoryObj<typeof Badge>

export const Default: Story = {}

export const Approved: Story = {
  args: { label: '承認済', variant: 'approved' },
}

export const Rejected: Story = {
  args: { label: '却下', variant: 'rejected' },
}

/** 要求ステータスに対応する全バリアントの一覧。 */
export const AllStatuses: Story = {
  render: () => (
    <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', maxWidth: 480 }}>
      {(Object.keys(STATUS_LABELS) as RequirementStatus[]).map((s) => (
        <Badge key={s} label={STATUS_LABELS[s]} variant={s} />
      ))}
    </div>
  ),
}
