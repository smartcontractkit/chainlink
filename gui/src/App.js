import React, { PureComponent } from 'react'
import Layout from 'Layout'
import createStore from 'connectors/redux'
import { Provider } from 'react-redux'
import { hot } from 'react-hot-loader'

class App extends PureComponent {
  // Remove the server-side injected CSS.
  componentDidMount () {
    const jssStyles = document.getElementById('jss-server-side')
    if (jssStyles && jssStyles.parentNode) {
      jssStyles.parentNode.removeChild(jssStyles)
    }
  }

  render () {
    return (
      <Provider store={createStore()}>
        <Layout />
      </Provider>
    )
  }
}

export default hot(module)(App)
