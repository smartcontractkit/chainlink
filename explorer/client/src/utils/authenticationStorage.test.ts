import { setAdminAllowed, getAdminAllowed } from './authenticationStorage'

describe('authenticationStorage', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  describe('setAdminAllowed', () => {
    it('saves the value to local storage', () => {
      setAdminAllowed(true)
      expect(localStorage.getItem('explorer.adminAllowed')).toEqual('true')
    })
  })

  describe('getAdminAllowed', () => {
    it('returns the value from local storage as a boolean', () => {
      expect(getAdminAllowed()).toEqual(false)

      localStorage.setItem('explorer.adminAllowed', 'true')
      expect(getAdminAllowed()).toEqual(true)
    })
  })
})
