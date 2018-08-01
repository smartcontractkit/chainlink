import { get, set } from 'utils/sessionStorage'

describe('utils/sessionStorage', () => {
  beforeEach(() => {
    global.localStorage.clear()
  })

  describe('get', () => {
    it('returns a JS object for JSON stored as "chainlink.session" in localStorage', () => {
      global.localStorage.setItem('chainlink.session', '{"foo":"FOO"}')
      expect(get()).toEqual({foo: 'FOO'})
    })
  })

  describe('set', () => {
    it('saves the JS object as JSON under the key "chainlink.session" in localStorage', () => {
      set({foo: 'FOO'})
      expect(global.localStorage.getItem('chainlink.session')).toEqual('{"foo":"FOO"}')
    })
  })
})
