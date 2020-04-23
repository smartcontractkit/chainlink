import storage from 'local-storage-fallback'

const ADMIN_ALLOWED = 'explorer.adminAllowed'

export function getAdminAllowed(): boolean {
  const value: string | null = storage.getItem(ADMIN_ALLOWED)

  if (value) {
    try {
      return JSON.parse(value)
    } catch (e) {
      // Fall through
      console.error(
        'could not parse local storage key: %o, value: %o, error: %o',
        ADMIN_ALLOWED,
        value,
        e,
      )
    }
  }

  return false
}

export function setAdminAllowed(allowed: boolean): void {
  storage.setItem(ADMIN_ALLOWED, JSON.stringify(allowed))
}
