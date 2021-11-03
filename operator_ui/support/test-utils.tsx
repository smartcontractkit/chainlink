import React from 'react'
import { Provider as ReduxProvider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import { render, RenderOptions } from '@testing-library/react'

import { MuiThemeProvider } from '@material-ui/core/styles'

import createStore from 'src/createStore'
import Notifications from 'pages/Notifications'
import { theme } from 'src/theme'

// AllTheProviders wraps the ui with our providers
const AllTheProviders: React.FC = ({ children }) => {
  return (
    <MuiThemeProvider theme={theme}>
      <ReduxProvider store={createStore()}>{children}</ReduxProvider>
    </MuiThemeProvider>
  )
}

// customRender wraps the render method with the providers.
//
// https://testing-library.com/docs/react-testing-library/setup/#custom-render
const customRender = (
  ui: React.ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>,
) => render(ui, { wrapper: AllTheProviders, ...options })

interface RenderWithRouterProps {
  initialEntries?: string[]
}

// renderWithRouter behaves like 'render' except it wraps the provided ui in a
// Router.
//
// Use this when you need a router in your tests for page components and any
// component which uses react-router hooks.
//
// TODO - Could not figure out how to type the return.
const renderWithRouter: any = (
  ui: React.ReactElement,
  { initialEntries }: RenderWithRouterProps = { initialEntries: ['/'] },
  options?: Omit<RenderOptions, 'wrapper'>,
) => {
  return {
    ...customRender(
      <MemoryRouter initialEntries={initialEntries}>
        {/* Notifications are added here because they can display a link */}
        <Notifications />
        {ui}
      </MemoryRouter>,
      options,
    ),
  }
}

export * from '@testing-library/react'
export { customRender as render, renderWithRouter }
