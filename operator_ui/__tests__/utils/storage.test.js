import { get, set } from 'utils/storage'

describe('utils/storage', () => {
  beforeEach(() => {
    global.localStorage.clear()
  })

  describe('get', () => {
    it('returns a JS object for JSON keyed under "chainlink." in localStorage', () => {
      global.localStorage.setItem('chainlink.foo', '{"foo":"FOO"}')
      expect(get('foo')).toEqual({ foo: 'FOO' })
    })

    it('returns an empty JS object when not valid JSON', () => {
      global.localStorage.setItem('chainlink.foo', '{"foo"}')
      expect(get('foo')).toEqual({})
    })

    it('returns an empty JS object when the key does not exist', () => {
      expect(get('foo')).toEqual({})
    })
  })

  describe('set', () => {
    it('saves the JS object as JSON keyed under "chainlink." in localStorage', () => {
      set('foo', { foo: 'FOO' })
      expect(global.localStorage.getItem('chainlink.foo')).toEqual(
        '{"foo":"FOO"}',
      )
    })
  })
})
