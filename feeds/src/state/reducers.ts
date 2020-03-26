import { combineReducers } from 'redux'
import { createFilter } from 'redux-persist-transform-filter'
import { persistReducer, createMigrate } from 'redux-persist'
import storage from 'redux-persist/lib/storage'
import * as reducers from './ducks'

const rootReducer = combineReducers({
  ...reducers,
})

export const INITIAL_STATE = rootReducer(undefined, { type: 'initial_state' })
export type AppState = typeof INITIAL_STATE

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

export const reducer = persistReducer(persistConfig, rootReducer)
