import React, { PureComponent } from 'react'
import { Provider } from 'react-redux'
import createStore from './createStore'
import './index.css'
import Layout from './Layout'
import { setPersistUrl } from './utils/storage'

const SIGNIN_PATH = '/signin'

const store = createStore()

store.subscribe(() => {
  const prevURL = store.getState().notifications.currentUrl
  if (prevURL && prevURL !== SIGNIN_PATH) {
    setPersistUrl(prevURL)
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
