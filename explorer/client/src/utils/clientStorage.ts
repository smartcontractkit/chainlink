import storage from 'local-storage-fallback'

export function get(key: string): any {
  const localStorageItem = storage.getItem(`explorer.${key}`)

  if (localStorageItem) {
    try {
      return JSON.parse(localStorageItem)
    } catch (e) {
      // continue regardless of error
    }
  }

  return {}
}

export function set(key: string, obj: any): void {
  storage.setItem(`explorer.${key}`, JSON.stringify(obj))
}
