import React, { PureComponent } from 'react'
import { Provider } from 'react-redux'
import { SnackbarProvider } from 'material-ui-snackbar-provider'
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
      <SnackbarProvider
        SnackbarProps={{
          autoHideDuration: 4000,
          anchorOrigin: { vertical: 'bottom', horizontal: 'right' },
        }}
      >
        <Provider store={store}>
          <Layout />
        </Provider>
      </SnackbarProvider>
    )
  }
}

export default App
