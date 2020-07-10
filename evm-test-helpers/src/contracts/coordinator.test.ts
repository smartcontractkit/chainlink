import { generateSAID, ServiceAgreement } from './coordinator'

describe('generateSAID test', () => {
  it('should return the hashed result of the abi encoded service agreement', () => {
    const sa: ServiceAgreement = {
      payment: '1000000000000000000',
      expiration: 300,
      oracles: ['0x7fc66c61f88A61DFB670627cA715Fe808057123e'],
      endAt: Math.round(new Date('2020-10-19T22:17:19Z').getTime() / 1000),
      aggregator: '0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF',
      aggInitiateJobSelector: '0xe16c4c94',
      aggFulfillSelector: '0x9760168f',
      requestDigest:
        '0xbadc0de5badc0de5badc0de5badc0de5badc0de5badc0de5badc0de5badc0de5',
    }

    const expected =
      '0x31e6113ed6267498e525eb904421e8d2e2a90289553334c692f07505c8c059a6'
    const actual = generateSAID(sa)
    expect(actual).toEqual(expected)
  })
})
