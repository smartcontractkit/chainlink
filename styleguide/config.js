import { configure } from '@storybook/react'

const loadStories = () => {
  require('./stories/index.js')
}

configure(loadStories, module)
