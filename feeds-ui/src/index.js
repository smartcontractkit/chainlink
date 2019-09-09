import React from 'react'
import ReactDOM from 'react-dom'
import AppRoutes from './pages'
import * as serviceWorker from './serviceWorker'
import { Provider as ReduxProvider } from 'react-redux'
import { store, persistor } from './state/store'
import { PersistGate } from 'redux-persist/integration/react'
import { Header } from 'components/header'
import './App.css'

ReactDOM.render(
  <ReduxProvider store={store}>
    <PersistGate loading={null} persistor={persistor}>
      <Header />
      <AppRoutes />
    </PersistGate>
  </ReduxProvider>,
  document.getElementById('root')
)

serviceWorker.unregister()
