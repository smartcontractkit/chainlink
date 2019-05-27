import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { mount } from 'enzyme'
import { Provider } from 'react-redux'
import { MuiThemeProvider } from '@material-ui/core/styles'
import theme from 'theme'
import createStore from 'connectors/redux'

export default (children, opts = {}) =>
  mount(
    <MuiThemeProvider theme={theme}>
      <Provider store={createStore()}>
        <MemoryRouter initialEntries={['/']}>{children}</MemoryRouter>
      </Provider>
    </MuiThemeProvider>
  )
