import { Middleware } from 'redux'
import { createStore } from '@chainlink/redux'
import thunkMiddleware from 'redux-thunk'
import { createQueryMiddleware } from './middleware'
import reducer from './reducers'

const middleware: Middleware[] = [thunkMiddleware, createQueryMiddleware()]

export default () => createStore(reducer, middleware)
