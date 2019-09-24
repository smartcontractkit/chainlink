import { createStore, applyMiddleware, combineReducers } from 'redux'
import thunkMiddleware from 'redux-thunk'
import { createLogger } from 'redux-logger'
import { persistStore, persistReducer, createMigrate } from 'redux-persist'
import storage from 'redux-persist/lib/storage'
import { createFilter } from 'redux-persist-transform-filter'
import * as reducers from './ducks'
import aggregatorMiddleware from './middlewares/aggregatorMiddleware'

const logger = createLogger({
  collapsed: true
})

/**
 * Redux persis config
 */

const migrations = {
  1: state => {
    return {}
  }
}

const persistConfig = {
  key: 'heartbeat',
  version: 1,
  storage,
  whitelist: [''],
  transforms: [createFilter('aggregation', [''])],
  migrate: createMigrate(migrations, {
    debug: process.env.NODE_ENV !== 'production'
  })
}

const rootReducer = combineReducers({
  ...reducers
})

const persistedReducer = persistReducer(persistConfig, rootReducer)

const initialState = {}

const developmentMiddlewares = applyMiddleware(
  thunkMiddleware,
  aggregatorMiddleware,
  logger
)

const productionMiddlewares = applyMiddleware(
  thunkMiddleware,
  aggregatorMiddleware
)

const middlewares =
  process.env.NODE_ENV === 'production'
    ? productionMiddlewares
    : developmentMiddlewares

const store = createStore(persistedReducer, initialState, middlewares)
const persistor = persistStore(store)

export { store, persistor }
