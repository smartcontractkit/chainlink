import { combineReducers } from 'redux'
import { createMigrate, persistReducer } from 'redux-persist'
import { createFilter } from 'redux-persist-transform-filter'
import storage from 'redux-persist/lib/storage'
import * as reducers from './ducks'
import { InitialStateAction } from 'state/actions'

const rootReducer = combineReducers(reducers)

const initialAction: InitialStateAction = { type: 'INITIAL_STATE' }
export const INITIAL_STATE = rootReducer(undefined, initialAction)
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
  migrate: createMigrate(migrations),
}

export const reducer = persistReducer(persistConfig, rootReducer)
