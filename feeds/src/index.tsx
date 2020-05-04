import React from 'react'
import ReactDOM from 'react-dom'
import { Provider as ReduxProvider } from 'react-redux'
import * as serviceWorker from './serviceWorker'
import App from './App'
import { Config } from './config'
import createStore from './state/createStore'
import { PersistGate } from 'redux-persist/integration/react'
import ReactGA from 'react-ga'
import './theme.css'

ReactGA.initialize(Config.gaId() ?? '')

const { store, persistor } = createStore()

ReactDOM.render(
  <ReduxProvider store={store}>
    <PersistGate loading={null} persistor={persistor}>
      <App />
    </PersistGate>
  </ReduxProvider>,
  document.getElementById('root'),
)

serviceWorker.unregister()
