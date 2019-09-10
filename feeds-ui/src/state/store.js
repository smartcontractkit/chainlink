import { createStore, applyMiddleware, combineReducers } from 'redux'
import thunkMiddleware from 'redux-thunk'
import { createLogger } from 'redux-logger'
import { persistStore, persistReducer, createMigrate } from 'redux-persist'
import storage from 'redux-persist/lib/storage'
import { createFilter } from 'redux-persist-transform-filter'
import * as reducers from './ducks'

const logger = createLogger({
  collapsed: true
})

/**
 * Redux persis config
 */

const migrations = {
  1: state => {
    return {}
  },
  2: state => {
    return {}
  }
}

const persistConfig = {
  key: 'hearbeat',
  version: 2,
  storage,
  whitelist: ['aggregation'],
  transforms: [createFilter('aggregation', ['oracles'])],
  migrate: createMigrate(migrations, {
    debug: process.env.NODE_ENV !== 'production'
  })
}

const rootReducer = combineReducers({
  ...reducers
})

const persistedReducer = persistReducer(persistConfig, rootReducer)

const initialState = {}

const developmentMiddlewares = applyMiddleware(thunkMiddleware, logger)

const productionMiddlewares = applyMiddleware(thunkMiddleware)

const middlewares =
  process.env.NODE_ENV === 'production'
    ? productionMiddlewares
    : developmentMiddlewares

const store = createStore(persistedReducer, initialState, middlewares)
const persistor = persistStore(store)

export { store, persistor }
