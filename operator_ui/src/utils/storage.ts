import * as storage from '@chainlink/local-storage'

const PERSIST_URL = 'persistURL'
const PERSIST_JOB_SPEC = 'persistJobSpec'

export function getPersistUrl(): string {
  return storage.get(PERSIST_URL) || ''
}

export function setPersistUrl(url: string): void {
  storage.set(PERSIST_URL, url)
}

export function setPersistJobSpec(spec: string): void {
  storage.set(PERSIST_JOB_SPEC, spec)
}

export function getPersistJobSpec(): string {
  return storage.get(PERSIST_JOB_SPEC) || ''
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
