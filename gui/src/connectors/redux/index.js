import { createStore, applyMiddleware, compose } from 'redux'
import thunkMiddleware from 'redux-thunk'
import logger from 'redux-logger'
import reducer from './reducers'

if (typeof window === 'undefined') {
  global.window = {}
}

let middleware = [thunkMiddleware]
if (process.env.LOG_REDUX === 'true') {
  middleware = middleware.concat(logger)
}
const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose

export default () =>
  createStore(reducer, {}, composeEnhancers(applyMiddleware(...middleware)))
