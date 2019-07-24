import React from 'react'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router'
import createStore from '../../src/connectors/redux'

interface IProps {
  children: React.ReactNode
  initialEntries?: string[]
}

export default ({ children, initialEntries }: IProps) => {
  return (
    <Provider store={createStore()}>
      <MemoryRouter initialEntries={initialEntries}>{children}</MemoryRouter>
    </Provider>
  )
}
