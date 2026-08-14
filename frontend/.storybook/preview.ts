import type { Preview } from '@storybook/react'
import '../src/styles.css'

const preview: Preview = {
  parameters: {
    layout: 'centered',
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
    backgrounds: {
      default: 'app',
      values: [
        { name: 'app', value: '#f5f6fa' },
        { name: 'surface', value: '#ffffff' },
        { name: 'dark', value: '#111827' },
      ],
    },
  },
}

export default preview
