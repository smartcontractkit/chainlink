import React from 'react'
import ReactDOM from 'react-dom'
import promiseFinally from 'promise.prototype.finally'
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles'

import App from './App'
import theme from './theme'

promiseFinally.shim(Promise)

window.JavascriptTimeAgo = JavascriptTimeAgo
JavascriptTimeAgo.locale(en)

// Export top level component as JSX (for static rendering)
export default App

// Render app
if (typeof document !== 'undefined') {
  const renderMethod = module.hot ? ReactDOM.render : ReactDOM.hydrate
  const muiTheme = createMuiTheme(theme)

  const render = Comp => {
    renderMethod(
      <MuiThemeProvider theme={muiTheme}>
        <Comp />
      </MuiThemeProvider>,
      document.getElementById('root')
    )
  }

  render(App)
}
