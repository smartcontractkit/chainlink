import { theme } from 'theme'
import {
  createGenerateClassName,
  createMuiTheme,
  MuiThemeProvider,
} from '@material-ui/core/styles'
import { mount } from 'enzyme'
import React from 'react'
import { SheetsRegistry } from 'react-jss/lib/jss'
import JssProvider from 'react-jss/lib/JssProvider'
import StoreAndMemoryRouter from './StoreAndMemoryRouter'

const sheetsRegistry = new SheetsRegistry()
const muiTheme = createMuiTheme(theme)
const generateClassName = createGenerateClassName()

export const mountWithProviders = (elements, { initialEntries = ['/'] } = {}) =>
  mount(
    <JssProvider
      registry={sheetsRegistry}
      generateClassName={generateClassName}
    >
      <MuiThemeProvider theme={muiTheme} sheetsManager={new Map()}>
        <StoreAndMemoryRouter initialEntries={initialEntries}>
          {elements}
        </StoreAndMemoryRouter>
      </MuiThemeProvider>
    </JssProvider>,
  )

export default mountWithProviders
