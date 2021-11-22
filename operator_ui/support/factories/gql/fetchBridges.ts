// buildFeedsManager builds a feeds manager for the FetchFeedsManagers query.
export function buildBridge(
  overrides?: Partial<BridgesPayload_ResultsFields>,
): BridgesPayload_ResultsFields {
  return {
    __typename: 'Bridge',
    name: 'bridge-api',
    url: 'http://bridge.com',
    confirmations: 1,
    minimumContractPayment: '0',
    ...overrides,
  }
}

// buildsBridges builds a list of bridges.
export function buildBridges(): ReadonlyArray<BridgesPayload_ResultsFields> {
  return [
    buildBridge({
      name: 'bridge-api1',
      url: 'http://bridge1.com',
      confirmations: 1,
      minimumContractPayment: '100',
    }),
    buildBridge({
      name: 'bridge-api2',
      url: 'http://bridge2.com',
      confirmations: 2,
      minimumContractPayment: '200',
    }),
  ]
}
