import { Networks } from './networks'

export function networkName(networkId: Networks) {
  switch (networkId) {
    case Networks.MAINNET:
      return 'mainnet'
    case Networks.ROPSTEN:
      return 'ropsten'
    default:
      return 'mainnet'
  }
}
