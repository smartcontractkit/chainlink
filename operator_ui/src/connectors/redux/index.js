import { createStore, applyMiddleware, compose } from 'redux'
import thunkMiddleware from 'redux-thunk'
import reducer from './reducers'

if (typeof window === 'undefined') {
  global.window = {}
}

let middleware = [thunkMiddleware]
const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose

export default () =>
  createStore(reducer, {}, composeEnhancers(applyMiddleware(...middleware)))
