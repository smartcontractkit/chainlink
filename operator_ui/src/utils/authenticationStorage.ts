import * as storage from 'utils/storage';

export const get = () => storage.get('authentication')

// CHECK ME
export const set = (obj: any) => storage.set('authentication', obj)
