import { Networks } from './networks'

export function networkName(networkId: Networks) {
  switch (networkId) {
    case Networks.MAINNET:
      return 'mainnet'
    case Networks.ROPSTEN:
      return 'ropsten'
    case Networks.KOVAN:
      return 'kovan'
    default:
      return 'mainnet'
  }
}
