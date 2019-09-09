import React from 'react'
import { mount } from 'enzyme'
import StoreAndMemoryRouter from './StoreAndMemoryRouter'

export default (children: React.ReactNode, initialEntries?: string[]) => {
  return mount(
    <StoreAndMemoryRouter initialEntries={initialEntries}>
      {children}
    </StoreAndMemoryRouter>,
  )
}
