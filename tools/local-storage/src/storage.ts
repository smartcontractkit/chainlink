import storage from 'local-storage-fallback'

export function get(key: string): string | null {
  return storage.getItem(`chainlink.${key}`)
}

export function set(key: string, val: string): void {
  storage.setItem(`chainlink.${key}`, val)
}

export function remove(key: string): void {
  storage.removeItem(`chainlink.${key}`)
}

export function getJson(key: string): any {
  const stored = get(key)
  const obj = {}

  if (stored) {
    try {
      return JSON.parse(stored)
    } catch (e) {
      // continue regardless of error
    }
  }

  return obj
}

export function setJson(key: string, obj: any): void {
  const json = JSON.stringify(obj)
  set(key, json)
}
