import React from 'react'
import ReactDOM from 'react-dom'
import promiseFinally from 'promise.prototype.finally'
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import { AppContainer } from 'react-hot-loader'
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles'

import App from 'App'
import theme from 'theme'

promiseFinally.shim(Promise)

window.JavascriptTimeAgo = JavascriptTimeAgo
JavascriptTimeAgo.locale(en)

export default App

if (typeof document !== 'undefined') {
  const renderMethod = module.hot ? ReactDOM.render : ReactDOM.hydrate
  const muiTheme = createMuiTheme(theme)

  const render = Comp => {
    renderMethod(
      <AppContainer>
        <MuiThemeProvider theme={muiTheme}>
          <Comp />
        </MuiThemeProvider>
      </AppContainer>,
      document.getElementById('root')
    )
  }

  render(App)
  // Hot Module Replacement
  if (module.hot) {
    module.hot.accept('./App', () => render(require('./App').default))
  }
}
