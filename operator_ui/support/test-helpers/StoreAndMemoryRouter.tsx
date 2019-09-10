import React from 'react'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router'
import createStore from '../../src/connectors/redux'

interface Props {
  children: React.ReactNode
  initialEntries?: string[]
}

export default ({ children, initialEntries }: Props) => {
  return (
    <Provider store={createStore()}>
      <MemoryRouter initialEntries={initialEntries}>{children}</MemoryRouter>
    </Provider>
  )
}
