import React, { PureComponent } from 'react'
import { Provider } from 'react-redux'
import createStore from './connectors/redux'
import Layout from './Layout'
import './index.css'
import { set } from './utils/storage'

const store = createStore()

store.subscribe(() => {
  const prevURL = store.getState().notifications.currentUrl
  if(prevURL !== '/signin') set("persistURL", prevURL)
})

class App extends PureComponent {
  // Remove the server-side injected CSS.
  componentDidMount () {
    const jssStyles = document.getElementById('jss-server-side')
    if (jssStyles && jssStyles.parentNode) {
      jssStyles.parentNode.removeChild(jssStyles)
    }
  }

  render (): JSX.Element {
    return (
      <Provider store={store}>
        <Layout />
      </Provider>
    )
  }
}

export default (App)
