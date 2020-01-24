import { combineReducers, Middleware } from 'redux'
import { createStore } from '@chainlink/redux'
import thunkMiddleware from 'redux-thunk'
import aggregatorMiddleware from './middlewares/aggregatorMiddleware'
import * as reducers from './ducks'

const middleware: Middleware[] = [thunkMiddleware, aggregatorMiddleware]

const rootReducer = combineReducers({
  ...reducers,
})

// const migrations = {
//   1: () => {
//     return {}
//   },
// }

// const persistConfig = {
//   key: 'heartbeat',
//   version: 1,
//   storage,
//   whitelist: [''],
//   transforms: [createFilter('aggregation', [''])],
//   migrate: createMigrate(migrations, {
//     debug: process.env.NODE_ENV !== 'production',
//   }),
// }

// const persistedReducer = persistReducer(persistConfig, rootReducer)
// const persistor = persistStore(store)

export default () => createStore(rootReducer, middleware)
