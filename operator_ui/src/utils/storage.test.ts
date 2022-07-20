import { getAuthentication, setAuthentication } from '../../src/utils/storage'

describe('utils/storage', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  describe('getAuthentication', () => {
    it('returns a JS object for JSON stored as "chainlink.authentication" in localStorage', () => {
      localStorage.setItem('chainlink.authentication', '{"allowed":true}')
      expect(getAuthentication()).toEqual({ allowed: true })
    })
  })

  describe('setAuthentication', () => {
    it('saves the JS object as JSON under the key "chainlink.authentication" in localStorage', () => {
      setAuthentication({ allowed: true })
      expect(localStorage.getItem('chainlink.authentication')).toEqual(
        '{"allowed":true}',
      )
    })
  })
})
