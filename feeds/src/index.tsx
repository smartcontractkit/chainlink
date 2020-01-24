import React from 'react'
import ReactDOM from 'react-dom'
import { Provider as ReduxProvider } from 'react-redux'
import * as serviceWorker from './serviceWorker'
import App from './App'
/* import { store, persistor } from './state/store' */
import createStore from './state/createStore'
/* import { PersistGate } from 'redux-persist/integration/react' */
import ReactGA from 'react-ga'
import './theme.css'

ReactGA.initialize(process.env['REACT_APP_GA_ID'] || '')

const store = createStore()

ReactDOM.render(
  <ReduxProvider store={store}>
    {/* <PersistGate loading={null} persistor={persistor}> */}
    <App />
    {/* </PersistGate> */}
  </ReduxProvider>,
  document.getElementById('root'),
)

serviceWorker.unregister()
