import storage from 'local-storage-fallback'

export const get = key => {
  const localStorageItem = storage.getItem(`chainlink.${key}`)
  const obj = {}

  if (localStorageItem) {
    try {
      return JSON.parse(localStorageItem)
    } catch (e) {
      // continue regardless of error
    }
  }

  return obj
}

export const set = (key, obj) => {
  storage.setItem(`chainlink.${key}`, JSON.stringify(obj))
}
