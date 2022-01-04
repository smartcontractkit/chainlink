// buildEthTx builds a eth transactions for the FetchEthTransactions query.
export function buildEthTx(
  overrides?: Partial<EthTransactionsPayload_ResultsFields>,
): EthTransactionsPayload_ResultsFields {
  return {
    __typename: 'EthTransaction',
    chain: {
      id: '42',
    },
    from: '0x0000000000000000000000000000000000000001',
    hash: '0x1111111111111111',
    to: '0x0000000000000000000000000000000000000002',
    nonce: '0',
    sentAt: '1000',
    ...overrides,
  }
}

// buildEthTxs builds a list of eth keys.
export function buildEthTxs(): ReadonlyArray<EthTransactionsPayload_ResultsFields> {
  return [
    buildEthTx(),
    buildEthTx({
      from: '0x0000000000000000000000000000000000000003',
      hash: '0x2222222222222222',
      to: '0x0000000000000000000000000000000000000004',
      nonce: '1',
      sentAt: '1001',
    }),
  ]
}
