import React from 'react'
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

const App = () => {
  return (
    <Provider store={store}>
      <Layout />
    </Provider>
  )
}

export default App
