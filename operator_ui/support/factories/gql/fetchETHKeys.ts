import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

// buildETHKey builds a eth key for the FetchETHKeys query.
export function buildETHKey(
  overrides?: Partial<EthKeysPayload_ResultsFields>,
): EthKeysPayload_ResultsFields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return {
    __typename: 'EthKey',
    address: '0x0000000000000000000000000000000000000001',
    chain: {
      id: '42',
    },
    createdAt: minuteAgo,
    ethBalance: '0.100000000000000000',
    isFunding: false,
    linkBalance: '1000000000000000000',
    ...overrides,
  }
}

// buildETHKeys builds a list of eth keys.
export function buildETHKeys(): ReadonlyArray<EthKeysPayload_ResultsFields> {
  return [
    buildETHKey(),
    buildETHKey({
      address: '0x0000000000000000000000000000000000000002',
    }),
  ]
}
