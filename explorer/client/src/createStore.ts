import { Middleware } from 'redux'
import { createStore } from '@chainlink/redux'
import thunkMiddleware from 'redux-thunk'
import reducer from './reducers'

const middleware: Middleware[] = [thunkMiddleware]

export default () => createStore(reducer, middleware)
