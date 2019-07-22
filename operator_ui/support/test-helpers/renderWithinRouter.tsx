import React from 'react'
import { MemoryRouter } from 'react-router'
import { render } from 'enzyme'

export default (component: React.ReactNode) => {
  return render(<MemoryRouter>{component}</MemoryRouter>)
}
