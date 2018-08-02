import * as storage from 'utils/storage'

export const get = () => storage.get('authentication')

export const set = obj => storage.set('authentication', obj)
