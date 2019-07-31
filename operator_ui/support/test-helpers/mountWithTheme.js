import React from 'react'
import { mount } from 'enzyme'
import JssProvider from 'react-jss/lib/JssProvider'
import { SheetsRegistry } from 'react-jss/lib/jss'
import {
  MuiThemeProvider,
  createMuiTheme,
  createGenerateClassName
} from '@material-ui/core/styles'
import theme from '@chainlink/styleguide/src/theme'
import StoreAndMemoryRouter from './StoreAndMemoryRouter'

const sheetsRegistry = new SheetsRegistry()
const muiTheme = createMuiTheme(theme)
const generateClassName = createGenerateClassName()

export default (children, opts = {}) =>
  mount(
    <JssProvider
      registry={sheetsRegistry}
      generateClassName={generateClassName}
    >
      <MuiThemeProvider theme={muiTheme} sheetsManager={new Map()}>
        <StoreAndMemoryRouter initialEntries={['/']}>
          {children}
        </StoreAndMemoryRouter>
      </MuiThemeProvider>
    </JssProvider>
  )
