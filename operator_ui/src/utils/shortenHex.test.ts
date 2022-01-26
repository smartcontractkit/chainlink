import { shortenHex } from './shortenHex'

describe('shortenHex', () => {
  it('shortens the hex with default values', () => {
    expect(shortenHex('0x123456789abcdef')).toEqual('0x1234...cdef')
  })

  it('shortens the hex with custom values', () => {
    expect(shortenHex('0x123456789abcdef', { start: 2, end: 2 })).toEqual(
      '0x...ef',
    )
  })
})
