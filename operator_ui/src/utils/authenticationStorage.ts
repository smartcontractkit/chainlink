import * as storage from '@chainlink/local-storage'

interface Auth {
  allowed?: boolean
}

export function get(): Auth {
  return storage.getJson('authentication')
}

export function set(auth: Auth): void {
  storage.setJson('authentication', auth)
}
