// buildBridge builds a Bridge for the FetchBridges query.
export function buildBridge(
  overrides?: Partial<BridgesPayload_ResultsFields>,
): BridgesPayload_ResultsFields {
  return {
    __typename: 'Bridge',
    id: 'bridge-api',
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
      id: 'bridge-api1',
      name: 'bridge-api1',
      url: 'http://bridge1.com',
      confirmations: 1,
      minimumContractPayment: '100',
    }),
    buildBridge({
      id: 'bridge-api2',
      name: 'bridge-api2',
      url: 'http://bridge2.com',
      confirmations: 2,
      minimumContractPayment: '200',
    }),
  ]
}
