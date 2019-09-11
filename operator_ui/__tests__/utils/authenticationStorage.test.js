import { get, set } from 'utils/authenticationStorage'

describe('utils/authenticationStorage', () => {
  beforeEach(() => {
    global.localStorage.clear()
  })

  describe('get', () => {
    it('returns a JS object for JSON stored as "chainlink.authentication" in localStorage', () => {
      global.localStorage.setItem('chainlink.authentication', '{"foo":"FOO"}')
      expect(get()).toEqual({ foo: 'FOO' })
    })
  })

  describe('set', () => {
    it('saves the JS object as JSON under the key "chainlink.authentication" in localStorage', () => {
      set({ foo: 'FOO' })
      expect(global.localStorage.getItem('chainlink.authentication')).toEqual(
        '{"foo":"FOO"}',
      )
    })
  })
})
