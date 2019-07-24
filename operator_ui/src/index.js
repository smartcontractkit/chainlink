import React from 'react'
import ReactDOM from 'react-dom'
import promiseFinally from 'promise.prototype.finally'
import JavascriptTimeAgo from 'javascript-time-ago'
import moment from 'moment'
import en from 'javascript-time-ago/locale/en'
import { AppContainer } from 'react-hot-loader'
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles'
import theme from '@chainlink/styleguide/src/theme'
import App from './App'

promiseFinally.shim(Promise)

window.JavascriptTimeAgo = JavascriptTimeAgo
JavascriptTimeAgo.locale(en)
moment.defaultFormat = 'YYYY-MM-DD h:mm:ss A'

export default App

if (typeof document !== 'undefined') {
  const renderMethod = module.hot ? ReactDOM.render : ReactDOM.hydrate

  const render = Comp => {
    renderMethod(
      <AppContainer>
        <MuiThemeProvider theme={createMuiTheme(theme)}>
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
