import { generateSAID } from './helpersV2'

describe('generateSAID test', () => {
  it('should the hashed result of the abi encoded service agreement', () => {
    const sa = {
      initiators: [{ type: 'execagreement' }],
      tasks: [
        {
          type: 'HttpGet',
          params: { get: 'https://bitstamp.net/api/ticker/' },
        },
        { type: 'JsonParse', params: { path: ['last'] } },
        { type: 'EthBytes32' },
        {
          type: 'EthTx',
          params: {
            address: '0x356a04bce728ba4c62a30294a55e6a8600a320b3',
            functionSelector: '0x609ff1bd',
          },
        },
      ],
      payment: '1000000000000000000',
      expiration: 300,
      oracles: ['0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f'],
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
