import React, { PureComponent } from 'react'
import { Provider } from 'react-redux'
import createStore from './connectors/redux'
import './index.css'
import Layout from './Layout'
import { set } from './utils/storage'

const store = createStore()

store.subscribe(() => {
  const prevURL = store.getState().notifications.currentUrl
  if (prevURL !== '/signin') {
    set('persistURL', prevURL)
  }
})

class App extends PureComponent {
  // Remove the server-side injected CSS.
  public componentDidMount() {
    const jssStyles = document.getElementById('jss-server-side')
    if (jssStyles && jssStyles.parentNode) {
      jssStyles.parentNode.removeChild(jssStyles)
    }
  }

  public render(): JSX.Element {
    return (
      <Provider store={store}>
        <Layout />
      </Provider>
    )
  }
}

export default App
