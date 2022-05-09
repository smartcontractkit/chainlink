import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

// buildNodePayloadFields builds the node fields.
export function buildNodePayloadFields(
  overrides?: Partial<NodePayload_Fields>,
): NodePayload_Fields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return {
    __typename: 'Node',
    id: '1',
    name: 'node1',
    httpURL: 'https://node1.com',
    wsURL: 'wss://node1.com',
    createdAt: minuteAgo,
    chain: {
      id: '42',
    },
    state: '',
    ...overrides,
  }
}
