import React from 'react'
import { render } from 'react-snapshot'
import './index.css'
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles'
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import theme from './theme'
import { Provider } from 'react-redux'
import createStore from './createStore'
import App from './App'
import * as serviceWorker from './serviceWorker'

JavascriptTimeAgo.locale(en)

const muiTheme = createMuiTheme(theme)
const store = createStore()

render(
  <MuiThemeProvider theme={muiTheme}>
    <Provider store={store}>
      <App />
    </Provider>
  </MuiThemeProvider>,
  document.getElementById('root')
)

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister()
