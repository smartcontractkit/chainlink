import storage from 'local-storage-fallback'

const ADMIN_ALLOWED = 'explorer.adminAllowed'

export function getAdminAllowed(): boolean {
  const adminAllowed: string | null = storage.getItem(ADMIN_ALLOWED)

  if (adminAllowed) {
    try {
      return JSON.parse(adminAllowed)
    } catch {
      // Fall through
    }
  }

  return false
}

export function setAdminAllowed(allowed: boolean): void {
  storage.setItem(ADMIN_ALLOWED, JSON.stringify(allowed))
}
