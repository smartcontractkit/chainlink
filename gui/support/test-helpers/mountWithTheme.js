import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import JssProvider from 'react-jss/lib/JssProvider'
import { SheetsRegistry } from 'react-jss/lib/jss'
import {
  MuiThemeProvider,
  createMuiTheme,
  createGenerateClassName
} from '@material-ui/core/styles'
import theme from 'theme'
import createStore from 'connectors/redux'

const sheetsRegistry = new SheetsRegistry()
const muiTheme = createMuiTheme(theme)
const generateClassName = createGenerateClassName()

export default (children, opts = {}) => (
  mount(
    <JssProvider registry={sheetsRegistry} generateClassName={generateClassName}>
      <MuiThemeProvider theme={muiTheme} sheetsManager={new Map()}>
        <Provider store={createStore()}>
          <MemoryRouter initialEntries={['/']}>
            {children}
          </MemoryRouter>
        </Provider>
      </MuiThemeProvider>
    </JssProvider>
  )
)
