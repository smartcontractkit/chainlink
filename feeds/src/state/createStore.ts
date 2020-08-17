import { createStore } from '@chainlink/redux'
import { Middleware } from 'redux'
import { persistStore } from 'redux-persist'
import thunkMiddleware from 'redux-thunk'
import { reducer } from './reducers'

const middleware: Middleware[] = [thunkMiddleware]

export default () => {
  const store = createStore(reducer, middleware)
  const persistor = persistStore(store)
  return { store, persistor }
}
