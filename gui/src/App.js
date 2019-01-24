import React, { PureComponent } from 'react'
import Layout from 'Layout'
import createStore from 'connectors/redux'
import { Provider } from 'react-redux'

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

export default (App)
