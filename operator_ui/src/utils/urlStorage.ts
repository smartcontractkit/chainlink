import * as storage from '@chainlink/local-storage'

const PERSIST_URL = 'persistURL'

export function getPersistUrl(): string {
  return storage.get(PERSIST_URL) || ''
}

export function setPersistUrl(url: string): void {
  storage.set(PERSIST_URL, url)
}
