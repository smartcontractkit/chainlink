// buildEthTx builds a eth transaction for the FetchEthTransaction query.
export function buildEthTx(
  overrides?: Partial<EthTransactionPayloadFields>,
): EthTransactionPayloadFields {
  return {
    __typename: 'EthTransaction',
    chain: {
      id: '42',
    },
    data: '0x',
    from: '0x0000000000000000000000000000000000000001',
    gasLimit: '21000',
    gasPrice: '2500000008',
    hash: '0x1111111111111111',
    hex: '0xf',
    nonce: '0',
    sentAt: '1000',
    state: 'confirmed',
    to: '0x0000000000000000000000000000000000000002',
    value: '0.020000000000000000',
    ...overrides,
  }
}
