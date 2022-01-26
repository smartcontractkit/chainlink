import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

// buildChains builds a chain for the FetchChains query.
export function buildChain(
  overrides?: Partial<ChainsPayload_ResultsFields>,
): ChainsPayload_ResultsFields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return {
    __typename: 'Chain',
    id: '5',
    enabled: true,
    createdAt: minuteAgo,
    ...overrides,
  }
}

// buildsChains builds a list of chains.
export function buildChains(): ReadonlyArray<ChainsPayload_ResultsFields> {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return [
    buildChain({
      id: '5',
      enabled: true,
      createdAt: minuteAgo,
    }),
    buildChain({
      id: '42',
      enabled: true,
      createdAt: minuteAgo,
    }),
  ]
}
