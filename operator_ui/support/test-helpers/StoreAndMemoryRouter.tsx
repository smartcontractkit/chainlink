import React from 'react'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import createStore from 'createStore'

interface Props {
  initialEntries?: string[]
}

const StoreAndMemoryRouter: React.FC<Props> = ({
  children,
  initialEntries,
}) => {
  return (
    <Provider store={createStore()}>
      <MemoryRouter initialEntries={initialEntries}>{children}</MemoryRouter>
    </Provider>
  )
}

export default StoreAndMemoryRouter
