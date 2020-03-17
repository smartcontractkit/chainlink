import { combineReducers, Middleware } from 'redux'
import thunkMiddleware from 'redux-thunk'
import { createFilter } from 'redux-persist-transform-filter'
import { persistStore, persistReducer, createMigrate } from 'redux-persist'
import storage from 'redux-persist/lib/storage'
import { createStore } from '@chainlink/redux'
import aggregatorMiddleware from './middlewares/aggregatorMiddleware'
import * as reducers from './ducks'

const middleware: Middleware[] = [thunkMiddleware, aggregatorMiddleware]

const rootReducer = combineReducers({
  ...reducers,
})

const migrations = {
  1: () => {
    return {}
  },
}

const persistConfig = {
  key: 'heartbeat',
  version: 1,
  storage,
  whitelist: [''],
  transforms: [createFilter('aggregation', [''])],
  migrate: createMigrate(migrations, {
    debug: process.env.NODE_ENV !== 'production',
  }),
}

const persistedReducer = persistReducer(persistConfig, rootReducer)

export default () => {
  const store = createStore(persistedReducer, middleware)
  const persistor = persistStore(store)
  return { store, persistor }
}
