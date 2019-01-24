import React, { PureComponent } from 'react'
import { Provider } from 'react-redux'
import createStore from 'connectors/redux'
import Layout from 'Layout'
import './index.css'

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
