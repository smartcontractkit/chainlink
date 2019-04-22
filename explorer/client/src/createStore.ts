import { createStore, applyMiddleware } from 'redux'
import { composeWithDevTools } from 'redux-devtools-extension'
import thunkMiddleware from 'redux-thunk'
import reducer from './reducers'

let middleware = [thunkMiddleware]
const composeEnhancers = composeWithDevTools({})

export default () =>
  createStore(reducer, composeEnhancers(applyMiddleware(...middleware)))
