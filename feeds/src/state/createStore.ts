import { Middleware } from 'redux'
import thunkMiddleware from 'redux-thunk'
import { persistStore } from 'redux-persist'
import { createStore } from '@chainlink/redux'
import { reducer } from './reducers'

const middleware: Middleware[] = [thunkMiddleware]

export default () => {
  const store = createStore(reducer, middleware)
  const persistor = persistStore(store)
  return { store, persistor }
}
