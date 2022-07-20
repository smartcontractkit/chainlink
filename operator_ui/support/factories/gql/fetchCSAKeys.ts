// buildCSAKey builds a CSA Key for the FetchCSAKeys query.
export function buildCSAKey(
  overrides?: Partial<CsaKeysPayload_ResultsFields>,
): CsaKeysPayload_ResultsFields {
  return {
    __typename: 'CSAKey',
    id: 'aa67b61969793d51a3008cffba147bf57f1c89c423e32ce93ec9471d21e4231d',
    publicKey:
      'aa67b61969793d51a3008cffba147bf57f1c89c423e32ce93ec9471d21e4231d',
    ...overrides,
  }
}

// buildCSAKeys builds a list of csa keys.
export function buildCSAKeys(): ReadonlyArray<CsaKeysPayload_ResultsFields> {
  return [
    buildCSAKey({
      id: 'aa67b61969793d51a3008cffba147bf57f1c89c423e32ce93ec9471d21e4231d',
      publicKey:
        'aa67b61969793d51a3008cffba147bf57f1c89c423e32ce93ec9471d21e4231d',
    }),
    buildCSAKey({
      id: 'e09c2e1444322d91cfb9b8576ce5895e54dc5caef37c5aff4accca9272412f5b',
      publicKey:
        'e09c2e1444322d91cfb9b8576ce5895e54dc5caef37c5aff4accca9272412f5b',
    }),
  ]
}
