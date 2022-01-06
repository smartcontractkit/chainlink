import { MuiThemeProvider } from '@material-ui/core/styles'
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import moment from 'moment'
import React from 'react'
import ReactDOM from 'react-dom'
import { AppContainer } from 'react-hot-loader'
import App from './App'
import { theme } from './theme'
import { ApolloProvider } from '@apollo/client'
import { client } from './apollo'

JavascriptTimeAgo.locale(en)
moment.defaultFormat = 'YYYY-MM-DD h:mm:ss A'

export default App

if (typeof document !== 'undefined') {
  const renderMethod = module.hot ? ReactDOM.render : ReactDOM.hydrate

  const render = (Comp) => {
    renderMethod(
      <ApolloProvider client={client}>
        <AppContainer>
          <MuiThemeProvider theme={theme}>
            <Comp />
          </MuiThemeProvider>
        </AppContainer>
      </ApolloProvider>,
      document.getElementById('root'),
    )
  }

  render(App)
  // Hot Module Replacement
  if (module.hot) {
    module.hot.accept('./App', () => render(require('./App').default))
  }
}
