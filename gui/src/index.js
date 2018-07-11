import React from 'react'
import ReactDOM from 'react-dom'
import createStore from 'connectors/redux'
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles'
import { Provider } from 'react-redux'

// Your top level component
import App from './App'

// Your Material UI Custom theme
import theme from './theme'

// Export your top level component as JSX (for static rendering)
export default App

// Render your app
if (typeof document !== 'undefined') {
  const renderMethod = module.hot ? ReactDOM.render : ReactDOM.hydrate
  const muiTheme = createMuiTheme(theme)

  const render = Comp => {
    renderMethod(
      <MuiThemeProvider theme={muiTheme}>
        <Provider store={createStore()}>
          <Comp />
        </Provider>
      </MuiThemeProvider>,
      document.getElementById('root')
    )
  }

  // Render!
  render(App)
}
