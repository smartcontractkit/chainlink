import * as storage from 'utils/storage'

export const get = () => storage.get('session')

export const set = obj => storage.set('session', obj)
