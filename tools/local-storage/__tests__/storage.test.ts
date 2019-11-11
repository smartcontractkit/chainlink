import storage from 'local-storage-fallback'
import { get, set, getJson, setJson } from '../src/storage'

beforeEach(() => {
  storage.clear()
})

describe('get', () => {
  it('returns a string keyed under "chainlink." in localStorage', () => {
    storage.setItem('chainlink.foo', 'FOO')
    expect(get('foo')).toEqual('FOO')
  })

  it('returns null when the key does not exist', () => {
    expect(get('foo')).toEqual(null)
  })
})

describe('getJson', () => {
  it('returns a JS object for JSON keyed under "chainlink." in localStorage', () => {
    storage.setItem('chainlink.foo', '{"foo":"FOO"}')
    expect(getJson('foo')).toEqual({ foo: 'FOO' })
  })

  it('returns an empty JS object when it retrieves invalid JSON from storage', () => {
    storage.setItem('chainlink.foo', '{"foo"}')
    expect(getJson('foo')).toEqual({})
  })

  it('returns an empty JS object when the key does not exist', () => {
    expect(getJson('foo')).toEqual({})
  })
})

describe('set', () => {
  it('saves the string keyed under "chainlink." in localStorage', () => {
    set('foo', 'FOO')

    const stored = storage.getItem('chainlink.foo')
    expect(stored).toEqual('FOO')
  })
})

describe('setJson', () => {
  it('saves the JS object as JSON keyed under "chainlink." in localStorage', () => {
    setJson('foo', { foo: 'FOO' })

    const stored = storage.getItem('chainlink.foo')
    expect(stored).toEqual('{"foo":"FOO"}')
  })
})
