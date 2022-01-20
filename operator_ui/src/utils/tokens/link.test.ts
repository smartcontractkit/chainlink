import { fromJuels } from './link'

describe('fromJuels', () => {
  it('converts juels to LINK', () => {
    expect(fromJuels('1010000000000000001')).toEqual('1.01000000')
  })
})
