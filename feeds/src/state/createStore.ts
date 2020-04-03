import { Middleware } from 'redux'
import thunkMiddleware from 'redux-thunk'
import { persistStore } from 'redux-persist'
import { createStore } from '@chainlink/redux'
import aggregatorMiddleware from './middlewares/aggregatorMiddleware'
import { reducer } from './reducers'

const middleware: Middleware[] = [thunkMiddleware, aggregatorMiddleware]

export default () => {
  const store = createStore(reducer, middleware)
  const persistor = persistStore(store)
  return { store, persistor }
}
