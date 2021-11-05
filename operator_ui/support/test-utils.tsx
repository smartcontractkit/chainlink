import React from 'react'
import { Provider as ReduxProvider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import { render, RenderOptions, RenderResult } from '@testing-library/react'

import { MuiThemeProvider } from '@material-ui/core/styles'

import createStore from 'src/createStore'
import { theme } from 'src/theme'

const AllTheProviders: React.FC = ({ children }) => {
  return (
    <MuiThemeProvider theme={theme}>
      <ReduxProvider store={createStore()}>{children}</ReduxProvider>
    </MuiThemeProvider>
  )
}

const customRender = (
  ui: React.ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>,
): RenderResult =>
  render(ui, {
    wrapper: AllTheProviders,
    ...options,
  })
interface RenderWithRouterProps {
  initialEntries?: string[]
}

// renderWithRouter behaves like 'render' except it wraps the provided ui in a
// Router.
//
// Use this when you need a router in your tests for page components and any
// component which uses react-router hooks.
const renderWithRouter = (
  ui: React.ReactElement,
  { initialEntries }: RenderWithRouterProps = { initialEntries: ['/'] },
  options?: Omit<RenderOptions, 'wrapper'>,
): RenderResult => {
  return {
    ...customRender(
      <MemoryRouter initialEntries={initialEntries}>{ui}</MemoryRouter>,
      options,
    ),
  }
}

export * from '@testing-library/react'
export { customRender as render, renderWithRouter }
