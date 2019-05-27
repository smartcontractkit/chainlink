import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles'
import theme from 'theme'
import createStore from 'connectors/redux'
import { SheetsRegistry } from 'jss'
import { StylesProvider, createGenerateClassName } from '@material-ui/styles'

const sheetsRegistry = new SheetsRegistry()
const muiTheme = createMuiTheme(theme)
const generateClassName = createGenerateClassName()
export default (children, opts = {}) =>
  mount(
    <StylesProvider
      registry={sheetsRegistry}
      generateClassName={generateClassName}
    >
      <MuiThemeProvider theme={muiTheme} sheetsManager={new Map()}>
        <Provider store={createStore()}>
          <MemoryRouter initialEntries={['/']}>{children}</MemoryRouter>
        </Provider>
      </MuiThemeProvider>
    </StylesProvider>
  )
