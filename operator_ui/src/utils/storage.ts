import * as storage from 'utils/local-storage'

const PERSIST_URL = 'persistURL'

export function getPersistUrl(): string {
  return storage.get(PERSIST_URL) || ''
}

export function setPersistUrl(url: string): void {
  storage.set(PERSIST_URL, url)
}

export interface Auth {
  allowed?: boolean
}

export function getAuthentication(): Auth {
  return storage.getJson('authentication')
}

export function setAuthentication(auth: Auth): void {
  storage.setJson('authentication', auth)
}
