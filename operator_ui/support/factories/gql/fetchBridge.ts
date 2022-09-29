// buildBridgePayloadFields builds the bridge fields.
export function buildBridgePayloadFields(
  overrides?: Partial<BridgePayload_Fields>,
): BridgePayload_Fields {
  return {
    __typename: 'Bridge',
    id: 'bridge-api',
    name: 'bridge-api',
    url: 'http://bridge.com',
    confirmations: 1,
    minimumContractPayment: '0',
    outgoingToken: 'outgoing1',
    ...overrides,
  }
}
