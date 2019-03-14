import React from 'react'
import { render } from 'react-snapshot'
import './index.css'
import App from './App'
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles'
import theme from './theme'
import * as serviceWorker from './serviceWorker'

const muiTheme = createMuiTheme(theme)

render(
  <MuiThemeProvider theme={muiTheme}>
    <App />
  </MuiThemeProvider>,
  document.getElementById('root')
)

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister()
