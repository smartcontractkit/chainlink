// buildP2PKey builds a p2p key for the FetchP2PKeys query.
export function buildP2PKey(
  overrides?: Partial<P2PKeysPayload_ResultsFields>,
): P2PKeysPayload_ResultsFields {
  return {
    __typename: 'P2PKey',
    id: '12D3KooWQTF8qHapWg89jsVucDvZittNUAGdkBph8ZgTuMAq7Ftk',
    peerID: 'p2p_12D3KooWQTF8qHapWg89jsVucDvZittNUAGdkBph8ZgTuMAq7Ftk',
    publicKey:
      'd976279be41a66b1192c1ba065d8bc6ba95a3777271009de99de177ce559fd41',
    ...overrides,
  }
}

// buildP2PKeys builds a list of p2p keys.
export function buildP2PKeys(): ReadonlyArray<P2PKeysPayload_ResultsFields> {
  return [
    buildP2PKey(),
    buildP2PKey({
      id: '12D3KooWNkPvkVkT3tRB179fjxdudwW6JRf4EZs8gJM6sDXPazEy',
      peerID: 'p2p_12D3KooWNkPvkVkT3tRB179fjxdudwW6JRf4EZs8gJM6sDXPazEy',
      publicKey:
        'c023986ca3ad70b4f06893d0a9b5b0a338578160dce173aacfb159dda3a54876',
    }),
  ]
}
