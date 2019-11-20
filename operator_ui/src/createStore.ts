import { createStore, applyMiddleware, Middleware } from 'redux'
import { composeWithDevTools } from 'redux-devtools-extension'
import thunkMiddleware from 'redux-thunk'
import logger from 'redux-logger'
import { createExplorerConnectionMiddleware } from './middleware'
import reducer from './reducers'

let middleware: Middleware[] = [
  thunkMiddleware,
  createExplorerConnectionMiddleware(),
]
if (process.env.LOG_REDUX === 'true') {
  middleware = middleware.concat(logger)
}
const composeEnhancers = composeWithDevTools({})

export default () =>
  createStore(reducer, composeEnhancers(applyMiddleware(...middleware)))
