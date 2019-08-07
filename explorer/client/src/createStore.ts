import { createStore, applyMiddleware, Middleware } from 'redux'
import { composeWithDevTools } from 'redux-devtools-extension'
import thunkMiddleware from 'redux-thunk'
import logger from 'redux-logger'
import reducer from './reducers'

let middleware: Middleware[] = [thunkMiddleware]
if (process.env.LOG_REDUX === 'true') {
  middleware = middleware.concat(logger)
}
const composeEnhancers = composeWithDevTools({})

export default () =>
  createStore(reducer, composeEnhancers(applyMiddleware(...middleware)))
