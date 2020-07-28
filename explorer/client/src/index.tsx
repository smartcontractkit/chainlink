import { createMuiTheme, MuiThemeProvider } from '@material-ui/core/styles'
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import moment from 'moment'
import React from 'react'
import { Provider } from 'react-redux'
import { render } from 'react-snapshot'
import App from './App'
import createStore from './createStore'
import './index.css'
import * as serviceWorker from './serviceWorker'
import { theme } from '@chainlink/styleguide'
import ReactGA from 'react-ga'
import * as api from './api/index'

JavascriptTimeAgo.locale(en)
moment.defaultFormat = 'YYYY-MM-DD h:mm:ss A'

const muiTheme = createMuiTheme(theme)
const store = createStore()

async function setGa(): Promise<void> {
  try {
    const { ga } = await api.v1.config.getConfig()
    ReactGA.initialize(ga)
    ReactGA.pageview(window.location.pathname + window.location.search)
  } catch {
    console.error('Could not set GA')
  }
}

setGa()

render(
  <MuiThemeProvider theme={muiTheme}>
    <Provider store={store}>
      <App />
    </Provider>
  </MuiThemeProvider>,
  document.getElementById('root'),
)

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister()
