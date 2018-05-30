import { createStore, applyMiddleware, compose } from 'redux'
import thunkMiddleware from 'redux-thunk'
import { createLogger } from 'redux-logger'
import reducer from './reducers'

if (typeof window === 'undefined') {
  global.window = {}
}

const middleware = [thunkMiddleware]
if (process.env.NODE_ENV !== 'test') {
  middleware.push(createLogger())
}

const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose

export default () => (
  createStore(
    reducer,
    {},
    composeEnhancers(
      applyMiddleware(...middleware)
    )
  )
)
