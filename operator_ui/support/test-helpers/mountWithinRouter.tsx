import React from 'react'
import { MemoryRouter } from 'react-router'
import { mount } from 'enzyme'

export default (children: React.ReactNode, initialEntries?: string[]) => {
  return mount(
    <MemoryRouter initialEntries={initialEntries}>{children}</MemoryRouter>
  )
}
