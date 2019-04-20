import { createStore, applyMiddleware } from 'redux'
import { composeWithDevTools } from 'redux-devtools-extension'
import thunkMiddleware from 'redux-thunk'
import * as createLogger from 'redux-logger'
import reducer from './reducers'

let middleware = [thunkMiddleware]
if (process.env.LOG_REDUX === 'true') {
  const logger = (createLogger as any)()
  middleware = middleware.concat(logger)
}
const composeEnhancers = composeWithDevTools({})

export default () =>
  createStore(reducer, composeEnhancers(applyMiddleware(...middleware)))
