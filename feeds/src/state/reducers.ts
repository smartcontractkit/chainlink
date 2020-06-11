import { combineReducers } from 'redux'
import { createMigrate, persistReducer } from 'redux-persist'
import { createFilter } from 'redux-persist-transform-filter'
import storage from 'redux-persist/lib/storage'
import * as reducers from './ducks'

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
  migrate: createMigrate(migrations),
}

const rootReducer = combineReducers(reducers)
export type AppState = ReturnType<typeof rootReducer>
export const reducer = persistReducer(persistConfig, rootReducer)
