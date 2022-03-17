import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

// buildNode builds a node for the FetchNodes query.
export function buildNode(
  overrides?: Partial<NodesPayload_ResultsFields>,
): NodesPayload_ResultsFields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return {
    __typename: 'Node',
    id: '1',
    name: 'node1',
    chain: {
      id: '42',
    },
    createdAt: minuteAgo,
    state: '',
    ...overrides,
  }
}

// buildNodes builds a list of nodes.
export function buildNodes(): ReadonlyArray<NodesPayload_ResultsFields> {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return [
    buildNode({
      id: '1',
      name: 'node1',
      chain: {
        id: '42',
      },
      createdAt: minuteAgo,
    }),
    buildNode({
      id: '2',
      name: 'node2',
      chain: {
        id: '5',
      },
      createdAt: minuteAgo,
    }),
  ]
}
